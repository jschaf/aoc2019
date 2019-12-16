package geom

func NewPoint(x, y int) Point {
	return Point{
		X: x,
		Y: y,
	}
}

type Point struct{ X, Y int }
