package main

import (
	// "errors"
	"flag"
	"fmt"
	pb3 "goNotes/grpcnotes/echo"
	pb "goNotes/grpcnotes/hello"
	pb2 "goNotes/grpcnotes/world"
	"io"
	"log"
	"net"

	"github.com/pkg/errors"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// const (
// 	port = ":9000"
// )

// server is used to implement helloworld.GreeterServer.
type serverHello struct{}
type serverWorld struct{}
type serverEcho struct{}

// SayHello implements helloworld.GreeterServer
func (s *serverHello) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// return &pb.HelloReply{Message: strings.Repeat("Hello ", int(in.Num))}, nil
	return &pb.HelloReply{Message: "Hello "}, nil
}

func (s *serverWorld) SayWorld(ctx context.Context, in *pb2.WorldRequest) (*pb2.WorldReply, error) {
	// return &pb2.WorldReply{Message: strings.Repeat("World ", int(in.Num))}, nil
	return &pb2.WorldReply{Message: "World "}, nil
}

func (s *serverEcho) SayEcho(ctx context.Context, in *pb3.EchoRequest) (*pb3.EchoReply, error) {
	return &pb3.EchoReply{Message: "echo "}, nil
}

func (s *serverEcho) SayEchoS(stream pb3.Echo_SayEchoSServer) error {
	// return &pb2.WorldReply{Message: strings.Repeat("World ", int(in.Num))}, nil
	// return &pb2.WorldReply{Message: "World "}, nil
	for {
		_, err := stream.Recv()
		if err == io.EOF {

			return errors.New("stream.Recv() io.EOF")
		}
		if err != nil {
			return errors.Wrap(err, "receive errors")
		}

		stream.Send(&pb3.EchoReply{
			Message: "echo",
		})
	}

}

var port = flag.Int("port", 19000, "port")

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{}),
		grpc.MaxConcurrentStreams(10000))
	pb.RegisterHelloServer(s, &serverHello{})
	pb2.RegisterWorldServer(s, &serverWorld{})
	pb3.RegisterEchoServer(s, &serverEcho{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
