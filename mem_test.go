// Copyright © 2017 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Test suite for mem module.

// +build linux,arm

package gpio

import (
	"testing"
)

func TestOpen(t *testing.T) {
	if err := Open(); err != nil {
		t.Fatal("Open returned error", err)
	}
	defer Close()
}

func TestOpenOpened(t *testing.T) {
	if err := Open(); err != nil {
		t.Fatal("Open returned error", err)
	}
	defer Close()
	if err := Open(); err == nil {
		t.Fatal("Open when opened didn't return error")
	}
}

func TestReOpen(t *testing.T) {
	if err := Open(); err != nil {
		t.Fatal("Open returned error", err)
	}
	Close()
	if err := Open(); err != nil {
		t.Fatal("Open returned error", err)
	}
	defer Close()
}
