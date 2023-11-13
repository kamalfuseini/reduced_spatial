package simple_point

import (
	"math"
	"testing"
)

func TestMagnitude(t *testing.T) {
	m := NewSimplePoint(3, 0, 1).Magnitude()
	expected := 3.162
	eps := 0.001
	diff := math.Abs(m - expected)
	if diff > eps {
		t.Errorf("Expected: %f(+-%f)\nGot: %f, diff: %f\n", expected, eps, m, diff)
	}
}

func TestDistance(t *testing.T) {
	d := NewSimplePoint(1, 1, 0).Distance(NewSimplePoint(2, 1, 2))
	expected := 2.236
	eps := 0.001
	diff := math.Abs(d - expected)
	if diff > eps {
		t.Errorf("Expected: %f(+-%f)\nGot: %f, diff: %f\n", expected, eps, d, diff)
	}
}

func TestCross(t *testing.T) {
	c := NewSimplePoint(-4, -7, -1).Cross(NewSimplePoint(3, 0, 1))
	expected := NewSimplePoint(-7, 1, 21)

	isEqual := func (a SimplePoint, b SimplePoint, eps float64) bool {
		return math.Abs(a.X - b.X) < eps && math.Abs(a.Y - b.Y) < eps && math.Abs(a.Z - b.Z) < eps
	}
	
	eps := 0.001
	if !isEqual(c, expected, eps) {
		t.Errorf("Expected: %v(+-%f)\nGot: %v\n", expected, eps, c)
	}
}

func TestPerpendicularDistance(t *testing.T) {
	d := PerpendicularDistance(NewSimplePoint(-3, -4, 0), NewSimplePoint(1, 3, 1), NewSimplePoint(4, 3, 2))
	expected := 7.007
	eps := 0.001
	diff := math.Abs(d - expected)
	if diff > eps {
		t.Errorf("Expected: %f(+-%f)\nGot: %f, diff: %f\n", expected, eps, d, diff)
	}
}

func TestShortestDistance(t *testing.T) {
	d := ShortestDistance(NewSimplePoint(5, 4, 0), NewSimplePoint(1, 3, 0), NewSimplePoint(4, 3, 0))
	expected := 1.414
	eps := 0.001
	diff := math.Abs(d - expected)
	if diff > eps {
		t.Errorf("Expected: %f(+-%f)\nGot: %f, diff: %f\n", expected, eps, d, diff)
	}
}
