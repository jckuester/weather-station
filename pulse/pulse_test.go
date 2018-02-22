package pulse

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatches(t *testing.T) {
	s := &Signal{
		Lengths: []int{516, 2116, 4152, 9112},
		Seq:     "0102020101020201020101020102010202020202020202010201010202010202020202020103",
	}
	p := &Protocol{
		SeqLength: []int{76},
		Lengths:   []int{496, 2048, 4068, 8960},
	}

	assert.Equal(t, true, matches(s, p))
}

func TestMatches_pulseSeqTooShort(t *testing.T) {
	s := &Signal{
		Lengths: []int{516, 2116, 4152, 9112},
		Seq:     "010202010102020102",
	}
	p := &Protocol{
		SeqLength: []int{76},
		Lengths:   []int{496, 2048, 4068, 8960},
	}

	assert.Equal(t, false, matches(s, p))
}

func TestMatches_pulseLengthDeviationTooHigh(t *testing.T) {
	s := &Signal{
		Lengths: []int{516, 2116, 4152, 9112},
		Seq:     "0102020101020201020101020102010202020202020202010201010202010202020202020103",
	}
	p := &Protocol{
		SeqLength: []int{76},
		Lengths:   []int{496, 2048, 2000, 8960},
	}

	assert.Equal(t, false, matches(s, p))
}

func TestMatches_numberOfPuleLengthDiffer(t *testing.T) {
	s := &Signal{
		Lengths: []int{516, 2116, 4152, 9112},
		Seq:     "0102020101020201020101020102010202020202020202010201010202010202020202020103",
	}
	p := &Protocol{
		SeqLength: []int{76},
		Lengths:   []int{496, 2048, 2000},
	}

	assert.Equal(t, false, matches(s, p))
}

func TestConvert(t *testing.T) {
	seq := "020101020201020201010103"
	bits := "10011011000"

	var m = map[string]string{
		"01": "0",
		"02": "1",
		"03": "",
	}

	mapped, err := convert(seq, m)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, mapped, bits)
}

func TestPrepare(t *testing.T) {
	input := "255 2904 1388 771 11346 0 0 0 0100020002020000020002020000020002000202000200020002000200000202000200020000020002000200020002020002000002000200000002000200020002020002000200020034"

	p, _ := Prepare(input)

	assert.Equal(t, []int{255, 771, 1388, 2904, 11346}, p.Lengths)
	assert.Equal(t, "0300020002020000020002020000020002000202000200020002000200000202000200020000020002000200020002020002000002000200000002000200020002020002000200020014", p.Seq)
}

func TestPrepare_InvalidCharacters(t *testing.T) {
	pC := "544 4128 2100 100 140 320 808 188 01020202010202020202020202020202020101020102010101010J�G_YJ�Üxx�1��Nz�8��&[��"

	_, err := Prepare(pC)

	assert.Error(t, err)
}

func TestSortIndices(t *testing.T) {
	a := []int{516, 9112, 4152, 2116}

	sortedIndices := sortIndices(a)

	assert.Equal(t, sortedIndices, map[string]string{
		"0": "0",
		"1": "3",
		"2": "2",
		"3": "1",
	})
}

func TestSortSignal(t *testing.T) {
	s := &Signal{
		Lengths: []int{516, 9112, 4152, 2116},
		Seq:     "01020203",
	}

	expectedSignal := &Signal{
		Lengths: []int{516, 2116, 4152, 9112},
		Seq:     "03020201",
	}

	sortedSignal, _ := sortSignal(s)
	require.Equal(t, expectedSignal, sortedSignal)
}

func TestSortSignal_lengthsAlreadySorted(t *testing.T) {
	s := &Signal{
		Lengths: []int{516, 2116, 4152, 9112},
		Seq:     "0102020101020201020101020102010202020202020202010201010202010202020202020103",
	}

	sortedSignal, _ := sortSignal(s)
	require.Equal(t, s, sortedSignal)
}
