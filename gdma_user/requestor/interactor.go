package requestor

import (
	"context"
	"log"
	"time"

	pb "gdma_pb"

	"google.golang.org/grpc"
)

var (
	conn    *grpc.ClientConn
	address = "localhost:50051"
)

func init() {

	var err error
	conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
}

func mockGDMARequestsInteractor() {

	if conn == nil {
		log.Println("connection error")
		return
	}

	c := pb.NewGDMAClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	go c.GetDistance(context.WithValue(ctx, "KEY", "ctx1"), &pb.GDMARequest{Sources: []string{}, Destinations: []string{}})
	time.Sleep(time.Second * 3)
	_, err := c.GetDistance(context.WithValue(ctx, "KEY", "ctx2"), &pb.GDMARequest{Sources: []string{}, Destinations: []string{}})
	if err != nil {
		log.Fatalf("gdma error: %v", err)
	}

}
