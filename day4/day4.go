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

func convertFloatToIntArray(current float64) []int{
	scurrent := strconv.FormatFloat(current, 'f', -1, 64)
	sacurrent := strings.Split(scurrent, "")

	var acurrent []int;
	for _, v := range sacurrent {
		i, err := strconv.Atoi(v)
		check(err)
		acurrent = append(acurrent, i)
	}

	return acurrent
}

func convertIntArrayToFloat(current []int) float64{
	val, err := strconv.ParseFloat(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(current)), ""),"[]"), 8);
	check(err)
	return val;
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

func isValidPwd(current float64, mode int) (isValid bool){
	if mode == 1{
		return isValidPwdP1(current)
	}else if mode == 2{
		return isValidPwdP2(current)
	}else{
		panic("Uh, what?")
	}
}

func isValidPwdP1(current float64)(isValid bool){
	pwd := convertFloatToIntArray(current)
	_,mini := firstMin(pwd)
	if mini > 0 {
		return false
	}
	_,maxi := lastMax(pwd)
	if maxi < len(pwd) - 1  {
		return false
	}

	for i := 1; i < len(pwd); i++ {
		if pwd[i] == pwd[i-1] {
			return true
		}
	}
	return false
}

func isValidPwdP2(current float64)(isValid bool){
	pwd := convertFloatToIntArray(current)
	_,mini := firstMin(pwd)
	if mini > 0 {
		return false
	}
	_,maxi := lastMax(pwd)
	if maxi < len(pwd) - 1  {
		return false
	}

	for i := 0; i < len(pwd); i++ {
		consec_counter := 1
		//fmt.Println(consec_counter, pwd, i)
		for j := 1; j + i < len(pwd); j++ {
			//fmt.Println(consec_counter, pwd, i,j )
			if pwd[i] == pwd[i + j] {
				consec_counter += 1
				if i + j == len(pwd) - 1 && consec_counter % 2 != 0 {
					return false;
				}
			} else {
				break
			}
		}
		if consec_counter > 1 && consec_counter == 2 {
			return true
		} else {
			i += consec_counter - 1
		}

	}

	return false
}

func generateNext(fcurrent float64)(fnext float64, changed bool){

	current := convertFloatToIntArray(fcurrent)
	changed = false;

	next := make([]int, len(current))
	copy(next, current)

	min, mini := firstMin(current)
	log.WithFields(log.Fields{
		"min":        min,
		"mini":       mini,
		"current":    current,
		"next": next,
	}).Trace("After min")

	if mini != 0 {
		changed = true
		for i := mini + 1; i < len(next); i++ {
			next[i] = min
		}
	}

	log.WithFields(log.Fields{
		"min":        min,
		"mini":       mini,
		"current":    current,
		"next": next,
	}).Trace("After change min")

	max, maxi := lastMax(next)
	log.WithFields(log.Fields{
		"max":        max,
		"maxi":       maxi,
		"current":    current,
		"next": next,
	}).Trace("After max")

	if maxi != len(next)-1 {
		changed = true
		for i := maxi + 1; i < len(next); i++ {
			next[i] = max
		}
	}
	log.WithFields(log.Fields{
		"max": max,
		"maxi": maxi,
		"current": current,
		"next": next,
	}).Trace("After change max")

	fnext = convertIntArrayToFloat(next)
	return
}

//271973
func nextValidPwd(current float64, mode int) (next float64) {

	var isValid bool = false
	var changed bool = false
	next = current
	for true{
		next, changed = generateNext(next)
		isValid = isValidPwd(next,mode)
		log.WithFields(log.Fields{
			"current": current,
			"next": next,
			"isValid": isValid,
			"changed": changed,
		}).Debug("IsValid/Changed")
		if isValid && changed {
			return next
		} else if isValidPwd(next+1,mode){
			return next + 1
		} else {
			next += 1
		}
	}

	return next
}

func part1(lower float64, upper float64) int{
	current := lower
	counter := 0
	for current <= upper {
		log.WithFields(log.Fields{
			"current": current,
			"counter": counter,
		}).Debug("current")
		current = nextValidPwd(current,1)
		counter += 1
	}
	return counter - 1
}

func part2(lower float64, upper float64) int{
	current := lower
	counter := 0
	for current <= upper {
	//for i :=0; i<100;i++ {
		log.WithFields(log.Fields{
			"current": current,
			"counter": counter,
		}).Debug("current")
		current = nextValidPwd(current,2)
		counter += 1
	}
	return counter - 1
}

func main() {
	log.SetLevel(log.InfoLevel)

	log.WithFields(log.Fields{
		"Valid Count": part1(271973,785961),
	}).Info("Part 1")

	log.WithFields(log.Fields{
		"Valid Count": part2(271973,785961),
	}).Info("Part 2")

}
