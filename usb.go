// usb - Self contained USB and HID library for Go
// Copyright 2019 The library Authors
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

// Package usb provide interfaces for generic USB devices.
package usb

// DeviceType represents the type of a USB device (generic or HID)
type DeviceType int

// List of supported device types
const (
	DeviceTypeGeneric DeviceType = 0
	DeviceTypeHID     DeviceType = 1
)

// Enumerate returns a list of all the HID devices attached to the system which
// match the vendor and product id:
//  - If the vendor id is set to 0 then any vendor matches.
//  - If the product id is set to 0 then any product matches.
//  - If the vendor and product id are both 0, all HID devices are returned.
// func Enumerate(vendorID uint16, productID uint16) []DeviceInfo {
// }

// DeviceInfo is a generic libusb info interface
type DeviceInfo interface {
	// Type returns the type of the device (generic or HID)
	Type() DeviceType

	// Platform-specific device path
	GetPath() string

	// IDs returns the vendor and product IDs for the device,
	// as well as the endpoint id and the usage page.
	IDs() (uint16, uint16, int, uint16)

	// Open tries to open the USB device represented by the current DeviceInfo
	Open() (Device, error)
}

// Device is a generic libusb device interface
type Device interface {
	Close() error

	Write(b []byte) (int, error)

	Read(b []byte) (int, error)

	// Type returns the type of the device (generic or HID)
	Type() DeviceType
}
