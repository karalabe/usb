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

// +build freebsd,cgo linux,cgo darwin,!ios,cgo windows,cgo

package usb

// #include "./libusb/libusb/libusb.h"
import "C"

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

// enumerateRaw returns a list of all the USB devices attached to the system which
// match the vendor and product id:
//  - If the vendor id is set to 0 then any vendor matches.
//  - If the product id is set to 0 then any product matches.
//  - If the vendor and product id are both 0, all USB devices are returned.
func enumerateRaw(vendorID uint16, productID uint16, skipHid bool) ([]DeviceInfo, error) {
	// Create a context to interact with USB devices through
	var ctx *C.libusb_context
	errCode := int(C.libusb_init((**C.libusb_context)(&ctx)))
	if errCode < 0 {
		return nil, fmt.Errorf("Error while initializing libusb: %d", errCode)
	}
	// Retrieve all the available USB devices and wrap them in Go
	var deviceList **C.libusb_device
	count := C.libusb_get_device_list(ctx, &deviceList)
	if count < 0 {
		return nil, rawError(count)
	}
	defer C.libusb_free_device_list(deviceList, 1)

	var devices []*C.libusb_device
	*(*reflect.SliceHeader)(unsafe.Pointer(&devices)) = reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(deviceList)),
		Len:  int(count),
		Cap:  int(count),
	}
	//
	var infos []DeviceInfo
	for devnum, dev := range devices {
		// Retrieve the libusb device descriptor and skip non-queried ones
		var desc C.struct_libusb_device_descriptor
		if err := fromRawErrno(C.libusb_get_device_descriptor(dev, &desc)); err != nil {
			return nil, fmt.Errorf("failed to get device %d descriptor: %v", devnum, err)
		}
		if (vendorID > 0 && uint16(desc.idVendor) != vendorID) || (productID > 0 && uint16(desc.idProduct) != productID) {
			continue
		}
		// Skip HID devices if requested, they will be handled later
		if skipHid && desc.bDeviceClass == C.LIBUSB_CLASS_HID {
			continue
		}
		// Iterate over all the configurations and find raw interfaces
		for cfgnum := 0; cfgnum < int(desc.bNumConfigurations); cfgnum++ {
			// Retrieve the all the possible USB configurations of the device
			var cfg *C.struct_libusb_config_descriptor
			if err := fromRawErrno(C.libusb_get_config_descriptor(dev, C.uint8_t(cfgnum), &cfg)); err != nil {
				return nil, fmt.Errorf("failed to get device %d config %d: %v", devnum, cfgnum, err)
			}
			var ifaces []C.struct_libusb_interface
			*(*reflect.SliceHeader)(unsafe.Pointer(&ifaces)) = reflect.SliceHeader{
				Data: uintptr(unsafe.Pointer(cfg._interface)),
				Len:  int(cfg.bNumInterfaces),
				Cap:  int(cfg.bNumInterfaces),
			}
			// Drill down into each advertised interface
			for ifacenum, iface := range ifaces {
				if iface.num_altsetting == 0 {
					continue
				}
				var alts []C.struct_libusb_interface_descriptor
				*(*reflect.SliceHeader)(unsafe.Pointer(&alts)) = reflect.SliceHeader{
					Data: uintptr(unsafe.Pointer(iface.altsetting)),
					Len:  int(iface.num_altsetting),
					Cap:  int(iface.num_altsetting),
				}
				for _, alt := range alts {
					// Find the endpoints that can speak libusb interrupts
					var ends []C.struct_libusb_endpoint_descriptor
					*(*reflect.SliceHeader)(unsafe.Pointer(&ends)) = reflect.SliceHeader{
						Data: uintptr(unsafe.Pointer(alt.endpoint)),
						Len:  int(alt.bNumEndpoints),
						Cap:  int(alt.bNumEndpoints),
					}
					var reader, writer *uint8
					for _, end := range ends {
						switch {
						case end.bEndpointAddress&C.LIBUSB_ENDPOINT_OUT == C.LIBUSB_ENDPOINT_OUT && end.bmAttributes == C.LIBUSB_TRANSFER_TYPE_INTERRUPT:
							writer = new(uint8)
							*writer = uint8(end.bEndpointAddress)
						case end.bEndpointAddress&C.LIBUSB_ENDPOINT_IN == C.LIBUSB_ENDPOINT_IN && end.bmAttributes == C.LIBUSB_TRANSFER_TYPE_INTERRUPT:
							reader = new(uint8)
							*reader = uint8(end.bEndpointAddress)
						}
					}
					// If both in and out interrupts are available, match the device
					if reader != nil && writer != nil {
						infos = append(infos, DeviceInfo{
							Path:      fmt.Sprintf("%x:%x:%d", vendorID, uint16(desc.idProduct), uint8(C.libusb_get_port_number(dev))),
							VendorID:  uint16(desc.idVendor),
							ProductID: uint16(desc.idProduct),
							Interface: ifacenum,
							rawDevice: dev,
							rawReader: reader,
							rawWriter: writer,
						})
					}
				}
			}
		}
	}
	return infos, nil
}

// openRaw connects to a low level libusb device by its path name.
func openRaw(info DeviceInfo) (*RawDevice, error) {
	var handle *C.struct_libusb_device_handle
	if err := fromRawErrno(C.libusb_open(info.rawDevice.(*C.libusb_device), (**C.struct_libusb_device_handle)(&handle))); err != nil {
		return nil, fmt.Errorf("failed to open device: %v", err)
	}
	return &RawDevice{
		DeviceInfo: info,
		handle:     handle,
	}, nil
}

// RawDevice is a live low level USB connected device handle.
type RawDevice struct {
	DeviceInfo // Embed the infos for easier access

	handle *C.struct_libusb_device_handle // Low level USB device to communicate through
	lock   sync.Mutex
}

// Close releases the raw USB device handle.
func (dev *RawDevice) Close() error {
	dev.lock.Lock()
	defer dev.lock.Unlock()

	if dev.handle != nil {
		C.libusb_close(dev.handle)
		dev.handle = nil
	}
	return nil
}

// Write sends a binary blob to a low level USB device.
func (dev *RawDevice) Write(b []byte) (int, error) {
	dev.lock.Lock()
	defer dev.lock.Unlock()

	var transferred C.int
	if err := fromRawErrno(C.libusb_interrupt_transfer(dev.handle, (C.uchar)(*dev.rawWriter), (*C.uchar)(&b[0]), (C.int)(len(b)), &transferred, (C.uint)(0))); err != nil {
		return 0, err
	}
	return int(transferred), nil
}

// Read retrieves a binary blob from a low level USB device.
func (dev *RawDevice) Read(b []byte) (int, error) {
	dev.lock.Lock()
	defer dev.lock.Unlock()

	var transferred C.int
	if err := fromRawErrno(C.libusb_interrupt_transfer(dev.handle, (C.uchar)(*dev.rawReader), (*C.uchar)(&b[0]), (C.int)(len(b)), &transferred, (C.uint)(0))); err != nil {
		return 0, err
	}
	return int(transferred), nil
}
