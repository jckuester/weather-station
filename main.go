package main

import (
	"log"
	"net/http"

	"github.com/jckuester/weather-station/sensor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	device     = kingpin.Arg("device", "Arduino connected to USB, such as /dev/ttyUSB0").Required().String()
	listenAddr = kingpin.Arg("listen-address", "The address to listen on for HTTP requests.").
			Default(":8080").String()
)

var (
	temperature = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_temperature_celsius",
		Help: "Current temperature in Celsius",
	})
	humidity = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "meter_humidity_percent",
		Help: "Current humidity level in %",
	})
)

func init() {
	prometheus.MustRegister(temperature)
	prometheus.MustRegister(humidity)
}

func main() {
	kingpin.Parse()
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
			temperature.Set(result.Temperature)
			humidity.Set(float64(result.Humidity))
		}
	}
}
