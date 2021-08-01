package geography

import (
	"math"
)

var emptyBound = Bound{Min: Point{1, 1}, Max: Point{-1, -1}}

type Bound struct {
	Min Point
	Max Point
}

func (b Bound) ToGeom() Geom {
	return b
}

func (b Bound) Clip(bound Bound) Geom {
	return Bound{
		Min: Point{
			math.Max(b.Min[0], bound.Min[0]),
			math.Max(b.Min[1], bound.Min[1]),
		},
		Max: Point{
			math.Min(b.Max[0], bound.Max[0]),
			math.Min(b.Max[1], bound.Max[1]),
		},
	}
}

func (b Bound) Cap() int {
	return b.AsPolygon().Cap()
}

func (b Bound) Geometry() []uint32 {
	return b.AsPolygon().Geometry()
}

func (b Bound) Project(transform Transform) Geom {
	return Bound{
		Min: transform(b.Min),
		Max: transform(b.Max),
	}
}

func (Bound) Type() string {
	return "Polygon"
}

func (b Bound) IsEmpty() bool {
	return b.Min[0] > b.Max[0] || b.Min[1] > b.Max[1]
}

func (b Bound) Equal(g Geom) bool {
	switch bound := g.(type) {
	case Bound:
		return b.Min.Equal(bound.Min) && bound.Max.Equal(bound.Max)
	}
	return false
}

func (b Bound) Bound() Bound {
	return b
}

func (b Bound) AsPolygon() *Polygon {
	return &Polygon{{
		b.Min,
		{b.Max[0], b.Min[1]},
		b.Max,
		{b.Min[0], b.Max[1]},
		b.Min,
	}}
}

func (Bound) DataType(driverName string) string {
	return "Polygon"
}

func (b Bound) Top() float64 {
	return b.Max[1]
}

func (b Bound) Bottom() float64 {
	return b.Min[1]
}

func (b Bound) Right() float64 {
	return b.Max[0]
}

func (b Bound) Left() float64 {
	return b.Min[0]
}

func (b Bound) LeftTop() Point {
	return Point{b.Left(), b.Top()}
}

func (b Bound) RightBottom() Point {
	return Point{b.Right(), b.Bottom()}
}

func (b Bound) Intersects(bound Bound) bool {
	return !((b.Max[0] < bound.Min[0]) || (b.Min[0] > bound.Max[0]) || (b.Max[1] < bound.Min[1]) || (b.Min[1] > bound.Max[1]))
}

func (b Bound) Contains(point Point) bool {
	if point[1] < b.Min[1] || b.Max[1] < point[1] {
		return false
	}
	if point[0] < b.Min[0] || b.Max[0] < point[0] {
		return false
	}
	return true
}

func (b Bound) Extend(point Point) Bound {
	if b.Contains(point) {
		return b
	}

	return Bound{
		Min: Point{
			math.Min(b.Min[0], point[0]),
			math.Min(b.Min[1], point[1]),
		},
		Max: Point{
			math.Max(b.Max[0], point[0]),
			math.Max(b.Max[1], point[1]),
		},
	}
}

func (b Bound) Union(other Bound) Bound {
	if other.IsEmpty() {
		return b
	}

	nextB := b.Extend(other.Min)
	nextB = b.Extend(other.Max)
	nextB = b.Extend(other.LeftTop())
	nextB = b.Extend(other.RightBottom())

	return nextB
}

func (b Bound) Center() Point {
	return Point{
		(b.Min[0] + b.Max[0]) / 2.0,
		(b.Min[1] + b.Max[1]) / 2.0,
	}
}

func (b Bound) Pad(d float64) Bound {
	b.Min[0] -= d
	b.Min[1] -= d

	b.Max[0] += d
	b.Max[1] += d
	return b
}
