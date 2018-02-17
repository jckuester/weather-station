// Package protocol contains supported protocols that
// can be used to decode signals received via the
// https://github.com/pimatic/RFControl library.
package protocol

import (
	"github.com/jckuester/weather-station/binary"
)

// Protocol defines a protocol that is used to match
// received signals and decode them.
type Protocol struct {
	Device    string                            // the type of device that uses the protocol
	SeqLength []int                             // allowed lengths of the sequence of pulses
	Lengths   []int                             // pulse lengths
	Mapping   map[string]string                 // maps the pulse sequence into binary representation (i.e. 0s and 1s)
	Decode    func(string) (interface{}, error) // decodes the binary representation into a human-readable struct
}

// Measurement is the result of a decoded pulse
// of the protocol for the "Globaltronics GT-WT-01 variant".
type Measurement struct {
	Id          int
	Channel     int
	Temperature float64
	Humidity    int
	LowBattery  bool
}

// Supported returns a list of all the currently supported
// protocols that can be used for trying to decode received signals.
func Supported() []*Protocol {
	return []*Protocol{
		// Only one protocol supported right now
		// (see: https://github.com/pimatic/rfcontroljs/blob/master/src/protocols/weather15.coffee)
		{
			Device:    "Globaltronics GT-WT-01 variant",
			SeqLength: []int{76},
			Lengths:   []int{496, 2048, 4068, 8960},
			Mapping: map[string]string{
				"01": "0",
				"02": "1",
				"03": "",
			},
			Decode: func(binSeq string) (interface{}, error) {
				id, err := binary.ToNumber(binSeq, 0, 12)
				if err != nil {
					return nil, err
				}

				channel, err := binary.ToNumber(binSeq, 14, 16)
				if err != nil {
					return nil, err
				}

				temp, err := binary.ToSignedNumber(binSeq, 16, 27)
				if err != nil {
					return nil, err
				}

				humidity, err := binary.ToNumber(binSeq, 28, 36)
				if err != nil {
					return nil, err
				}

				return &Measurement{
					Id:          id,
					Channel:     channel + 1,
					Temperature: float64(temp) / 10,
					Humidity:    humidity,
					LowBattery:  binary.ToBoolean(binSeq, 12),
				}, nil
			},
		},
	}
}
