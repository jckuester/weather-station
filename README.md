<img alt="Weather Station" src="img/logo_banner.png" />
 
---
 
<p align="right">
  <a href="https://github.com/jckuester/weather-station/releases/latest">
    <img alt="Release" src="https://img.shields.io/github/release/jckuester/weather-station.svg?style=flat-square">
  </a>
  <a href="https://github.com/jckuester/weather-station/master">
    <img alt="Travis" src="https://img.shields.io/travis/jckuester/weather-station/master.svg?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/jckuester/weather-station">
    <img alt="Go Report" src="https://goreportcard.com/badge/github.com/jckuester/weather-station?style=flat-square" />
  </a>
  <a href="https://codecov.io/gh/jckuester/weather-station">
    <img alt="Codecov branch" src="https://codecov.io/gh/jckuester/weather-station/branch/master/graph/badge.svg?style=flat-square" />
  </a>
  <a href="https://godoc.org/github.com/jckuester/weather-station">
    <img alt="Go Doc" src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square" />
  </a>
  <a href="https://github.com/jckuester/weather-station/blob/master/LICENSE">
    <img alt="Software License" src="https://img.shields.io/github/license/jckuester/weather-station.svg?style=flat-square" />
  </a>
</p>

<p>
  <img src="img/hardware.jpg" alt="Hardware of the weather station">
  <em>Figure 1: Hardware in use: Raspberry Pi, Arduino Nano, RXB6 433Mhz receiver,
  and as many GT-WT-01 temperature/humidity sensors as you like.</em>
</p>

<p>
  <img src="img/gauges.png" alt="Grafana dashboard">
  <em>Figure 2: Grafana dashboard showing an overview of current
  temperatures and humidities around my house as well as the status of alerts
  (e.g., I will get notified when my piano is too cold or humid so that I can assure it stays longer tuned).</em>
</p>

<p>
  <img src="img/fridge.png" alt="Fridge temperature and humidity (last 24h)">
  <em>Figure 3: Temperature and humidity of the fridge within the last 24h (the upper red line defines a threshold of 
  9° Celsius; so if my old fridge gets too warm I will get notified in Slack).</em>
</p>

<p>
  <img src="img/humidity.png" alt="Humidity inside (last 24h)">
  <em>Figure 4: Humidity inside rooms and piano.</em>
</p>

  
This is an opinionated and affordable setup to measure and log temperature and humidity around the house. Opinionated 
because I like Go, Prometheus, and Grafana. Affordable because each sensor costs around 10 Euros.
 
In a nutshell, this repo offers you a prometheus exporter for 433 MHz temperature/humidity sensors, where
signals are received via an Arduino ([with this software flashed to it](https://github.com/pimatic/homeduino#flashing))
connected to (in my case) a Raspberry Pi.

Happy measuring!

## Hardware

You see it all in Figure 1:

* Raspberry Pi
* Arduino Nano
* [RF unit](https://www.amazon.de/gp/product/B06XHJMC82/ref=oh_aui_detailpage_o00_s00?ie=UTF8&psc=1) for the Arduino
 (don't try to save a few dollars by buying a cheap one; it has a very tiny range)
* [GT-WT-01 temperature/humidity sensors](https://www.teknihall.be/en/node/1430)
(get as many as you like; I bought six on [eBay](https://www.ebay.com/itm/361435018543): one for outside, one for every room, one for the fridge, one for inside my piano, etc.)

***Note:*** I realised that the GT-WT-01 sensors seem to be purchasable only in the EU and UK.
However, [protocols for other sensors](https://github.com/pimatic/rfcontroljs/blob/master/protocols.md) 
can easily be added to the [supported protocols](pulse/protocol.go) of this project, too. I would be happy to see your contribution.

## Software

* [Prometheus](https://prometheus.io/)
* [Grafana](https://grafana.com/) (you can download the newest Deb packages for the Pi [here](https://github.com/fg2it/grafana-on-raspberry/releases))
* [Flash Arduino to recieve signals](https://github.com/pimatic/homeduino#flashing)
* [Slack](https://slack.com/) (if you want to get notified about alerts)

## Download

Get the binary for the Raspberry Pi (ARM) or other platforms [here](https://github.com/jckuester/weather-station/releases).

## Usage

To see options available run `$ ./weather-station --help`:
```
usage: weather-station [<flags>] [<config.yaml>]

Flags:
  --help                    Show context-sensitive help (also try --help-long and --help-man).
  --device="/dev/ttyUSB0"   Arduino connected to USB
  --listen-address=":8080"  The address to listen on for HTTP requests

Args:
  [<config.yaml>]  Path to config file.
```

### Scanning mode (find your sensors)

If this is your first time using the exporter, simply start it without any parameters (i.e. in scanning mode). 
In this case, logs will show all signals that the exporter is able to decode, but nothing is exported yet:

```
$ ./weather-station
2019/10/06 21:29:08 Config File "weather-station" Not Found in "[/etc/weather-station /home/jan/.config/weather-station /home/jan/git/github.com/weather-station]" Error reading config file. Running in scanning mode.
2019/10/06 21:29:08 Setting up serial port to use 115200 baud
2019/10/06 21:29:08 Serial port prepared
2019/10/06 21:29:08 Device '/dev/ttyUSB0' opened
2019/10/06 21:29:08 Wrote command 'RESET' to '/dev/ttyUSB0' (6 bytes)
2019/10/06 21:29:09 Line Scanned: ready
2019/10/06 21:29:09 Serving metrics at ':8080/metrics'
2019/10/06 21:29:09 Wrote command 'RF receive 0' to '/dev/ttyUSB0' (13 bytes)
2019/10/06 21:29:09 Write RF receive 0
2019/10/06 21:29:09 Line Scanned: ACK
2019/10/06 21:29:11 Line Scanned: RF receive 1008 576 7880 4064 1856 16004 0 0 00000000121313141313131413141413141414141414131413141313141314141314131314131414141414141315
2019/10/06 21:29:12 Unsupported protocol or error decoding the pulse
2019/10/06 21:29:28 Line Scanned: RF receive 560 4112 2068 9080 0 0 0 0 0102020102020201020202020102020102020202010102020201020202020101020101020103
2019/10/06 21:29:28 Sensor has no matching configuration, potential protocols:
2019/10/06 21:29:28 weather12: {ID:145 Name:145 Channel:1 Temperature:-178 Humidity:33 LowBattery:false}
2019/10/06 21:29:28 weather15: {ID:2320 Name:2320 Channel:2 Temperature:19.6 Humidity:54 LowBattery:true}
...
```

### Export mode (use your sensors)

1) To find your sensors in the logs, look for lines with `Sensor has no matching configuration, potential protocols`. The exporter tries to 
   decode every received signal with all currently supported protocols. However, not all decodings make sense. For example,
   the decoded signal using protocol `weather12` makes no sense (`Temperature: -178`), but using `weather15` does
   (as these decoded temperature/humidity values are also shown on the little sensor's display in front of me).

2) Then, write all the IDs of sensors you'd like export to Prometheus into a config file (e.g. `my-sensors.yml`):

    ```
    sensors:
      2320:
        location: kitchen
        protocol: weather15
      <another-id>:
        location: <some-location>
        protocol: <a-supported-protocol>        
    ```

3) Restart the weather station and use the config: `$ ./weather-station my-sensors.yml`

## Currently supported devices

* GT-WT-01 temperature/humidity sensor (use `weather15` protocol)
* GT-WT-01 non-variant temperature/humidity sensor (use `weather12` protocol)

Feel free to support more devices!

## Dashboards

Here is [my version of the Grafana dashboard](./grafana-dashboard.json) that you see above. Feel free to use
 it as a starting point for customizing yours (of course, sensor IDs need to be adapted accordingly).

## Credits

* Thanks for contributing the logo, [ssnjrthegr8](https://github.com/ssnjrthegr8).

* Idea of this [CO2 Prometheus exporter](https://github.com/larsp/co2monitor)
has been an inspiration and starting point for this project.

* All the hard work of reverse engineering various protocols and finding the correct decodings has been already
done by the [Pimatic project](https://github.com/pimatic/rfcontroljs).
