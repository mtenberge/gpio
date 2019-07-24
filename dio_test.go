// Copyright © 2017 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

//
//  Test suite for dio module.
//
//	Tests use J8 pins 7 (mostly) and 15 and 16 (for looped tests)

// +build linux,arm

package gpio

import (
	"testing"
)

func setupDIO(t *testing.T) {
	if err := Open(); err != nil {
		t.Fatal("Open returned error", err)
	}
}

func teardownDIO() {
	Close()
}

func TestUninitialisedPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("NewPin did not panic")
		}
	}()
	p := NewPin(J8p7)
	_ = p
}

func TestClosedPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("NewPin did not panic")
		}
	}()
	setupDIO(t)
	teardownDIO()
	p := NewPin(J8p7)
	_ = p
}

func TestRead(t *testing.T) {
	setupDIO(t)
	defer teardownDIO()
	pin := NewPin(J8p7)
	// A basic read test - assuming the pin is input and pulled high
	// which is the default state for this pin on a Pi.
	mode := pin.Mode()
	if mode != Input {
		t.Error("Not an input pin")
	}
	val := pin.Read()
	// Assumes pin is initially pulled up and set as an input.
	if val != High {
		t.Error("Not pulled high")
	}
}

func TestMode(t *testing.T) {
	setupDIO(t)
	defer teardownDIO()
	pin := NewPin(J8p7)
	mode := pin.Mode()
	if mode != Input {
		t.Fatal("Not an input pin")
	}
	pin.SetMode(Output)
	mode = pin.Mode()
	if mode != Output {
		t.Error("Failed to set output")
	}
	pin.SetMode(Input)
	mode = pin.Mode()
	if mode != Input {
		t.Error("Failed to set input")
	}
	pin.Output()
	mode = pin.Mode()
	if mode != Output {
		t.Error("Failed to set output")
	}
	pin.Input()
	mode = pin.Mode()
	if mode != Input {
		t.Error("Failed to set input")
	}
}

func TestPull(t *testing.T) {
	setupDIO(t)
	defer teardownDIO()
	pin := NewPin(J8p7)
	defer pin.PullUp()
	// A basic read test - using the pull up/down to drive he pin.
	mode := pin.Mode()
	if mode != Input {
		t.Error("Not an input pin")
	}
	pin.PullUp()
	val := pin.Read()
	if val != High {
		t.Error("Not pulled up")
	}
	pin.PullDown()
	val = pin.Read()
	if val != Low {
		t.Error("Not pulled down")
	}
	pin.SetPull(PullUp)
	val = pin.Read()
	if val != High {
		t.Error("Not pulled up by SetPull")
	}
	// no real way of testing this, but to trick coverage...
	pin.PullNone()
}

func TestWrite(t *testing.T) {
	setupDIO(t)
	defer teardownDIO()
	pin := NewPin(J8p7)
	mode := pin.Mode()
	if mode != Input {
		t.Error("Not an input pin")
	}
	defer pin.SetMode(Input)
	pin.Write(Low)
	pin.SetMode(Output)
	mode = pin.Mode()
	if mode != Output {
		t.Fatal("Failed to set output")
	}
	if pin.Read() != Low {
		t.Error("Failed to init Low")
	}
	pin.Write(High)
	if pin.Shadow() != High {
		t.Error("Failed to shadow write High")
	}
	if pin.Read() != High {
		t.Error("Failed to write High")
	}
	pin.Write(Low)
	if pin.Shadow() != Low {
		t.Error("Failed to shadow write Low")
	}
	if pin.Read() != Low {
		t.Error("Failed to write Low")
	}
	pin.High()
	if pin.Shadow() != High {
		t.Error("Failed to shadow write High")
	}
	if pin.Read() != High {
		t.Error("Failed to write High")
	}
	pin.Low()
	if pin.Shadow() != Low {
		t.Error("Failed to shadow write Low")
	}
	if pin.Read() != Low {
		t.Error("Failed to write Low")
	}
}

// Looped tests require a jumper across Raspberry Pi J8 pins 15 and 16.
func TestWriteLooped(t *testing.T) {
	setupDIO(t)
	defer teardownDIO()
	pinIn := NewPin(J8p15)
	pinOut := NewPin(J8p16)
	pinIn.SetMode(Input)
	defer pinOut.SetMode(Input)
	pinOut.Write(Low)
	pinOut.SetMode(Output)
	if pinIn.Read() != Low {
		t.Error("Failed to init Low")
	}
	pinOut.Write(High)
	if pinIn.Read() != High {
		t.Error("Failed to write High")
	}
	pinOut.Write(Low)
	if pinIn.Read() != Low {
		t.Error("Failed to write Low")
	}
}

func TestToggle(t *testing.T) {
	setupDIO(t)
	defer teardownDIO()
	pin := NewPin(J8p7)
	defer pin.SetMode(Input)
	pin.Write(Low)
	pin.SetMode(Output)
	mode := pin.Mode()
	if mode != Output {
		t.Fatal("Failed to set output")
	}
	if pin.Read() != Low {
		t.Error("Failed to init Low")
	}
	pin.Toggle()
	if pin.Shadow() != High {
		t.Error("Failed to shadow toggle High")
	}
	if pin.Read() != High {
		t.Error("Failed to toggle High")
	}
	pin.Toggle()
	if pin.Shadow() != Low {
		t.Error("Failed to shadow toggle Low")
	}
	if pin.Read() != Low {
		t.Error("Failed to toggle Low")
	}
}

// Looped tests require a jumper across Raspberry Pi J8 pins 15 and 16.
func TestToggleLooped(t *testing.T) {
	setupDIO(t)
	defer teardownDIO()
	pinIn := NewPin(J8p15)
	pinOut := NewPin(J8p16)
	pinIn.SetMode(Input)
	defer pinOut.SetMode(Input)
	pinOut.Write(Low)
	pinOut.SetMode(Output)
	mode := pinOut.Mode()
	if mode != Output {
		t.Fatal("Failed to set output")
	}
	if pinIn.Read() != Low {
		t.Error("Failed to init Low")
	}
	pinOut.Toggle()
	if pinIn.Read() != High {
		t.Error("Failed to toggle High")
	}
	pinOut.Toggle()
	if pinIn.Read() != Low {
		t.Error("Failed to toggle Low")
	}
}

func BenchmarkRead(b *testing.B) {
	err := Open()
	if err != nil {
		b.Fatal("Open returned error", err)
	}
	defer Close()
	pin := NewPin(J8p7)
	for i := 0; i < b.N; i++ {
		_ = pin.Read()
	}
}

func BenchmarkWrite(b *testing.B) {
	err := Open()
	if err != nil {
		b.Fatal("Open returned error", err)
	}
	defer Close()
	pin := NewPin(J8p7)
	for i := 0; i < b.N; i++ {
		pin.Write(High)
	}
}

func BenchmarkToggle(b *testing.B) {
	err := Open()
	if err != nil {
		b.Fatal("Open returned error", err)
	}
	defer Close()
	pin := NewPin(J8p7)
	defer pin.SetMode(Input)
	pin.Write(Low)
	pin.SetMode(Output)
	for i := 0; i < b.N; i++ {
		pin.Toggle()
	}
}

func BenchmarkSysfsRead(b *testing.B) {
	err := Open()
	if err != nil {
		b.Fatal("Open returned error", err)
	}
	defer Close()
	pin := NewPin(J8p7)
	// setup sysfs
	err = export(pin)
	if err != nil {
		b.Fatal("Couldn't export pin.")
	}
	defer unexport(pin)
	f, err := openValue(pin)
	if err != nil {
		b.Fatal("Couldn't open value.")
	}
	defer f.Close()
	r := make([]byte, 1)
	r[0] = 0
	for i := 0; i < b.N; i++ {
		f.Read(r)
	}
}

func BenchmarkSysfsWrite(b *testing.B) {
	err := Open()
	if err != nil {
		b.Fatal("Open returned error", err)
	}
	defer Close()
	pin := NewPin(J8p7)
	// setup sysfs
	err = export(pin)
	if err != nil {
		b.Fatal("Couldn't export pin.")
	}
	defer unexport(pin)
	f, err := openValue(pin)
	if err != nil {
		b.Fatal("Couldn't open value.")
	}
	defer f.Close()
	r := "0"
	for i := 0; i < b.N; i++ {
		f.WriteString(r)
	}
}

func BenchmarkSysfsToggle(b *testing.B) {
	err := Open()
	if err != nil {
		b.Fatal("Open returned error", err)
	}
	defer Close()
	pin := NewPin(J8p7)
	// setup sysfs
	err = export(pin)
	if err != nil {
		b.Fatal("Couldn't export pin.")
	}
	defer unexport(pin)
	f, err := openValue(pin)
	if err != nil {
		b.Fatal("Couldn't open value.")
	}
	defer f.Close()
	r := "0"
	for i := 0; i < b.N; i++ {
		if r == "0" {
			r = "1"
		} else {
			r = "0"
		}
		f.WriteString(r)
	}
}
