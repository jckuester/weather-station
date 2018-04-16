// Package arduino helps to talk to
// an Arduino connected via USB (to a Raspberry Pi).
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"context"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"strings"
	"sync"
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

	PingCmd  = "PING pong"
	ResetCmd = "RESET"

	ResetCmdResponse = "ready"
)

type ProcessorFunc func(s string)

// Device represents the device file of an Arduino
// connected to the USB port.
type Device struct {
	afero.File
	sync.Mutex
	open bool
}

// Open opens the named device file for reading.
func OpenDevice(name string) (*Device, error) {
	file, err := AppFs.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open '%v'", name)
	}
	log.Printf("Device '%v' opened", file.Name())

	d := &Device{
		File:  file,
		open:  true,
		Mutex: sync.Mutex{},
	}

	return d, nil
}

func (d *Device) Reset() error {
	err := d.Write(ResetCmd)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	err = d.Process(ctx, func(s string) {
		for !strings.EqualFold(s, ResetCmdResponse){
			log.Printf("Not %v msg: %v\n", ResetCmdResponse, s)
			return
		}
		cancel()
	})

	if err != nil && err != context.Canceled {
		return err
	}

	return nil
}

// ReadProcess reads the next line from a device file in a loop and
// applies a Processor to it. If the Processor returns false,
// reading is stopped.
// Before ReadProcess can be used the device file needs to be opened via Open.
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

		handle(line)

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
