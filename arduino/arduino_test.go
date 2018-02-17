package arduino

import (
	"testing"

	"strings"

	"log"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var device = "/dev/ttyUSB0"

type MockedProcessor struct {
	mock.Mock
}

func (m MockedProcessor) Process(s string) bool {
	log.Println(s)
	args := m.Called(s)
	log.Println(args)
	return args.Bool(0)
}

func TestOpen(t *testing.T) {
	writeFile(device)

	d := Device{}
	err := d.Open(device)
	defer d.Close()

	require.NoError(t, err)
}

func TestOpen_DeviceNotExist(t *testing.T) {
	writeFile("/dev/someotherdevice")

	d := Device{}
	err := d.Open(device)

	require.Error(t, err)
}

func TestRead_DeviceNotOpened(t *testing.T) {
	d := Device{}
	err := d.Read(MockedProcessor{})

	require.Error(t, err)
}

func TestRead_DeviceClosed(t *testing.T) {
	writeFile(device)

	d := Device{}
	d.Open(device)
	d.Close()

	err := d.Read(MockedProcessor{})

	require.Error(t, err)
}

func TestRead(t *testing.T) {
	l1 := "some line to read process"
	l2 := "some other line to read and process"
	writeFile(device, l1, l2)

	m := MockedProcessor{}
	m.On("Process", l1).Return(true).Once()
	m.On("Process", l2).Return(true).Once()

	d := Device{}
	d.Open(device)

	d.Read(m)

	m.AssertExpectations(t)
}

func TestRead_stopProcessing(t *testing.T) {
	l1 := "some line to read and process"
	l2 := "some other line to read and process"
	l3 := "this line shouldn't be processed"

	writeFile(device, l1, l2, l3)

	m := MockedProcessor{}
	m.On("Process", l1).Return(true).Once()
	m.On("Process", l2).Return(false).Once()

	d := Device{}
	d.Open(device)

	d.Read(m)

	m.AssertExpectations(t)
	m.AssertNotCalled(t, "process", l3)
}

func writeFile(name string, content ...string) {
	AppFs = afero.NewMemMapFs()
	afero.WriteFile(AppFs, name, []byte(strings.Join(content, "\n")), 0644)
}
