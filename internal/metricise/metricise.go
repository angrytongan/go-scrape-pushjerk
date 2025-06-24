package metricise

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

func addKilosWords(s string) string {
	re := regexp.MustCompile(` (\d+) pounds`)
	matches := re.FindSubmatch([]byte(s))
	if len(matches) == 2 {
		kg := 0

		pounds, err := strconv.Atoi(string(matches[1]))
		if err == nil {
			kg = poundsToKilos(pounds)
		}

		return fmt.Sprintf("%s (%d kg)", s, kg)
	}

	return s
}

func addKilosRange(s string) string {
	re := regexp.MustCompile(`(\d+)-(\d+) pounds`)
	matches := re.FindSubmatch([]byte(s))
	if len(matches) == 3 {
		aKg := 0
		bKg := 0

		aPounds, err := strconv.Atoi(string(matches[1]))
		if err == nil {
			aKg = poundsToKilos(aPounds)
		}

		bPounds, err := strconv.Atoi(string(matches[2]))
		if err == nil {
			bKg = poundsToKilos(bPounds)
		}

		return fmt.Sprintf("%s (%d-%d kg)", s, aKg, bKg)
	}

	return s
}

func Metricise(workout string) string {
	kilos := regexp.MustCompile(`(\d+)/(\d+)#`)
	cm := regexp.MustCompile(`(\d+)/(\d+)”`)
	pounds := regexp.MustCompile(` (\d+) pounds`)
	poundsRange := regexp.MustCompile(`(\d+)-(\d+) pounds`)

	out := kilos.ReplaceAllStringFunc(workout, addKilos)
	out = cm.ReplaceAllStringFunc(out, addCm)
	out = pounds.ReplaceAllStringFunc(out, addKilosWords)
	out = poundsRange.ReplaceAllStringFunc(out, addKilosRange)

	return out
}
