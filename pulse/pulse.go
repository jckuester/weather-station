package pulse

import (
	"log"
	"strconv"
	"strings"

	"sort"

	"fmt"

	"github.com/bradfitz/slice"
)

type PulseInfo struct {
	Lengths []int
	Pulses  string
}

type Tuple struct {
	first  int
	second int
}

func PrepareCompressedPulses(input string) *PulseInfo {
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
