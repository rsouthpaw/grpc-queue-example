package handler

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	pb "gdma_pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}
type PacketReq struct {
	Request pb.GDMARequest
	Key     string
}
type PacketRes struct {
	Request pb.GDMAResponse
	Key     string
}

var (
	bchanInput   = make(chan PacketReq, 2)
	ubchanOutput = make(chan PacketRes)
)
var (
	ri = 0
)

func init() {
	//initialize channel
}

func startServerInteractor(port string) {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGDMAServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) GetDistance(ctx context.Context, in *pb.GDMARequest) (*pb.GDMAResponse, error) {

	ri++
	log.Println("New Request!", ri)

	var key string
	if ctx.Value("KEY") == nil {
		key = "ctx" + strconv.Itoa(ri)
	}
	bchanInput <- PacketReq{
		Request: *in,
		Key:     key,
	}

	go waitForData()

	for {
		val := <-ubchanOutput
		log.Println("Received:", val.Key)
		if val.Key == key {
			return &pb.GDMAResponse{}, nil
		}
	}

	return &pb.GDMAResponse{}, nil
}
func waitForData() {
	if len(bchanInput) == cap(bchanInput) {
		log.Println("Calling GDMA...")
		callDistanceMatrixApi()
	} else {
		log.Println("Waiting for more data...")
	}
}
func callDistanceMatrixApi() {
	time.Sleep(time.Second * 1)
	batch := createNewBatch()
	/**
	  GDMA WILL BE CALLED HERE WITH GIVEN BATCH...
	*/
	log.Println("Response received from GDMA for", len(batch), "requests...")
	ubchanOutput <- PacketRes{Key: "ctx1"}
	ubchanOutput <- PacketRes{Key: "ctx2"}
}
func createNewBatch() []pb.GDMARequest {
	var request PacketReq
	var batch []pb.GDMARequest
	request = <-bchanInput
	batch = append(batch, request.Request)

	request = <-bchanInput
	batch = append(batch, request.Request)
	log.Println("out")
	return batch
}
