package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Orbits struct {
	object   string
	parent   string
	children []string
}

var MAP = map[string]Orbits{}

func updateMap(record string) {
	parts := strings.Split(record, ")")
	if len(parts) > 2 {
		log.WithFields(log.Fields{
			"parts":  parts,
			"record": record,
		}).Error("Unexpected number of parts..")
	}

	orbiting := parts[0]
	object := parts[1]

	log.WithFields(log.Fields{"object": object, "orbiting": orbiting, "map": MAP}).Debug("Before adding object")

	// Put the object in the map..
	orbit := Orbits{object, orbiting, make([]string, 1)}
	if _, ok := MAP[object]; !ok {
		MAP[object] = orbit
	} else {
		objectTmp := MAP[object]
		objectTmp.parent = orbiting
		MAP[object] = objectTmp
	}

	log.WithFields(log.Fields{"object": object, "orbiting": orbiting, "map": MAP}).Debug("Before adding parent")
	// Update the object it's orbiting
	if _, ok := MAP[orbiting]; !ok {
		orbits := append(make([]string, 1), object)
		MAP[orbiting] = Orbits{orbiting, "", orbits}
	} else {
		objectTmp := MAP[orbiting]
		objectTmp.children = append(objectTmp.children, object)
		MAP[orbiting] = objectTmp
	}
	log.WithFields(log.Fields{"object": object, "orbiting": orbiting, "map": MAP}).Debug("After updating")
}

func provideParentChain(object string) (parents []string) {
	parents = make([]string, 0)
	for parent := MAP[object].parent; object != "COM"; {
		parents = append(parents, parent)
		object = parent
		parent = MAP[object].parent
	}
	return
}

func calculateOrbits(object string, depth int) int {
	orbits := 0
	if len(MAP[object].children) > 0 {
		for i := 0; i < len(MAP[object].children); i++ {
			orbits += calculateOrbits(MAP[object].children[i], depth+1)
		}
		return orbits
	} else {
		return depth - 1
	}
}

func find(needle string, haystack []string) (exists bool, index int, val string) {
	for i, v := range haystack {
		if needle == v {
			return true, i, v
		}
	}
	return false, -1, ""
}

func loadAndRun(file io.ReadSeeker) {
	_, err := file.Seek(0, io.SeekStart)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		updateMap(scanner.Text())
	}

	fmt.Println("Part 1:", calculateOrbits("COM", 0))

	san := provideParentChain("SAN")
	you := provideParentChain("YOU")

	for i, v := range you {
		exists, index, val := find(v, san)

		if exists {
			fmt.Printf("\nPart 2: %d\nFound closest common planet: %s\nFrom YOU: %d\nFrom SAN: %d\nTotal 'orbital transfers': %d", i+index, val, i, index, i+index)
			break
		}
	}
}

func main() {
	log.SetLevel(log.InfoLevel)

	file, err := os.Open("./challenge.txt")
	check(err)
	defer func() {
		err := file.Close()
		check(err)
	}()

	loadAndRun(file)
}
