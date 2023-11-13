package simple_point

import (
	"math"

	pb "github.com/kfuseini/reduced_spatial/reduced_spatial"
)

type SimplePoint struct {
	X float64
	Y float64
	Z float64
}

func NewSimplePoint(X float64, Y float64, Z float64) SimplePoint {
	return SimplePoint{ X, Y, Z }
}

func (a SimplePoint) Sub(b SimplePoint) SimplePoint {
	return SimplePoint{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func (a SimplePoint) Cross(b SimplePoint) SimplePoint {
	return SimplePoint{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func (a SimplePoint) Dot(b SimplePoint) float64 {
	return a.X * b.X + a.Y * b.Y + a.Z * b.Z
}

func (a SimplePoint) Magnitude() float64 {
	return math.Sqrt(a.Dot(a))
}

func (a SimplePoint) Distance(b SimplePoint) float64 {
	return math.Sqrt(math.Pow(a.X - b.X, 2) + math.Pow(a.Y - b.Y, 2) + math.Pow(a.Z - b.Z, 2))
}

func SimplePointFromPoint(point *pb.Point) SimplePoint {
	return SimplePoint{
		X: point.X,
		Y: point.Y,
		Z: point.Z,
	}
}

func PerpendicularDistance(P SimplePoint, A SimplePoint, B SimplePoint) float64 {
	AB := B.Sub(A)
	AP := P.Sub(A)
	APXAB := AP.Cross(AB)

	return APXAB.Magnitude() / AB.Magnitude()
}

func ShortestDistance(P SimplePoint, A SimplePoint, B SimplePoint) float64 {
	AB := B.Sub(A)

	AP := P.Sub(A)
	if AP.Dot(AB) <= 0 {
		return AP.Magnitude()
	}

	BP := P.Sub(B)
	if BP.Dot(AB) >= 0 {
		return BP.Magnitude()
	}

	return AP.Cross(AB).Magnitude() / AB.Magnitude()
}

