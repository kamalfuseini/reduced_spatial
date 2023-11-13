package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	pb "github.com/kfuseini/reduced_spatial/reduced_spatial"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func NewPoint(id string, x float64, y float64, z float64) *pb.Point {
	return &pb.Point{
		ID: id,
		T:  time.Now().UnixMilli() + rand.Int63n(100),
		X: x,
		Y: y,
		Z: z,
	}
}

func SendPointsReplyToString(r *pb.SendPointsReply) string {
	return fmt.Sprintf("SendPointsReply { NumPoints: %d, NumReducedPoints: %d }", r.NumPoints, r.NumReducedPoints)
}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewReducedSpatialClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	id := uuid.NewString()
	arrOfPoints := [][]*pb.Point{
		{},
		{{ID: uuid.NewString(), X: 9, Y: 2, Z: -1}},
		{NewPoint(id, 1, 1, 1), NewPoint(id, 1, 1, 1)},
		{NewPoint(id, 1, 1, 1), NewPoint(id, 1, 1, 1), NewPoint(id, 3, 3, 3), NewPoint(id, 4, 4, 4)},
		{NewPoint(id, 1, 5, 0), NewPoint(id, 2, 4, 0), NewPoint(id, 3, 2, 0), NewPoint(id, 4, 2.5, 0), NewPoint(id, 5, 4, 0), NewPoint(id, 6, 3.5, 0), NewPoint(id, 7, 2, 0), NewPoint(id, 8, 5, 0)},
	}
	noDb := false
	for _, points := range arrOfPoints {
		r, err := c.SendPoints(ctx, &pb.SendPointsReq{Points: points, NoDb: &noDb})
		if err != nil {
			log.Fatalf("Got Err: %q\n", err)
		}
		log.Printf("response: %v", r)
	}
}
