// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package c

import (
	"errors"
	"fmt"

	stdcontext "context"

	"github.com/DioneProtocol/coreth/plugin/delta"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/DioneProtocol/odysseygo/database"
	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/utils/crypto/keychain"
	"github.com/DioneProtocol/odysseygo/utils/crypto/secp256k1"
	"github.com/DioneProtocol/odysseygo/utils/hashing"
	"github.com/DioneProtocol/odysseygo/utils/set"
	"github.com/DioneProtocol/odysseygo/vms/components/dione"
	"github.com/DioneProtocol/odysseygo/vms/components/verify"
	"github.com/DioneProtocol/odysseygo/vms/secp256k1fx"
)

const version = 0

var (
	_ Signer = (*txSigner)(nil)

	errUnknownInputType      = errors.New("unknown input type")
	errUnknownCredentialType = errors.New("unknown credential type")
	errUnknownOutputType     = errors.New("unknown output type")
	errInvalidUTXOSigIndex   = errors.New("invalid UTXO signature index")

	emptySig [secp256k1.SignatureLen]byte
)

type Signer interface {
	SignUnsignedAtomic(ctx stdcontext.Context, tx delta.UnsignedAtomicTx) (*delta.Tx, error)
	SignAtomic(ctx stdcontext.Context, tx *delta.Tx) error
}

type EthKeychain interface {
	// The returned Signer can provide a signature for [addr]
	GetEth(addr ethcommon.Address) (keychain.Signer, bool)
	// Returns the set of addresses for which the accessor keeps an associated
	// signer
	EthAddresses() set.Set[ethcommon.Address]
}

type SignerBackend interface {
	GetUTXO(ctx stdcontext.Context, chainID, utxoID ids.ID) (*dione.UTXO, error)
}

type txSigner struct {
	dioneKC keychain.Keychain
	ethKC   EthKeychain
	backend SignerBackend
}

func NewSigner(dioneKC keychain.Keychain, ethKC EthKeychain, backend SignerBackend) Signer {
	return &txSigner{
		dioneKC: dioneKC,
		ethKC:   ethKC,
		backend: backend,
	}
}

func (s *txSigner) SignUnsignedAtomic(ctx stdcontext.Context, utx delta.UnsignedAtomicTx) (*delta.Tx, error) {
	tx := &delta.Tx{UnsignedAtomicTx: utx}
	return tx, s.SignAtomic(ctx, tx)
}

func (s *txSigner) SignAtomic(ctx stdcontext.Context, tx *delta.Tx) error {
	switch utx := tx.UnsignedAtomicTx.(type) {
	case *delta.UnsignedImportTx:
		signers, err := s.getImportSigners(ctx, utx.SourceChain, utx.ImportedInputs)
		if err != nil {
			return err
		}
		return sign(tx, true, signers)
	case *delta.UnsignedExportTx:
		signers := s.getExportSigners(utx.Ins)
		return sign(tx, true, signers)
	default:
		return fmt.Errorf("%w: %T", errUnknownTxType, tx)
	}
}

func (s *txSigner) getImportSigners(ctx stdcontext.Context, sourceChainID ids.ID, ins []*dione.TransferableInput) ([][]keychain.Signer, error) {
	txSigners := make([][]keychain.Signer, len(ins))
	for credIndex, transferInput := range ins {
		input, ok := transferInput.In.(*secp256k1fx.TransferInput)
		if !ok {
			return nil, errUnknownInputType
		}

		inputSigners := make([]keychain.Signer, len(input.SigIndices))
		txSigners[credIndex] = inputSigners

		utxoID := transferInput.InputID()
		utxo, err := s.backend.GetUTXO(ctx, sourceChainID, utxoID)
		if err == database.ErrNotFound {
			// If we don't have access to the UTXO, then we can't sign this
			// transaction. However, we can attempt to partially sign it.
			continue
		}
		if err != nil {
			return nil, err
		}

		out, ok := utxo.Out.(*secp256k1fx.TransferOutput)
		if !ok {
			return nil, errUnknownOutputType
		}

		for sigIndex, addrIndex := range input.SigIndices {
			if addrIndex >= uint32(len(out.Addrs)) {
				return nil, errInvalidUTXOSigIndex
			}

			addr := out.Addrs[addrIndex]
			key, ok := s.dioneKC.Get(addr)
			if !ok {
				// If we don't have access to the key, then we can't sign this
				// transaction. However, we can attempt to partially sign it.
				continue
			}
			inputSigners[sigIndex] = key
		}
	}
	return txSigners, nil
}

func (s *txSigner) getExportSigners(ins []delta.DELTAInput) [][]keychain.Signer {
	txSigners := make([][]keychain.Signer, len(ins))
	for credIndex, input := range ins {
		inputSigners := make([]keychain.Signer, 1)
		txSigners[credIndex] = inputSigners

		key, ok := s.ethKC.GetEth(input.Address)
		if !ok {
			// If we don't have access to the key, then we can't sign this
			// transaction. However, we can attempt to partially sign it.
			continue
		}
		inputSigners[0] = key
	}
	return txSigners
}

// TODO: remove [signHash] after the ledger supports signing all transactions.
func sign(tx *delta.Tx, signHash bool, txSigners [][]keychain.Signer) error {
	unsignedBytes, err := delta.Codec.Marshal(version, &tx.UnsignedAtomicTx)
	if err != nil {
		return fmt.Errorf("couldn't marshal unsigned tx: %w", err)
	}
	unsignedHash := hashing.ComputeHash256(unsignedBytes)

	if expectedLen := len(txSigners); expectedLen != len(tx.Creds) {
		tx.Creds = make([]verify.Verifiable, expectedLen)
	}

	sigCache := make(map[ids.ShortID][secp256k1.SignatureLen]byte)
	for credIndex, inputSigners := range txSigners {
		credIntf := tx.Creds[credIndex]
		if credIntf == nil {
			credIntf = &secp256k1fx.Credential{}
			tx.Creds[credIndex] = credIntf
		}

		cred, ok := credIntf.(*secp256k1fx.Credential)
		if !ok {
			return errUnknownCredentialType
		}
		if expectedLen := len(inputSigners); expectedLen != len(cred.Sigs) {
			cred.Sigs = make([][secp256k1.SignatureLen]byte, expectedLen)
		}

		for sigIndex, signer := range inputSigners {
			if signer == nil {
				// If we don't have access to the key, then we can't sign this
				// transaction. However, we can attempt to partially sign it.
				continue
			}
			addr := signer.Address()
			if sig := cred.Sigs[sigIndex]; sig != emptySig {
				// If this signature has already been populated, we can just
				// copy the needed signature for the future.
				sigCache[addr] = sig
				continue
			}

			if sig, exists := sigCache[addr]; exists {
				// If this key has already produced a signature, we can just
				// copy the previous signature.
				cred.Sigs[sigIndex] = sig
				continue
			}

			var sig []byte
			if signHash {
				sig, err = signer.SignHash(unsignedHash)
			} else {
				sig, err = signer.Sign(unsignedBytes)
			}
			if err != nil {
				return fmt.Errorf("problem signing tx: %w", err)
			}
			copy(cred.Sigs[sigIndex][:], sig)
			sigCache[addr] = cred.Sigs[sigIndex]
		}
	}

	signedBytes, err := delta.Codec.Marshal(version, tx)
	if err != nil {
		return fmt.Errorf("couldn't marshal tx: %w", err)
	}
	tx.Initialize(unsignedBytes, signedBytes)
	return nil
}
