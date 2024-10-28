package server

import (
	"github.com/nuzur/extension-sdk/client"
	pb "github.com/nuzur/extension-sdk/idl/gen"
	"go.uber.org/fx"
)

type server struct {
	pb.UnimplementedNuzurExtensionServer
	client *client.Client
}

type Params struct {
	fx.In
	Client *client.Client
}

func New(params Params) (pb.NuzurExtensionServer, error) {
	return &server{
		client: params.Client,
	}, nil
}
