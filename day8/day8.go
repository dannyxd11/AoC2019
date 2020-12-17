package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
)

const ( // iota is reset to 0
	BLACK = iota
	WHITE
	TRANSPARENT
)

type Image struct {
	encoding [][][]int
	height   int
	width    int
	depth    int
	length   int
}

func (a Image) LayerSize() int { return a.height * a.width }
func (a Image) IndexToPoints(i int) (int, int, int) {
	return i % a.width, i / a.width % a.height, i / a.width / a.height
}
func (a Image) GetElement(x int, y int, z int) int    { return a.encoding[x][y][z] }
func (a Image) SetElement(x int, y int, z int, v int) { a.encoding[x][y][z] = v }
func (a Image) GetElementByInd(i int) int             { return a.GetElement(a.IndexToPoints(i)) }
func (a Image) SetElementByInd(i int, v int) {
	x, y, z := a.IndexToPoints(i)
	a.SetElement(x, y, z, v)
}
func (a Image) CountDigitsOnLayer(layer int, digit int) (count int) {
	count = 0
	for i := 0; i < a.LayerSize(); i++ {
		if a.GetElementByInd(a.LayerSize()*layer+i) == digit {
			count += 1
		}
	}
	return
}
func (a Image) GetColourAtPoint(x int, y int) (pixel int) {
	for z := 0; z < a.depth; z++ {
		pixel = a.GetElement(x, y, z)
		if pixel != TRANSPARENT {
			return pixel
		}
	}
	return
}
func (a Image) Print() {
	for i := 0; i < a.LayerSize(); i++ {
		x, y, _ := a.IndexToPoints(i)
		colour := a.GetColourAtPoint(x, y)

		reset := "\033[0m"
		white := "\033[31m" //"\033[36m"
		black := "\033[34m"

		if i%a.width == 0 {
			fmt.Println()
		}

		if colour == WHITE {
			fmt.Print(white, "#", reset)
		} else if colour == BLACK {
			fmt.Print(black, "-", reset)
		} else {
			fmt.Print(" ")
		}

	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func make3D(m, n, p int) [][][]int {
	buf := make([]int, m*n*p)

	x := make([][][]int, m)
	for i := range x {
		x[i] = make([][]int, n)
		for j := range x[i] {
			x[i][j] = buf[:p:p]
			buf = buf[p:]
		}
	}
	return x
}

func load(file io.ReadSeeker) (image Image) {
	_, err := file.Seek(0, io.SeekStart)
	check(err)

	fs := bufio.NewScanner(file)
	fs.Scan()
	rawEncoding := fs.Text()

	width, height, length := 25, 6, len(fs.Bytes())
	depth := length / width / height
	encoding := make3D(width, height, depth)
	image = Image{encoding: encoding, height: height, width: width, depth: depth, length: length}

	log.WithFields(log.Fields{
		"width":  width,
		"height": height,
		"length": len(fs.Bytes()),
		"depth":  depth,
	}).Debug("Creating 3D array")

	for i, v := range rawEncoding {
		vi, _ := strconv.Atoi(string(v))
		image.SetElementByInd(i, vi)
	}

	return image
}

func part1(file io.ReadSeeker) {
	image := load(file)

	mLayer := -1
	mZeros := image.length
	count := 0
	for z := 0; z < image.depth; z++ {
		count = image.CountDigitsOnLayer(z, 0)

		log.WithFields(log.Fields{"Layer": z, "Count": count}).Debug("Layer count")
		if count < mZeros {
			mLayer = z
			mZeros = count
		}
	}

	ones := image.CountDigitsOnLayer(mLayer, 1)
	twos := image.CountDigitsOnLayer(mLayer, 2)
	log.WithFields(log.Fields{
		"Layer":    mLayer,
		"Count 0s": mZeros,
		"Count 1s": ones,
		"Count 2s": twos,
		"Answer":   twos * ones,
	}).Info("Part 1")
}

func part2(file io.ReadSeeker) {
	image := load(file)
	image.Print()
}

func main() {
	log.SetLevel(log.InfoLevel)

	file, err := os.Open("./challenge.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := file.Close()
		check(err)
	}()

	part1(file)
	part2(file)
}
