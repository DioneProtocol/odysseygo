// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package poll

import (
	"fmt"
	"strings"
	"time"

	"github.com/ava-labs/avalanchego/utils/linkedhashmap"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/utils/metric"
)

type pollHolder interface {
	GetPoll() Poll
	StartTime() time.Time
}

type poll struct {
	Poll
	start time.Time
}

func (p poll) GetPoll() Poll {
	return p
}

func (p poll) StartTime() time.Time {
	return p.start
}

type set struct {
	log      logging.Logger
	numPolls prometheus.Gauge
	durPolls prometheus.Histogram
	factory  Factory
	// maps requestID -> poll
	polls linkedhashmap.LinkedHashmap
}

// NewSet returns a new empty set of polls
func NewSet(
	factory Factory,
	log logging.Logger,
	namespace string,
	registerer prometheus.Registerer,
) Set {
	numPolls := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "polls",
		Help:      "Number of pending network polls",
	})
	if err := registerer.Register(numPolls); err != nil {
		log.Error("failed to register polls statistics due to %s", err)
	}

	durPolls := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "poll_duration",
		Help:      "Length of time the poll existed in milliseconds",
		Buckets:   metric.MillisecondsBuckets,
	})
	if err := registerer.Register(durPolls); err != nil {
		log.Error("failed to register poll_duration statistics due to %s", err)
	}

	return &set{
		log:      log,
		numPolls: numPolls,
		durPolls: durPolls,
		factory:  factory,
		polls:    linkedhashmap.New(),
	}
}

// Add to the current set of polls
// Returns true if the poll was registered correctly and the network sample
//         should be made.
func (s *set) Add(requestID uint32, vdrs ids.ShortBag) bool {
	if _, exists := s.polls.Get(requestID); exists {
		s.log.Debug("dropping poll due to duplicated requestID: %d", requestID)
		return false
	}

	s.log.Verbo("creating poll with requestID %d and validators %s",
		requestID,
		&vdrs)

	s.polls.Put(requestID, poll{
		Poll:  s.factory.New(vdrs), // create the new poll
		start: time.Now(),
	})
	s.numPolls.Inc() // increase the metrics
	return true
}

// Vote registers the connections response to a query for [id]. If there was no
// query, or the response has already be registered, nothing is performed.
func (s *set) Vote(requestID uint32, vdr ids.ShortID, vote ids.ID) ([]ids.Bag, bool) {
	pollHolderIntf, exists := s.polls.Get(requestID)
	if !exists {
		s.log.Verbo("dropping vote from %s to an unknown poll with requestID: %d",
			vdr,
			requestID)
		return []ids.Bag{}, false
	}

	pollHolder := pollHolderIntf.(pollHolder)
	poll := pollHolder.GetPoll()

	s.log.Verbo("processing vote from %s in the poll with requestID: %d with the vote %s",
		vdr,
		requestID,
		vote)

	poll.Vote(vdr, vote)
	if !poll.Finished() {
		return []ids.Bag{}, false
	}

	s.log.Verbo("poll with requestID %d finished as %s", requestID, poll)

	s.polls.Delete(requestID) // remove the poll from the current set
	s.durPolls.Observe(float64(time.Since(pollHolder.StartTime()).Milliseconds()))
	s.numPolls.Dec() // decrease the metrics

	result := poll.Result()
	return []ids.Bag{result}, true
}

// Drop registers the connections response to a query for [id]. If there was no
// query, or the response has already be registered, nothing is performed.
func (s *set) Drop(requestID uint32, vdr ids.ShortID) ([]ids.Bag, bool) {
	pollHolderIntf, exists := s.polls.Get(requestID)
	if !exists {
		s.log.Verbo("dropping vote from %s to an unknown poll with requestID: %d",
			vdr,
			requestID)
		return []ids.Bag{}, false
	}

	s.log.Verbo("processing dropped vote from %s in the poll with requestID: %d",
		vdr,
		requestID)

	pollHolder := pollHolderIntf.(pollHolder)
	poll := pollHolder.GetPoll()

	poll.Drop(vdr)
	if !poll.Finished() {
		return []ids.Bag{}, false
	}

	s.log.Verbo("poll with requestID %d finished as %s", requestID, poll)

	s.polls.Delete(requestID) // remove the poll from the current set
	s.durPolls.Observe(float64(time.Since(pollHolder.StartTime()).Milliseconds()))
	s.numPolls.Dec() // decrease the metrics
	return []ids.Bag{poll.Result()}, true
}

// Len returns the number of outstanding polls
func (s *set) Len() int { return s.polls.Len() }

func (s *set) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("current polls: (Size = %d)", s.polls.Len()))
	iter := s.polls.NewIterator()
	for iter.Next() {
		requestID := iter.Key()
		poll := iter.Value().(Poll)
		sb.WriteString(fmt.Sprintf("\n    %d: %s", requestID, poll.PrefixedString("    ")))
	}
	return sb.String()
}
