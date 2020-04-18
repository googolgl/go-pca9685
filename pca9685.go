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
	"time"

	i2c "github.com/d2r2/go-i2c"
)

const (
	// Address default for controller
	Address byte = 0x40

	// Registers
	Mode1    byte = 0x00
	Prescale byte = 0xFE
	Led0On   byte = 0x06

	// The internal reference clock is 25mhz but may vary slightly with
	// environmental conditions and manufacturing variances. Providing a more precise
	// "ReferenceClockSpeed" can improve the accuracy of the frequency and duty_cycle computations.
	ReferenceClockSpeed float32 = 25000000.0 // 25MHz
	StepCount           float32 = 4096.0     // 12-bit
	DefaultPWMFrequency float32 = 50.0       // 50Hz
)

// PCA9685 is a Driver for the PCA9685 16-channel 12-bit PWM/Servo controller
type PCA9685 struct {
	Conn *i2c.I2C
	Optn *Options
}

// Options for controller
type Options struct {
	Name       string
	Frequency  float32
	ClockSpeed float32
}

// PCANew creates a new driver with specified i2c interface
func PCANew(i2c *i2c.I2C, optn *Options) *PCA9685 {
	adr := i2c.GetAddr()
	pca := &PCA9685{
		Conn: i2c,
		Optn: &Options{
			Name:       "Controller" + fmt.Sprintf("-0x%x", adr),
			Frequency:  DefaultPWMFrequency,
			ClockSpeed: ReferenceClockSpeed,
		},
	}
	if optn != nil {
		pca.Optn = optn
	}
	return pca
}

// Init initialize the PCA9685
func (pca *PCA9685) Init() (err error) {
	if pca.Conn.GetAddr() == 0 {
		return fmt.Errorf(`device %v is not initiated`, pca.Optn.Name)
	}
	return pca.SetFreq(pca.Optn.Frequency)
}

// SetFreq sets the PWM frequency in Hz for controller
func (pca *PCA9685) SetFreq(freq float32) (err error) {
	prescaleVal := ReferenceClockSpeed/StepCount/freq + 0.5
	if prescaleVal < 3.0 {
		return fmt.Errorf("PCA9685 cannot output at the given frequency")
	}
	oldMode, err := pca.Conn.ReadRegU8(Mode1)
	if err != nil {
		return err
	}
	newMode := (oldMode & 0x7F) | 0x10 // Mode 1, sleep
	if err := pca.Conn.WriteRegU8(Mode1, newMode); err != nil {
		return err
	}
	if err := pca.Conn.WriteRegU8(Prescale, byte(prescaleVal)); err != nil {
		return err
	}
	if err := pca.Conn.WriteRegU8(Mode1, oldMode); err != nil {
		return err
	}
	time.Sleep(5 * time.Millisecond)
	return pca.Conn.WriteRegU8(Mode1, oldMode|0xA1) // Mode 1, autoincrement on)
}

// GetFreq returns frequency value
func (pca *PCA9685) GetFreq() float32 {
	return pca.Optn.Frequency
}

// DeInit reset the chip
func (pca *PCA9685) DeInit() (err error) {
	return pca.Conn.WriteRegU8(Mode1, 0x00)
}

// SetChannel sets a single PWM channel
func (pca *PCA9685) SetChannel(chn, on, off int) (err error) {
	if chn < 0 || chn > 15 {
		return fmt.Errorf("invalid [channel] value")
	}
	if on < 0 || on > int(StepCount) {
		return fmt.Errorf("invalid [on] value")
	}
	if off < 0 || off > int(StepCount) {
		return fmt.Errorf("invalid [off] value")
	}

	buf := []byte{Led0On + byte(4*chn), byte(on) & 0xFF, byte(on >> 8), byte(off) & 0xFF, byte(off >> 8)}
	_, err = pca.Conn.WriteBytes(buf)
	return err
}
