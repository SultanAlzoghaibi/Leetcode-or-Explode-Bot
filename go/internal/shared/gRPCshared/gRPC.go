// monitor_rpc.go
// One file that exposes BOTH a gRPC server (listener) and a simple client sender.
// Drop this in a shared package and import where needed.

package gRPCshared

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "Leetcode-or-Explode-Bot/proto" // generated from your proto
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ----- Server -----

// MonitorServer implements the gRPC service.

// SendStatus logs the incoming boolean and echoes it back.
// (Change behavior as you like; 1/0 is just a bool here.)
type MonitorServer struct {
	pb.UnimplementedMonitorServiceServer
	OnReceive func(status Status)
}

func (s *MonitorServer) SendStatus(ctx context.Context, req *pb.Status) (*pb.Status, error) {
	log.Printf("📥 Received status: %v", req.ServerStatus)
	if s.OnReceive != nil {
		s.OnReceive(Status{
			PodName:  req.PodName,
			IsOn:     req.ServerStatus,
			Password: req.Password,
		})
	}
	// echo back (or customize)
	return &pb.Status{ServerStatus: req.ServerStatus}, nil
}

type Status struct {
	PodName  string
	IsOn     bool
	Password uint32
}

// StartMonitorServer starts a gRPC server bound to bindAddr, e.g. "0.0.0.0:4000" or "10.0.0.5:4000".
func StartListen(bindAddr string) (Status, error) {
	if bindAddr == "" {
		bindAddr = "0.0.0.0:50051"
	}

	lis, err := net.Listen("tcp", bindAddr)
	if err != nil {
		return Status{}, fmt.Errorf("failed to listen on %s: %w", bindAddr, err)
	}

	// Channel for receiving full Status struct
	received := make(chan Status, 1)

	/*
		→ Makes your own struct that actually implements
		the MonitorService interface generated from your .proto file.
		That means it has a SendStatus method with the right signature.
	*/
	server := &MonitorServer{
		UnimplementedMonitorServiceServer: pb.UnimplementedMonitorServiceServer{},
		OnReceive: func(s Status) { // Now OnReceive takes the full Status
			fmt.Println("Server &MonitorServer")
			received <- s
		},
	}

	grpcServer := grpc.NewServer()
	/*
		→ Makes an empty gRPC server container —
		it’s just a framework object that can host services,
		but it doesn’t know what services exist yet.
	*/

	// For requests to the MonitorService service,
	//call methods on this specific server object.
	pb.RegisterMonitorServiceServer(grpcServer, server)

	go func() {
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			log.Printf("gRPC server stopped: %v", serveErr)
		}
	}()

	return <-received, nil
}

// ----- Client -----

// SendStatusTo dials targetAddr (e.g. "10.0.0.5:4000") and sends a single boolean.
// Returns the server's echoed boolean (or other response you implement).
func SendStatusTo(targetAddr string, value bool) (bool, error) {
	conn, err := grpc.Dial(
		targetAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // swap to TLS for public exposure
		grpc.WithBlock(),
	)
	if err != nil {
		return false, fmt.Errorf("dial %s: %w", targetAddr, err)
	}
	defer conn.Close()

	client := pb.NewMonitorServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.SendStatus(ctx, &pb.Status{ServerStatus: value})
	if err != nil {
		return false, fmt.Errorf("SendStatus RPC: %w", err)
	}
	log.Printf("📤 Sent %v, server replied %v", value, resp.ServerStatus)
	return resp.ServerStatus, nil
}
