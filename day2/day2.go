package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Operation struct {
	opcode int
	opa  int
	opb int
	opc int
}

type pair [2]int

func cartesian(a, b []int) []pair {
	p := make([]pair, len(a)*len(b))
	i := 0
	for _, a := range a {
		for _, b := range b {
			p[i] = pair{a, b}
			i++
		}
	}
	return p
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func parseOp(ops []int) (op Operation, isTerminated bool){
	opcode := ops[0]
	opa := ops[1]
	opb := ops[2]
	opc := ops[3]

	if opcode == 99 {
		isTerminated = true
	} else {
		isTerminated = false
	}

	op = Operation{opcode, opa, opb, opc}
	return
}

func execInstruction(program *[]int, op Operation){
	if op.opcode == 1 {
		(*program)[op.opc] = (*program)[op.opa] + (*program)[op.opb]
	} else if op.opcode == 2 {
		(*program)[op.opc] = (*program)[op.opa] * (*program)[op.opb]
	} else if op.opcode == 99 {
		return
	} else {
		panic(fmt.Sprintf("Unrecognised opcode: %d", op.opcode))
	}
}

func execute(program *[]int) (output []int){

	for i := 0; i < len(*program); i += 4 {
		op, isTerminated := parseOp((*program)[i:i+4])

		if isTerminated {
			break
		}

		execInstruction(program, op)
	}

	return
}

func part1(file io.ReadSeeker){
	file.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(file)

	// Move to first line
	scanner.Scan()

	program_s := strings.Split(scanner.Text(), ",")

	// Cast elements to int
	var program = []int{}
	for _, i := range program_s {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		program = append(program, j)
	}

	//fmt.Println("[Part 1] Original Program:",program);
	program[1] = 12
	program[2] = 2

	execute(&program)

	fmt.Println("[Part 1] Result:", program[0])
}

func part2(file io.ReadSeeker){
	file.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(file)

	// Move to first line
	scanner.Scan()

	program_s := strings.Split(scanner.Text(), ",")

	// Cast elements to int
	var program = []int{}
	for _, i := range program_s {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		program = append(program, j)
	}

	inputs := makeRange(0, 99)
	input_pairs := cartesian(inputs, inputs)

	//fmt.Println("[Part 2] Original Program:",program);
	for _, v := range input_pairs {
		var i_program = make([]int, len(program))
		copy(i_program, program)

		noun, verb := v[0], v[1]
		i_program[1] = noun
		i_program[2] = verb

		execute(&i_program)

		result := i_program[0]

		//fmt.Printf("[Part 2] Noun: %d, Verb: %d, Result: %d\n", noun, verb, result);
		if result == 19690720 {
			fmt.Printf("[Part 2] 100 * noun [%d] + verb [%d] = %d", noun, verb, 100 * noun + verb)
			break
		}
	}

}


func main() {

	file, err := os.Open("./challenge.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	part1(file)
	part2(file)

}
