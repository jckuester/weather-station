package arduino

// Package arduino helps to talk to
// an Arduino connected via USB (to a Raspberry Pi).

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync/atomic"

	"github.com/jckuester/weather-station/protocol"
	"github.com/jckuester/weather-station/pulse"
	"github.com/pkg/errors"
)

// Arduino represents the device file of and Arduino
// connected to the USB port.
type Arduino struct {
	file   *os.File
	opened int32
}

// Open opens the named device file for reading,
// which is usually /dev/ttyUSB0 for an Arduino
// connected to a Raspberry Pi.
func (m *Arduino) Open(name string) (err error) {
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
func (m *Arduino) Read() (*protocol.Measurement, error) {
	if atomic.LoadInt32(&m.opened) != 1 {
		return nil, errors.New("Device needs to be opened")
	}

	scanner := bufio.NewScanner(m.file)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)

		if strings.HasPrefix(line, "RF receive") {
			pulseTrimmed := strings.TrimPrefix(line, "RF receive ")

			p, err := pulse.PrepareCompressed(pulseTrimmed)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to prepare compressed pulse '%s'", pulseTrimmed)
			}
			log.Printf("%+v\n", *p)

			m, err := pulse.Decode(p)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed decode pulse info '%s'", p)
			}
			if m != nil {
				return m.(*protocol.Measurement), nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Wrapf(err, "Error while scanning device file '%s'", m.file.Name())
	}

	return nil, nil
}

// Close closes the opened device file.
func (m *Arduino) Close() error {
	log.Printf("Closing '%v'", m.file.Name())
	atomic.StoreInt32(&m.opened, 0)
	return m.file.Close()
}
