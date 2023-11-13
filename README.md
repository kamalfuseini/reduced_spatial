# Reduced Spatial
GRPC Server that receives a list of Points, spatially reduces points
and optionally saves the points to a Cassandra database

## How to run
1. clone the repo
```
git clone https://github.com/kamalfuseini/reduced_spatial.git
```
2. Start the server and Cassandra db
```
docker compose up
```
By default the service should be running at `localhost:50051`

## Run Tests
There's the main grpc test in `@Todo link to serrver_test.go` and some tests for the points used in the reduction process in `@Todo: link to simple_point_ttestt.go`
Run at the top level for the main server test
```
go test
```
or all tests
```
go test ./...
```

## Architecture
This repo implements a grpc service `ReducedSpatial` described in `@Todo: link to reduced_spatial.proto` with one procedure `SendPoints`. A grpc service was chosen for its simplicity compared to a message broker like kafka and its high performance and reduced latency compared to REST. The request payload for the procedure is the message `SendPointsReq` that contains the track points(`points`), to reduce and save, an opional reduction distance dimension(`eps, default=1`) and an optional boolean (`noDb, default=false`) to control saving to Cassandra. The `noDb` option is mainly used for testing.

#### A Track Point
```
Point {
  X: float64 (position in X)
  Y: float64 (position in Y)
  Z: float64 (position in Z)
  T: int64 (timestamp) in milliseconds
  ID: string (uuid format) a unique identifier to identify the track
}

```
The points are reduceed using the Ramer–Douglas–Peucker algorithm before being saved to a Cassandra database in the keyspace `reduced_spatial`. The service responds with the `SendPointsReply` message that contains the total number of points recevied in `numPoints` and the number of points remaining after reducing in the field `numReducedPoints`

Cassandra is selected as the db since is a high performance database optimized for fast writes and capable of handling large amounts of data which is expected of this service. With multiple nodes Cassandra can be highly available, keeping the grpc service highly available.

## Todo:
- Add authentication to the grpc server
- Add test that uses Cassandra
- Add Authenication to Cassandra
- Save and retrieve secret keys with vault
- Configure the Cassandra instance better