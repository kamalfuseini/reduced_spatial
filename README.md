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
There's the main grpc test in [server_test.go](https://github.com/kamalfuseini/reduced_spatial/blob/main/server_test.go) and some tests for the point type used in the reduction process in [simple_point_test.go](https://github.com/kamalfuseini/reduced_spatial/blob/main/simple_point/simple_point_test.go)
Run at the top level for the main server test
```
go test
```
or all tests
```
go test ./...
```

## Architecture
This repo implements a grpc service `ReducedSpatial` described in [reduced_spatial.proto](https://github.com/kamalfuseini/reduced_spatial/blob/main/reduced_spatial/reduced_spatial.proto) with one procedure `SendPoints`. A grpc service was chosen for its simplicity compared to a message broker like kafka and its high performance and reduced latency compared to REST. The request payload for the procedure is the message `SendPointsReq` that contains the track points(`points`) to reduce and save, an optional reduction distance dimension(`eps, default=1`) and an optional boolean (`noDb, default=false`) to control saving to Cassandra. The `noDb` option is mainly used for testing.

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
The points are reduceed using the [Ramer–Douglas–Peucker algorithm](https://en.wikipedia.org/wiki/Ramer%E2%80%93Douglas%E2%80%93Peucker_algorithm) before being saved to a Cassandra database in the keyspace `reduced_spatial`. The shortest distance to the line segment is used in place of the perpendicular distance described in the Wikipiedia article of the algorithm. This should be more accurate for points the lie outside the line segment created by the start and end points. The service responds with the `SendPointsReply` message that contains the total number of points recevied in the field `numPoints` and the number of points remaining after reducing in the field `numReducedPoints`

Cassandra is selected for the database since it is a high performance, optimized for fast writes and capable of handling large amounts of data which is expected of this service. With multiple nodes Cassandra can be highly available, keeping the grpc service highly available.

## Possible Improvements:
- Add authentication to the grpc server
- Add test that uses Cassandra
- Add Authenication to Cassandra
- Save and retrieve secret keys with vault
- Configure the Cassandra instance better