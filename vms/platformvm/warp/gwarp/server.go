// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package gwarp

import (
	"context"

	"github.com/DioneProtocol/odysseygo/ids"
	"github.com/DioneProtocol/odysseygo/vms/platformvm/warp"

	pb "github.com/DioneProtocol/odysseygo/proto/pb/warp"
)

var _ pb.SignerServer = (*Server)(nil)

type Server struct {
	pb.UnsafeSignerServer
	signer warp.Signer
}

func NewServer(signer warp.Signer) *Server {
	return &Server{signer: signer}
}

func (s *Server) Sign(_ context.Context, unsignedMsg *pb.SignRequest) (*pb.SignResponse, error) {
	sourceChainID, err := ids.ToID(unsignedMsg.SourceChainId)
	if err != nil {
		return nil, err
	}

	msg, err := warp.NewUnsignedMessage(
		unsignedMsg.NetworkId,
		sourceChainID,
		unsignedMsg.Payload,
	)
	if err != nil {
		return nil, err
	}

	sig, err := s.signer.Sign(msg)
	return &pb.SignResponse{
		Signature: sig,
	}, err
}
