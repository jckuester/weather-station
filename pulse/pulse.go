package pulse

import (
	"strconv"
	"strings"

	"sort"

	"fmt"

	"errors"

	"math"

	"github.com/bradfitz/slice"
	"github.com/jckuester/weather-station/protocol"
)

// Package pulse decodes a compressed pulse sequence
// received via the Arduino library https://github.com/pimatic/RFControl.
//
// For decoding details see also https://github.com/pimatic/rfcontroljs#details.

type PulseInfo struct {
	Lengths []int  // length of pulses
	Seq     string // sequence of pulses
}

type Pair struct {
	first  int
	second int
}

func Decode(p *PulseInfo, pc *protocol.Protocol) (interface{}, error) {
	if protocolMatches(p, pc) {
		binary, err := Map(p.Seq, pc.Mapping)
		if err != nil {
			return nil, err
		}
		return pc.Decode(binary)
	}
	return nil, nil
}

// protocolMatches checks whether a received pulseInfo matches
// a protocol.
func protocolMatches(p *PulseInfo, pc *protocol.Protocol) bool {
	var i int
	var maxDelta float64

	// length of the pulse sequence must match
	if !contains(pc.SeqLength, len(p.Seq)) {
		return false
	}

	// number of pulse length must match
	if len(p.Lengths) != len(pc.Lengths) {
		return false
	}

	// pulse length must be in a certain range
	for i < len(p.Lengths) {
		maxDelta = float64(float64(p.Lengths[i]) * float64(0.4))
		if math.Abs(float64(p.Lengths[i]-pc.Lengths[i])) > maxDelta {
			return false
		}
		i++
	}
	return true
}

// PrepareCompressed takes an input, such as
// "516 4156 2116 9116 0 0 0 0 01020201010202010201010201020102020",
// splits it into the pulse lenghts (removes  pulse lenght that are 0),
// and sorts the pulse lengths in ascending order
func PrepareCompressed(input string) (*PulseInfo, error) {
	parts := strings.Split(input, " ")
	if len(parts) < 8 {
		return nil, errors.New(fmt.Sprintf("Incorrect number of pulse lengths: %s", input))
	}
	pulseLengths := parts[:8]
	seq := parts[8]

	pulseLengths = Filter(pulseLengths, func(s string) bool {
		return s != "0"
	})

	return sortCompressed(&PulseInfo{
		StringToIntArray(pulseLengths),
		seq,
	})
}

// sortCompressed sorts the given pulse lengths in ascending order
// and changes their representation in the pulse sequence accordingly to the new indices.
func sortCompressed(p *PulseInfo) (*PulseInfo, error) {
	sortedIndices := sortIndices(p.Lengths)
	sort.Ints(p.Lengths)

	seq, err := changeRepresentation(p.Seq, sortedIndices)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to change the represenation of: %s", p.Seq))
	}

	return &PulseInfo{
		p.Lengths,
		seq,
	}, nil
}

// sortIndices sorts the indicies of a
// given array a, i.e. iff the array is
// [200, 600, 500], then it returns [0, 2, 1].
func sortIndices(a []int) []int {
	pairs := make([]Pair, len(a))

	for i, e := range a {
		pairs[i] = Pair{e, i}
	}
	slice.Sort(pairs[:], func(l, r int) bool {
		return pairs[l].first < pairs[r].first
	})

	indices := make([]int, len(a))

	for j, p := range pairs {
		indices[p.second] = j
	}
	return indices
}

// Instead of using the pulse lengths themselves, each character of the pulse sequence
// represents a pulse length by its index in the array of pulse lengths.
// Therefore, after sorting the pulse lengths array, the representation in the pulse sequence needs to
// be changed accordingly.
func changeRepresentation(data string, mapping []int) (string, error) {
	var result string
	var d int

	for d < len(data) {
		i, err := strconv.ParseInt(string(data[d]), 10, 0)
		if err != nil {
			return "", err
		}

		result = fmt.Sprintf("%s%d", result, mapping[i])
		d++
	}

	return result, nil
}

// Map maps a pulse sequence to a binary representation, using a given mapping.
func Map(seq string, mapping map[string]string) (string, error) {
	var hadMatch bool
	var i int
	var result string

	for i < len(seq) {
		hadMatch = false
		for search, replace := range mapping {
			if len(seq)-i >= len(search) {
				if string(seq[i:i+len(search)]) == search {
					result += replace
					i += len(search)
					hadMatch = true
					break
				}
			}
		}
		if !hadMatch {
			return "", errors.New(fmt.Sprintf("Unable to apply mapping to pulse sequence %s", seq))
		}
	}

	return result, nil
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

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
