package main

import (
	"fmt"
	"strconv"
)

//DeviceType for enumeration
//go:generate stringer -type DeviceType
type DeviceType uint8

const (
	GT_WT_01 DeviceType = iota
	GT_WT_01_variant
)

// Protocol defines a protocol that can be used to match
// received signals and decode them.
type Protocol struct {
	Device    string            // the type of device that uses the protocol
	SeqLength []int             // allowed lengths of the sequence of pulses
	Lengths   []int             // pulse lengths
	Mapping   map[string]string // maps the pulse sequence into binary representation (i.e. 0s and 1s)
	Type      DeviceType
	Decode    func(string) (interface{}, error) // decodes the binary representation into a human-readable struct
}

// Protocols returns a map of all the currently supported
// protocols that can be used for trying to decode received signals.
func Protocols() map[string]*Protocol {
	return map[string]*Protocol{
		// (see: https://github.com/pimatic/rfcontroljs/blob/master/src/protocols/weather15.coffee)
		"weather15": {
			Device:    "Globaltronics GT-WT-01 variant",
			SeqLength: []int{76},
			Lengths:   []int{496, 2048, 4068, 8960},
			Mapping: map[string]string{
				"01": "0",
				"02": "1",
				"03": "",
			},
			Type: GT_WT_01_variant,
			Decode: func(binSeq string) (interface{}, error) {
				id, err := strconv.ParseUint(binSeq[0:12], 2, 0)
				if err != nil {
					return nil, err
				}

				channel, err := strconv.ParseUint(binSeq[14:16], 2, 0)
				if err != nil {
					return nil, err
				}

				temp, err := parse12BitSignedInt(binSeq[16:28])
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
					ID:          int(id),
					Name:        fmt.Sprint(id),
					Channel:     int(channel) + 1,
					Temperature: float64(temp) / 10,
					Humidity:    int(humidity),
					LowBattery:  lowBattery,
				}, nil
			},
		},
		// (see: https://github.com/pimatic/rfcontroljs/blob/master/src/protocols/weather12.coffee)
		"weather12": {
			Device:    "Globaltronics GT-WT-01 non-variant",
			SeqLength: []int{76},
			Lengths:   []int{496, 2048, 4068, 8960},
			Mapping: map[string]string{
				"01": "0",
				"02": "1",
				"03": "",
			},
			Type: GT_WT_01,
			Decode: func(binSeq string) (interface{}, error) {
				id, err := strconv.ParseUint(binSeq[0:8], 2, 0)
				if err != nil {
					return nil, err
				}

				channel, err := strconv.ParseUint(binSeq[10:12], 2, 0)
				if err != nil {
					return nil, err
				}

				temp, err := parse12BitSignedInt(binSeq[12:24])
				if err != nil {
					return nil, err
				}

				humidity, err := strconv.ParseInt(binSeq[24:31], 2, 0)
				if err != nil {
					return nil, err
				}

				lowBattery, err := strconv.ParseBool(string(binSeq[8]))
				if err != nil {
					return nil, err
				}

				return &GTWT01Result{
					ID:          int(id),
					Name:        fmt.Sprint(id),
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
// for the "GT-WT-01 variant and non-variant".
type GTWT01Result struct {
	ID          int
	Name        string
	Channel     int
	Temperature float64
	Humidity    int
	LowBattery  bool
}

func parse12BitSignedInt(s string) (int64, error) {
	if s[0] == '1' {
		v, err := strconv.ParseInt(s, 2, 0)
		if err != nil {
			return 0, err
		}
		return v - 4096, nil // 2**12
	}
	return strconv.ParseInt(s, 2, 0)
}
