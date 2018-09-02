package main

import (
	"fmt"
	"strconv"
)

//go:generate stringer -type ProtocolType

type ProtocolType uint8

const (
	// Unknown is the type of protocol returned if the protocol used to decode a signal is (currently) unknown.
	Unknown ProtocolType = iota
	// Weather15 is the only protocol type supported right now
	// (see: https://github.com/pimatic/rfcontroljs/blob/master/src/protocols/weather15.coffee)
	Weather15
)

// Protocol defines a protocol that can be used to match
// received signals and decode them.
type Protocol struct {
	Devices   string                            // comma-separated list describing all the devices using this protocol
	SeqLength []int                             // allowed lengths of the sequence of pulses
	Lengths   []int                             // pulse lengths
	Mapping   map[string]string                 // maps the pulse sequence into binary representation (i.e. 0s and 1s)
	Type      ProtocolType                      // the type defining this protocol
	Decode    func(string) (interface{}, error) // decodes the binary representation into a human-readable struct
}

// Protocols returns a list of all the currently supported
// protocols that can be used for trying to decode received signals.
func Protocols() []*Protocol {
	return []*Protocol{
		{
			Devices:   "Globaltronics GT-WT-01 variant",
			SeqLength: []int{76},
			Lengths:   []int{496, 2048, 4068, 8960},
			Mapping: map[string]string{
				"01": "0",
				"02": "1",
				"03": "",
			},
			Type: Weather15,
			Decode: func(binSeq string) (interface{}, error) {
				id, err := strconv.ParseUint(binSeq[0:12], 2, 0)
				if err != nil {
					return nil, err
				}

				channel, err := strconv.ParseUint(binSeq[14:16], 2, 0)
				if err != nil {
					return nil, err
				}

				temp, err := strconv.ParseInt(binSeq[16:28], 2, 0)
				if err != nil {
					return nil, err
				}

				humidity, err := strconv.ParseInt(binSeq[28:36], 2, 0)
				if err != nil {
					return nil, err
				}

				lowBattery, err := strconv.ParseBool(string(binSeq[12]))
				if err != nil {
					return nil, err
				}

				return &GTWT01Result{
					ID:          fmt.Sprint(id),
					Channel:     int(channel) + 1,
					Temperature: float64(temp) / 10,
					Humidity:    int(humidity),
					LowBattery:  lowBattery,
				}, nil
			},
		},
	}
}

// GTWT01Result is the human-readable result of a decoded pulse
// for the "GT-WT-01 variant".
type GTWT01Result struct {
	ID          string
	Channel     int
	Temperature float64
	Humidity    int
	LowBattery  bool
}
