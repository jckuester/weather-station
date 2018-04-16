package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/http"
	"strings"
)

var (
	device = kingpin.Flag("device", "Arduino connected to USB").
		Default("/dev/ttyUSB0").String()
	listenAddr = kingpin.Flag("listen-address", "The address to listen on for HTTP requests.").
			Default(":8080").String()
	ids = kingpin.Arg("ids", "Sensor IDs that will be exported").StringMap()

	temperature *prometheus.GaugeVec
	humidity    *prometheus.GaugeVec
	sensors     map[string]string
)

const (
	SENSOR_ID       = "id"
	SENSOR_LOCATION = "location"
)

func main() {
	kingpin.Parse()

	sensors = *ids

	temperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "meter_temperature_celsius",
		Help: "Current temperature in Celsius",
	}, []string{
		SENSOR_ID,
		SENSOR_LOCATION,
	})
	humidity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "meter_humidity_percent",
		Help: "Current humidity level in %",
	}, []string{
		SENSOR_ID,
		SENSOR_LOCATION,
	})

	prometheus.MustRegister(temperature)
	prometheus.MustRegister(humidity)

	http.Handle("/metrics", promhttp.Handler())

	dev, err := OpenDevice(*device)
	if err != nil {
		log.Fatalf("Could not open '%v'", *device)
	}
	defer dev.Close()

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
	err = a.Process(ctx, DecodedSignal)
	if err != nil {
		log.Println(err)
	}
}

// Process decodes a compressed signal read from the Arduino
// by trying all currently supported protocols.
func DecodedSignal(line string) {
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
		case GT_WT_01:
			m := result.(*GTWT01Result)
			log.Printf("%v: %+v\n", device, *m)
			temperature.With(prometheus.Labels{
				SENSOR_ID:       m.Name,
				SENSOR_LOCATION: sensors[m.Name],
			}).Set(m.Temperature)

			humidity.With(prometheus.Labels{
				SENSOR_ID:       m.Name,
				SENSOR_LOCATION: sensors[m.Name],
			}).Set(float64(m.Humidity))
		default:
			log.Println("Device", device)
		}
	}
	return
}
