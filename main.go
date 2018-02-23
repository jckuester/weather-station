package main

import (
	"log"
	"net/http"

	"fmt"

	"strings"

	"github.com/jckuester/weather-station/arduino"
	"github.com/jckuester/weather-station/pulse"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	device = kingpin.Flag("device", "Arduino connected to USB").
		Default("/dev/ttyUSB0").String()
	listenAddr = kingpin.Flag("listen-address", "The address to listen on for HTTP requests.").
			Default(":8080").String()
	ids = kingpin.Arg("ids", "Sensor IDs that will be exported").Ints()

	temperature = make(map[int]prometheus.Gauge)
	humidity    = make(map[int]prometheus.Gauge)
)

func main() {
	kingpin.Parse()

	for _, i := range *ids {
		temperature[i] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: fmt.Sprintf("meter_temperature_celsius_%d", i),
			Help: "Current temperature in Celsius",
		})
		humidity[i] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: fmt.Sprintf("meter_humidity_percent_%d", i),
			Help: "Current humidity level in %",
		})

		prometheus.MustRegister(temperature[i])
		prometheus.MustRegister(humidity[i])
	}

	http.Handle("/metrics", promhttp.Handler())

	go receive(device)

	log.Printf("Serving metrics at '%v/metrics'", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func receive(device *string) {
	a := &arduino.Device{}

	err := a.Open(*device)
	if err != nil {
		log.Fatalf("Could not open '%v'", *device)
	}

	// wait until the Arduino is ready to accept commands
	err = a.Read(Ready{})
	if err != nil {
		log.Fatalf("Device is not ready to take commands: %s", err)
	}

	// tell the Arduino to start receiving signals
	err = a.Write(arduino.ReceiveCmd)
	if err != nil {
		log.Fatalf("Could not write to '%v'", *device)
	}

	// read and decode received signals forever
	err = a.Read(DecodedSignal{})
	if err != nil {
		log.Println(err)
	}
}

// Ready implements a Processor that waits and returns
// as soon as the Arduino is ready to accept commands.
type Ready struct{}

// Process reads from the Arduino until it returns "ready".
func (Ready) Process(s string) bool {
	return strings.Contains(s, arduino.Ready)
}

// DecodedSignal implements a Processor that
// decodes received raw signals and sets the Prometheus Gauge
// based on the decoded values.
type DecodedSignal struct{}

// Process decodes a compressed signal read from the Arduino
// by trying all currently supported protocols.
func (DecodedSignal) Process(line string) bool {
	if strings.HasPrefix(line, arduino.ReceivePrefix) {
		trimmed := strings.TrimPrefix(line, arduino.ReceivePrefix)

		p, err := pulse.Prepare(trimmed)
		if err != nil {
			log.Println(err)
			return true
		}

		result, err := pulse.Decode(p)
		if err != nil {
			log.Println(err)
			return true
		}

		if result != nil {
			m := result.(*pulse.GTWT01Result)
			log.Printf("%+v\n", *m)

			if t, ok := temperature[m.ID]; ok {
				t.Set(m.Temperature)
			}
			if h, ok := humidity[m.ID]; ok {
				h.Set(float64(m.Humidity))
			}
		}
	}
	return true
}
