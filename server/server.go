package server

import (
	"context"

	"github.com/h0n9/petit-chat/control"
	"github.com/h0n9/petit-chat/p2p"
	"github.com/h0n9/petit-chat/types"
	"github.com/h0n9/petit-chat/util"
)

type Server struct {
	ctx context.Context
	cfg *util.Config

	node   *p2p.Node
	center *control.Center
}

func NewServer(ctx context.Context, cfg *util.Config) (*Server, error) {
	node, err := p2p.NewNode(ctx, cfg)
	if err != nil {
		return nil, err
	}
	center, err := control.NewCenter(ctx, node.GetHostID())
	if err != nil {
		return nil, err
	}
	return &Server{
		ctx: ctx,
		cfg: cfg,

		node:   node,
		center: center,
	}, nil
}

func (s *Server) Close() error {
	return s.node.Close()
}

func (s *Server) GetID() types.ID {
	return s.node.GetHostID()
}

func (s *Server) DiscoverPeers() error {
	return s.node.DiscoverPeers(s.cfg.BootstrapNodes)
}

func (s *Server) GetPeers() []types.ID {
	return s.node.GetPeers()
}

func (s *Server) PrintInfo() {
	s.node.Info()
}

func (s *Server) GetCenter() *control.Center {
	return s.center
}

func (s *Server) CreateMsgBox(tStr, nickname string, pub bool) (*control.Box, error) {
	topic, err := s.node.Join(tStr)
	if err != nil {
		return nil, err
	}
	// TODO: get metdata from parameters
	persona, err := types.NewPersona(nickname, []byte{}, s.node.PubKey)
	if err != nil {
		return nil, err
	}
	return s.center.CreateBox(topic, pub, s.node.PrivKey, persona)
}

func (s *Server) LeaveMsgBox(topicStr string) error {
	return s.center.LeaveBox(topicStr)
}

func (s *Server) GetMsgBox(topicStr string) (*control.Box, bool) {
	return s.center.GetBox(topicStr)
}
