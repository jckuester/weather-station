package main

import (
	"bytes"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	vip = viper.New()

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	loadConfig("")

	assert.Equal(t, len(vip.GetStringMap("sensors")), 0, "Should initialize empty sensor list")
	assert.Contains(t, buf.String(), "Not Found")

}

func TestConfigWithSensors(t *testing.T) {
	vip = viper.New()

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	loadSampleConfig()
	sensors := vip.GetStringMap("sensors")
	assert.Equal(t, len(sensors), 2, "Should initialize sensor list from config file")
	assert.NotNil(t, sensors["91"], "Sensor 91 exists in config")
	assert.Nil(t, sensors["0815"], "Sensor 0185 does not exist")
	assert.Equal(t, vip.GetString("sensors.1235.location"), "kitchen", "Sensor location is kitchen")
	assert.Equal(t, vip.GetString("sensors.1235.maeh"), "", "Sensor location is kitchen")
	assert.Equal(t, buf.String(), "")
}

func loadSampleConfig() {
	vip = viper.New()

	initConfig("")
	vip.SetConfigName("sample-weather-station")
	readConfig()
}
