package main

import (
	"log"
	"net/http"

	"fmt"

	"github.com/jckuester/weather-station/sensor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	device     = kingpin.Flag("device", "Arduino connected to USB, such as /dev/ttyUSB0").Required().String()
	listenAddr = kingpin.Flag("listen-address", "The address to listen on for HTTP requests.").
			Default(":8080").String()
	ids = kingpin.Arg("ids", "Device ids of the sensors").Required().Ints()
)

var (
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
	go measure()

	log.Printf("Serving metrics at '%v/metrics'", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func measure() {
	sensor := &sensor.Sensor{}
	err := sensor.Open(*device)
	if err != nil {
		log.Fatalf("Could not open '%v'", *device)
		return
	}

	for {
		result, err := sensor.Read()
		if err != nil {
			log.Fatalf("Something went wrong: '%v'", err)
		}

		if result != nil {
			if t, ok := temperature[result.Id]; ok {
				t.Set(result.Temperature)
			}
			if h, ok := humidity[result.Id]; ok {
				h.Set(float64(result.Humidity))
			}
		}
	}
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
