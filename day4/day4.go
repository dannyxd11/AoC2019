package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func convertFloatToIntArray(current float64) []int {
	scurrent := strconv.FormatFloat(current, 'f', -1, 64)
	sacurrent := strings.Split(scurrent, "")

	var acurrent []int
	for _, v := range sacurrent {
		i, err := strconv.Atoi(v)
		check(err)
		acurrent = append(acurrent, i)
	}

	return acurrent
}

func convertIntArrayToFloat(current []int) float64 {
	val, err := strconv.ParseFloat(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(current)), ""), "[]"), 8)
	check(err)
	return val
}

func lastMax(array []int) (max int, maxi int) {
	max = array[0]
	maxi = 0

	for i, value := range array {
		if max <= value {
			max = value
			maxi = i
		}
	}
	return
}

func firstMin(array []int) (min int, mini int) {
	min = array[0]
	mini = 0
	for i, value := range array {
		if min > value {
			min = value
			mini = i
		}
	}
	return
}

func isValidPwd(current float64, mode int) (isValid bool) {
	switch {
	case mode == 1:
		return isValidPwdP1(current)
	case mode == 2:
		return isValidPwdP2(current)
	default:
		panic("Uh, what?")
	}
}

func isValidPwdP1(current float64) (isValid bool) {
	pwd := convertFloatToIntArray(current)
	_, mini := firstMin(pwd)
	if mini > 0 {
		return false
	}
	_, maxi := lastMax(pwd)
	if maxi < len(pwd)-1 {
		return false
	}

	for i := 1; i < len(pwd); i++ {
		if pwd[i] == pwd[i-1] {
			return true
		}
	}
	return false
}

func isValidPwdP2(current float64) (isValid bool) {
	pwd := convertFloatToIntArray(current)
	_, mini := firstMin(pwd)
	if mini > 0 {
		return false
	}
	_, maxi := lastMax(pwd)
	if maxi < len(pwd)-1 {
		return false
	}

	for i := 0; i < len(pwd); i++ {
		consecCounter := 1
		for j := 1; j+i < len(pwd); j++ {
			if pwd[i] == pwd[i+j] {
				consecCounter += 1
				if i+j == len(pwd)-1 && consecCounter%2 != 0 {
					return false
				}
			} else {
				break
			}
		}
		if consecCounter > 1 && consecCounter == 2 {
			return true
		} else {
			i += consecCounter - 1
		}

	}

	return false
}

func generateNext(fcurrent float64) (fnext float64, changed bool) {

	current := convertFloatToIntArray(fcurrent)
	changed = false

	next := make([]int, len(current))
	copy(next, current)

	min, mini := firstMin(current)

	if mini != 0 {
		changed = true
		for i := mini + 1; i < len(next); i++ {
			next[i] = min
		}
	}

	max, maxi := lastMax(next)

	if maxi != len(next)-1 {
		changed = true
		for i := maxi + 1; i < len(next); i++ {
			next[i] = max
		}
	}

	fnext = convertIntArrayToFloat(next)
	return
}

func nextValidPwd(current float64, mode int) (next float64) {

	var isValid = false
	var changed = false
	next = current
	for true {
		next, changed = generateNext(next)
		isValid = isValidPwd(next, mode)

		if isValid && changed {
			return next
		} else if isValidPwd(next+1, mode) {
			return next + 1
		} else {
			next++
		}
	}

	return next
}

func part1(lower float64, upper float64) int {
	current := lower
	counter := 0
	for current <= upper {
		log.WithFields(log.Fields{
			"current": current,
			"counter": counter,
		}).Debug("current")
		current = nextValidPwd(current, 1)
		counter += 1
	}
	return counter - 1
}

func part2(lower float64, upper float64) int {
	current := lower
	counter := 0
	for current <= upper {
		//for i :=0; i<100;i++ {
		log.WithFields(log.Fields{
			"current": current,
			"counter": counter,
		}).Debug("current")
		current = nextValidPwd(current, 2)
		counter++
	}
	return counter - 1
}

func main() {
	log.SetLevel(log.InfoLevel)

	log.WithFields(log.Fields{
		"Valid Count": part1(271973, 785961),
	}).Info("Part 1")

	log.WithFields(log.Fields{
		"Valid Count": part2(271973, 785961),
	}).Info("Part 2")

}
