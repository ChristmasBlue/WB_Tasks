package point

import (
	"math"
)

type Point struct {
	x float64
	y float64
}

// NewPoint конструктор
func NewPoint(x, y float64) *Point {
	return &Point{x: x, y: y}
}

// Distance функция подсчёта растояния между точками на плоскости,
// использовал теорему пифагора и формулу: a² + b² = c²
func (p *Point) Distance(secondPoint *Point) float64 {
	cat1 := p.x - secondPoint.x
	cat2 := p.y - secondPoint.y
	return math.Sqrt(cat1*cat1 + cat2*cat2)
}
