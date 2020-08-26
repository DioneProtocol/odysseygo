// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package snowstorm

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/snow"
	"github.com/ava-labs/gecko/snow/choices"
	"github.com/ava-labs/gecko/utils/formatting"

	sbcon "github.com/ava-labs/gecko/snow/consensus/snowball"
)

// DirectedFactory implements Factory by returning a directed struct
type DirectedFactory struct{}

// New implements Factory
func (DirectedFactory) New() Consensus { return &Directed{} }

// Directed is an implementation of a multi-color, non-transitive, snowball
// instance
type Directed struct {
	common

	// Key: Transaction ID
	// Value: Node that represents this transaction in the conflict graph
	txs map[[32]byte]*directedTx

	// Key: UTXO ID
	// Value: IDs of transactions that consume the UTXO specified in the key
	utxos map[[32]byte]ids.Set
}

type directedTx struct {
	snowball

	// pendingAccept identifies if this transaction has been marked as accepted
	// once its transitive dependencies have also been accepted
	pendingAccept bool

	// ins is the set of txIDs that this tx conflicts with that are less
	// preferred than this tx
	ins ids.Set

	// outs is the set of txIDs that this tx conflicts with that are more
	// preferred than this tx
	outs ids.Set

	// tx is the actual transaction this node represents
	tx Tx
}

// Initialize implements the Consensus interface
func (dg *Directed) Initialize(
	ctx *snow.Context,
	params sbcon.Parameters,
) error {
	dg.txs = make(map[[32]byte]*directedTx)
	dg.utxos = make(map[[32]byte]ids.Set)

	return dg.common.Initialize(ctx, params)
}

// IsVirtuous implements the Consensus interface
func (dg *Directed) IsVirtuous(tx Tx) bool {
	txID := tx.ID()
	// If the tx is currently processing, we should just return if was registed
	// as rogue or not.
	if node, exists := dg.txs[txID.Key()]; exists {
		return !node.rogue
	}

	// The tx isn't processing, so we need to check to see if it conflicts with
	// any of the other txs that are currently processing.
	for _, input := range tx.InputIDs().List() {
		if _, exists := dg.utxos[input.Key()]; exists {
			// A currently processing tx names the same input as the provided
			// tx, so the provided tx would be rogue.
			return false
		}
	}

	// This tx is virtuous as far as this consensus instance knows.
	return true
}

// Conflicts implements the Consensus interface
func (dg *Directed) Conflicts(tx Tx) ids.Set {
	conflicts := ids.Set{}
	if node, exists := dg.txs[tx.ID().Key()]; exists {
		// If the tx is currently processing, the conflicting txs is just the
		// union of the inbound conflicts and the outbound conflicts.
		conflicts.Union(node.ins)
		conflicts.Union(node.outs)
	} else {
		// If the tx isn't currently processing, the conflicting txs is the
		// union of all the txs that spend an input that this tx spends.
		for _, input := range tx.InputIDs().List() {
			if spends, exists := dg.utxos[input.Key()]; exists {
				conflicts.Union(spends)
			}
		}
	}
	return conflicts
}

// Add implements the Consensus interface
func (dg *Directed) Add(tx Tx) error {
	if dg.Issued(tx) {
		// If the tx was previously inserted, nothing should be done here.
		return nil
	}

	txID := tx.ID()
	bytes := tx.Bytes()

	// Notify the IPC socket that this tx has been issued.
	dg.ctx.DecisionDispatcher.Issue(dg.ctx.ChainID, txID, bytes)

	// Notify the metrics that this transaction was just issued.
	dg.metrics.Issued(txID)

	inputs := tx.InputIDs()

	// If this tx doesn't have any inputs, it's impossible for there to be any
	// conflicting transactions. Therefore, this transaction is treated as
	// vacuously accepted.
	if inputs.Len() == 0 {
		// Accept is called before notifying the IPC so that acceptances that
		// cause fatal errors aren't sent to an IPC peer.
		if err := tx.Accept(); err != nil {
			return err
		}

		// Notify the IPC socket that this tx has been accepted.
		dg.ctx.DecisionDispatcher.Accept(dg.ctx.ChainID, txID, bytes)

		// Notify the metrics that this transaction was just accepted.
		dg.metrics.Accepted(txID)
		return nil
	}

	txNode := &directedTx{tx: tx}

	// For each UTXO consumed by the tx:
	// * Add edges between this tx and txs that consume this UTXO
	// * Mark this tx as attempting to consume this UTXO
	for _, inputID := range inputs.List() {
		inputKey := inputID.Key()

		// Get the set of txs that are currently processing that also consume
		// this UTXO
		spenders := dg.utxos[inputKey]

		// Add all the txs that spend this UTXO to this txs conflicts that are
		// preferred over this tx. We know all these txs are preferred over
		// this tx, because this tx currently has a bias of 0 and the tie goes
		// to the tx whose bias was updated first.
		txNode.outs.Union(spenders)

		// Update txs conflicting with tx to account for its issuance
		for _, conflictID := range spenders.List() {
			conflictKey := conflictID.Key()

			// Get the node that contains this conflicting tx
			conflict := dg.txs[conflictKey]

			// This conflicting tx can't be virtuous anymore. So we remove this
			// conflicting tx from any of the virtuous sets if it was previously
			// in them.
			dg.virtuous.Remove(conflictID)
			dg.virtuousVoting.Remove(conflictID)

			// This tx should be set to rogue if it wasn't rogue before.
			conflict.rogue = true

			// This conflicting tx is preferred over the tx being inserted, as
			// described above. So we add the conflict to the inbound set.
			conflict.ins.Add(txID)
		}

		// Add this tx to list of txs consuming the current UTXO
		spenders.Add(txID)

		// Because this isn't a pointer, we should re-map the set.
		dg.utxos[inputKey] = spenders
	}

	// Mark this transaction as rogue if had any conflicts registered above
	txNode.rogue = txNode.outs.Len() != 0

	if !txNode.rogue {
		// If this tx is currently virtuous, add it to the virtuous sets
		dg.virtuous.Add(txID)
		dg.virtuousVoting.Add(txID)

		// If a tx is virtuous, it must be preferred.
		dg.preferences.Add(txID)
	}

	// Add this tx to the set of currently processing txs
	dg.txs[txID.Key()] = txNode

	// This tx can be accepted only if all the txs it depends on are also
	// accepted. If any txs that this tx depends on are rejected, reject it.
	toReject := &rejector{
		g:    dg,
		errs: &dg.errs,
		txID: txID,
	}

	// Register all of this txs dependencies as possibilities to reject this tx.
	for _, dependency := range tx.Dependencies() {
		if dependency.Status() != choices.Accepted {
			// If the dependency isn't accepted, then it must be processing. So,
			// this tx should be rejected if any of these processing txs are
			// rejected. Note that the dependencies can't be rejected, because
			// it is assumped that this tx is currently considered valid.
			toReject.deps.Add(dependency.ID())
		}
	}

	// Register these dependencies
	dg.pendingReject.Register(toReject)

	// Registering the rejector can't result in an error, so we can safely
	// return nil here.
	return nil
}

// Issued implements the Consensus interface
func (dg *Directed) Issued(tx Tx) bool {
	// If the tx is either Accepted or Rejected, then it must have been issued
	// previously.
	if tx.Status().Decided() {
		return true
	}

	// If the tx is currently processing, then it must have been issued.
	_, ok := dg.txs[tx.ID().Key()]
	return ok
}

// RecordPoll implements the Consensus interface
func (dg *Directed) RecordPoll(votes ids.Bag) (bool, error) {
	// Increase the vote ID. This is updated here and is used to reset the
	// confidence values of transactions lazily.
	dg.currentVote++

	// Changed tracks if the Avalanche instance needs to recompute its
	// frontiers. Frontiers only need to be recalculated if preferences change
	// or if a tx was accepted.
	changed := false

	// We only want to iterate over txs that received alpha votes
	votes.SetThreshold(dg.params.Alpha)
	// Get the set of IDs that meet this alpha threshold
	metThreshold := votes.Threshold()
	for _, txID := range metThreshold.List() {
		txKey := txID.Key()

		// Get the node this tx represents
		txNode, exist := dg.txs[txKey]
		if !exist {
			// This tx may have already been accepted because of tx
			// dependencies. If this is the case, we can just drop the vote.
			continue
		}

		txNode.RecordSuccessfulPoll(dg.currentVote)

		dg.ctx.Log.Verbo("Updated TxID=%s to have consensus state=%s",
			txID, &txNode.snowball)

		// If the tx should be accepted, then we should defer its acceptance
		// until its dependencies are decided. However, if this tx was
		// already marked to be accepted, we shouldn't register it again.
		if !txNode.pendingAccept &&
			txNode.Finalized(dg.params.BetaVirtuous, dg.params.BetaRogue) {
			dg.deferAcceptance(txNode)
			if dg.errs.Errored() {
				return changed, dg.errs.Err
			}
		}

		if txNode.tx.Status() != choices.Accepted {
			// If this tx wasn't accepted, then this instance is only changed if
			// preferences changed.
			changed = dg.redirectEdges(txNode) || changed
		} else {
			// By accepting a tx, the state of this instance has changed.
			changed = true
		}
	}
	return changed, dg.errs.Err
}

func (dg *Directed) String() string {
	nodes := make([]*directedTx, 0, len(dg.txs))
	for _, tx := range dg.txs {
		nodes = append(nodes, tx)
	}
	// Sort the nodes so that the string representation is canonical
	sortTxNodes(nodes)

	sb := strings.Builder{}
	sb.WriteString("DG(")

	format := fmt.Sprintf(
		"\n    Choice[%s] = ID: %%50s %%s",
		formatting.IntFormat(len(dg.txs)-1))
	for i, txNode := range nodes {
		sb.WriteString(fmt.Sprintf(format,
			i, txNode.tx.ID(), txNode.snowball.CurrentString(dg.currentVote)))
	}

	if len(nodes) > 0 {
		sb.WriteString("\n")
	}
	sb.WriteString(")")
	return sb.String()
}

// deferAcceptance attempts to mark this tx as accepted now or in the future
// once dependencies are accepted
func (dg *Directed) deferAcceptance(txNode *directedTx) {
	// Mark that this tx is pending acceptance so this function won't be called
	// again
	txNode.pendingAccept = true

	toAccept := &directedAccepter{
		dg:     dg,
		txNode: txNode,
	}
	for _, dependency := range txNode.tx.Dependencies() {
		if dependency.Status() != choices.Accepted {
			// If the dependency isn't accepted, then it must be processing. So,
			// this tx should be accepted after all of these processing txs are
			// accepted.
			toAccept.deps.Add(dependency.ID())
		}
	}

	// This tx is no longer being voted on, so we remove it from the voting set.
	// This ensures that virtuous txs built on top of rogue txs don't force the
	// node to treat the rogue tx as virtuous.
	dg.virtuousVoting.Remove(txNode.tx.ID())
	dg.pendingAccept.Register(toAccept)
}

// reject all the named txIDs and remove them from the graph
func (dg *Directed) reject(conflictIDs ...ids.ID) error {
	for _, conflictID := range conflictIDs {
		conflictKey := conflictID.Key()
		conflict := dg.txs[conflictKey]

		// This tx is not longer an option for consuming the UTXOs from its
		// inputs, so we should remove their reference to this tx.
		for _, inputID := range conflict.tx.InputIDs().List() {
			inputKey := inputID.Key()
			txIDs, exists := dg.utxos[inputKey]
			if !exists {
				// This UTXO may no longer exist because it was removed due to
				// the acceptance of a tx. If that is the case, there is nothing
				// left to remove from memory.
				continue
			}
			txIDs.Remove(conflictID)
			if txIDs.Len() == 0 {
				// If this tx was the last tx consuming this UTXO, we should
				// prune the UTXO from memory entirely.
				delete(dg.utxos, inputKey)
			} else {
				// If this UTXO still has txs consuming it, then we should make
				// sure this update is written back to the UTXOs map.
				dg.utxos[inputKey] = txIDs
			}
		}

		// We are rejecting the tx, so we should remove it from the graph
		delete(dg.txs, conflictKey)

		// While it's statistically unlikely that something being rejected is
		// preferred, it is handled for completion.
		dg.preferences.Remove(conflictID)

		// remove the edge between this node and all its neighbors
		dg.removeConflict(conflictID, conflict.ins.List()...)
		dg.removeConflict(conflictID, conflict.outs.List()...)

		// Reject is called before notifying the IPC so that rejections that
		// cause fatal errors aren't sent to an IPC peer.
		if err := conflict.tx.Reject(); err != nil {
			return err
		}

		// Notify the IPC that the tx was rejected
		dg.ctx.DecisionDispatcher.Reject(dg.ctx.ChainID, conflict.tx.ID(), conflict.tx.Bytes())

		// Update the metrics to account for this transaction's rejection
		dg.metrics.Rejected(conflictID)

		// If there is a tx that was accepted pending on this tx, the ancestor
		// tx can't be accepted.
		dg.pendingAccept.Abandon(conflictID)
		// If there is a tx that was issued pending on this tx, the ancestor tx
		// must be rejected.
		dg.pendingReject.Fulfill(conflictID)
	}
	return nil
}

// redirectEdges attempts to turn outbound edges into inbound edges if the
// preferences have changed
func (dg *Directed) redirectEdges(tx *directedTx) bool {
	changed := false
	for _, conflictID := range tx.outs.List() {
		changed = dg.redirectEdge(tx, conflictID) || changed
	}
	return changed
}

// Change the direction of this edge if needed. Returns true if the direction
// was switched.
func (dg *Directed) redirectEdge(txNode *directedTx, conflictID ids.ID) bool {
	conflict := dg.txs[conflictID.Key()]
	if txNode.numSuccessfulPolls <= conflict.numSuccessfulPolls {
		return false
	}

	// Because this tx has a higher preference than the conflicting tx, we must
	// ensure that the edge is directed towards this tx.
	nodeID := txNode.tx.ID()

	// Change the edge direction according to the conflict tx
	conflict.ins.Remove(nodeID)
	conflict.outs.Add(nodeID)
	dg.preferences.Remove(conflictID) // This conflict has an outbound edge

	// Change the edge direction according to this tx
	txNode.ins.Add(conflictID)
	txNode.outs.Remove(conflictID)
	if txNode.outs.Len() == 0 {
		// If this tx doesn't have any outbound edges, it's preferred
		dg.preferences.Add(nodeID)
	}
	return true
}

func (dg *Directed) removeConflict(txID ids.ID, neighborIDs ...ids.ID) {
	for _, neighborID := range neighborIDs {
		neighborKey := neighborID.Key()
		neighbor, exists := dg.txs[neighborKey]
		if !exists {
			// If the neighbor doesn't exist, they may have already been
			// rejected, so this mapping can be skipped.
			continue
		}

		// Remove any edge to this tx.
		neighbor.ins.Remove(txID)
		neighbor.outs.Remove(txID)

		if neighbor.outs.Len() == 0 {
			// If this tx should now be preferred, make sure its status is
			// updated.
			dg.preferences.Add(neighborID)
		}
	}
}

type directedAccepter struct {
	dg       *Directed
	deps     ids.Set
	rejected bool
	txNode   *directedTx
}

func (a *directedAccepter) Dependencies() ids.Set { return a.deps }

func (a *directedAccepter) Fulfill(id ids.ID) {
	a.deps.Remove(id)
	a.Update()
}

func (a *directedAccepter) Abandon(id ids.ID) { a.rejected = true }

func (a *directedAccepter) Update() {
	// If I was rejected or I am still waiting on dependencies to finish or an
	// error has occurred, I shouldn't do anything.
	if a.rejected || a.deps.Len() != 0 || a.dg.errs.Errored() {
		return
	}

	txID := a.txNode.tx.ID()
	// We are accepting the tx, so we should remove the node from the graph.
	delete(a.dg.txs, txID.Key())

	// This tx is consuming all the UTXOs from its inputs, so we can prune them
	// all from memory
	for _, inputID := range a.txNode.tx.InputIDs().List() {
		delete(a.dg.utxos, inputID.Key())
	}

	// This tx is now accepted, so it shouldn't be part of the virtuous set or
	// the preferred set. Its status as Accepted implies these descriptions.
	a.dg.virtuous.Remove(txID)
	a.dg.preferences.Remove(txID)

	// Reject all the txs that conflicted with this tx.
	if err := a.dg.reject(a.txNode.ins.List()...); err != nil {
		a.dg.errs.Add(err)
		return
	}
	// While it is typically true that a tx this is being accepted is preferred,
	// it is possible for this to not be the case. So this is handled for
	// completeness.
	if err := a.dg.reject(a.txNode.outs.List()...); err != nil {
		a.dg.errs.Add(err)
		return
	}

	// Accept is called before notifying the IPC so that acceptances that cause
	// fatal errors aren't sent to an IPC peer.
	if err := a.txNode.tx.Accept(); err != nil {
		a.dg.errs.Add(err)
		return
	}

	// Notify the IPC socket that this tx has been accepted.
	a.dg.ctx.DecisionDispatcher.Accept(a.dg.ctx.ChainID, txID, a.txNode.tx.Bytes())

	// Update the metrics to account for this transaction's acceptance
	a.dg.metrics.Accepted(txID)

	// If there is a tx that was accepted pending on this tx, the ancestor
	// should be notified that it doesn't need to block on this tx anymore.
	a.dg.pendingAccept.Fulfill(txID)
	// If there is a tx that was issued pending on this tx, the ancestor tx
	// doesn't need to be rejected because of this tx.
	a.dg.pendingReject.Abandon(txID)
}

type sortTxNodeData []*directedTx

func (tnd sortTxNodeData) Less(i, j int) bool {
	return bytes.Compare(
		tnd[i].tx.ID().Bytes(),
		tnd[j].tx.ID().Bytes()) == -1
}
func (tnd sortTxNodeData) Len() int      { return len(tnd) }
func (tnd sortTxNodeData) Swap(i, j int) { tnd[j], tnd[i] = tnd[i], tnd[j] }

func sortTxNodes(nodes []*directedTx) { sort.Sort(sortTxNodeData(nodes)) }
