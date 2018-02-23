package pulse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode_weather15(t *testing.T) {
	bits := "1001100101100010000011001000010000111"

	result, _ := Protocols()["protocol1"].Decode(bits)
	m := result.(*GTWT01Result)

	assert.Equal(t, 67, m.Humidity, "Humidity")
	assert.Equal(t, 20.0, m.Temperature, "Temperature")
	assert.Equal(t, false, m.LowBattery, "LowBattery")
	assert.Equal(t, 3, m.Channel, "Channel")
	assert.Equal(t, 2454, m.ID, "Id")
}
