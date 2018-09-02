package main

import (
	"context"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var deviceFile = "/dev/ttyUSB0"

func MockedProcessor(s string) bool {
	return false
}

func TestOpen(t *testing.T) {
	writeFile(deviceFile)

	d, err := OpenDevice(deviceFile)
	defer d.Close()

	require.NoError(t, err)
}

func TestOpen_DeviceNotExist(t *testing.T) {
	writeFile("/dev/someotherdevice")

	_, err := OpenDevice(deviceFile)

	require.Error(t, err)
}

//  TODO Reset
func TestReset(t *testing.T) {
	l1 := "some line to read and process"
	l2 := "some other line to read and process"
	l3 := "this line shouldn't be processed"

	writeFile(deviceFile, l1, l2, l3)

	m := mock.Mock{}
	m.On("func1", l1).Once()
	m.On("func1", l2).Once()

	d, _ := OpenDevice(deviceFile)
	counter := 0
	d.Reset()

	m.AssertExpectations(t)
	m.AssertNotCalled(t, "func1", l3)
}

func TestRead_DeviceNotOpened(t *testing.T) {
	d := Device{}
	err := d.Process(context.Background(), MockedProcessor)

	require.Error(t, err)
}

func TestRead_DeviceClosed(t *testing.T) {
	writeFile(deviceFile)

	d, _ := OpenDevice(deviceFile)
	d.Close()

	err := d.Process(context.Background(), MockedProcessor)

	require.Error(t, err)
}

func TestRead(t *testing.T) {
	lines := []string{
		"some line to read process",
		"some other line to read and process",
	}

	writeFile(deviceFile, lines...)

	m := mock.Mock{}
	m.On("func1", lines[0]).Once()
	m.On("func1", lines[1]).Once()

	d, _ := OpenDevice(deviceFile)
	counter := 0
	err := d.Process(context.Background(), func(s string) bool {
		counter++
		m.Called(s)
		if counter == 2 {
			return true
		}
		return false
	})
	assert.NoError(t, err)
	m.AssertNumberOfCalls(t, "func1", 2)

}

func TestRead_stopProcessing(t *testing.T) {
	l1 := "some line to read and process"
	l2 := "some other line to read and process"
	l3 := "this line shouldn't be processed"

	writeFile(deviceFile, l1, l2, l3)

	m := mock.Mock{}
	m.On("func1", l1).Once()
	m.On("func1", l2).Once()

	d, _ := OpenDevice(deviceFile)
	counter := 0
	d.Process(context.Background(), func(s string) bool {
		counter++
		m.Called(s)
		if counter == 2 {
			return true
		}
		return false
	})

	m.AssertExpectations(t)
	m.AssertNotCalled(t, "func1", l3)
}

func writeFile(name string, content ...string) {
	AppFs = afero.NewMemMapFs()
	afero.WriteFile(AppFs, name, []byte(strings.Join(content, "\n")), 0644)
}
