package arduino

import (
	"log"
	"os"
	"sync/atomic"

	"strings"

	"bufio"

	"github.com/jckuester/weather-station/protocol"
	"github.com/jckuester/weather-station/pulse"
	"github.com/pkg/errors"
)

// Arduino represents an Arduino
// connected to the USB port.
type Arduino struct {
	file   *os.File
	opened int32
}

// Open opens the named device file for reading,
// which is usually /dev/ttyUSB0 for an Arduino connected to a Raspberry Pi.
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
			log.Println(p)

			// only one protocol that we need for the weather station is supported right now
			pc := protocol.Weather15
			m, err := pulse.Decode(p, pc)
			return m.(*protocol.Measurement), errors.Wrapf(err, "Failed decode pulse info '%s' for protocol '%s'", p, pc)
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
