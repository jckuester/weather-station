# Weather station

<p align="right">
  <a href="https://goreportcard.com/report/github.com/jckuester/weather-station">
  <img src="https://goreportcard.com/badge/github.com/cloudetc/awsweeper" /></a>
  <a href="https://godoc.org/github.com/jckuester/weather-station">
  <img src="https://godoc.org/github.com/cloudetc/awsweeper?status.svg" /></a>
</p>

<p>
 <img src="img/hardware.jpg" alt="Hardware used for the weather station">
 <em>Figure 1: Hardware used: Raspberry Pi, Arduino Nano, RXB6 433Mhz receiver,
 and as many GT-WT-01 temperature/humidity sensors as you like.</em>
</p>

<p>
 <img src="img/gauges.png" alt="Grafana dashboard">
 <em>Figure 2: Grafana dashboard showing the current values of some temperature/humidity sensors as well as the status of the alerts.</em>
</p>

<p>
 <img src="img/fridge.png" alt="Temerperature/humidity of the fridge within the last 24h">
 <em>Figure 3: A graph in Grafana showing the temerperature/humidity of the fridge within the last 24h.</em>
</p>

<p>
 <img src="img/humidity.png" alt="Humidity (24h)">
 <em>Figure 4: A graph in Grafana showing the humidity inside (the upper red line is defines a threshold;
 crossing it will raise an alert in Slack).</em>
</p>

## Stuff you need

Hardware (see Fig. 1):

* Raspberry Pi
* Arduino Nano
* [RF unit](https://www.amazon.de/gp/product/B06XHJMC82/ref=oh_aui_detailpage_o00_s00?ie=UTF8&psc=1) for the Arduino
 (don't try to safe a few dollars by buying the cheap one; it has a very tiny range)
* [GT-WT-01 temperature/humidity sensors](https://www.ebay.com/itm/361435018543)
(get as many as you like; I bought six: one for outside, one for every room, one for the fridge, one inside my piano, etc.)

Software:
* [Prometheus](https://prometheus.io/)
* [Grafana](https://grafana.com/) (you can download the newest Deb packages for the Pi [here](https://github.com/fg2it/grafana-on-raspberry/releases))
* [Flash Arduino to recieve signals](https://github.com/pimatic/homeduino#flashing)
* [Slack](https://slack.com/) (if you want to get notified about alerts)

## Usage

```
$ ./weather-station --help
  
Usage: weather-station --device=DEVICE [<flags>] <ids>...
  
Flags:
  --help                    Show context-sensitive help (also try --help-long and --help-man).
  --device=DEVICE           Arduino connected to USB, such as /dev/ttyUSB0
  --listen-address=":8080"  The address to listen on for HTTP requests.
  
Args:
  <ids>  Device ids of the sensors
...
```

Starting the exporter:

```
$ ./weather-station --device="/dev/ttyUSB0" 2321 2454
2018/02/16 19:51:05 RF receive 544 4124 2088 9080 0 0 0 0 0102020101020202020101010202020102020202010201010102020102010202020101020103
2018/02/16 19:51:05 {Id:2439 Channel:2 Temperature:18.5 Humidity:70 LowBattery:false}
...
```

## Dashboards

Here is [my version of the Grafana dashboard](./grafana-dashboard.json) that you see above. Feel free to use
 it as a starting point for customizing yours (of course, sensor IDs need to be adapted accordingly).

## Credits

* The idea and code of this [CO2 Prometheus exporter](https://github.com/larsp/co2monitor), also written in Go,
has been inspiration and starting point for this project.

* All that hard work of reverser engineering the protocols and finding the correct decodings has been
done [here](https://github.com/pimatic/rfcontroljs).
