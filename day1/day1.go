package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Fuel required to launch a given module is based on its mass. Specifically,
// to find the fuel required for a module,
// take its mass, divide by three, round down, and subtract 2.
func calculateFuel(mass int) (fuel int){
	fuel = int(math.Max(math.Floor(float64(mass / 3.0) - 2), 0));
	return;
}

func part1(file *os.File){
	file.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(file)

	var totalFuelRequired = 0;
	for scanner.Scan() {
		mass, err := strconv.Atoi(scanner.Text());
		check(err);

		totalFuelRequired += calculateFuel(mass);
	}

	fmt.Println("Part 1: ", totalFuelRequired);
}

func part2(file *os.File){
	file.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(file)

	var totalFuelRequired = 0;
	for scanner.Scan() {
		mass, err := strconv.Atoi(scanner.Text());
		check(err);


		fuel := calculateFuel(mass);
		totalFuelRequired += fuel;
		for fuel > 0 {
			fuel = calculateFuel(fuel);
			totalFuelRequired += fuel;
		}
	}


	fmt.Println("Part 2: ", totalFuelRequired);
}

func main() {

	file, err := os.Open("./challenge.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	part1(file);
	part2(file);

}
