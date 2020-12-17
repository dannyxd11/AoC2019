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

type Output struct {
	val       int
	hasOutput bool
}

type ProgramState struct {
	ip          int
	program     *[]int
	initialised bool
}

type Operation struct {
	opcode  int
	params  []Parameter
	nParams int
}

// https://stackoverflow.com/questions/30226438/generate-all-permutations-in-go
func permutations(arr []int) [][]int {
	var helper func([]int, int)
	var res [][]int

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
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

func getInput() int {
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("> ")

	input, err := buf.ReadByte()
	check(err)
	iInput, err := strconv.Atoi(string(input))
	check(err)

	return iInput
}

func getVal(program *[]int, param Parameter) int {
	if param.mode == 0 {
		log.WithFields(log.Fields{
			"param":              param,
			"program[param.val]": (*program)[param.val],
		}).Trace("Getting Value (Position)")
		return (*program)[param.val]
	} else if param.mode == 1 {
		log.WithFields(log.Fields{
			"param":     param,
			"param.val": param.val,
		}).Trace("Getting Value (Immediate)")
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

func execInstruction(program *[]int, op Operation, ip int, input int) (i int, output Output) {
	log.WithFields(log.Fields{"op": op, "ip": ip}).Debug("Executing Operation")
	ip += op.nParams + 1
	output = Output{0, false}
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
		setVal(program,
			op.params[0].val,
			input,
		)
	} else if op.opcode == 4 {
		output = Output{getVal(program, op.params[0]), true}
		log.WithFields(log.Fields{
			"out": getVal(program, op.params[0]),
		}).Debug("Operation 4 Output")
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
	return ip, output
}

func execute(program *[]int, inputs []int, mode string, ip int) (output []int, i int, paused bool) {
	nInput := 0
	for i = ip; i < len(*program); {
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

		var input = 0
		if op.opcode == 3 {
			if nInput >= len(inputs) {
				input = getInput()
			} else {
				input = inputs[nInput]
				nInput++
			}
		}

		ip, out := execInstruction(program, op, i, input)
		i = ip
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
			//if (*program)[i] == 99 {
			//	return output, i, false
			//}
			return output, i, true
		}
	}

	return output, i, false
}

func load(file io.ReadSeeker) *[]int {
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

	return &program
}

func loadAndRun(file io.ReadSeeker, inputs []int, mode string, ip int) ([]int, int, bool) {
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

	return execute(&program, inputs, mode, ip)
}

func part1(file io.ReadSeeker) {
	P1PhaseCodes := []int{0, 1, 2, 3, 4}
	P1perms := permutations(P1PhaseCodes)
	P1highestOutput, P1highestIndex := -1, -1
	for i, v := range P1perms {
		output := []int{0}
		for _, p := range v {
			// Amp A
			inputs := []int{p}
			inputs = append(inputs, output...)
			output, _, _ = loadAndRun(file, inputs, "", 0)

			if output[0] > P1highestOutput {
				P1highestOutput = output[0]
				P1highestIndex = i
			}
		}
	}
	log.WithFields(log.Fields{
		"Highest Index":  P1highestIndex,
		"Highest Output": P1highestOutput,
		"Highest Perms":  P1perms[P1highestIndex],
	}).Info("Part 1 Output")
}

func part2(file io.ReadSeeker) {
	P2PhaseCodes := []int{5, 6, 7, 8, 9}
	P2perms := permutations(P2PhaseCodes)
	//P2perms := [][]int{{5,6,7,8,9}}
	P2highestOutput, P2highestIndex := -1, -1

	for i, v := range P2perms {
		var output []int
		P2amplifierState := map[int]ProgramState{}
		for p := 0; ; p++ {
			var inputs []int

			state := P2amplifierState[p%5]
			if !state.initialised {
				inputs = append(inputs, v[p%5])
				if p%5 == 0 {
					inputs = append(inputs, 0)
				}
				state.program = load(file)
			}

			inputs = append(inputs, output...)

			log.WithFields(log.Fields{
				"program": p % 5,
				"v":       v,
				"i":       i,
				"inputs":  inputs,
				"state":   P2amplifierState[p%5],
			}).Debug("Starting program..")

			o, ip, paused := execute(state.program, inputs, "FEEDBACK", state.ip)
			output = o
			P2amplifierState[p%5] = ProgramState{ip, state.program, true}

			log.WithFields(log.Fields{
				"program": p % 5,
				"v":       v,
				"i":       i,
				"inputs":  inputs,
				"output":  output,
				"paused":  paused,
				"state":   P2amplifierState[p%5],
			}).Debug("Stopped program..")

			if len(output) == 0 {
				output = append(output, inputs...)
			}

			if paused == false {
				break
			}
		}

		if output[0] > P2highestOutput {
			P2highestOutput = output[0]
			P2highestIndex = i
		}
	}
	log.WithFields(log.Fields{
		"Highest Index":  P2highestIndex,
		"Highest Output": P2highestOutput,
		"Highest Perms":  P2perms[P2highestIndex],
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
