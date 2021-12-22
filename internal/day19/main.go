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

var rotations = []mat.Matrix{
	mat.NewDense(3, 3, []float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}),
	mat.NewDense(3, 3, []float64{
		1, 0, 0,
		0, 0, -1,
		0, 1, 0,
	}),
	mat.NewDense(3, 3, []float64{
		1, 0, 0,
		0, -1, 0,
		0, 0, -1,
	}),
	mat.NewDense(3, 3, []float64{
		1, 0, 0,
		0, 0, 1,
		0, -1, 0,
	}),
	mat.NewDense(3, 3, []float64{
		0, -1, 0,
		1, 0, 0,
		0, 0, 1,
	}),
	mat.NewDense(3, 3, []float64{
		0, 0, 1,
		1, 0, 0,
		0, 1, 0,
	}),
	mat.NewDense(3, 3, []float64{
		0, 1, 0,
		1, 0, 0,
		0, 0, -1,
	}),
	mat.NewDense(3, 3, []float64{
		0, 0, -1,
		1, 0, 0,
		0, -1, 0,
	}),
	mat.NewDense(3, 3, []float64{
		-1, 0, 0,
		0, -1, 0,
		0, 0, 1,
	}),
	mat.NewDense(3, 3, []float64{
		-1, 0, 0,
		0, 0, -1,
		0, -1, 0,
	}),
	mat.NewDense(3, 3, []float64{
		-1, 0, 0,
		0, 1, 0,
		0, 0, -1,
	}),
	mat.NewDense(3, 3, []float64{
		-1, 0, 0,
		0, 0, 1,
		0, 1, 0,
	}),
	mat.NewDense(3, 3, []float64{
		0, 1, 0,
		-1, 0, 0,
		0, 0, 1,
	}),
	mat.NewDense(3, 3, []float64{
		0, 0, 1,
		-1, 0, 0,
		0, -1, 0,
	}),
	mat.NewDense(3, 3, []float64{
		0, -1, 0,
		-1, 0, 0,
		0, 0, -1,
	}),
	mat.NewDense(3, 3, []float64{
		0, 0, -1,
		-1, 0, 0,
		0, 1, 0,
	}),
	mat.NewDense(3, 3, []float64{
		0, 0, -1,
		0, 1, 0,
		1, 0, 0,
	}),
	mat.NewDense(3, 3, []float64{
		0, 1, 0,
		0, 0, 1,
		1, 0, 0,
	}),
	mat.NewDense(3, 3, []float64{
		0, 0, 1,
		0, -1, 0,
		1, 0, 0,
	}),
	mat.NewDense(3, 3, []float64{
		0, -1, 0,
		0, 0, -1,
		1, 0, 0,
	}),
	mat.NewDense(3, 3, []float64{
		0, 0, -1,
		0, -1, 0,
		-1, 0, 0,
	}),
	mat.NewDense(3, 3, []float64{
		0, -1, 0,
		0, 0, 1,
		-1, 0, 0,
	}),
	mat.NewDense(3, 3, []float64{
		0, 0, 1,
		0, 1, 0,
		-1, 0, 0,
	}),
	mat.NewDense(3, 3, []float64{
		0, 1, 0,
		0, 0, -1,
		-1, 0, 0,
	}),
}

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

type Scanner struct {
	id      string
	beacons []*Vector
}

func (s *Scanner) VectorDistances() map[float64][]*Vector {
	distances := make(map[float64][]*Vector)
	cs := combin.Combinations(len(s.beacons), 2)
	for _, c := range cs {
		b1 := s.beacons[c[0]]
		b2 := s.beacons[c[1]]
		distances[Distance(b1, b2)] = []*Vector{b1, b2}
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

// func (s *Scanner) VectorSet() mapset.Set {
// 	beacons := make([]interface{}, len(s.beacons))
// 	for i := range s.beacons {
// 		beacons[i] = s.beacons[i]
// 	}
// 	beaconSet := mapset.NewSetFromSlice(beacons)
// 	return beaconSet
// }

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
func FindTransform(v0 *Vector, v1 *Vector) (mat.Matrix, bool) {
	// fmt.Println("Try to match", v0)
	for _, r := range rotations {
		vector := Rotate(r, v1)
		if vector.Equals(v0) {
			// fmt.Println("Found vector match:", vector)
			return r, true
		}
	}
	return nil, false
}

// func Rotations(r func(float64) mat.Matrix) []mat.Matrix {
// 	var rotations []mat.Matrix
// 	for theta := float64(0); theta < float64(2)*math.Pi; theta += math.Pi / 2 {
// 		rotations = append(rotations, r(theta))
// 	}
// 	return rotations
// }

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

func FindRotation(s1 *Scanner, s2 *Scanner) (mat.Matrix, [][]*Vector) {
	beaconPairs := VectorPairs(s1, s2)
	foundRotations := make(map[mat.Matrix][][][]*Vector)
	for _, pair := range beaconPairs {
		v0 := pair[0][1].Subtract(pair[0][0])
		v1 := pair[1][1].Subtract(pair[1][0])
		r, found := FindTransform(&v0, &v1)
		if found {
			// fmt.Println("found rotation:", r, pair)
			foundRotations[r] = append(foundRotations[r], pair)
		}
	}
	maxPairCount := 0
	var maxR mat.Matrix
	for r, pairs := range foundRotations {
		if len(pairs) > maxPairCount {
			maxPairCount = len(pairs)
			maxR = r
		}
	}
	fmt.Println("chose rotation:", maxR, foundRotations[maxR][0])
	return maxR, foundRotations[maxR][0]
}

// type ScannerParent struct {
// 	id        string
// 	rotations [][]mat.Matrix
// }

// func (sp *ScannerParent) Apply(v *Vector) Vector {
// 	newVector := *v
// 	for ri := len(sp.rotations) - 1; ri >= 0; ri-- {
// 		newVector = newVector.Apply(sp.rotations[ri])
// 	}
// 	return newVector
// }

type ScannerNode struct {
	scanner          *Scanner
	beaconSet        mapset.Set
	parent           *ScannerNode
	children         []*ScannerNode
	scannerVectors   map[string]*Vector
	scannerRotations map[string]mat.Matrix
	position         *Vector
}

func BuildScannerTree(scanners []*Scanner, graph map[string][]string) *ScannerNode {
	root := &ScannerNode{
		scanner:          scanners[0],
		beaconSet:        mapset.NewSet(),
		scannerVectors:   make(map[string]*Vector),
		scannerRotations: make(map[string]mat.Matrix),
		position:         &Vector{0, 0, 0},
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
					scanner:          scanners[i],
					beaconSet:        mapset.NewSet(),
					parent:           parent,
					scannerVectors:   make(map[string]*Vector),
					scannerRotations: make(map[string]mat.Matrix),
				}
				parent.children = append(parent.children, child)
				q.PushBack(child)
			}

		}
	}
	return root
}

func TransformBeacons(s1 *Scanner, s2 *Scanner) (mat.Matrix, *Vector) {
	// Find list of rotations from s2 -> s1
	r, beaconPair := FindRotation(s1, s2)
	fmt.Println("Rotation:")
	fmt.Println(r)
	s1_b0, _ := beaconPair[0][0], beaconPair[0][1]
	s2_b0, _ := beaconPair[1][0], beaconPair[1][1]
	// Invert the first beacon's coordinates to get the position of the target scanner relative to the beacon
	inverse := &Vector{-1, -1, -1}
	s2_b0_rel_s2 := inverse.Multiply(s2_b0)
	// Transform the beacon->scanner vector so its relative to the first beacon of the base scanner
	s2_b0_rel_s1 := Rotate(r, &s2_b0_rel_s2)
	// Add the transformed vector to the first beacon of the base scanner to get the position of the target scanner relative to the base scanner
	s1_rel_s0 := s1_b0.Add(&s2_b0_rel_s1)
	// fmt.Println("scanner", s2.id, "relative to", s1.id, s1_rel_s0)
	return r, &s1_rel_s0
}

func TraverseScannerTree(seen mapset.Set, node *ScannerNode) {
	seen.Add(node)
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
		if !seen.Contains(child) {
			TraverseScannerTree(seen, child)
			r, s1_rel_s0 := TransformBeacons(node.scanner, child.scanner)
			for _, el := range child.beaconSet.ToSlice() {
				b := el.(Vector)
				newVector := Rotate(r, &b)
				newVector = s1_rel_s0.Add(&newVector)
				// fmt.Println(b, "->", newVector)
				node.beaconSet.Add(newVector)
			}
			node.scannerRotations[child.scanner.id] = r
			node.scannerVectors[child.scanner.id] = s1_rel_s0
			for gcid, gcv := range child.scannerVectors {
				newVector := Rotate(r, gcv)
				newVector = s1_rel_s0.Add(&newVector)
				node.scannerVectors[gcid] = &newVector
			}
		}
	}
	// fmt.Println("total beacons for scanner", node.scanner.id, "=", node.beaconSet.Cardinality())
	fmt.Println("child scanners relative to", node.scanner.id, ":")
	for id, v := range node.scannerVectors {
		fmt.Println(id, v)
	}
}

// func PositionScanners(node *ScannerNode) {
// 	for _, child := range node.children {
// 		r := node.scannerRotations[child.scanner.id]
// 		// var m mat.Dense
// 		// m.Copy(r)
// 		// var inv mat.Dense
// 		// inv.Inverse(&m)
// 		newVector := Rotate(r, node.scannerVectors[child.scanner.id])
// 		position := node.position.Add(&newVector)
// 		child.position = &position
// 		fmt.Println("scanner", child.scanner.id, "new position", child.position)
// 		PositionScanners(child)
// 	}
// }

func Main() {
	// v1 := &Vector{68, -1246, -43}
	// v2 := &Vector{88, 113, -1104}
	// // v := v2.Subtract(v1)
	// r := mat.NewDense(3, 3, []float64{
	// 	-1, 0, 0,
	// 	0, 1, 0,
	// 	0, 0, -1,
	// })
	// // r := mat.NewDense(3, 3, []float64{
	// // 	0, 1, 0,
	// // 	0, 0, -1,
	// // 	-1, 0, 0,
	// // })
	// v := Rotate(r, v2)
	// v = v1.Add(&v)
	// fmt.Println(v)
	// return

	scanners := parseFile("test2")

	graph := make(map[string][]string)
	cs := combin.Combinations(len(scanners), 2)
	// cs := [][]int{{1, 4}}
	for _, sc := range cs {
		s1 := scanners[sc[0]]
		s2 := scanners[sc[1]]
		// fmt.Println("--- Comparing", s1.id, "to", s2.id)
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
			// fmt.Printf("Found %d pairs\n", s1set.Cardinality())
			graph[s1.id] = append(graph[s1.id], s2.id)
			graph[s2.id] = append(graph[s2.id], s1.id)
		}
	}

	for i := 0; i < len(scanners); i++ {
		id := strconv.Itoa(i)
		_, found := graph[id]
		if !found {
			log.Fatal("not found " + id)
		}
		fmt.Println(id, "=>", graph[id])
	}

	fmt.Println()
	fmt.Println()
	node := BuildScannerTree(scanners, graph)
	seen := mapset.NewSet()
	TraverseScannerTree(seen, node)
	// PositionScanners(node)
}
