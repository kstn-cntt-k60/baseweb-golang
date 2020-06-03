package kmeans

import (
	"fmt"
	"math"
)

type Point struct {
	Lat, Long float32
}

type Neighbor struct {
	Index int
	Point Point
}

func randomCenterPoints(points []Point, n int) []Point {
	if n > len(points) {
		n = len(points)
	}

	result := make([]Point, n)

	for i := 0; i < n; i++ {
		result[i] = points[i]
	}

	return result
}

func toRadians(degree float32) float64 {
	return float64(degree) * math.Pi / 180.0
}

func distance(a, b Point) float32 {
	lat1 := toRadians(a.Lat)
	long1 := toRadians(a.Long)

	lat2 := toRadians(b.Lat)
	long2 := toRadians(b.Long)

	dlong := long2 - long1
	dlat := lat2 - lat1

	ans := math.Pow(math.Sin(dlat/2), 2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Pow(math.Sin(dlong/2), 2)

	ans = 2 * math.Asin(math.Sqrt(ans))
	var R float64 = 6371.0
	ans = ans * R
	return float32(ans)
}

func add(a, b Point) Point {
	return Point{Lat: a.Lat + b.Lat, Long: a.Long + b.Long}
}

func minDistance(a Point, centerPoints []Point) int {
	minIndex := 0
	minD := distance(a, centerPoints[0])

	for i := 1; i < len(centerPoints); i++ {
		d := distance(a, centerPoints[i])
		if minD > d {
			minIndex = i
			minD = d
		}
	}

	return minIndex
}

func collectNeighbors(points []Point, centerPoints []Point) []Neighbor {
	neighbors := make([]Neighbor, 0)

	for _, p := range points {
		index := minDistance(p, centerPoints)
		n := Neighbor{Index: index, Point: p}
		neighbors = append(neighbors, n)
	}

	return neighbors
}

func kmeansStep(points []Point, centerPoints []Point) []Point {
	n := len(centerPoints)
	result := make([]Point, n)

	neighbors := collectNeighbors(points, centerPoints)

	for i := 0; i < n; i++ {
		sum := Point{Lat: 0, Long: 0}
		count := 0
		for _, nb := range neighbors {
			if nb.Index == i {
				sum = add(sum, nb.Point)
				count++
			}
		}
		result[i] = Point{Lat: sum.Lat / float32(count), Long: sum.Long / float32(count)}
	}

	fmt.Println("STEP", result)

	return result
}

var epsilon float32 = 0.00001

func abs(x float32) float32 {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func pointEquals(a, b Point) bool {
	return abs(a.Lat-b.Lat) < epsilon && abs(a.Long-b.Long) < epsilon
}

func pointsEqual(a, b []Point) bool {
	for i := 0; i < len(a); i++ {
		if !pointEquals(a[i], b[i]) {
			return false
		}
	}
	return true
}

func Kmeans(points []Point, n int) ([]Point, []Neighbor) {
	centerPoints := randomCenterPoints(points, n)

	for {
		newCenters := kmeansStep(points, centerPoints)
		if pointsEqual(newCenters, centerPoints) {
			break
		}
		centerPoints = newCenters
	}

	return centerPoints, collectNeighbors(points, centerPoints)
}
