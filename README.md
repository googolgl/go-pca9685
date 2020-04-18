PCA9685 16-Channel 12-Bit PWM Driver
============================================================

[![GoDoc](https://godoc.org/github.com/googolgl/go-pca9685?status.svg)](https://godoc.org/github.com/googolgl/go-pca9685)
[![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

PCA9685 ([pdf reference](https://raw.github.com/googolgl/go-pca9685/master/docs/pca9685.pdf)) is a popular controller among Arduino and Raspberry PI developers.
The 16-Channel 12-bit PWM/Servo Driver will drive up to 16 servos over I2C with only 2 pins.  The on-board PWM controller will drive all 16 channels.  What's more, you can chain up to 62 of them to control up to 992 servos - all with the same 2 pins!
![image](https://raw.github.com/googolgl/go-pca9685/master/docs/pca9685.jpg)

Here is a library written in [Go programming language](https://golang.org/) for Raspberry PI and counterparts.

Golang usage
------------


```go
package main

import (
	"log"
	"time"

	i2c "github.com/d2r2/go-i2c"
	"github.com/googolgl/go-pca9685"
)

func main() {
    // Create new connection to i2c-bus on 1 line with address 0x40.
    // Use i2cdetect utility to find device address over the i2c-bus
    i2c, err := i2c.NewI2C(pca9685.Address, 1)
    if err != nil {
        log.Fatal(err)
    }

    pca0 := pca9685.PCANew(i2c, nil)
    err = pca0.Init()
    if err != nil {
    log.Fatal(err)
    }

    // Sets frequency for channel 0
    pca0.SetChannel(0, 0, 130)
    time.Sleep(1 * time.Second)

    // Angle in degrees. Must be in the range `0` to `Range`
    // Rotates from 0 to 130 degrees
    servo1 := ServoNew(pca0, 0, nil)
    for i := 0; i < 130; i++ {
        servo1.Angle(i)
        time.Sleep(10 * time.Millisecond)
    }

    // Fraction as pulse width expressed between 0.0 `MinPulse` and 1.0 `MaxPulse`
    servo1.Fraction(0.5)

    pca0.DeInit()
}
```


Getting help
------------

GoDoc [documentation](http://godoc.org/github.com/googolgl/go-pca9685)

Installation
------------

```bash
$ go get -u github.com/googolgl/go-pca9685
```

Troubleshooting
--------------

- *How to obtain fresh Golang installation to RPi device (either any RPi clone):*
If your RaspberryPI golang installation taken by default from repository is outdated, you may consider
to install actual golang manually from official Golang [site](https://golang.org/dl/). Download
tar.gz file containing arm64 in the name. Follow installation instructions.

- *How to enable I2C bus on RPi device:*
If you employ RaspberryPI, use raspi-config utility to activate i2c-bus on the OS level.
Go to "Interfacing Options" menu, to active I2C bus.
Probably you will need to reboot to load i2c kernel module.
Finally you should have device like /dev/i2c-1 present in the system.

- *How to find I2C bus allocation and device address:*
Use i2cdetect utility in format "i2cdetect -y X", where X may vary from 0 to 5 or more,
to discover address occupied by peripheral device. To install utility you should run
`apt install i2c-tools` on debian-kind system. `i2cdetect -y 1` sample output:
	```
	     0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
	00:          -- -- -- -- -- -- -- -- -- -- -- -- --
	10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	40: 40 -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	60: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	70: -- -- -- -- -- -- 76 --    
	```

Contact
-------

Please use [Github issue tracker](https://github.com/googolgl/go-pca9685/issues) for filing bugs or feature requests.


License
-------

Go-pca9685 is licensed under MIT License.
