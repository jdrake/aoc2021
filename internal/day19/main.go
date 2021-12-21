package day19

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/combin"
)

type Vector struct {
	x, y, z int
}

func (v *Vector) ToMatrixVector() mat.Vector {
	return mat.NewDense(3, 1, []float64{
		float64(v.x), float64(v.y), float64(v.z),
	}).ColView(0)
}

func (v *Vector) Equals(v2 *Vector) bool {
	return v2.x == v.x && v2.y == v.y && v2.z == v.z
}

func (v *Vector) Multiply(v2 *Vector) Vector {
	return Vector{
		v.x * v2.x,
		v.y * v2.y,
		v.z * v2.z,
	}
}

type Beacon struct {
	x, y, z int
}

func Distance(b1 *Beacon, b2 *Beacon) float64 {
	return math.Sqrt(math.Pow(float64(b1.x-b2.x), 2) + math.Pow(float64(b1.y-b2.y), 2) + math.Pow(float64(b1.z-b2.z), 2))
}

func (b Beacon) String() string {
	return fmt.Sprintf("(%4d, %4d, %4d)", b.x, b.y, b.z)
}

func (b *Beacon) Add(v *Vector) Beacon {
	return Beacon{v.x + b.x, v.y + b.y, v.z + b.z}
}

func (b *Beacon) Subtract(b2 *Beacon) Vector {
	return Vector{b2.x - b.x, b2.y - b.y, b2.z - b.z}
}

type Scanner struct {
	id      string
	beacons []*Beacon
}

func (s *Scanner) BeaconDistances() map[float64][]*Beacon {
	distances := make(map[float64][]*Beacon)
	for _, b1 := range s.beacons {
		for _, b2 := range s.beacons {
			if b2 != b1 {
				distances[Distance(b1, b2)] = []*Beacon{b1, b2}
			}
		}
	}
	return distances
}

func (s Scanner) String() string {
	beacons := ""
	for _, b := range s.beacons {
		beacons += fmt.Sprintf("\n%s", b)
		var distances []float64
		for _, b2 := range s.beacons {
			if b2 != b {
				distances = append(distances, Distance(b, b2))
			}
		}
		sort.Float64s(distances)
		for _, d := range distances {
			beacons += fmt.Sprintf("%12f", d)
		}
	}
	return fmt.Sprintf("--- Scanner %s ---%s\n", s.id, beacons)
}

func (s *Scanner) BeaconSet() mapset.Set {
	beacons := make([]interface{}, len(s.beacons))
	for i := range s.beacons {
		beacons[i] = s.beacons[i]
	}
	beaconSet := mapset.NewSetFromSlice(beacons)
	return beaconSet
}

func parseFile(name string) []*Scanner {
	_, fileName, _, _ := runtime.Caller(0)
	prefixPath := filepath.Dir(fileName)
	file, err := os.Open(prefixPath + "/" + name + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	var scanners []*Scanner
	var s *Scanner
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "---") {
			fields := strings.Fields(line)
			s = &Scanner{id: fields[2]}
		} else if line == "" {
			scanners = append(scanners, s)
		} else {
			var values []int
			for _, num := range strings.Split(line, ",") {
				value, _ := strconv.Atoi(num)
				values = append(values, value)
			}
			b := &Beacon{values[0], values[1], values[2]}
			s.beacons = append(s.beacons, b)
		}
	}
	scanners = append(scanners, s)
	return scanners
}

type ScannerPair struct {
	base     *Scanner
	target   *Scanner
	rotation []mat.Matrix
}

func Rotate(rotation *mat.Matrix, vector *Vector) Vector {
	v := mat.NewDense(3, 1, []float64{
		float64(vector.x), float64(vector.y), float64(vector.z),
	})
	var c mat.Dense
	c.Mul(*rotation, v)
	return Vector{
		int(math.Round(c.At(0, 0))),
		int(math.Round(c.At(1, 0))),
		int(math.Round(c.At(2, 0))),
	}
}

func xAxisRotation(theta float64) mat.Matrix {
	return mat.NewDense(3, 3, []float64{
		1, 0, 0,
		0, math.Round(math.Cos(theta)), math.Round(-math.Sin(theta)),
		0, math.Round(math.Sin(theta)), math.Round(math.Cos(theta)),
	})
}

func yAxisRotation(theta float64) mat.Matrix {
	return mat.NewDense(3, 3, []float64{
		math.Round(math.Cos(theta)), 0, math.Round(math.Sin(theta)),
		0, 1, 0,
		math.Round(-math.Sin(theta)), 0, math.Round(math.Cos(theta)),
	})

}

func zAxisRotation(theta float64) mat.Matrix {
	return mat.NewDense(3, 3, []float64{
		math.Round(math.Cos(theta)), math.Round(-math.Sin(theta)), 0,
		math.Round(math.Sin(theta)), math.Round(math.Cos(theta)), 0,
		0, 0, 1,
	})
}

// var orientations = []*Vector{
// 	{1, 1, 1},
// 	{-1, 1, 1},
// }

// type Transformation struct {
// 	// orientation *Vector
// 	rotation []mat.Matrix
// }

// func (t Transformation) String() string {
// 	return fmt.Sprintf("orientation=(%d,%d,%d) rotation=%v", t.orientation.x, t.orientation.y, t.orientation.z, mat.Formatted(*t.rotation))
// }

// func (t *Transformation) Apply(v *Vector) Vector {
// 	vector := t.orientation.Multiply(v)
// 	vector = Rotate(t.rotation, &vector)
// 	return vector
// }

// type SortedVectors []*Vector

// func (s SortedVectors) Len() int      { return len(s) }
// func (s SortedVectors) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
// func (s SortedVectors) Less(i, j int) bool {
// 	return s[i].x < s[j].x && s[i].y < s[j].y && s[i].z < s[j].z
// }

// Find the orientation vector and rotation matrix from v1 -> v0
func FindTransform(rotations [][]mat.Matrix, v0 *Vector, v1 *Vector) ([]mat.Matrix, bool) {
	// fmt.Println("Try to match", v0)
	for _, r := range rotations {
		rx, ry, rz := r[0], r[1], r[2]
		vector := Rotate(&rx, v1)
		vector = Rotate(&ry, &vector)
		vector = Rotate(&rz, &vector)
		if vector.Equals(v0) {
			// fmt.Println("Found vector match:", vector)
			return r, true
		}
	}
	return nil, false
}

// func FindBeaconTransform(s0_b0 *Beacon, s0_b1 *Beacon, s1_b0 *Beacon, s1_b1 *Beacon) (*Transformation, bool) {
// 	v0 := s0_b1.Subtract(s0_b0)
// 	// fmt.Println(v0)
// 	v1 := s1_b1.Subtract(s1_b0)
// 	// fmt.Println(v1)
// 	t, found := FindTransform(rotations, &v0, &v1)
// 	if !found {
// 		return nil, false
// 	}
// 	// fmt.Println(t)
// 	// fmt.Println()

// 	// Invert the first beacon's coordinates to get the position of the target scanner relative to the beacon
// 	s1_rel_s1_b0 := orientations[1].Multiply((*Vector)(s1_b0))
// 	// Transform the beacon->scanner vector so its relative to the first beacon of the base scanner
// 	s1_rel_s0_b0 := t.Apply(&s1_rel_s1_b0)
// 	// Add the transformed vector to the first beacon of the base scanner to get the position of the target scanner relative to the base scanner
// 	s1_rel_s0 := s0_b0.Add(&s1_rel_s0_b0)
// 	fmt.Println(s1_rel_s0)
// 	return t, true
// }

func Rotations(r func(float64) mat.Matrix) []mat.Matrix {
	var rotations []mat.Matrix
	for theta := float64(0); theta < float64(2)*math.Pi; theta += math.Pi / 2 {
		rotations = append(rotations, r(theta))
	}
	return rotations
}

func Main() {
	var rotations [][]mat.Matrix
	for _, rx := range Rotations(xAxisRotation) {
		for _, ry := range Rotations(yAxisRotation) {
			for _, rz := range Rotations(zAxisRotation) {
				rotations = append(rotations, []mat.Matrix{rx, ry, rz})
			}
		}
	}
	fmt.Println("rotation count:", len(rotations))

	scanners := parseFile("test2")

	// s0_b0 := &Beacon{-618, -824, -621}
	// s0_b1 := &Beacon{-537, -823, -458}

	// s1_b0 := &Beacon{686, 422, 578}
	// s1_b1 := &Beacon{605, 423, 415}

	graph := make(map[string][]string)
	cs := combin.Combinations(len(scanners), 2)
	// cs := [][]int{{1, 4}}
	for _, sc := range cs {
		s1 := scanners[sc[0]]
		s2 := scanners[sc[1]]
		fmt.Println("--- Comparing", s1.id, "to", s2.id)
		s1d := s1.BeaconDistances()
		s2d := s2.BeaconDistances()
		s1set := mapset.NewSet()
		s2set := mapset.NewSet()
		var beaconPairs [][][]*Beacon
		for d, s2b := range s2d {
			if s1b, found := s1d[d]; found {
				// The distance between the beacons in s1 equals the distance
				// between the beacons in s2
				s1set.Add(s1b[0])
				s1set.Add(s1b[1])
				s2set.Add(s2b[0])
				s2set.Add(s2b[1])
				// fmt.Println(d, s1b, s2b)
				beaconPairs = append(beaconPairs, [][]*Beacon{s1b, s2b})
			}
		}
		if s1set.Cardinality() >= 12 {
			fmt.Printf("Found %d pairs\n", s1set.Cardinality())
			graph[s1.id] = append(graph[s1.id], s2.id)
			graph[s2.id] = append(graph[s2.id], s1.id)
		} else {
			fmt.Println("Found <12 pairs")
		}
		fmt.Println("")
	}

	// keep track of a master map of beacons from s0 perspective
	// keep track of list of rotations to perform with each new level, e.g. [0 <- 1, 1 <- 4]
	// for each level, apply transforms in reverse order and add to master map of beacons
	fmt.Println(graph)
	q := list.New()
	q.PushBack("0")
	seen := mapset.NewSet()
	for q.Len() > 0 {
		el := q.Front()
		q.Remove(el)
		id := el.Value.(string)
		i, _ := strconv.Atoi(id)
		s1 := scanners[i]
		seen.Add(id)
		for _, neighbor := range graph[id] {
			if !seen.Contains(neighbor) {
				q.PushBack(neighbor)
			}
		}
	}

	for _, pair := range beaconPairs {
		v0 := pair[0][1].Subtract(pair[0][0])
		v1 := pair[1][1].Subtract(pair[1][0])
		rotations, found := FindTransform(rotations, &v0, &v1)
		if found {
			for _, r := range rotations {
				fmt.Println(r)
			}
			sp := &ScannerPair{
				base:     s1,
				target:   s2,
				rotation: rotations,
			}
			break
		}
	}

	// beaconCombos1 := combin.Combinations(len(s1.beacons), 2)
	// beaconCombos2 := combin.Combinations(len(s2.beacons), 2)
	// count := 0
	// var transformation *Transformation
	// for _, bc1 := range beaconCombos1 {
	// 	v0 := s1.beacons[bc1[1]].Subtract(s1.beacons[bc1[0]])
	// 	for _, bc2 := range beaconCombos2 {
	// 		v1 := s2.beacons[bc2[1]].Subtract(s2.beacons[bc2[0]])
	// 		t, found := FindTransform(&v0, &v1)
	// 		if found {
	// 			count += 1
	// 			if transformation == nil {
	// 				transformation = t
	// 			}
	// 			if count >= 12 {
	// 				break
	// 			}
	// 			fmt.Println("found", s1.beacons[bc1[0]], s1.beacons[bc1[1]], s2.beacons[bc2[0]], s2.beacons[bc2[1]])
	// 			// fmt.Println(t)
	// 			// fmt.Println()
	// 		}
	// 	}
	// 	if count >= 12 {
	// 		break
	// 	}
	// }
	// if count >= 12 {
	// 	fmt.Println("Pair")
	// }
	// fmt.Println()

	// count := 0
	// for _, scanner := range scanners {
	// 	bset := scanner.BeaconSet()
	// 	for _, sp := range pairs {
	// 		if sp.scanners[0] == scanner {
	// 			bset = bset.Difference(sp.beacons[0])
	// 		} else if sp.scanners[1] == scanner {
	// 			bset = bset.Difference(sp.beacons[1])
	// 		}
	// 	}
	// 	fmt.Println("unique", scanner.id, bset)
	// 	count += bset.Cardinality()
	// }

	// for _, sp := range pairs {
	// 	count += sp.beacons[0].Cardinality()
	// }

	// fmt.Println("beacon count", count)
}
