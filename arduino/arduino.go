// Package arduino helps to talk to
// an Arduino connected via USB (to a Raspberry Pi).
package arduino

import (
	"bufio"
	"log"
	"os"
	"sync/atomic"

	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

var (
	// AppFs is an abstraction of the file system
	// to allow mocking in tests.
	AppFs = afero.NewOsFs()
)

const (
	// ReceiveCmd is the command that needs to b e sent to the Arduino to start
	// receiving on Pin 0.
	ReceiveCmd = "RF receive 0"
	// ReceivePrefix is the prefix of raw signals read from the device file of the Arduino.
	ReceivePrefix = "RF receive "
	// Ready is the statement returned by the Arduino when it's ready
	// to receive commands.
	Ready = "ready"
)

// Device represents the device file of an Arduino
// connected to the USB port.
type Device struct {
	file   afero.File
	opened int32
}

// Open opens the named device file for reading.
func (d *Device) Open(name string) (err error) {
	atomic.StoreInt32(&d.opened, 1)

	d.file, err = AppFs.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		return errors.Wrapf(err, "Failed to open '%v'", name)
	}

	log.Printf("Device '%v' opened", d.file.Name())
	return nil
}

// Read reads the next line from a device file in a loop and
// applies a Processor to it. If the Processor returns false,
// reading is stopped.
// Before Read can be used the device file needs to be opened via Open.
func (d *Device) Read(p Processor) error {
	if atomic.LoadInt32(&d.opened) != 1 {
		return errors.New("Device needs to be open")
	}

	scanner := bufio.NewScanner(d.file)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)

		if !p.Process(line) {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return errors.Wrapf(err, "Error while scanning device file '%s'", d.file.Name())
	}

	return nil
}

// Write writes a command to the device file
// (e.g. to tell the Arduino to start receiving signals).
func (d *Device) Write(cmd string) error {
	if atomic.LoadInt32(&d.opened) != 1 {
		return errors.New("Device needs to be opened")
	}
	b, err := d.file.WriteString(fmt.Sprintf("%s\n", cmd))
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Wrote command '%s' to '%s' (%d bytes)\n", cmd, d.file.Name(), b)

	d.file.Sync()

	return nil
}

// Close closes the opened device file.
func (d *Device) Close() error {
	log.Printf("Closing '%v'", d.file.Name())
	atomic.StoreInt32(&d.opened, 0)
	return d.file.Close()
}

// Processor defines a function that can be applied
// to each line read from the device file.
// Let the function return false if reading should be stopped.
type Processor interface {
	Process(s string) bool
}
