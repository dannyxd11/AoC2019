package main

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Coordinate struct {
	x int
	y int
	i int
}

func abs(val int) (out int){
	out = int(math.Abs(float64(val)))
	return
}

type Coordinates []Coordinate
type CoordinatesMHSort Coordinates

func (a CoordinatesMHSort) Len() int           { return len(a) }
func (a CoordinatesMHSort) Less(i, j int) bool {
	if abs(a[i].x) + abs(a[i].y) <  abs(a[j].x) + abs(a[j].y) {
		return true
	} else if abs(a[i].x) + abs(a[i].y) ==  abs(a[j].x) + abs(a[j].y){
		if abs(a[i].x) <  abs(a[j].x) {
			return true
		} else if   abs(a[i].x) ==  abs(a[j].x) {
			return  abs(a[i].y) <  abs(a[j].y)
		}
	}
	return false
}
func (a CoordinatesMHSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type Intersect struct {
	x int
	y int
	w1i int
	w2i int
}
type IntersectStepSort []Intersect

func (a IntersectStepSort) Len() int           { return len(a) }
func (a IntersectStepSort) Less(i, j int) bool { return a[i].w1i + a[i].w2i < a[j].w1i + a[j].w2i  }
func (a IntersectStepSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }


func calcSteps(currentPos Coordinate, step string) (newPos Coordinate, steps []Coordinate){
	direction := string(step[0])
	stepCount, err := strconv.Atoi(step[1:])
	check(err)

	val := 1
	if direction == "L" || direction == "D"{
		val = -1
	}

	nextStep := Coordinate{currentPos.x, currentPos.y, currentPos.i}
	switch {
	case direction == "R" || direction == "L":
		for i := 0; i < stepCount; i++ {
			nextStep = Coordinate{nextStep.x + val, nextStep.y, nextStep.i + 1}
			log.WithFields(log.Fields{
				"direction": direction,
				"val":       val,
				"stepCount": stepCount,
				"i":         i,
				"nextStep":  nextStep,
			}).Trace("calcSteps:73")
			steps = append(steps, nextStep)
		}
	case direction == "U" || direction == "D":
		for i := 0; i < stepCount; i++ {
			nextStep = Coordinate{nextStep.x, nextStep.y + val, nextStep.i + 1}
			log.WithFields(log.Fields{
				"direction": direction,
				"val":       val,
				"stepCount": stepCount,
				"i":         i,
				"nextStep":  nextStep,
			}).Trace("calcSteps:86")
			steps = append(steps, nextStep)
		}
	default:
		panic("Uhh.. what direction is this?!")
	}

	newPos = nextStep
	return
}

func plotRoute(wire []string)(route []Coordinate){
	route = []Coordinate{}
	pos := Coordinate{0,0, 0}

	for _, v := range wire {
		var steps []Coordinate
		pos, steps = calcSteps(pos, v)
		log.WithFields(log.Fields{
			"pos": pos,
			"steps": steps,
			"instruction": v,
		}).Debug("plotRoute:107")
		route = append(route, steps...)
	}

	return
}

func find(needle Coordinate, haystack Coordinates) (exists bool, val Coordinate){
	for _, v := range haystack {
		if needle.x == v.x && needle.y == v.y {
			exists = true
			val = v
			return
		}

		if abs(needle.x) + abs(needle.y) < abs(v.x) + abs(v.y){
			exists = false
			return
		}
	}
	return
}

func part1and2(file *os.File){
	_, err := file.Seek(0, io.SeekStart)
	check(err)
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	wire1 := strings.Split(scanner.Text(), ",")
	wire1route := plotRoute(wire1)

	sortedWire1route := make(Coordinates, len(wire1route))
	copy(sortedWire1route, wire1route)

	sort.Sort(CoordinatesMHSort(sortedWire1route))
	log.WithFields(log.Fields{
		"sortedWire2route": sortedWire1route,
	}).Debug("Sorted Route")
	log.WithFields(log.Fields{
		"len": len(sortedWire1route),
		"cap": cap(sortedWire1route),
	}).Debug("part1:155" )
	for i, v := range sortedWire1route {
		log.WithFields(log.Fields{
			"index": i,
			"value": v,
		}).Trace("part1:160")
	}

	scanner.Scan()
	wire2 := strings.Split(scanner.Text(), ",")
	wire2route := plotRoute(wire2)

	sortedWire2route := make(Coordinates, len(wire2route))
	copy(sortedWire2route, wire2route)
	sort.Sort(CoordinatesMHSort(sortedWire2route))
	log.WithFields(log.Fields{
		"sortedWire2route": sortedWire2route,
	}).Debug("Sorted Route")
	log.WithFields(log.Fields{
		"len": len(sortedWire2route),
		"cap": cap(sortedWire2route),
	}).Debug("part1:176")
	for i, v := range sortedWire2route {
		log.WithFields(log.Fields{
			"index": i,
			"value": v,
		}).Trace("part1:181")
	}

	var intersections []Intersect
	for _, v := range sortedWire1route {
		exists, found := find(v, sortedWire2route)
		if exists {
			intersect := Intersect{v.x,v.y,v.i,found.i}
			log.Info("[Part 1] Found! ", intersect)
			intersections = append(intersections, intersect)
			//break; // Can return here since they're sorted based on manhattan length
		}
	}

	if len(intersections) < 1{
		log.Info("No Intersections found")
		return
	}

	log.WithFields(log.Fields{
		"Manhattan Length": abs(intersections[0].x) + abs(intersections[0].y),
		"Intersect": intersections[0],
	}).Info("Part 1 Result")

	sort.Sort(IntersectStepSort(intersections));


	log.WithFields(log.Fields{
		"Combined Step Count": intersections[0].w1i + intersections[0].w2i,
		"Intersect": intersections[0],
	}).Info("Part 2 Result")
}

func part2only(file *os.File){
	_, err := file.Seek(0, io.SeekStart)
	check(err)

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	wire1 := strings.Split(scanner.Text(), ",")
	wire1route := plotRoute(wire1)

	sortedWire1route := make(Coordinates, len(wire1route))
	copy(sortedWire1route, wire1route)

	sort.Sort(CoordinatesMHSort(sortedWire1route))
	log.WithFields(log.Fields{
		"sortedWire2route": sortedWire1route,
	}).Debug("sortedroute")
	log.WithFields(log.Fields{
		"len": len(sortedWire1route),
		"cap": cap(sortedWire1route),
	}).Debug("part1:148" )
	for i, v := range sortedWire1route {
		log.WithFields(log.Fields{
			"index": i,
			"value": v,
		}).Trace("part1:145")
	}

	scanner.Scan()
	wire2 := strings.Split(scanner.Text(), ",")
	wire2route := plotRoute(wire2)

	sortedWire2route := make(Coordinates, len(wire2route))
	copy(sortedWire2route, wire2route)
	sort.Sort(CoordinatesMHSort(sortedWire2route))
	log.WithFields(log.Fields{
		"sortedWire2route": sortedWire2route,
	}).Debug("sortedroute")
	log.WithFields(log.Fields{
		"len": len(sortedWire2route),
		"cap": cap(sortedWire2route),
	}).Debug("part1:590")
	for i, v := range sortedWire2route {
		log.WithFields(log.Fields{
			"index": i,
			"value": v,
		}).Trace("part1:164")
	}

	var intersections []Intersect
	for _, v := range sortedWire1route {
		exists, found := find(v, sortedWire2route)
		if exists {
			intersect := Intersect{v.x,v.y,v.i,found.i}
			log.Info("[Part 2] Found! ", intersect)
			intersections = append(intersections, intersect)
		}
	}

	sort.Sort(IntersectStepSort(intersections));

	log.WithFields(log.Fields{
		"Combined Step Count": intersections[0].w1i + intersections[0].w2i,
		"Intersect": intersections[0],
	}).Info("Part 2 Result")
}

func main() {
	log.SetLevel(log.DebugLevel)

	file, err := os.Open("./challenge.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	part1and2(file)
	// part2only(file);

}
