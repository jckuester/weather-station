package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"context"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

var (
	// AppFs is an abstraction of the file system
	// to allow mocking in tests.
	AppFs = afero.NewOsFs()
)

const (
	// ReceiveCmd is the command to be sent to the Arduino to start receiving (on Pin 0).
	ReceiveCmd = "RF receive 0"
	// ReceivePrefix is the prefix of raw signals read from the Arduino.
	ReceivePrefix = "RF receive "
	// ResetCmd resets the Arduino (or better the homeduino software installed on it) into an expected state.
	ResetCmd = "RESET"
	// ResetCmdResponse is the response to ResetCmd when resetting is finished.
	ResetCmdResponse = "ready"
)

// ProcessorFunc are implementations that can be applied
// to each line read from the device file.
// Let the function return true if otherwise continuous reading should be stopped.
type ProcessorFunc func(s string) bool

// Device represents the device file of an Arduino
// connected to the USB port.
type Device struct {
	afero.File
	sync.Mutex
	open bool
}

// OpenDevice opens the named device file for reading.
func OpenDevice(name string) (*Device, error) {
	file, err := AppFs.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open '%v'", name)
	}
	log.Printf("Device '%v' is open", file.Name())

	d := &Device{
		File:  file,
		open:  true,
		Mutex: sync.Mutex{},
	}

	return d, nil
}

// Reset sends ResetCmd to the device to set it in an expected state
// (i.e. waits until a ResetCmdResponse is received).
func (d *Device) Reset() error {
	err := d.Write(ResetCmd)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = d.Process(ctx, func(s string) bool {
		if !strings.Contains(s, ResetCmdResponse) {
			log.Printf("Not %v msg: %v\n", ResetCmdResponse, s)
			return false
		}
		return true
	})

	if err != nil && err != context.Canceled {
		return err
	}

	return nil
}

// Process reads the next line from a device file in a loop and
// applies a ProcessorFunc to it. The Context is used to stop reading.
func (d *Device) Process(ctx context.Context, handle ProcessorFunc) error {
	d.Lock()
	defer d.Unlock()
	if !d.open {
		return errors.New("File already closed")
	}

	scanner := bufio.NewScanner(d)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println("Line Scanned:", line)

		stop := handle(line)
		if stop {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	if err := scanner.Err(); err != nil {
		return errors.Wrapf(err, "Error while scanning device file '%s'", d.Name())
	}

	return nil
}

// Write writes a command to the device file
// (e.g. to tell the Arduino to start receiving signals).
func (d *Device) Write(cmd string) error {
	d.Lock()
	if !d.open {
		d.Unlock()
		return errors.New("File already closed")
	}
	d.Unlock()

	b, err := d.WriteString(fmt.Sprintf("%s\n", cmd))
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Wrote command '%s' to '%s' (%d bytes)\n", cmd, d.Name(), b)

	d.Sync()

	return nil
}

// Close closes the opened device file.
func (d *Device) Close() error {
	d.Lock()
	if !d.open {
		d.Unlock()
		return errors.New("File already closed")
	}
	log.Printf("Closing '%v'", d.Name())
	d.open = false
	d.Unlock()
	return d.Close()
}
