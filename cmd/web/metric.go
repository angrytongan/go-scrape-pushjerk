package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

func poundsToKilos(pounds int) int {
	return int(math.Ceil(float64(pounds) / 2.2))
}

func inchesToCm(inches int) int {
	return int(math.Trunc(2.54 * float64(inches)))
}

func addKilos(s string) string {
	re := regexp.MustCompile(`(\d+)/(\d+)#`)
	matches := re.FindSubmatch([]byte(s))
	if len(matches) == 3 {
		aKilos := 0
		bKilos := 0

		mensPounds, err := strconv.Atoi(string(matches[1]))
		if err == nil {
			aKilos = poundsToKilos(mensPounds)
		}

		womensPounds, err := strconv.Atoi(string(matches[2][:]))
		if err == nil {
			bKilos = poundsToKilos(womensPounds)
		}

		return fmt.Sprintf("%s (%d/%dkg)", s, aKilos, bKilos)
	}

	return s
}

func addCm(s string) string {
	re := regexp.MustCompile(`(\d+)/(\d+)”`)
	matches := re.FindSubmatch([]byte(s))
	if len(matches) == 3 {
		aCm := 0
		bCm := 0

		aInches, err := strconv.Atoi(string(matches[1]))
		if err == nil {
			aCm = inchesToCm(aInches)
		}

		bInches, err := strconv.Atoi(string(matches[2][:]))
		if err == nil {
			bCm = inchesToCm(bInches)
		}

		return fmt.Sprintf("%s (%d/%dcm)", s, aCm, bCm)
	}

	return s
}

func Metricise(workout string) string {
	kilos := regexp.MustCompile(`(\d+)/(\d+)#`)
	cm := regexp.MustCompile(`(\d+)/(\d+)”`)

	out := kilos.ReplaceAllStringFunc(workout, addKilos)
	out = cm.ReplaceAllStringFunc(out, addCm)

	return out
}
