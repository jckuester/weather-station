package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	device = kingpin.Flag("device", "Arduino connected to USB").
		Default("/dev/ttyUSB0").String()
	listenAddr = kingpin.Flag("listen-address", "The address to listen on for HTTP requests").
			Default(":8080").String()
	configFile = kingpin.Arg("config.yaml", "Path to config file.").String()

	temperature *prometheus.GaugeVec
	humidity    *prometheus.GaugeVec
)

const (
	// SensorID is the unique identifier of the sensor
	SensorID = "id"
	// SensorLocation is the location where the sensor is placed
	SensorLocation = "location"
)

func main() {
	kingpin.Parse()

	loadConfig(*configFile)

	setupMetrics()

	http.Handle("/metrics", promhttp.Handler())
	SetupDevice(*device)
	dev, err := OpenDevice(*device)
	if err != nil {
		log.Fatalf("Could not open '%v'", *device)
	}
	defer dev.Close()

	err = dev.Reset()
	if err != nil {
		log.Fatalf("Could not reset '%v'", *device)
	}

	go receive(dev)

	log.Printf("Serving metrics at '%v/metrics'", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func setupMetrics() {
	temperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "meter_temperature_celsius",
		Help: "Current temperature in Celsius",
	}, []string{
		SensorID,
		SensorLocation,
	})
	humidity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "meter_humidity_percent",
		Help: "Current humidity level in %",
	}, []string{
		SensorID,
		SensorLocation,
	})
	prometheus.MustRegister(temperature)
	prometheus.MustRegister(humidity)
}

func receive(a *Device) {
	// tell the Arduino to start receiving signals
	err := a.Write(ReceiveCmd)
	if err != nil {
		log.Fatalf("Could not write to '%v'", a)
	}
	log.Println("Write", ReceiveCmd)

	ctx := context.Background()
	// read and decode received signals forever
	err = a.Process(ctx, DecodedSignal)
	if err != nil {
		log.Println(err)
	}
}

// DecodedSignal decodes a compressed signal read from the Arduino
// by trying all currently supported protocols and stores result for Prometheus scraping
func DecodedSignal(line string) (stop bool) {
	stop = false

	if strings.HasPrefix(line, ReceivePrefix) {
		trimmed := strings.TrimPrefix(line, ReceivePrefix)

		pulse, err := PreparePulse(trimmed)
		if err != nil {
			log.Println(err)
			return
		}

		matchingProtocols := MatchingProtocols(pulse)

		if !processedWithMatchingConfig(matchingProtocols, pulse) {
			printAllMatchingProtocols(matchingProtocols, pulse)
		}

	}
	return
}

func printAllMatchingProtocols(matchingProtocols []string, pulse *Signal) {
	firstMatch := true
	for _, p := range matchingProtocols {
		result, err := DecodePulse(pulse, p)
		if err != nil {
			log.Println(err)
			continue
		}
		m := result.(*GTWT01Result)
		if firstMatch {
			log.Println("Sensor has no matching configuration, potential protocols:")
			firstMatch = false
		}
		log.Printf("%v: %+v\n", p, *m)
	}
	if firstMatch {
		log.Println("Unsupported protocol or error decoding the pulse")
	} else {
		log.Println("Add to configuration with appropriate protocol")
	}
}

func processedWithMatchingConfig(matchingProtocols []string, pulse *Signal) bool {
	protocolMatch := false
	configuredSensors := vip.GetStringMap("sensors")
	for id := range configuredSensors {
		location := vip.GetString(fmt.Sprintf("sensors.%s.location", id))
		if location == "" {
			panic(fmt.Errorf("fatal error sensor id %s has no location specified in config file", id))
		}
		protocol := vip.GetString(fmt.Sprintf("sensors.%s.protocol", id))
		if protocol == "" {
			panic(fmt.Errorf("fatal error sensor id %s has no protocol specified in config file", id))
		}

		for _, p := range matchingProtocols {
			if p == protocol {
				result, err := DecodePulse(pulse, protocol)
				if err != nil {
					log.Println(err)
					break
				}
				m := result.(*GTWT01Result)
				if m.Name == id {
					temperature.With(prometheus.Labels{
						SensorID:       m.Name,
						SensorLocation: location,
					}).Set(m.Temperature)

					humidity.With(prometheus.Labels{
						SensorID:       m.Name,
						SensorLocation: location,
					}).Set(float64(m.Humidity))
					log.Printf("%v: %+v\n", location, *m)
					protocolMatch = true
					break
				}

			}
		}
	}

	return protocolMatch
}
