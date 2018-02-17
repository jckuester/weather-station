// Package pulse decodes a compressed pulse sequence
// received via the Arduino library https://github.com/pimatic/RFControl.
//
// For decoding details see also https://github.com/pimatic/rfcontroljs#details.
package pulse

import (
	"strconv"
	"strings"

	"sort"

	"fmt"

	"math"

	"github.com/bradfitz/slice"
	"github.com/jckuester/weather-station/protocol"
	"github.com/pkg/errors"
)

// Signal implements a received 433 MHz signal of compressed raw time series
// that consists of pulse lengths and a sequence of pulses.
type Signal struct {
	Lengths []int
	Seq     string
}

// Pair simply implements a tuple fo two values
// (first, second).
type Pair struct {
	first  int
	second int
}

// Decode tries to decode a received Signal
// based on all currently supported protocols.
func Decode(s *Signal) (interface{}, error) {
	for _, p := range protocol.Supported() {
		if matches(s, p) {
			binary, err := Map(s.Seq, p.Mapping)
			if err != nil {
				return nil, err
			}
			return p.Decode(binary)
		}
	}
	return nil, nil
}

// matches checks whether a received Signal matches
// a protocol.
func matches(s *Signal, p *protocol.Protocol) bool {
	var i int
	var maxDelta float64

	// length of the pulse sequence must match
	if !contains(p.SeqLength, len(s.Seq)) {
		return false
	}

	// number of pulse length must match
	if len(s.Lengths) != len(p.Lengths) {
		return false
	}

	// pulse length must be in a certain range
	for i < len(s.Lengths) {
		maxDelta = float64(float64(s.Lengths[i]) * float64(0.4))
		if math.Abs(float64(s.Lengths[i]-p.Lengths[i])) > maxDelta {
			return false
		}
		i++
	}
	return true
}

// Prepare takes an compressed signal as input,
// 1) splits it into pulse lengths and pulse sequence,
// 2) removes pulse lengths that are 0,
// 3) sorts the pulse lengths in ascending order, and
// 4) rearranges the pulse sequence, which characters each is a pulse length
// represented by its index in the array of pulse lengths.
func Prepare(input string) (*Signal, error) {
	parts := strings.Split(input, " ")
	if len(parts) < 8 {
		return nil, errors.New(fmt.Sprintf("Incorrect number of pulse lengths: %s", input))
	}
	lengths := parts[:8]
	seq := parts[8]

	lengths = filter(lengths, func(s string) bool {
		return s != "0"
	})

	lengthsInts, err := toIntArray(lengths)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot convert pulse lengths to integers: %s", lengths))
	}

	return sortSignal(
		&Signal{
			lengthsInts,
			seq,
		})
}

// sortSignal sorts the given pulse lengths in ascending order
// and changes the pulse sequence, where each character is a pulse length
// represented by its index in the array of pulse lengths,
// according to the new order of indices.
func sortSignal(s *Signal) (*Signal, error) {
	sortedIndices := sortIndices(s.Lengths)
	sort.Ints(s.Lengths)

	seq, err := changeRepresentation(s.Seq, sortedIndices)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to change the representation of '%s'", s.Seq)
	}

	return &Signal{
		s.Lengths,
		seq,
	}, nil
}

// sortIndices sorts the indicies of a
// given array a, i.e. if the array is
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
func changeRepresentation(seq string, mapping []int) (string, error) {
	var result string
	var d int

	for d < len(seq) {
		i, err := strconv.ParseInt(string(seq[d]), 10, 0)
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

func filter(a []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range a {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func toIntArray(a []string) ([]int, error) {
	var intArray = []int{}

	for _, i := range a {
		j, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}
		intArray = append(intArray, j)
	}

	return intArray, nil
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
