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

func Distance(b1 *Vector, b2 *Vector) float64 {
	return math.Sqrt(math.Pow(float64(b1.x-b2.x), 2) + math.Pow(float64(b1.y-b2.y), 2) + math.Pow(float64(b1.z-b2.z), 2))
}

func (v Vector) String() string {
	return fmt.Sprintf("(%4d, %4d, %4d)", v.x, v.y, v.z)
}

func (v *Vector) Add(v2 *Vector) Vector {
	return Vector{v2.x + v.x, v2.y + v.y, v2.z + v.z}
}

func (v *Vector) Subtract(v2 *Vector) Vector {
	return Vector{v2.x - v.x, v2.y - v.y, v2.z - v.z}
}

func (v *Vector) Apply(rotations []mat.Matrix) Vector {
	vector := *v
	for _, r := range rotations {
		vector = Rotate(r, &vector)
	}
	return vector
}

type Scanner struct {
	id      string
	beacons []*Vector
}

func (s *Scanner) VectorDistances() map[float64][]*Vector {
	distances := make(map[float64][]*Vector)
	for _, b1 := range s.beacons {
		for _, b2 := range s.beacons {
			if b2 != b1 {
				distances[Distance(b1, b2)] = []*Vector{b1, b2}
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

func (s *Scanner) VectorSet() mapset.Set {
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
			b := &Vector{values[0], values[1], values[2]}
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

func Rotate(rotation mat.Matrix, vector *Vector) Vector {
	v := mat.NewDense(3, 1, []float64{
		float64(vector.x), float64(vector.y), float64(vector.z),
	})
	var c mat.Dense
	c.Mul(rotation, v)
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
		vector := Rotate(rx, v1)
		vector = Rotate(ry, &vector)
		vector = Rotate(rz, &vector)
		if vector.Equals(v0) {
			// fmt.Println("Found vector match:", vector)
			return r, true
		}
	}
	return nil, false
}

func Rotations(r func(float64) mat.Matrix) []mat.Matrix {
	var rotations []mat.Matrix
	for theta := float64(0); theta < float64(2)*math.Pi; theta += math.Pi / 2 {
		rotations = append(rotations, r(theta))
	}
	return rotations
}

func VectorPairs(s1 *Scanner, s2 *Scanner) [][][]*Vector {
	fmt.Println("--- Comparing", s1.id, "to", s2.id)
	s1d := s1.VectorDistances()
	s2d := s2.VectorDistances()
	var beaconPairs [][][]*Vector
	for d, s2b := range s2d {
		if s1b, found := s1d[d]; found {
			// The distance between the beacons in s1 equals the distance
			// between the beacons in s2
			beaconPairs = append(beaconPairs, [][]*Vector{s1b, s2b})
		}
	}
	return beaconPairs
}

func FindRotation(rotations [][]mat.Matrix, s1 *Scanner, s2 *Scanner) ([]mat.Matrix, [][]*Vector) {
	beaconPairs := VectorPairs(s1, s2)
	for _, pair := range beaconPairs {
		v0 := pair[0][1].Subtract(pair[0][0])
		v1 := pair[1][1].Subtract(pair[1][0])
		r, found := FindTransform(rotations, &v0, &v1)
		if found {
			fmt.Println(pair)
			return r, pair
		}
	}
	return nil, nil
}

type ScannerParent struct {
	id        string
	rotations [][]mat.Matrix
}

func (sp *ScannerParent) Apply(v *Vector) Vector {
	newVector := *v
	for ri := len(sp.rotations) - 1; ri >= 0; ri-- {
		newVector = newVector.Apply(sp.rotations[ri])
	}
	return newVector
}

type ScannerNode struct {
	scanner   *Scanner
	beaconSet mapset.Set
	parent    *ScannerNode
	children  []*ScannerNode
}

func BuildScannerTree(scanners []*Scanner, graph map[string][]string) *ScannerNode {
	root := &ScannerNode{
		scanner:   scanners[0],
		beaconSet: mapset.NewSet(),
	}
	q := list.New()
	q.PushBack(root)
	seen := mapset.NewSet()
	for q.Len() > 0 {
		el := q.Front()
		q.Remove(el)
		parent := el.Value.(*ScannerNode)
		seen.Add(parent.scanner.id)
		for _, childId := range graph[parent.scanner.id] {
			if !seen.Contains(childId) {
				i, _ := strconv.Atoi(childId)
				child := &ScannerNode{
					scanner:   scanners[i],
					beaconSet: mapset.NewSet(),
					parent:    parent,
				}
				parent.children = append(parent.children, child)
				q.PushBack(child)
			}

		}
	}
	return root
}

func TransformBeacons(rotations [][]mat.Matrix, s1 *Scanner, s2 *Scanner) ([]mat.Matrix, *Vector) {
	// Find list of rotations from s2 -> s1
	r, beaconPair := FindRotation(rotations, s1, s2)
	fmt.Println("Rotations:")
	for _, rot := range r {
		fmt.Println(rot)
	}
	s1_b0, _ := beaconPair[0][0], beaconPair[0][1]
	s2_b0, _ := beaconPair[1][0], beaconPair[1][1]
	// Invert the first beacon's coordinates to get the position of the target scanner relative to the beacon
	inverse := &Vector{-1, -1, -1}
	s2_b0_rel_s2 := inverse.Multiply(s2_b0)
	// Transform the beacon->scanner vector so its relative to the first beacon of the base scanner
	s2_b0_rel_s1 := s2_b0_rel_s2.Apply(r)
	// Add the transformed vector to the first beacon of the base scanner to get the position of the target scanner relative to the base scanner
	s1_rel_s0 := s1_b0.Add(&s2_b0_rel_s1)
	fmt.Println("scanner", s2.id, "relative to", s1.id, s1_rel_s0)
	return r, &s1_rel_s0
}

func TraverseScannerTree(rotations [][]mat.Matrix, node *ScannerNode) {
	for _, b := range node.scanner.beacons {
		node.beaconSet.Add(*b)
	}
	if len(node.children) == 0 {
		return
	}
	if node.parent == nil {
		fmt.Println(node.scanner.id)
	} else {
		fmt.Println(node.scanner.id, "parent:", node.parent.scanner.id)
	}
	for _, child := range node.children {
		TraverseScannerTree(rotations, child)
		r, s1_rel_s0 := TransformBeacons(rotations, node.scanner, child.scanner)
		for _, el := range child.beaconSet.ToSlice() {
			b := el.(Vector)
			newVector := b.Apply(r)
			newVector = s1_rel_s0.Add(&newVector)
			fmt.Println(b, "->", newVector)
			node.beaconSet.Add(newVector)
		}
	}
	fmt.Println("total beacons for scanner", node.scanner.id, "=", node.beaconSet.Cardinality())
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

	scanners := parseFile("input")

	graph := make(map[string][]string)
	cs := combin.Combinations(len(scanners), 2)
	// cs := [][]int{{1, 4}}
	for _, sc := range cs {
		s1 := scanners[sc[0]]
		s2 := scanners[sc[1]]
		fmt.Println("--- Comparing", s1.id, "to", s2.id)
		s1d := s1.VectorDistances()
		s2d := s2.VectorDistances()
		s1set := mapset.NewSet()
		for d := range s2d {
			if s1b, found := s1d[d]; found {
				// The distance between the beacons in s1 equals the distance
				// between the beacons in s2
				s1set.Add(s1b[0])
				s1set.Add(s1b[1])
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

	fmt.Println()
	fmt.Println()
	node := BuildScannerTree(scanners, graph)
	TraverseScannerTree(rotations, node)

	// // keep track of a master map of beacons from s0 perspective
	// // keep track of list of rotations to perform with each new level, e.g. [0 <- 1, 1 <- 4]
	// // for each level, apply transforms in reverse order and add to master map of beacons
	// beaconMap := make(map[Vector]int)
	// fmt.Println(graph)
	// q := list.New()
	// q.PushBack(&ScannerParent{id: "0"})
	// seen := mapset.NewSet()
	// for q.Len() > 0 {
	// 	el := q.Front()
	// 	q.Remove(el)
	// 	sp := el.Value.(*ScannerParent)
	// 	i, _ := strconv.Atoi(sp.id)
	// 	s1 := scanners[i]
	// 	seen.Add(sp.id)
	// 	for _, neighbor := range graph[sp.id] {
	// 		if !seen.Contains(neighbor) {
	// 			j, _ := strconv.Atoi(neighbor)
	// 			s2 := scanners[j]
	// 			// Find list of rotations from s2 -> s1
	// 			r, beaconPair := FindRotation(rotations, s1, s2)
	// 			fmt.Println("Rotations:")
	// 			for _, rot := range r {
	// 				fmt.Println(rot)
	// 			}
	// 			s1_b0, _ := beaconPair[0][0], beaconPair[0][1]
	// 			s2_b0, _ := beaconPair[1][0], beaconPair[1][1]
	// 			// Invert the first beacon's coordinates to get the position of the target scanner relative to the beacon
	// 			inverse := &Vector{-1, -1, -1}
	// 			s2_b0_rel_s2 := inverse.Multiply(s2_b0)
	// 			// Transform the beacon->scanner vector so its relative to the first beacon of the base scanner
	// 			s2_b0_rel_s1 := s2_b0_rel_s2.Apply(r)
	// 			// s2_b0_rel_s1 = sp.Apply(&s2_b0_rel_s1)
	// 			// Add the transformed vector to the first beacon of the base scanner to get the position of the target scanner relative to the base scanner
	// 			s1_rel_s0 := s1_b0.Add(&s2_b0_rel_s1)
	// 			fmt.Println("scanner", s2.id, "relative to 0:", s1_rel_s0)
	// 			for _, b := range s2.beacons {
	// 				newVector := b.Apply(r)
	// 				newVector = sp.Apply(&newVector)
	// 				newVector = s1_rel_s0.Add(&newVector)
	// 				fmt.Println(b, "->", newVector)
	// 				beaconMap[newVector] += 1
	// 			}
	// 			parentRotations := append(sp.rotations, r)
	// 			q.PushBack(&ScannerParent{
	// 				id:        s2.id,
	// 				rotations: parentRotations,
	// 			})
	// 		}
	// 	}
	// }

	// fmt.Println("---")
	// fmt.Println("Beacons from s0:")
	// for b := range beaconMap {
	// 	fmt.Println(b)
	// }

}
