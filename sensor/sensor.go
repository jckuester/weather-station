package sensor

import (
	"log"
	"os"
	"sync/atomic"

	"strconv"

	"strings"

	"bufio"

	"github.com/jckuester/weather-station/pulse"
	"github.com/pkg/errors"
)

type Sensor struct {
	file   *os.File
	opened int32
}

// Measurement is the result of a Read operation.
type Measurement struct {
	Id          int
	Channel     int
	Temperature float64
	Humidity    int
	LowBattery  bool
}

// Open opens the named device file for reading,
// which is for an Arduino connected to the Raspberry Pi usually /dev/ttyUSB0.
func (m *Sensor) Open(name string) (err error) {
	atomic.StoreInt32(&m.opened, 1)

	m.file, err = os.Open(name)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Device '%v' opened", m.file.Name())
	return nil
}

// Read will read from the device file and decodes the received sensor information.
// Before it can be used the device file needs to be opened via Open.
func (m *Sensor) Read() (*Measurement, error) {
	if atomic.LoadInt32(&m.opened) != 1 {
		return nil, errors.New("Device needs to be opened")
	}

	scanner := bufio.NewScanner(m.file)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)

		if strings.HasPrefix(line, "RF receive") {
			pulseTrimmed := strings.TrimPrefix(line, "RF receive ")

			p := pulse.PrepareCompressedPulses(pulseTrimmed)
			log.Println(p)

			if p != nil {
				measurement := m.decode(mapPulse(p.Pulses))
				log.Println(measurement)

				return measurement, nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil, nil
}

// Decode decodes the pulse received from the "Globaltronics GT-WT-01 variant"
// and returns the human-readable temperature and humidity. See:
// https://github.com/pimatic/rfcontroljs/blob/master/src/protocols/weather15.coffee
func (m *Sensor) decode(pulses string) *Measurement {
	return &Measurement{
		Id:          binaryToNumber(pulses, 0, 12),
		Channel:     binaryToNumber(pulses, 14, 16) + 1,
		Temperature: float64(binaryToSignedNumber(pulses, 16, 27)) / 10,
		Humidity:    binaryToNumber(pulses, 28, 36),
		LowBattery:  binaryToBoolean(pulses, 12),
	}
}

func binaryToNumber(data string, b int, e int) int {
	i, err := strconv.ParseInt(data[b:e], 2, 0)
	if err != nil {
		log.Fatal(err)
	}

	return int(i)

}

func binaryToSignedNumber(data string, b int, e int) int {
	signedPos := b
	b++

	i, err := strconv.ParseInt(string(data[signedPos]), 2, 0)
	if err != nil {
		log.Fatal(err)
	}

	if i == 1 {
		return _binaryToSignedNumberMSBLSB(data, b, e)
	} else {
		return binaryToNumberMSBLSB(data, b, e)
	}
}

func _binaryToSignedNumberMSBLSB(data string, b int, e int) int {
	number := ^0
	i := b

	for i <= e {
		s, _ := strconv.ParseInt(string(data[i]), 2, 0)
		number <<= 1
		number |= int(s)
		i++
	}
	return number
}

func binaryToNumberMSBLSB(data string, b int, e int) int {
	number := 0
	i := b

	for i <= e {
		s, _ := strconv.ParseInt(string(data[i]), 2, 0)
		number <<= 1
		number |= int(s)
		i++
	}
	return number
}

func binaryToBoolean(data string, i int) bool {
	return string(data[i]) == "1"
}

func mapPulse(data string) string {
	var hadMatch bool
	var i int
	var result string

	mapping := map[string]string{
		"01": "0",
		"02": "1",
		"03": "",
	}

	for i < len(data) {
		hadMatch = false
		for search, replace := range mapping {
			if len(data)-i >= len(search) {
				if string(data[i:i+len(search)]) == search {
					result += replace
					i += len(search)
					hadMatch = true
					break
				}
			}
		}
	}

	if !hadMatch {
		return ""
	}

	return result
}

// Close closes the opened device file.
func (m *Sensor) Close() error {
	log.Printf("Closing '%v'", m.file.Name())
	atomic.StoreInt32(&m.opened, 0)
	return m.file.Close()
}
