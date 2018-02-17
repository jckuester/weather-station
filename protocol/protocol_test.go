package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode_weather15(t *testing.T) {
	bits := "1001100101100010000011001000010000111"

	d, _ := Supported()[0].Decode(bits)
	m := d.(*Measurement)

	assert.Equal(t, 67, m.Humidity, "Humidity")
	assert.Equal(t, 20.0, m.Temperature, "Temperature")
	assert.Equal(t, false, m.LowBattery, "LowBattery")
	assert.Equal(t, 3, m.Channel, "Channel")
	assert.Equal(t, 2454, m.Id, "Id")
}
