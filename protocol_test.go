package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode_weather15(t *testing.T) {
	p, _ := PreparePulse("564 4116 2068 9112 0 0 0 0 0102020101020201020101020202020102020202020201020102010102010202010101020103")
	bits, _ := convert(p.Seq, Protocols()["protocol1"].Mapping)
	//bits := "1001100101100010000011001000010000111"

	result, _ := Protocols()["protocol1"].Decode(bits)
	m := result.(*GTWT01Result)

	assert.Equal(t, 78, m.Humidity, "Humidity")
	assert.Equal(t, 4.3, m.Temperature, "Temperature")
	assert.Equal(t, false, m.LowBattery, "LowBattery")
	assert.Equal(t, 2, m.Channel, "Channel")
	assert.Equal(t, 2454, m.ID, "Id")
}

func TestDecode_weather12(t *testing.T) {
	p, _ := PreparePulse("568 4040 2008 8952 0 0 0 0 0102020102010202020201020202020201010201020202020201010101020201020102010103")
	bits, _ := convert(p.Seq, Protocols()["protocol1"].Mapping)
	//bits := "0110101111011111001011111000011010100"

	result, _ := Protocols()["protocol1"].Decode(bits)
	m := result.(*GTWT01Result)

	assert.Equal(t, 60, m.Humidity, "Humidity")
	assert.Equal(t, 20.8, m.Temperature, "Temperature")
	assert.Equal(t, false, m.LowBattery, "LowBattery")
	assert.Equal(t, 3, m.Channel, "Channel")
	assert.Equal(t, 148, m.ID, "Id")
}
