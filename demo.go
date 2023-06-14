package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

// RAWINPUTDEVICELIST structure
type rawInputDeviceList struct {
	DeviceHandle uintptr
	Type         uint32
}

var (
	user32                    = syscall.NewLazyDLL("user32.dll")
	getRawInputDeviceListProc = user32.NewProc("GetRawInputDeviceList")
)

func main1() {
	dl := rawInputDeviceList{}
	size := uint32(unsafe.Sizeof(dl))

	// First I determine how many input devices are on the system, which
	// gets assigned to `devCount`
	var devCount uint32
	_ = getRawInputDeviceList(nil, &devCount, size)

	if devCount > 0 {
		devices := make([]rawInputDeviceList, size*devCount) // <- This is definitely wrong

		for i := 0; i < int(devCount); i++ {
			devices[i] = rawInputDeviceList{}
		}

		// Here is where I get the "The parameter is incorrect." error:
		err := getRawInputDeviceList(&devices[0], &devCount, size)
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
		for i := 0; i < int(devCount); i++ {
			fmt.Printf("Type: %v", devices[i].Type)
		}

	}
}

// Enumerates the raw input devices attached to the system.
func getRawInputDeviceList(rawInputDeviceList *rawInputDeviceList, numDevices *uint32, size uint32) error {
	_, _, err := getRawInputDeviceListProc.Call(
		uintptr(unsafe.Pointer(rawInputDeviceList)),
		uintptr(unsafe.Pointer(numDevices)),
		uintptr(size))
	if err != syscall.Errno(0) {
		return err
	}

	return nil
}
