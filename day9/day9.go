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
	9:  1, // adjust bp
	99: 0, // halt
}

type Parameter struct {
	val  float64
	mode int
}

type Output struct {
	val       float64
	hasOutput bool
}

type Operation struct {
	opcode  int
	params  []Parameter
	nParams int
}

func parseOp(program []float64) (op Operation, isTerminated bool) {
	op = Operation{}

	sop := fmt.Sprintf("%05.0f", program[0])
	l := len(sop)
	opcode, err := strconv.ParseFloat(sop[l-2:l], 64)
	check(err)

	op.opcode = int(opcode)
	op.nParams = OPCODES[int(opcode)]

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

func getInput() float64 {
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("> ")

	input, err := buf.ReadByte()
	check(err)
	iInput, err := strconv.ParseFloat(string(input), 64)
	check(err)

	return iInput
}

func getVal(program *[]float64, param Parameter, bp int) float64 {
	if param.mode == 0 {
		log.WithFields(log.Fields{
			"param":              param,
			"program[param.val]": (*program)[int(param.val)],
		}).Trace("Getting Value (Position)")
		return (*program)[int(param.val)]
	} else if param.mode == 1 {
		log.WithFields(log.Fields{
			"param":     param,
			"param.val": param.val,
		}).Trace("Getting Value (Immediate)")
		return param.val
	} else if param.mode == 2 {
		log.WithFields(log.Fields{
			"param":     param,
			"param.val": param.val,
		}).Trace("Getting Value (Immediate)")
		return (*program)[int(param.val)+bp]
	} else {
		panic(fmt.Sprintf("Unrecognised mode: %d", param.mode))
	}
}

func setVal(program *[]float64, param Parameter, bp int, val float64) float64 {
	log.WithFields(log.Fields{
		"param": param,
		"val":   val,
	}).Trace("Setting Value")

	if param.mode == 0 {
		(*program)[int(param.val)] = val
	} else if param.mode == 2 {
		(*program)[int(param.val)+bp] = val
	}
	return 0
}

func execInstruction(program *[]float64, op Operation, ip int, bp int, input float64) (i int, b int, output Output) {
	log.WithFields(log.Fields{"op": op, "ip": ip, "bp": bp}).Debug("Executing Operation")
	ip += op.nParams + 1
	output = Output{0, false}
	if op.opcode == 1 {
		setVal(program,
			//op.params[2].val,
			op.params[2],
			bp,
			getVal(program, op.params[0], bp)+getVal(program, op.params[1], bp),
		)
	} else if op.opcode == 2 {
		setVal(program,
			//op.params[2].val,
			op.params[2],
			bp,
			getVal(program, op.params[0], bp)*getVal(program, op.params[1], bp),
		)
	} else if op.opcode == 3 {
		setVal(program,
			op.params[0],
			bp,
			input,
		)
		log.WithFields(log.Fields{
			"in": input,
		}).Debug("Operation 3 Input")
	} else if op.opcode == 4 {
		output = Output{getVal(program, op.params[0], bp), true}
		log.WithFields(log.Fields{
			"out": getVal(program, op.params[0], bp),
		}).Debug("Operation 4 Output")
	} else if op.opcode == 5 {
		if getVal(program, op.params[0], bp) != 0 {
			ip = int(getVal(program, op.params[1], bp))
		}
	} else if op.opcode == 6 {
		if getVal(program, op.params[0], bp) == 0 {
			ip = int(getVal(program, op.params[1], bp))
		}
	} else if op.opcode == 7 {
		var v = -1.0
		if getVal(program, op.params[0], bp) < getVal(program, op.params[1], bp) {
			v = 1.0
		} else {
			v = 0.0
		}
		setVal(program, op.params[2], bp, v)
	} else if op.opcode == 8 {
		var v = -1.0
		if getVal(program, op.params[0], bp) == getVal(program, op.params[1], bp) {
			v = 1.0
		} else {
			v = 0.0
		}
		setVal(program, op.params[2], bp, v)
	} else if op.opcode == 9 {
		bp += int(getVal(program, op.params[0], bp))
	} else {
		panic(fmt.Sprintf("Unrecognised opcode: %d", op.opcode))
	}
	log.WithFields(log.Fields{"op": op, "ip": ip, "bp": bp}).Debug("Executed Operation")
	return ip, bp, output
}

func execute(program *[]float64, inputs []float64, mode string, ip int, bp int) (output []float64, i int, b int, paused bool) {
	nInput := 0
	for i, b := ip, bp; i < len(*program); {
		log.WithFields(log.Fields{
			"i":              i,
			"bp":             b,
			"program[i:i+4]": (*program)[i : i+4],
		}).Trace("Parsing op")
		op, isTerminated := parseOp((*program)[i:])

		log.WithFields(log.Fields{
			"i":              i,
			"bp":             b,
			"program[i:i+4]": (*program)[i : i+4],
			"op":             op,
			"isTerminated":   isTerminated,
		}).Trace("Parsed op")
		if isTerminated {
			break
		}

		var input = 0.0
		if op.opcode == 3 {
			if nInput >= len(inputs) {
				input = getInput()
			} else {
				input = inputs[nInput]
				nInput++
			}
		}

		ip, bp, out := execInstruction(program, op, i, b, input)
		i, b = ip, bp
		if out.hasOutput {
			output = append(output, out.val)
		}

		if op.opcode == 4 && mode == "FEEDBACK" {
			log.WithFields(log.Fields{
				"i":              i,
				"program[i:i+4]": (*program)[i : i+4],
				"program[i]":     (*program)[i],
				"op":             op,
				"isTerminated":   isTerminated,
				"output":         output,
			}).Debug("Pausing program after output")

			return output, i, bp, true
		}
	}

	return output, i, b, false
}

func loadAndRun(file io.ReadSeeker, inputs []float64, mode string, ip int) ([]float64, int, int, bool) {
	_, err := file.Seek(0, io.SeekStart)
	check(err)

	scanner := bufio.NewScanner(file)

	// Move to first line
	scanner.Scan()

	sProgram := strings.Split(scanner.Text(), ",")

	// Cast elements to float64
	var program []float64
	for _, i := range sProgram {
		j, err := strconv.ParseFloat(i, 64)
		if err != nil {
			panic(err)
		}
		program = append(program, j)
	}

	buffer := make([]float64, len(program)*10)
	program = append(program, buffer...)

	bp := 0
	return execute(&program, inputs, mode, ip, bp)
}

func part1(file io.ReadSeeker) {
	inputs := []float64{1}
	output, _, _, _ := loadAndRun(file, inputs, "", 0)

	log.WithFields(log.Fields{
		"Output": fmt.Sprintf("%f", output),
	}).Info("Part 1 Output")
}

func part2(file io.ReadSeeker) {
	inputs := []float64{2}
	output, _, _, _ := loadAndRun(file, inputs, "", 0)

	log.WithFields(log.Fields{
		"Output": fmt.Sprintf("%f", output),
	}).Info("Part 2 Output")
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
