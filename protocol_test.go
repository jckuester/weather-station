package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode_weather15(t *testing.T) {
	p, _ := PreparePulse("564 4116 2068 9112 0 0 0 0 0102020101020201020101020202020102020202020201020102010102010202010101020103")
	bits, _ := convert(p.Seq, Protocols()["weather15"].Mapping)
	//bits := "1001100101100010000011001000010000111"

	result, _ := Protocols()["weather15"].Decode(bits)
	m := result.(*GTWT01Result)

	assert.Equal(t, 78, m.Humidity, "Humidity")
	assert.Equal(t, 4.3, m.Temperature, "Temperature")
	assert.Equal(t, false, m.LowBattery, "LowBattery")
	assert.Equal(t, 2, m.Channel, "Channel")
	assert.Equal(t, 2454, m.ID, "Id")
}

func TestDecode_weather12(t *testing.T) {
	p, _ := PreparePulse("616 1996 4048 9044 0 0 0 0 0102010202010202010101010101010102010202020102020102020102010202020102020103")
	bits, _ := convert(p.Seq, Protocols()["weather12"].Mapping)
	//bits := "1010010011111111010001001001010001001"

	result, _ := Protocols()["weather12"].Decode(bits)
	m := result.(*GTWT01Result)

	assert.Equal(t, 53, m.Humidity, "Humidity")
	assert.Equal(t, 18.7, m.Temperature, "Temperature")
	assert.Equal(t, false, m.LowBattery, "LowBattery")
	assert.Equal(t, 1, m.Channel, "Channel")
	assert.Equal(t, 91, m.ID, "Id")
}
