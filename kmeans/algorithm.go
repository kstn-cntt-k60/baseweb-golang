package kmeans

import (
	"fmt"
	"log"
	"math"
	"math/rand"
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

	copy := make([]Point, 0)
	for _, p := range points {
		copy = append(copy, p)
	}

	// rand.Seed(1000)
	rand.Shuffle(len(copy), func(i, j int) {
		copy[i], copy[j] = copy[j], copy[i]
	})

	for i := 0; i < n; i++ {
		result[i] = copy[i]
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

func distance2(a, b Point) float32 {
	ans := math.Pow(float64(a.Lat-b.Lat), 2) + math.Pow(float64(a.Long-b.Long), 2)
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
		if count != 0 {
			result[i] = Point{Lat: sum.Lat / float32(count), Long: sum.Long / float32(count)}
		}
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

	count := 0
	for {
		count++
		newCenters := kmeansStep(points, centerPoints)
		if count >= 100 || pointsEqual(newCenters, centerPoints) {
			break
		}
		centerPoints = newCenters
	}
	log.Println(count)

	return centerPoints, collectNeighbors(points, centerPoints)
}
