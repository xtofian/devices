// Copyright 2024 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package ezo

import (
	"testing"
	"time"

	"periph.io/x/conn/v3/i2c/i2ctest"
)

func TestCommand_success(t *testing.T) {
	bus := i2ctest.Playback{
		Ops: []i2ctest.IO{
			{Addr: 0x63, W: []uint8{0x69}, R: []uint8(nil)},
			{Addr: 0x63, W: []uint8(nil),
				R: []uint8{0x1, 0x3f, 0x49, 0x2c, 0x70, 0x48, 0x2c, 0x32, 0x2e, 0x31, 0x34, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
		},
	}
	d, err := NewEzo(&bus, EzoPhI2CAddr)
	if err != nil {
		t.Fatal(err)
	}
	got, err := d.Command("i", 300*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	want := "?I,pH,2.14"
	if got != want {
		t.Errorf("Command(\"i\") = %v; want %v", got, want)
	}
}
