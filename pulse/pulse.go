package pulse

import (
	"log"
	"strconv"
	"strings"

	"sort"

	"fmt"

	"math"

	"github.com/bradfitz/slice"
)

type PulseInfo struct {
	pulseLengths []int
	pulses       string
}

type Tuple struct {
	first  int
	second int
}

func prepareCompressedPulses(input string) *PulseInfo {
	parts := strings.Split(input, " ")
	pulseLengths := parts[:8]
	pulses := parts[8]

	pulseLengths = Filter(pulseLengths, func(s string) bool {
		return s != "0"
	})

	return sortCompressedPulses(StringToIntArray(pulseLengths), pulses)
}

func sortCompressedPulses(pulseLengths []int, pulses string) *PulseInfo {
	sortedIndices := sortIndices(pulseLengths)

	sort.Ints(pulseLengths)

	pulses = mapByArray(pulses, sortedIndices)

	return &PulseInfo{
		pulseLengths,
		pulses,
	}
}

func fixPulses(pulseLengths []int, pulses string) *PulseInfo {
	if len(pulseLengths) <= 3 {
		return nil
	}

	i := 1
	newPulseLengths := pulseLengths

	for i < len(pulseLengths) {
		if pulseLengths[i-1]*2 < pulseLengths[i] {
			i++
			continue
		}

		newPulseLength := math.Floor(float64(pulseLengths[i-1]+pulseLengths[i]) / 2)
		newPulseLengths2 := append(newPulseLengths[:i-1], int(newPulseLength))
		newPulseLengths = append(newPulseLengths2, newPulseLengths[i+1:]...)
		break
	}

	if i == len(pulseLengths) {
		return nil
	}
	newPulses := pulses

	for i < len(pulseLengths) {
		newPulses = strings.Replace(newPulses,
			strconv.FormatUint(uint64(i), 10),
			strconv.FormatUint(uint64(i-1), 10),
			-1)
		i++
	}
	return &PulseInfo{
		pulseLengths: newPulseLengths,
		pulses:       newPulses,
	}
}

func sortIndices(array []int) []int {
	tuples := make([]Tuple, len(array))

	for i, e := range array {
		tuples[i] = Tuple{e, i}
	}
	slice.Sort(tuples[:], func(l, r int) bool {
		return tuples[l].first < tuples[r].first
	})

	indices := make([]int, len(array))

	for j, t := range tuples {
		indices[t.second] = j
	}
	return indices
}

func mapByArray(data string, mapping []int) string {
	var result string
	var d int

	for d < len(data) {
		i, err := strconv.ParseInt(string(data[d]), 10, 0)
		if err != nil {
			log.Fatal(err)
		}

		result = fmt.Sprintf("%s%d", result, mapping[i])
		d++
	}
	return result
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func StringToIntArray(a []string) []int {
	var t2 = []int{}

	for _, i := range a {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		t2 = append(t2, j)
	}

	return t2
}
