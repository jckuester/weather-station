package main

import (
	"log"
	"net/http"

	"fmt"

	"strings"

	"github.com/jckuester/weather-station/arduino"
	"github.com/jckuester/weather-station/protocol"
	"github.com/jckuester/weather-station/pulse"
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
	a := &arduino.Arduino{}
	err := a.Open(*device)
	if err != nil {
		log.Fatalf("Could not open '%v'", *device)
		return
	}

	line, err := a.Read()
	if err != nil {
		log.Println(err)
	}
	if strings.Contains(line, "ready") {
		err = a.Write("RF receive 0")
		if err != nil {
			log.Fatalf("Could not write to '%v'", *device)
			return
		}
	}

	for {
		line, err := a.Read()
		if err != nil {
			log.Println(err)
			continue
		}

		if strings.HasPrefix(line, "RF receive") {
			pulseTrimmed := strings.TrimPrefix(line, "RF receive ")

			p, err := pulse.PrepareCompressed(pulseTrimmed)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Printf("%+v\n", *p)

			m, err := pulse.Decode(p)
			if err != nil {
				log.Println(err)
				continue
			}

			if m != nil {
				m := m.(*protocol.Measurement)
				log.Printf("%+v\n", *m)

				if t, ok := temperature[m.Id]; ok {
					t.Set(m.Temperature)
				}
				if h, ok := humidity[m.Id]; ok {
					h.Set(float64(m.Humidity))
				}
			}
		}
	}
}
