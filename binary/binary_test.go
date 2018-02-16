package binary

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToSignedNumber_positive(t *testing.T) {
	i, _ := ToSignedNumber("0011110", 1, 5)
	require.Equal(t, 15, i)
}

func TestToSignedNumber_negative(t *testing.T) {
	i, _ := ToSignedNumber("1111101", 1, 5)
	require.Equal(t, -2, i)
}

func TestToSignedNumber_invalidCharacter(t *testing.T) {
	_, err := ToSignedNumber("1112", 0, 4)
	require.Error(t, err)
}
