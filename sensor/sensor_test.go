package sensor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapPulse(t *testing.T) {
	pulse := "020101020201020201010103"
	bits := "10011011000"

	assert.Equal(t, mapPulse(pulse), bits)
	//"RF receive 512 4152 2124 9112 0 0 0 0 0102020101020201020101020202010202020202010102020201010202010202020102020103"
}

func TestDecode(t *testing.T) {
	bits := "1001100101100010000011001000010000111"

	s := &Sensor{}

	m := s.decode(bits)

	assert.Equal(t, 67, m.Humidity, "Humidity")
	assert.Equal(t, 20.0, m.Temperature, "Temperature")
	assert.Equal(t, false, m.LowBattery, "LowBattery")
	assert.Equal(t, 3, m.Channel, "Channel")
	assert.Equal(t, 2454, m.Id, "Id")
}

func TestDecode2(t *testing.T) {
	bits := "1001000100010001000011010111010000101"

	s := &Sensor{}

	m := s.decode(bits)

	assert.Equal(t, 66, m.Humidity, "Humidity")
	assert.Equal(t, 21.5, m.Temperature, "Temperature")
	assert.Equal(t, false, m.LowBattery, "LowBattery")
	assert.Equal(t, 2, m.Channel, "Channel")
	assert.Equal(t, 2321, m.Id, "Id")
}
