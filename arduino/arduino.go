package arduino

// Package arduino helps to talk to
// an Arduino connected via USB (to a Raspberry Pi).

import (
	"bufio"
	"log"
	"os"
	"sync/atomic"

	"fmt"

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

	m.file, err = os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Device '%v' opened", m.file.Name())
	return nil
}

// Read reads and returns the next line from a device file.
// Before it can be used the device file needs to be opened via Open.
func (m *Arduino) Read() (string, error) {
	if atomic.LoadInt32(&m.opened) != 1 {
		return "", errors.New("Device needs to be opened")
	}

	scanner := bufio.NewScanner(m.file)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)

		return line, nil
	}

	if err := scanner.Err(); err != nil {
		return "", errors.Wrapf(err, "Error while scanning device file '%s'", m.file.Name())
	}

	return "", nil
}

// Write writes a command to the device file (e.g., to tell the Arduino to start receiving signals).
func (m *Arduino) Write(cmd string) error {
	if atomic.LoadInt32(&m.opened) != 1 {
		return errors.New("Device needs to be opened")
	}
	b, err := m.file.WriteString(fmt.Sprintf("%s\n", cmd))
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Printf("Wrote command '%s' to '%s' (%d bytes)\n", cmd, m.file.Name(), b)

	m.file.Sync()

	return nil
}

// Close closes the opened device file.
func (m *Arduino) Close() error {
	log.Printf("Closing '%v'", m.file.Name())
	atomic.StoreInt32(&m.opened, 0)
	return m.file.Close()
}
