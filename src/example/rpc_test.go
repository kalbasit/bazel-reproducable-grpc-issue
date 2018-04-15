package example_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/gogo/protobuf/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/kalbasit/bazel-reproducable-grpc-issue/src/example"
)

type exampleService struct {
}

func (ls *exampleService) Run(ctx context.Context, e *types.Empty) (*types.Empty, error) {
	return &types.Empty{}, nil
}

const serverAddr = "127.0.0.1"

var serverPort = 2730

func newServer(t *testing.T) *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", serverAddr, serverPort))
	if err != nil {
		serverPort++
		return newServer(t)
	}

	s := grpc.NewServer()
	example.RegisterExampleServer(s, new(exampleService))
	go s.Serve(lis)
	return s
}

func TestEvent(t *testing.T) {
	server := newServer(t)
	defer server.Stop()

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverAddr, serverPort), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("unable to open a connection to the RPC server: %s", err)
	}
	defer conn.Close()
	client := example.NewExampleClient(conn)

	if _, err := client.Run(context.Background(), &types.Empty{}); err != nil {
		t.Fatalf("Event returned an unexpected error: %s", err)
	}
}

func TestError(t *testing.T) {
	server := newServer(t)
	defer server.Stop()

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", serverAddr, serverPort), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("unable to open a connection to the RPC server: %s", err)
	}
	defer conn.Close()
	client := example.NewExampleClient(conn)

	if _, err := client.Run(context.Background(), &types.Empty{}); err != nil {
		t.Fatalf("Event returned an unexpected error: %s", err)
	}
}
