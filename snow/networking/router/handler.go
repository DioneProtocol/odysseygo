// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package router

import (
	"math"
	"sync"
	"time"

	"github.com/ava-labs/gecko/utils/timer"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/snow"
	"github.com/ava-labs/gecko/snow/engine/common"
	"github.com/ava-labs/gecko/snow/networking/timeout"
	"github.com/ava-labs/gecko/snow/validators"
)

const (
	defaultStakerPortion float64 = 0.2
)

// Requirement: A set of nodes spamming messages (potentially costly) shouldn't
//              impact other node's queries.

// Requirement: The staked validators should be able to maintain liveness, even
//              if that requires sacrificing liveness of the non-staked
//              validators. This ensures the network keeps moving forwards.

// Idea: There is 1 second of cpu time per second, divide that out into the
//       stakers (potentially by staking weight).

// Idea: Non-stakers are treated as if they have the same identity.

// Idea: Beacons should receive special treatment.

// Problem: Our queues need to have bounded size, so we need to drop messages.
//          When should we be dropping messages? If we only drop based on the
//          queue being full, then a peer can spam messages to fill the queue
//          causing other peers' messages to drop.
// Answer: Drop messages if the peer has too many oustanding messages. (Could be
//         weighted by the size of the queue + stake amount)

// Problem: How should we prioritize peers? If we are already picking which
//          level of the queue the peer's messages are going into, then the
//          internal queue can just be FIFO.
// Answer: Somehow we track the cpu time of the peer (WMA). Based on that, we
//         place the message into a corresponding bucket. When pulling the
//         message from the bucket, we check to see if the message should be
//         moved to a lower bucket. If so move the message to the lower queue
//         and process the next message.

// Structure:
//  [000%-050%] P0: Chan msg    size = 1024    CPU time per iteration = 200ms
//  [050%-075%] P1: Chan msg    size = 1024    CPU time per iteration = 150ms
//  [075%-100%] P2: Chan msg    size = 1024    CPU time per iteration = 100ms
//  [100%-125%] P3: Chan msg    size = 1024    CPU time per iteration = 050ms
//  [125%-INF%] P4: Chan msg    size = 1024    CPU time per iteration = 001ms

// 20% resources for stakers. = RE_s
// 80% resources for non-stakers. = RE_n

// Each peer is going to have calculated their expected CPU utilization.
//
// E[Staker CPU Utilization] = RE_s * weight + RE_n / NumPeers
// E[Non-Staker CPU Utilization] = RE_n / NumPeers

// Each peer is going to have calculated their max number of outstanding
// messages.
//
// Max[Staker Messages] = (RE_s * weight + RE_n / NumPeers) * MaxMessages
// Max[Non-Staker Messages] = (RE_n / NumPeers) * MaxMessages

// Problem: If everyone is part of the P0 queue, except for a couple byzantine
//          nodes. Then the byzantine nodes can take up 80% of the CPU.

// Vars to tune:
// - % reserved for stakers.
// - CPU time per buckets
// - % range of buckets
// - number of buckets
// - size of buckets
// - how to track CPU utilization of a peer
// - "MaxMessages"

// Handler passes incoming messages from the network to the consensus engine
// (Actually, it receives the incoming messages from a ChainRouter, but same difference)
type Handler struct {
	metrics

	validators validators.Set

	// This is the channel of messages to process
	reliableMsgsSema chan struct{}
	reliableMsgsLock sync.Mutex
	reliableMsgs     []message
	closed           chan struct{}
	msgChan          <-chan common.Message

	clock              timer.Clock
	dropMessageTimeout time.Duration

	serviceQueue messageQueue
	msgSema      <-chan struct{}

	ctx    *snow.Context
	engine common.Engine

	toClose func()
	closing bool
}

// Initialize this consensus handler
// engine must be initialized before initializing the handler
func (h *Handler) Initialize(
	engine common.Engine,
	validators validators.Set,
	msgChan <-chan common.Message,
	bufferSize int,
	namespace string,
	metrics prometheus.Registerer,
) {
	h.metrics.Initialize(namespace, metrics)
	h.reliableMsgsSema = make(chan struct{}, 1)
	h.closed = make(chan struct{})
	h.msgChan = msgChan
	h.dropMessageTimeout = timeout.DefaultRequestTimeout

	h.ctx = engine.Context()

	// Defines the maximum current percentage of expected CPU utilization for
	// a message to be placed in the queue at the corresponding index
	consumptionRanges := []float64{
		0.5,
		0.75,
		1.5,
		math.MaxFloat64,
	}

	cpuInterval := float64(defaultCPUInterval)
	// Defines the percentage of CPU time allotted to processing messages
	// from the bucket at the corresponding index.
	consumptionAllotments := []float64{
		cpuInterval * 0.25,
		cpuInterval * 0.25,
		cpuInterval * 0.25,
		cpuInterval * 0.25,
	}

	h.serviceQueue, h.msgSema = newMultiLevelQueue(
		validators,
		h.ctx.Log,
		&h.metrics,
		consumptionRanges,
		consumptionAllotments,
		bufferSize,
		cpuInterval,
		defaultStakerPortion,
	)
	h.engine = engine
	h.validators = validators
}

// Context of this Handler
func (h *Handler) Context() *snow.Context { return h.engine.Context() }

// Engine returns the engine this handler dispatches to
func (h *Handler) Engine() common.Engine { return h.engine }

// SetEngine sets the engine for this handler to dispatch to
func (h *Handler) SetEngine(engine common.Engine) { h.engine = engine }

// Dispatch waits for incoming messages from the network
// and, when they arrive, sends them to the consensus engine
func (h *Handler) Dispatch() {
	defer h.shutdownDispatch()

	for {
		select {
		case _, ok := <-h.msgSema:
			if !ok {
				// the msgSema channel has been closed, so this dispatcher should exit
				return
			}

			msg, err := h.serviceQueue.PopMessage()
			if err != nil {
				h.ctx.Log.Warn("Could not pop messsage from service queue")
				continue
			}
			if !msg.deadline.IsZero() && h.clock.Time().After(msg.deadline) {
				h.ctx.Log.Verbo("Dropping message due to likely timeout: %s", msg)
				h.metrics.dropped.Inc()
				h.metrics.expired.Inc()
				continue
			}

			h.dispatchMsg(msg)
		case <-h.reliableMsgsSema:
			// get all the reliable messages
			h.reliableMsgsLock.Lock()
			msgs := h.reliableMsgs
			h.reliableMsgs = nil
			h.reliableMsgsLock.Unlock()

			// fire all the reliable messages
			for _, msg := range msgs {
				h.metrics.pending.Dec()
				h.dispatchMsg(msg)
			}
		case msg := <-h.msgChan:
			// handle a message from the VM
			h.dispatchMsg(message{messageType: notifyMsg, notification: msg})
		}

		if h.closing {
			return
		}
	}
}

// Dispatch a message to the consensus engine.
func (h *Handler) dispatchMsg(msg message) {
	if h.closing {
		h.ctx.Log.Debug("dropping message due to closing:\n%s", msg)
		h.metrics.dropped.Inc()
		return
	}

	startTime := h.clock.Time()

	h.ctx.Lock.Lock()
	defer h.ctx.Lock.Unlock()

	h.ctx.Log.Debug("Forwarding message to consensus: %s", msg)
	var (
		err error
	)
	switch msg.messageType {
	case notifyMsg:
		err = h.engine.Notify(msg.notification)
		h.notify.Observe(float64(h.clock.Time().Sub(startTime)))
	case gossipMsg:
		err = h.engine.Gossip()
		h.gossip.Observe(float64(h.clock.Time().Sub(startTime)))
	default:
		err = h.handleValidatorMsg(msg, startTime)
	}

	if err != nil {
		h.ctx.Log.Fatal("forcing chain to shutdown due to: %s", err)
		h.closing = true
	}
}

// GetAcceptedFrontier passes a GetAcceptedFrontier message received from the
// network to the consensus engine.
func (h *Handler) GetAcceptedFrontier(validatorID ids.ShortID, requestID uint32) bool {
	currentTime := h.clock.Time()
	return h.sendMsg(message{
		messageType: getAcceptedFrontierMsg,
		validatorID: validatorID,
		requestID:   requestID,
		received:    currentTime,
		deadline:    currentTime.Add(h.dropMessageTimeout),
	})
}

// AcceptedFrontier passes a AcceptedFrontier message received from the network
// to the consensus engine.
func (h *Handler) AcceptedFrontier(validatorID ids.ShortID, requestID uint32, containerIDs ids.Set) bool {
	return h.sendMsg(message{
		messageType:  acceptedFrontierMsg,
		validatorID:  validatorID,
		requestID:    requestID,
		containerIDs: containerIDs,
		received:     h.clock.Time(),
	})
}

// GetAcceptedFrontierFailed passes a GetAcceptedFrontierFailed message received
// from the network to the consensus engine.
func (h *Handler) GetAcceptedFrontierFailed(validatorID ids.ShortID, requestID uint32) {
	h.sendReliableMsg(message{
		messageType: getAcceptedFrontierFailedMsg,
		validatorID: validatorID,
		requestID:   requestID,
	})
}

// GetAccepted passes a GetAccepted message received from the
// network to the consensus engine.
func (h *Handler) GetAccepted(validatorID ids.ShortID, requestID uint32, containerIDs ids.Set) bool {
	currentTime := h.clock.Time()
	return h.sendMsg(message{
		messageType:  getAcceptedMsg,
		validatorID:  validatorID,
		requestID:    requestID,
		containerIDs: containerIDs,
		received:     currentTime,
		deadline:     currentTime.Add(h.dropMessageTimeout),
	})
}

// Accepted passes a Accepted message received from the network to the consensus
// engine.
func (h *Handler) Accepted(validatorID ids.ShortID, requestID uint32, containerIDs ids.Set) bool {
	return h.sendMsg(message{
		messageType:  acceptedMsg,
		validatorID:  validatorID,
		requestID:    requestID,
		containerIDs: containerIDs,
		received:     h.clock.Time(),
	})
}

// GetAcceptedFailed passes a GetAcceptedFailed message received from the
// network to the consensus engine.
func (h *Handler) GetAcceptedFailed(validatorID ids.ShortID, requestID uint32) {
	h.sendReliableMsg(message{
		messageType: getAcceptedFailedMsg,
		validatorID: validatorID,
		requestID:   requestID,
	})
}

// GetAncestors passes a GetAncestors message received from the network to the consensus engine.
func (h *Handler) GetAncestors(validatorID ids.ShortID, requestID uint32, containerID ids.ID) bool {
	currentTime := h.clock.Time()
	return h.sendMsg(message{
		messageType: getAncestorsMsg,
		validatorID: validatorID,
		requestID:   requestID,
		containerID: containerID,
		received:    currentTime,
		deadline:    currentTime.Add(h.dropMessageTimeout),
	})
}

// MultiPut passes a MultiPut message received from the network to the consensus engine.
func (h *Handler) MultiPut(validatorID ids.ShortID, requestID uint32, containers [][]byte) bool {
	return h.sendMsg(message{
		messageType: multiPutMsg,
		validatorID: validatorID,
		requestID:   requestID,
		containers:  containers,
		received:    h.clock.Time(),
	})
}

// GetAncestorsFailed passes a GetAncestorsFailed message to the consensus engine.
func (h *Handler) GetAncestorsFailed(validatorID ids.ShortID, requestID uint32) {
	h.sendReliableMsg(message{
		messageType: getAncestorsFailedMsg,
		validatorID: validatorID,
		requestID:   requestID,
	})
}

// Get passes a Get message received from the network to the consensus engine.
func (h *Handler) Get(validatorID ids.ShortID, requestID uint32, containerID ids.ID) bool {
	currentTime := h.clock.Time()
	return h.sendMsg(message{
		messageType: getMsg,
		validatorID: validatorID,
		requestID:   requestID,
		containerID: containerID,
		received:    currentTime,
		deadline:    currentTime.Add(h.dropMessageTimeout),
	})
}

// Put passes a Put message received from the network to the consensus engine.
func (h *Handler) Put(validatorID ids.ShortID, requestID uint32, containerID ids.ID, container []byte) bool {
	return h.sendMsg(message{
		messageType: putMsg,
		validatorID: validatorID,
		requestID:   requestID,
		containerID: containerID,
		container:   container,
		received:    h.clock.Time(),
	})
}

// GetFailed passes a GetFailed message to the consensus engine.
func (h *Handler) GetFailed(validatorID ids.ShortID, requestID uint32) {
	h.sendReliableMsg(message{
		messageType: getFailedMsg,
		validatorID: validatorID,
		requestID:   requestID,
	})
}

// PushQuery passes a PushQuery message received from the network to the consensus engine.
func (h *Handler) PushQuery(validatorID ids.ShortID, requestID uint32, blockID ids.ID, block []byte) bool {
	currentTime := h.clock.Time()
	return h.sendMsg(message{
		messageType: pushQueryMsg,
		validatorID: validatorID,
		requestID:   requestID,
		containerID: blockID,
		container:   block,
		received:    currentTime,
		deadline:    currentTime.Add(h.dropMessageTimeout),
	})
}

// PullQuery passes a PullQuery message received from the network to the consensus engine.
func (h *Handler) PullQuery(validatorID ids.ShortID, requestID uint32, blockID ids.ID) bool {
	currentTime := h.clock.Time()
	return h.sendMsg(message{
		messageType: pullQueryMsg,
		validatorID: validatorID,
		requestID:   requestID,
		containerID: blockID,
		received:    currentTime,
		deadline:    currentTime.Add(h.dropMessageTimeout),
	})
}

// Chits passes a Chits message received from the network to the consensus engine.
func (h *Handler) Chits(validatorID ids.ShortID, requestID uint32, votes ids.Set) bool {
	return h.sendMsg(message{
		messageType:  chitsMsg,
		validatorID:  validatorID,
		requestID:    requestID,
		containerIDs: votes,
		received:     h.clock.Time(),
	})
}

// QueryFailed passes a QueryFailed message received from the network to the consensus engine.
func (h *Handler) QueryFailed(validatorID ids.ShortID, requestID uint32) {
	h.sendReliableMsg(message{
		messageType: queryFailedMsg,
		validatorID: validatorID,
		requestID:   requestID,
	})
}

// Gossip passes a gossip request to the consensus engine
func (h *Handler) Gossip() {
	h.sendReliableMsg(message{
		messageType: gossipMsg,
	})
}

// Notify ...
func (h *Handler) Notify(msg common.Message) {
	h.sendReliableMsg(message{
		messageType:  notifyMsg,
		notification: msg,
	})
}

// Shutdown asynchronously shuts down the dispatcher.
// The handler should never be invoked again after calling
// Shutdown.
func (h *Handler) Shutdown() {
	h.serviceQueue.Shutdown()
}

func (h *Handler) shutdownDispatch() {
	h.ctx.Lock.Lock()
	defer h.ctx.Lock.Unlock()

	startTime := time.Now()
	if err := h.engine.Shutdown(); err != nil {
		h.ctx.Log.Error("Error while shutting down the chain: %s", err)
	}
	h.ctx.Log.Info("finished shutting down chain")
	if h.toClose != nil {
		go h.toClose()
	}
	h.closing = true
	h.shutdown.Observe(float64(time.Now().Sub(startTime)))
	close(h.closed)
}

func (h *Handler) handleValidatorMsg(msg message, startTime time.Time) error {
	var (
		err          error
		timeConsumed float64
	)
	switch msg.messageType {
	case getAcceptedFrontierMsg:
		err = h.engine.GetAcceptedFrontier(msg.validatorID, msg.requestID)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.getAcceptedFrontier.Observe(timeConsumed)
	case acceptedFrontierMsg:
		err = h.engine.AcceptedFrontier(msg.validatorID, msg.requestID, msg.containerIDs)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.acceptedFrontier.Observe(timeConsumed)
	case getAcceptedFrontierFailedMsg:
		err = h.engine.GetAcceptedFrontierFailed(msg.validatorID, msg.requestID)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.getAcceptedFrontierFailed.Observe(timeConsumed)
	case getAcceptedMsg:
		err = h.engine.GetAccepted(msg.validatorID, msg.requestID, msg.containerIDs)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.getAccepted.Observe(timeConsumed)
	case acceptedMsg:
		err = h.engine.Accepted(msg.validatorID, msg.requestID, msg.containerIDs)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.accepted.Observe(timeConsumed)
	case getAcceptedFailedMsg:
		err = h.engine.GetAcceptedFailed(msg.validatorID, msg.requestID)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.getAcceptedFailed.Observe(timeConsumed)
	case getAncestorsMsg:
		err = h.engine.GetAncestors(msg.validatorID, msg.requestID, msg.containerID)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.getAncestors.Observe(timeConsumed)
	case getAncestorsFailedMsg:
		err = h.engine.GetAncestorsFailed(msg.validatorID, msg.requestID)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.getAncestorsFailed.Observe(timeConsumed)
	case multiPutMsg:
		err = h.engine.MultiPut(msg.validatorID, msg.requestID, msg.containers)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.multiPut.Observe(timeConsumed)
	case getMsg:
		err = h.engine.Get(msg.validatorID, msg.requestID, msg.containerID)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.get.Observe(timeConsumed)
	case getFailedMsg:
		err = h.engine.GetFailed(msg.validatorID, msg.requestID)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.getFailed.Observe(timeConsumed)
	case putMsg:
		err = h.engine.Put(msg.validatorID, msg.requestID, msg.containerID, msg.container)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.put.Observe(timeConsumed)
	case pushQueryMsg:
		err = h.engine.PushQuery(msg.validatorID, msg.requestID, msg.containerID, msg.container)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.pushQuery.Observe(timeConsumed)
	case pullQueryMsg:
		err = h.engine.PullQuery(msg.validatorID, msg.requestID, msg.containerID)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.pullQuery.Observe(timeConsumed)
	case queryFailedMsg:
		err = h.engine.QueryFailed(msg.validatorID, msg.requestID)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.queryFailed.Observe(timeConsumed)
	case chitsMsg:
		err = h.engine.Chits(msg.validatorID, msg.requestID, msg.containerIDs)
		timeConsumed = float64(h.clock.Time().Sub(startTime))
		h.chits.Observe(timeConsumed)
	}

	h.serviceQueue.UtilizeCPU(msg.validatorID, timeConsumed)

	return err
}

func (h *Handler) sendMsg(msg message) bool {
	return h.serviceQueue.PushMessage(msg)
}

func (h *Handler) sendReliableMsg(msg message) {
	h.reliableMsgsLock.Lock()
	defer h.reliableMsgsLock.Unlock()

	h.metrics.pending.Inc()
	h.reliableMsgs = append(h.reliableMsgs, msg)
	select {
	case h.reliableMsgsSema <- struct{}{}:
	default:
	}
}

func (h *Handler) endInterval() { h.serviceQueue.EndInterval() }
