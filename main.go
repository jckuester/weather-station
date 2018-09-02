package main

import (
	"context"
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
	ids = kingpin.Arg("id=label ...", "List of all sensor IDs (e.g. 1234=kitchen 2353=piano)"+
		" that will be exported to prometheus. Each ID must be given a human-readable label.").StringMap()

	temperature     *prometheus.GaugeVec
	humidity        *prometheus.GaugeVec
	sensorLocations map[string]string
)

const (
	sensorID = "id"
	// sensorLocation is a label name used for prometheus' GaugeVec to partition sensor data by
	// where the sensor is located (e.g., kitchen, bedroom, outside, etc.)
	sensorLocation = "location"
)

func main() {
	kingpin.Parse()

	sensorLocations = *ids

	temperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "meter_temperature_celsius",
		Help: "Current temperature in Celsius",
	}, []string{
		sensorID,
		sensorLocation,
	})
	humidity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "meter_humidity_percent",
		Help: "Current humidity level in %",
	}, []string{
		sensorID,
		sensorLocation,
	})

	prometheus.MustRegister(temperature)
	prometheus.MustRegister(humidity)

	http.Handle("/metrics", promhttp.Handler())

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

func receive(a *Device) {
	// tell the Arduino to start receiving signals
	err := a.Write(ReceiveCmd)
	if err != nil {
		log.Fatalf("Could not write to '%v'", a)
	}
	log.Println("Write", ReceiveCmd)

	ctx := context.Background()
	// read and decode received signals forever
	err = a.Process(ctx, DecodeSignal)
	if err != nil {
		log.Println(err)
	}
}

// DecodeSignal decodes a compressed signal read from the Arduino
// by trying all currently supported protocols.
func DecodeSignal(line string) (stop bool) {
	stop = false

	if strings.HasPrefix(line, ReceivePrefix) {
		trimmed := strings.TrimPrefix(line, ReceivePrefix)

		p, err := PreparePulse(trimmed)
		if err != nil {
			log.Println(err)
			return
		}

		device, result, err := DecodePulse(p)
		if err != nil {
			log.Println(err)
			return
		}

		switch device {
		case Weather15:
			m := result.(*GTWT01Result)
			if loc, ok := sensorLocations[m.ID]; !ok || loc == "" {
				log.Printf("%v: %+v\n", device, *m)
				log.Println("Sensor hasn't set a location and won't be provided to Prometheus for monitoring")
				return
			}
			log.Printf("%v(%v): %+v\n", device, sensorLocations[m.ID], *m)

			temperature.With(prometheus.Labels{
				sensorID:       m.ID,
				sensorLocation: sensorLocations[m.ID],
			}).Set(m.Temperature)

			humidity.With(prometheus.Labels{
				sensorID:       m.ID,
				sensorLocation: sensorLocations[m.ID],
			}).Set(float64(m.Humidity))
		default:
			log.Println("Devices", device)
		}
	}
	return
}
