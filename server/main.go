package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/prasad-joshi-ntx/net_file_copy/file-copy"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type Server struct {
	pb.UnimplementedFileCopyServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) Write(
	ctx context.Context,
	args *pb.WriteArgs) (*pb.WriteResponse, error) {
	decoded, err := hex.DecodeString(args.GetData())
	if err != nil {
		return &pb.WriteResponse{ByesCopied: 0, Error: -1}, err
	}
	size := len(decoded)
	// fmt.Printf("Received data bytes ", size)
	return &pb.WriteResponse{ByesCopied: int64(size), Error: 0}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFileCopyServer(s, &Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
