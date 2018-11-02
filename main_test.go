package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestPrintAllMatchingProtocols(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	s := &Signal{
		Lengths: []int{516, 2116, 4152, 9112},
		Seq:     "0102020101020201020101020102010202020202020202010201010202010202020202020103",
	}

	printAllMatchingProtocols([]string{"weather12", "weather15"}, s)
	assert.Contains(t, buf.String(), "weather12")
	assert.Contains(t, buf.String(), "weather15")
}

func TestProcessedWithMatchingConfigNoMatch(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	s := &Signal{
		Lengths: []int{516, 2116, 4152, 9112},
		Seq:     "0102020101020201020101020102010202020202020202010201010202010202020202020103",
	}

	processed := processedWithMatchingConfig([]string{"weather12", "weather15"}, s)
	assert.False(t, processed)
}

func TestProcessedWithMatchingConfigMatch(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	loadSampleConfig()

	s := &Signal{
		Lengths: []int{616, 1996, 4048, 9044},
		Seq:     "0102010202010202010101010101010102010202020102020102020102010202020102020103",
	}

	setupMetrics()

	processed := processedWithMatchingConfig([]string{"weather12", "weather15"}, s)
	assert.True(t, processed)
	assert.Contains(t, buf.String(), "weather12: {ID:91 Name:91 Channel:1 Temperature:18.7 Humidity:53 LowBattery:false}")
}
