// usb - Self contained USB and HID library for Go
// Copyright 2017 The library Authors
//
// This library is free software: you can redistribute it and/or modify it under
// the terms of the GNU Lesser General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option) any
// later version.
//
// The library is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
// A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License along
// with the library. If not, see <http://www.gnu.org/licenses/>.

package usb

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"testing"
)

// Tests that HID enumeration can be called concurrently from multiple threads.
func TestThreadedEnumerateHid(t *testing.T) {
	var (
		threads         = 8
		errs    []error = make([]error, threads)
		pend    sync.WaitGroup
	)
	for i := 0; i < threads; i++ {
		pend.Add(1)

		go func(index int) {
			defer pend.Done()
			for j := 0; j < 512; j++ {
				if _, err := EnumerateHid(uint16(index), 0); err != nil {
					errs[index] = fmt.Errorf("thread %d, iter %d: failed to enumerate-hid: %v", index, j, err)
					break
				}
			}
		}(i)
	}
	pend.Wait()
	for _, err := range errs {
		if err != nil {
			t.Error(err)
		}
	}
}

// Tests that RAW enumeration can be called concurrently from multiple threads.
func TestThreadedEnumerateRaw(t *testing.T) {
	// Travis does not have usbfs enabled in the Linux kernel
	if os.Getenv("TRAVIS") != "" && runtime.GOOS == "linux" {
		t.Skip("Linux on Travis doesn't have usbfs, skipping test")
	}
	// Yay, we can actually test this
	var (
		threads         = 8
		errs    []error = make([]error, threads)
		pend    sync.WaitGroup
	)
	for i := 0; i < threads; i++ {
		pend.Add(1)

		go func(index int) {
			defer pend.Done()
			for j := 0; j < 512; j++ {
				if _, err := EnumerateRaw(uint16(index), 0); err != nil {
					errs[index] = fmt.Errorf("thread %d, iter %d: failed to enumerate-raw: %v", index, j, err)
					break
				}
			}
		}(i)
	}
	pend.Wait()
	for _, err := range errs {
		if err != nil {
			t.Error(err)
		}
	}
}

// Tests that generic enumeration can be called concurrently from multiple threads.
func TestThreadedEnumerate(t *testing.T) {
	// Travis does not have usbfs enabled in the Linux kernel
	if os.Getenv("TRAVIS") != "" && runtime.GOOS == "linux" {
		t.Skip("Linux on Travis doesn't have usbfs, skipping test")
	}
	var (
		threads         = 8
		errs    []error = make([]error, threads)
		pend    sync.WaitGroup
	)
	for i := 0; i < threads; i++ {
		pend.Add(1)

		go func(index int) {
			defer pend.Done()
			for j := 0; j < 512; j++ {
				if _, err := Enumerate(uint16(index), 0); err != nil {
					errs[index] = fmt.Errorf("thread %d, iter %d: failed to enumerate: %v", index, j, err)
					break
				}
			}
		}(i)
	}
	pend.Wait()
	for _, err := range errs {
		if err != nil {
			t.Error(err)
		}
	}
}
