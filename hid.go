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

import "errors"

// ErrDeviceClosed is returned for operations where the device closed before or
// during the execution.
var ErrDeviceClosed = errors.New("hid: device closed")

// ErrUnsupportedPlatform is returned for all operations where the underlying
// operating system is not supported by the library.
var ErrUnsupportedPlatform = errors.New("hid: unsupported platform")

// HidDeviceInfo is a hidapi info structure.
type HidDeviceInfo struct {
	Path         string // Platform-specific device path
	VendorID     uint16 // Device Vendor ID
	ProductID    uint16 // Device Product ID
	Release      uint16 // Device Release Number in binary-coded decimal, also known as Device Version Number
	Serial       string // Serial Number
	Manufacturer string // Manufacturer String
	Product      string // Product string
	UsagePage    uint16 // Usage Page for this Device/Interface (Windows/Mac only)
	Usage        uint16 // Usage for this Device/Interface (Windows/Mac only)

	// The USB interface which this logical device
	// represents. Valid on both Linux implementations
	// in all cases, and valid on the Windows implementation
	// only if the device contains more than one interface.
	Interface int
}

// GetPath returns the system-dependent path to the device
func (hdi *HidDeviceInfo) GetPath() string {
	return hdi.Path
}

// IDs returns the vendor and product id of the device
func (hdi *HidDeviceInfo) IDs() (uint16, uint16, int, uint16) {
	return hdi.VendorID, hdi.ProductID, hdi.Interface, hdi.UsagePage
}

// Type returns the type of the device (HID or generic)
func (hdi *HidDeviceInfo) Type() DeviceType {
	return DeviceTypeHID
}
