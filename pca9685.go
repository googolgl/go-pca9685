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
	//Address default for controller
	Address byte = 0x40

	Mode1      byte = 0x00
	Mode2      byte = 0x01
	SubAdr1    byte = 0x02
	SubAdr2    byte = 0x03
	SubAdr3    byte = 0x04
	Prescale   byte = 0xFE
	Led0OnL    byte = 0x06
	Led0OnH    byte = 0x07
	Led0OffL   byte = 0x08
	Led0OffH   byte = 0x09
	AllLedOnL  byte = 0xFA
	AllLedOnH  byte = 0xFB
	AllLedOffL byte = 0xFC
	AllLedOffH byte = 0xFD

	// Bits:
	Restart byte = 0x80
	Sleep   byte = 0x10
	AllCall byte = 0x01
	Invert  byte = 0x10
	OutDrv  byte = 0x04

	DefaultFreq float32 = 1
	OSCFreq     float32 = 25000000 // 25MHz
	StepCount   float32 = 4096     // 12-bit
)

// PCA9685 is a Driver for the PCA9685 16-channel 12-bit PWM/Servo controller
type PCA9685 struct {
	conn *i2c.I2C
	name string
}

// New creates a new driver with specified i2c interface
func New(i2c *i2c.I2C, name string) *PCA9685 {
	return &PCA9685{
		conn: i2c,
		name: name,
	}
}

// Init initialize the PCA9685
func (pca *PCA9685) Init() (err error) {
	if pca.conn.GetAddr() == 0 {
		return fmt.Errorf(`device %v is not initiated`, pca.name)
	}
	if err := pca.SetAllPWM(0, 0); err != nil {
		return err
	}
	if err := pca.conn.WriteRegU8(Mode2, OutDrv); err != nil {
		return err
	}
	if err := pca.conn.WriteRegU8(Mode1, AllCall); err != nil {
		return err
	}
	time.Sleep(5 * time.Millisecond)

	mode1, err := pca.conn.ReadRegU8(Mode1)
	if err != nil {
		return err
	}
	mode1 = mode1 & ^Sleep
	if err := pca.conn.WriteRegU8(Mode1, mode1); err != nil {
		return err
	}
	time.Sleep(5 * time.Millisecond)

	return
}

// SetPWMFreq sets the PWM frequency in Hz
func (pca *PCA9685) SetPWMFreq(freq float32) (err error) {
	var prescaleval float32 = OSCFreq / StepCount / freq
	prescaleval -= DefaultFreq
	prescale := byte(prescaleval + 0.5)

	mode1, err := pca.conn.ReadRegU8(Mode1)
	if err != nil {
		return err
	}
	newMode := (mode1 & 0x7F) | 0x10
	if err := pca.conn.WriteRegU8(Mode1, newMode); err != nil {
		return err
	}
	if err := pca.conn.WriteRegU8(Prescale, prescale); err != nil {
		return err
	}
	if err := pca.conn.WriteRegU8(Mode1, mode1); err != nil {
		return err
	}
	time.Sleep(5 * time.Millisecond)
	if err := pca.conn.WriteRegU8(Mode1, mode1|0x80); err != nil {
		return err
	}

	return
}

// SetPWM sets a single PWM channel
func (pca *PCA9685) SetPWM(chn, on, off int) (err error) {
	if err := pca.conn.WriteRegU8(Led0OnL+byte(4*chn), byte(on)&0xFF); err != nil {
		return err
	}
	if err := pca.conn.WriteRegU8(Led0OnH+byte(4*chn), byte(on>>8)); err != nil {
		return err
	}
	if err := pca.conn.WriteRegU8(AllLedOnL+byte(4*chn), byte(off)&0xFF); err != nil {
		return err
	}
	if err := pca.conn.WriteRegU8(AllLedOnH+byte(4*chn), byte(off>>8)); err != nil {
		return err
	}
	return
}

// SetAllPWM sets all PWM channels
func (pca *PCA9685) SetAllPWM(on, off int) (err error) {
	if err := pca.conn.WriteRegU8(AllLedOnL, byte(on)&0xFF); err != nil {
		return err
	}
	if err := pca.conn.WriteRegU8(AllLedOnH, byte(on>>8)); err != nil {
		return err
	}
	if err := pca.conn.WriteRegU8(AllLedOffL, byte(off)&0xFF); err != nil {
		return err
	}
	if err := pca.conn.WriteRegU8(AllLedOffH, byte(off>>8)); err != nil {
		return err
	}
	return
}
