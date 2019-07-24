// Copyright © 2017 Kent Gibson <warthog618@gmail.com>.
// Copyright © 2019 Matthijs ten Berge <m.tenberge@awl.nl>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Dummy implementation for 'interrupt' to make a build possible on other platforms.
// Note that such a build will not provide any meaningful functionality.

// +build !linux !arm

package gpio

import (
	"errors"
	"os"
	"sync"
)

const (
	// MaxGPIOInterrupt is the maximum pin number.
	MaxGPIOInterrupt = 54
)

// Edge represents the change in Pin level that triggers an interrupt.
type Edge string

const (
	// EdgeNone indicates no level transitions will trigger an interrupt
	EdgeNone Edge = "none"
	// EdgeRising indicates an interrupt is triggered when the pin transitions from low to high.
	EdgeRising Edge = "rising"
	// EdgeFalling indicates an interrupt is triggered when the pin transitions from high to low.
	EdgeFalling Edge = "falling"
	// EdgeBoth indicates an interrupt is triggered when the pin changes level.
	EdgeBoth Edge = "both"
)

type interrupt struct {
	pin       *Pin
	handler   func(*Pin)
	valueFile *os.File
}

// Watcher dummy
type Watcher struct {
	sync.Mutex // Guards the following, and sysfs interactions.
	Fd         int
	// Map from pin to value Fd.
	interruptFds map[uint8]int
	// Map from pin Fd to interrupt
	interrupts map[int]*interrupt
}

var defaultWatcher *Watcher

func getDefaultWatcher() *Watcher {
	return nil
}

// NewWatcher dummy
func NewWatcher() *Watcher {
	return nil
}

func closeInterrupts() {}

// Close dummy
func (watcher *Watcher) Close() {}

func waitExported(pin *Pin) error {
	return nil
}

func waitWriteable(path string) error {
	return nil
}

func export(pin *Pin) error {
	return nil
}

func unexport(pin *Pin) error {
	return nil
}

func openValue(pin *Pin) (*os.File, error) {
	return nil, nil
}

func setEdge(pin *Pin, edge Edge) error {
	return nil
}

// RegisterPin dummy
func (watcher *Watcher) RegisterPin(pin *Pin, edge Edge, handler func(*Pin)) (err error) {
	return nil
}

// UnregisterPin dummy
func (watcher *Watcher) UnregisterPin(pin *Pin) {}

// Watch dummy
func (pin *Pin) Watch(edge Edge, handler func(*Pin)) error {
	return nil
}

// Unwatch dummy
func (pin *Pin) Unwatch() {}

var (
	// ErrTimeout indicates the operation could not be performed within the
	// expected time.
	ErrTimeout = errors.New("timeout")
	// ErrBusy indicates the operation is already active on the pin.
	ErrBusy = errors.New("pin already in use")
)
