package pulse

import (
	"github.com/jckuester/weather-station/binary"
)

// Protocol defines a protocol that can be used to match
// received signals and decode them.
type Protocol struct {
	Device    string                            // the type of device that uses the protocol
	SeqLength []int                             // allowed lengths of the sequence of pulses
	Lengths   []int                             // pulse lengths
	Mapping   map[string]string                 // maps the pulse sequence into binary representation (i.e. 0s and 1s)
	Decode    func(string) (interface{}, error) // decodes the binary representation into a human-readable struct
}

// Protocols returns a list of all the currently supported
// protocols that can be used for trying to decode received signals.
func Protocols() map[string]*Protocol {
	return map[string]*Protocol{
		// Only one protocol supported right now
		// (see: https://github.com/pimatic/rfcontroljs/blob/master/src/protocols/weather15.coffee)
		"protocol1": {
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

				return &GTWT01Result{
					ID:          id,
					Channel:     channel + 1,
					Temperature: float64(temp) / 10,
					Humidity:    humidity,
					LowBattery:  binary.ToBoolean(binSeq, 12),
				}, nil
			},
		},
	}
}

// GTWT01Result is the human-readable result of a decoded pulse
// for the "GT-WT-01 variant".
type GTWT01Result struct {
	ID          int
	Channel     int
	Temperature float64
	Humidity    int
	LowBattery  bool
}
