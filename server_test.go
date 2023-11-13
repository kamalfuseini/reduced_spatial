package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kfuseini/reduced_spatial/config"
	"github.com/kfuseini/reduced_spatial/env"
	pb "github.com/kfuseini/reduced_spatial/reduced_spatial"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func getClientStartServer(ctx context.Context) (pb.ReducedSpatialClient, func()) {
	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)

	s := grpc.NewServer()
	config := config.Default()
	pb.RegisterReducedSpatialServer(s, &Server{
		env: &env.Env{
			Config: &config,
		},
	})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	stop := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		s.GracefulStop()
	}

	client := pb.NewReducedSpatialClient(conn)

	return client, stop
}

func SendPointsReplyToString(r *pb.SendPointsReply) string {
	return fmt.Sprintf("SendPointsReply { NumPoints: %d, NumReducedPoints: %d }", r.NumPoints, r.NumReducedPoints)
}

func generatePoints(n int) []*pb.Point {
	points := []*pb.Point{}
	for i := 0; i < n; i += 1 {
		points = append(points, &pb.Point{
			ID: uuid.NewString(),
		})
	}
	return points
}

func NewPoint(x float64, y float64, z float64) *pb.Point {
	return &pb.Point{
		// ID: uuid.NewString(),
		// T:  time.Now().Unix(),
		X: x,
		Y: y,
		Z: z,
	}
}

func TestSendPoints(t *testing.T) {
	client, stopServer := getClientStartServer(context.Background())
	defer stopServer()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	arrOfPoints := []struct {
		Points   []*pb.Point
		Eps      float64
		Expected *pb.SendPointsReply
	}{
		{
			Points:   []*pb.Point{},
			Expected: &pb.SendPointsReply{NumPoints: 0, NumReducedPoints: 0},
		},
		{
			Points:   []*pb.Point{{X: 9, Y: 2, Z: -1}},
			Expected: &pb.SendPointsReply{NumPoints: 1, NumReducedPoints: 1},
		},
		{
			Points:   []*pb.Point{NewPoint(1, 5, 0), NewPoint(2, 4, 0), NewPoint(3, 2, 0), NewPoint(4, 2.5, 0), NewPoint(5, 4, 0), NewPoint(6, 3.5, 0), NewPoint(7, 2, 0), NewPoint(8, 5, 0)},
			Eps:      1,
			Expected: &pb.SendPointsReply{NumPoints: 8, NumReducedPoints: 5},
		},
	}

	noDb := true
	for _, req := range arrOfPoints {
		points := req.Points
		expected := req.Expected
		eps := &req.Eps
		if *eps == 0 {
			eps = nil
		}
		r, err := client.SendPoints(ctx, &pb.SendPointsReq{Points: points, Eps: eps, NoDb: &noDb})
		if err != nil {
			t.Errorf("Expected: %q\nGot Err: %q\n", SendPointsReplyToString(expected), err)
		}
		if r.NumPoints != expected.NumPoints || r.NumReducedPoints != expected.NumReducedPoints {
			t.Errorf("Expected: %q\nGot: %q\n", SendPointsReplyToString(expected), SendPointsReplyToString(r))
		}
	}
}
