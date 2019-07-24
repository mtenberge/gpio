// Copyright © 2017 Kent Gibson <warthog618@gmail.com>.
// Copyright © 2019 Matthijs ten Berge <m.tenberge@awl.nl>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Dummy implementation for 'mem' to make a build possible on other platforms.
// Note that such a build will not provide any meaningful functionality.

// +build !linux !arm

package gpio

import (
	"errors"
	"sync"
)

// Arrays for 8 / 32 bit access to memory and a semaphore for write locking
var (
	// The memlock covers read/modify/write access to the mem block.
	// Individual reads and writes can skip the lock on the assumption that
	// concurrent register writes are atomic. e.g. Read, Write and Mode.
	memlock sync.Mutex
	mem     []uint32
	mem8    []uint8
)

// Open dummy for other platforms
func Open() (err error) {
	return nil
}

// Close dummy for other platforms
func Close() error {
	return nil
}

var (
	// ErrAlreadyOpen indicates the mem is already open.
	ErrAlreadyOpen = errors.New("already open")
)
