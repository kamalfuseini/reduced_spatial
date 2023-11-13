package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/gocql/gocql"
	"github.com/kfuseini/reduced_spatial/config"
	"github.com/kfuseini/reduced_spatial/env"
	pb "github.com/kfuseini/reduced_spatial/reduced_spatial"
	sp "github.com/kfuseini/reduced_spatial/simple_point"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Server struct {
	env        *env.Env
	grpcServer *grpc.Server
	pb.UnimplementedReducedSpatialServer
}

func (s *Server) SendPoints(ctx context.Context, in *pb.SendPointsReq) (*pb.SendPointsReply, error) {
	points := in.GetPoints()
	saveDb := in.NoDb == nil || !*in.NoDb
	eps := s.env.Config.Eps

	if in.Eps != nil {
		eps = *in.Eps
	}

	idPointsMap := make(map[string][]*pb.Point)
	for i, point := range points {
		if _, ok := idPointsMap[point.ID]; ok {
			idPointsMap[point.ID] = append(idPointsMap[point.ID], points[i])
		} else {
			idPointsMap[point.ID] = []*pb.Point{points[i]}
		}
	}

	numReduced := 0
	for _, points := range idPointsMap {
		reducedPoints := reducePoints(points, eps)
		numReduced += len(reducedPoints)

		if saveDb {
			batch := s.env.Cass.NewBatch(gocql.LoggedBatch)
			for _, point := range reducedPoints {
				query := "INSERT INTO reduced_spatial.track_points (id, t, x, y, z) VALUES(?,?,?,?,?)"
				batch.Query(query, point.ID, point.T, point.X, point.Y, point.Z)
			}

			if err := s.env.Cass.ExecuteBatch(batch); err != nil {
				return nil, errors.Wrap(err, "Failed to save to cassandra")
			}
		}
	}

	reply := &pb.SendPointsReply{
		NumPoints:        int32(len(points)),
		NumReducedPoints: int32(numReduced),
	}
	return reply, nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.env.Config.Port))
	if err != nil {
		return errors.Wrap(err, "Failed to listen")
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.grpcServer.Serve(lis); err != nil {
		return errors.Wrap(err, "Failed to serve")
	}

	return nil
}

func NewServer(e *env.Env) (*Server, error) {
	s := grpc.NewServer()
	server := Server{
		env:        e,
		grpcServer: s,
	}
	pb.RegisterReducedSpatialServer(s, &server)

	return &server, nil
}

func reducePoints(points []*pb.Point, eps float64) []*pb.Point {
	if len(points) <= 2 {
		if len(points) == 2 {
			d := sp.SimplePointFromPoint(points[0]).Distance(sp.SimplePointFromPoint(points[1]))
			if d <= eps {
				return []*pb.Point{points[0]}
			}
		}
		return append([]*pb.Point{}, points...)
	}

	keeping := make([]bool, len(points))

	var updateKeeping func(int, int)
	updateKeeping = func(startIdx int, endIdx int) {
		dMax := 0.0
		idx := 0
		for i, point := range points[startIdx+1 : endIdx] {
			A := sp.SimplePointFromPoint(points[startIdx])
			B := sp.SimplePointFromPoint(points[endIdx])
			P := sp.SimplePointFromPoint(point)
			d := sp.ShortestDistance(P, A, B)
			if d > dMax {
				dMax = d
				idx = i + startIdx + 1
			}
		}

		if dMax > eps {
			updateKeeping(startIdx, idx)
			updateKeeping(idx, endIdx)
		} else {
			keeping[startIdx] = true
			keeping[endIdx] = true
		}
	}

	updateKeeping(0, len(points)-1)

	reducedPoints := []*pb.Point{}
	for idx, keep := range keeping {
		if keep {
			reducedPoints = append(reducedPoints, points[idx])
		}
	}

	return reducedPoints
}

func main() {
	config, err := config.Init()
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}
	flag.Parse()

	e, err := env.NewEnv(config)
	if err != nil {
		log.Fatalf("Failed to create env: %v", err)
	}
	if e == nil {
		log.Fatalf("Failed to create env: env is nil")
	}

	server, err := NewServer(e)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
