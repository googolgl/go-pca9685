/*
Copyright (c) 2020
Author: Pavlo Lytvynoff

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package pca9685

import (
	"fmt"
)

const (
	// The specified pulse width range of a servo has historically been 1000-2000us,
	// for a 90 degree range of motion. But nearly all modern servos have a 170-180
	// degree range, and the pulse widths can go well out of the range to achieve this
	// extended motion. The default values here of `750` and `2250` typically give
	// 135 degrees of motion. You can set `Range` to correspond to the
	// actual range of motion you observe with your given `MinPulse` and `MaxPulse` values.
	ServoRangeDef    int     = 135
	ServoMinPulseDef float32 = 750.0
	ServoMaxPulseDef float32 = 2250.0
)

// Servo structure
type Servo struct {
	PWM     *PCA9685
	Channel uint8
	Options *ServOptions
}

// ServOptions for servo
type ServOptions struct {
	Range    int // actuation range
	MinPulse float32
	MaxPulse float32
}

// ServNew creates a new servo driver
func ServoNew(p *PCA9685, chn uint8, o *ServOptions) *Servo {
	s := &Servo{
		PWM:     p,
		Channel: chn,
		Options: &ServOptions{
			Range:    ServoRangeDef,
			MinPulse: ServoMinPulseDef,
			MaxPulse: ServoMaxPulseDef,
		},
	}
	if o != nil {
		s.Options = o
	}
	return s
}

// Angle in degrees. Must be in the range `0` to `Range`.
func (s *Servo) Angle(a int) (err error) {
	if a < 0 || a > s.Options.Range {
		return fmt.Errorf("Angle out of range")
	}
	return s.Fraction(float32(a) / float32(s.Options.Range))
}

// Fraction as pulse width expressed between 0.0 `MinPulse` and 1.0 `MaxPulse`.
// For conventional servos, corresponds to the servo position as a fraction
// of the actuation range.
func (s *Servo) Fraction(f float32) (err error) {
	if f < 0.0 || f > 1.0 {
		return fmt.Errorf("Must be 0.0 to 1.0")
	}

	freq := s.PWM.GetFreq()

	minDuty := s.Options.MinPulse * freq / 1000000 * 0xFFFF
	maxDuty := s.Options.MaxPulse * freq / 1000000 * 0xFFFF
	dutyRange := maxDuty - minDuty
	dutyCycle := (int(minDuty+f*dutyRange) + 1) >> 4

	return s.PWM.SetChannel(int(s.Channel), 0, dutyCycle)
}
