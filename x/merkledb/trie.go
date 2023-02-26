// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package merkledb

import (
	"context"
	"errors"

	"github.com/dioneprotocol/dionego/ids"
	"github.com/dioneprotocol/dionego/utils/set"
)

var errNoNewRoot = errors.New("there was no updated root in change list")

// Invariant: unexported methods (except lockStack) are only called when the
// trie's view stack is locked.
type ReadOnlyTrie interface {
	// Lock this trie and those under it.
	// If this is the Database (the bottom of the view stack) only grabs a read lock.
	// For all views, grabs a write lock.
	// Invariant: This must only be called by this trie, or a trie built atop this view.
	// Invariant: Views only modify the underlying Database by calling Commit.
	lockStack()

	// Unlock this trie and those under it.
	unlockStack()

	// get the value associated with the key
	// database.ErrNotFound if the key is not present
	GetValue(ctx context.Context, key []byte) ([]byte, error)

	// get the values associated with the keys
	// database.ErrNotFound if the key is not present
	GetValues(ctx context.Context, keys [][]byte) ([][]byte, []error)

	// get the value associated with the key in path form
	// database.ErrNotFound if the key is not present
	getValue(ctx context.Context, key path) ([]byte, error)

	// get the merkle root of the Trie
	GetMerkleRoot(ctx context.Context) (ids.ID, error)

	// get the node with the given key path
	getNode(ctx context.Context, key path) (*node, error)

	// generate a proof of the value associated with a particular key, or a proof of its absence from the trie
	GetProof(ctx context.Context, bytesPath []byte) (*Proof, error)

	// generate a proof of up to maxLength smallest key/values with keys between start and end
	GetRangeProof(ctx context.Context, start, end []byte, maxLength int) (*RangeProof, error)

	// GetKeyValues but doesn't grab any locks.
	getKeyValues(
		ctx context.Context,
		start []byte,
		end []byte,
		maxLength int,
		keysToIgnore set.Set[string],
	) ([]KeyValue, error)
}

type Trie interface {
	ReadOnlyTrie

	// Delete a key from the Trie
	Remove(ctx context.Context, key []byte) error

	// Get a new view on top of this Trie
	NewPreallocatedView(ctx context.Context, estimatedChanges int) (TrieView, error)

	// Get a new view on top of this Trie
	NewView(ctx context.Context) (TrieView, error)

	// Insert a key/value pair into the Trie
	Insert(ctx context.Context, key, value []byte) error
}

// Invariant: unexported methods (except lockStack) are only called when the
// trie's view stack is locked.
type TrieView interface {
	Trie

	// Commit the changes from this Trie into the database.
	// Any views that this Trie is built on will also be committed, starting at
	// the oldest.
	Commit(ctx context.Context) error

	// Insert key/value into the trie and get back the node associated with the
	// key.
	// Updates nodes in the trie, whereas Trie.Insert records the key/value
	// without updating any trie nodes.
	insertIntoTrie(ctx context.Context, keyPath path, value Maybe[[]byte]) (*node, error)

	// Remove the key's value from the trie.
	// Updates nodes in the trie, whereas Trie.Remove records the key without
	// updating any trie nodes.
	removeFromTrie(ctx context.Context, keyPath path) error
}
