// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package poll

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/utils/logging"
)

type set struct {
	log      logging.Logger
	numPolls prometheus.Gauge
	factory  Factory
	polls    map[uint32]Poll
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
		log.Error("failed to register av_polls statistics due to %s", err)
	}

	return &set{
		log:      log,
		numPolls: numPolls,
		factory:  factory,
		polls:    make(map[uint32]Poll),
	}
}

// Add to the current set of polls
// Returns true if the poll was registered correctly and the network sample
//         should be made.
func (s *set) Add(requestID uint32, vdrs ids.ShortSet) bool {
	if _, exists := s.polls[requestID]; exists {
		s.log.Debug("dropping poll due to duplicated requestID: %d", requestID)
		return false
	}

	s.log.Verbo("creating poll with requestID %d and validators %s",
		requestID,
		vdrs)

	s.polls[requestID] = s.factory.New(vdrs) // create the new poll
	s.numPolls.Inc()                         // increase the metrics
	return true
}

// Vote registers the connections response to a query for [id]. If there was no
// query, or the response has already be registered, nothing is performed.
func (s *set) Vote(
	requestID uint32,
	vdr ids.ShortID,
	votes []ids.ID,
) (ids.UniqueBag, bool) {
	poll, exists := s.polls[requestID]
	if !exists {
		s.log.Verbo("dropping vote from %s to an unknown poll with requestID: %d",
			vdr,
			requestID)
		return nil, false
	}

	s.log.Verbo("processing vote from %s in the poll with requestID: %d with the votes %v",
		vdr,
		requestID,
		votes)

	poll.Vote(vdr, votes)
	if !poll.Finished() {
		return nil, false
	}

	s.log.Verbo("poll with requestID %d finished as %s", requestID, poll)

	delete(s.polls, requestID) // remove the poll from the current set
	s.numPolls.Dec()           // decrease the metrics
	return poll.Result(), true
}

// Len returns the number of outstanding polls
func (s *set) Len() int { return len(s.polls) }

func (s *set) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("current polls: (Size = %d)", len(s.polls)))
	for requestID, poll := range s.polls {
		sb.WriteString(fmt.Sprintf("\n    %d: %s", requestID, poll))
	}
	return sb.String()
}
