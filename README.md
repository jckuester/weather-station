# Weather station

<p>
 <img src="img/hardware.jpg" alt="Hardware for the weather station">   
 <em>A home-made weather station using: Raspberry Pi, Arduino Nano, a RXB6 433Mhz receiver,
 and as many GT-WT-01 temperature/humidity sensors as you wish.</em>
</p>

## What you need

Hardware:

* Raspberry Pi
* Arduino Nano
* [RF unit](https://www.amazon.de/gp/product/B06XHJMC82/ref=oh_aui_detailpage_o00_s00?ie=UTF8&psc=1) (don't get the cheap one, it has a very limited range)
* [GT-WT-01 temperature/humidity sensors](https://www.ebay.com/itm/361435018543) (as many as you like; I have six: one outside, one in every room, one in the fridge, one inside my piano, etc.)

Software:
* Prometheus
* Slack
* Grafana
* [Flash Arduino to recieve signals](https://github.com/pimatic/homeduino#flashing)

## Dashboards

## Alerts

* [Create Slack workspace](https://slack.com/intl/de-de/get-started) 

## Credits

* The idea and code of this [CO2 prometheus-exporter](https://github.com/larsp/co2monitor) written in go 
was the initial starting point for this project.  

* All that hard work of decoding has been done [here](https://github.com/pimatic/rfcontroljs).
