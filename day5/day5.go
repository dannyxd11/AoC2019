package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var OPCODES = map[int]int{
	1:  3, // Add
	2:  3, // Multiply
	3:  1, // Input
	4:  1, // Output
	5:  2, // Jump-if-true
	6:  2, // Jump-if-false
	7:  3, // less than
	8:  3, // equals
	99: 0, // halt
}

type Parameter struct {
	val  int
	mode int
}

type Operation struct {
	opcode  int
	params  []Parameter
	nParams int
}

func parseOp(program []int) (op Operation, isTerminated bool) {
	op = Operation{}

	sop := fmt.Sprintf("%05d", program[0])
	l := len(sop)
	opcode, err := strconv.Atoi(sop[l-2 : l])
	check(err)

	op.opcode = opcode
	op.nParams = OPCODES[opcode]

	for p, n := 1, l-3; p <= op.nParams && n >= 0; p, n = p+1, n-1 {
		mode, err := strconv.Atoi(sop[n : n+1])
		check(err)
		op.params = append(op.params, Parameter{program[p], mode})
		log.WithFields(log.Fields{
			"params":  op.params,
			"opcode":  op.opcode,
			"nParams": op.nParams,
			"n":       n,
		}).Trace("Parsing parameters & modes")
	}

	if op.opcode == 99 {
		isTerminated = true
	} else {
		isTerminated = false
	}

	return
}

func getVal(program *[]int, param Parameter) int {
	if param.mode == 0 {
		return (*program)[param.val]
	} else if param.mode == 1 {
		return param.val
	} else {
		panic(fmt.Sprintf("Unrecognised mode: %d", param.mode))
	}
}

func setVal(program *[]int, index int, val int) int {
	log.WithFields(log.Fields{
		"index": index,
		"val":   val,
	}).Trace("Setting Value")
	(*program)[index] = val
	return 0
}

func execInstruction(program *[]int, op Operation, ip int) int {
	log.WithFields(log.Fields{"op": op, "ip": ip}).Debug("Executing Operation")
	ip += op.nParams + 1
	if op.opcode == 1 {
		setVal(program,
			op.params[2].val,
			getVal(program, op.params[0])+getVal(program, op.params[1]),
		)
	} else if op.opcode == 2 {
		setVal(program,
			op.params[2].val,
			getVal(program, op.params[0])*getVal(program, op.params[1]),
		)
	} else if op.opcode == 3 {
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("For Part 1 provide 1, for Part 2 provide 5. 1 digit only.\n> ")

		input, err := buf.ReadByte()
		check(err)
		iInput, err := strconv.Atoi(string(input))
		check(err)

		setVal(program,
			op.params[0].val,
			iInput,
		)
	} else if op.opcode == 4 {
		log.WithFields(log.Fields{
			"out": getVal(program, op.params[0]),
		}).Info("Operation 4 Output")
	} else if op.opcode == 5 {
		if getVal(program, op.params[0]) != 0 {
			ip = getVal(program, op.params[1])
		}
	} else if op.opcode == 6 {
		if getVal(program, op.params[0]) == 0 {
			ip = getVal(program, op.params[1])
		}
	} else if op.opcode == 7 {
		var v = -1
		if getVal(program, op.params[0]) < getVal(program, op.params[1]) {
			v = 1
		} else {
			v = 0
		}
		setVal(program, op.params[2].val, v)
	} else if op.opcode == 8 {
		var v = -1
		if getVal(program, op.params[0]) == getVal(program, op.params[1]) {
			v = 1
		} else {
			v = 0
		}
		setVal(program, op.params[2].val, v)
	} else {
		panic(fmt.Sprintf("Unrecognised opcode: %d", op.opcode))
	}
	log.WithFields(log.Fields{"op": op, "ip": ip}).Debug("Executed Operation")
	return ip
}

func execute(program *[]int) (output []int) {

	for i := 0; i < len(*program); { //i++{
		log.WithFields(log.Fields{
			"i":              i,
			"program[i:i+4]": (*program)[i : i+4],
		}).Trace("Parsing op")
		op, isTerminated := parseOp((*program)[i:])
		log.WithFields(log.Fields{
			"i":              i,
			"program[i:i+4]": (*program)[i : i+4],
			"op":             op,
			"isTerminated":   isTerminated,
		}).Trace("Parsed op")
		if isTerminated {
			break
		}

		i = execInstruction(program, op, i)
	}

	return
}

func loadAndRun(file io.ReadSeeker) {
	_, err := file.Seek(0, io.SeekStart)
	check(err)

	scanner := bufio.NewScanner(file)

	// Move to first line
	scanner.Scan()

	sProgram := strings.Split(scanner.Text(), ",")

	// Cast elements to int
	var program []int
	for _, i := range sProgram {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		program = append(program, j)
	}

	execute(&program)
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

	loadAndRun(file)
}
