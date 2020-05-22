package geom

func NewPoint(x, y int) Point {
	return Point{
		X: x,
		Y: y,
	}
}

type Point struct{ X, Y int }

func (p Point) North() Point {
	return NewPoint(p.X, p.Y+1)
}

func (p Point) South() Point {
	return NewPoint(p.X, p.Y-1)
}

func (p Point) East() Point {
	return NewPoint(p.X+1, p.Y)
}

func (p Point) West() Point {
	return NewPoint(p.X-1, p.Y)
}
