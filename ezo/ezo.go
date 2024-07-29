// Copyright 2024 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package ezo

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"periph.io/x/conn/v3/i2c"
	_ "periph.io/x/conn/v3/physic"
)

// EzoPhI2CAddr is the default I2C address for the EZO pH circuit.
const EzoPhI2CAddr uint16 = 0x63

// EzoOrpI2CAddr is the default I2C address for the EZO ORP circuit.
const EzoOrpI2CAddr uint16 = 0x62

// Maximum length of i2c response from the EZO circuit.
const maxRespLen = 40

// Dev is a handle to an EZO sensor.
type Dev struct {
	i2cDev  i2c.Dev
	name    string
	ezoType string
	mu      sync.Mutex
}

func NewEzo(bus i2c.Bus, addr uint16) (*Dev, error) {
	return &Dev{
		i2cDev: i2c.Dev{Bus: bus, Addr: addr},
		name:   "EZO",
	}, nil
}

type DeviceInfo struct {
	Device          string
	FirmwareVersion string
}

func (d *Dev) DeviceInfo() (DeviceInfo, error) {

	return DeviceInfo{}, nil
}

func (d *Dev) Command(cmd string, delay time.Duration) (string, error) {
	// Ensure exclusive use of sensor while request outstanding.
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.i2cDev.Tx([]byte(cmd), nil); err != nil {
		return "", d.wrap("Command: command i2c tx failed", err)
	}

	time.Sleep(delay)

	resp := make([]byte, maxRespLen+1)

	if err := d.i2cDev.Tx(nil, resp); err != nil {
		return "", d.wrap("Command: read-response i2c tx failed", err)
	}

	switch resp[0] {
	case 1:
		// success
	case 2:
		return "", d.errorf("Command: error status 2 (syntax error)")
	case 254:
		return "", d.errorf("Command: error status 254 (response not ready)")
	case 255:
		return "", d.errorf("Command: error status 255 (no data to send)")
	default:
		return "", d.errorf("Command: error status %d (unknown error)", resp[0])
	}

	eos := bytes.IndexByte(resp[1:], 0)
	var res string
	if eos < 0 {
		res = string(resp[1:])
	} else {
		res = string(resp[1 : eos+1])
	}
	return res, nil
}

func (d *Dev) errorf(msg string, a ...any) error {
	return fmt.Errorf(d.name+": "+msg, a...)
}

func (d *Dev) wrap(msg string, err error) error {
	return fmt.Errorf("%s: %s: %w", d.name, msg, err)
}
