package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	RID_INPUT                  = 0x10000003
	RID_DEVICE_INFO_MOUSE_TYPE = 0x02
	RID_HEADER                 = 0x10000005
	RID_DEVICE_NAME_SIZE       = 128

	MOUSE_WHEEL_ROUTED_EVENT = 0x040E
)

type RID_DEVICE_INFO struct {
	ID   uintptr
	Type uint32
}

type RID_DEVICE_INFO_MOUSE struct {
	Size       uint32
	Type       uint32
	ID         uintptr
	Buttons    uint32
	DataLength uint32
	Data       uintptr
	Attributes uint32
}

var (
	user32DLL = syscall.NewLazyDLL("user32.dll")
	// getRawInputDeviceListProc = user32DLL.NewProc("GetRawInputDeviceList")
	getRawInputDeviceInfo = user32DLL.NewProc("GetRawInputDeviceInfoW")
)

func getMouseDevices() ([]RID_DEVICE_INFO_MOUSE, error) {
	var deviceCount uint32
	ret, _, _ := getRawInputDeviceListProc.Call(uintptr(unsafe.Pointer(nil)), uintptr(unsafe.Pointer(&deviceCount)), unsafe.Sizeof(RID_DEVICE_INFO{}))
	if ret == 0xFFFFFFFF {
		return nil, fmt.Errorf("GetRawInputDeviceList failed")
	}

	rawInputDeviceList := make([]RID_DEVICE_INFO, deviceCount)
	deviceSize := unsafe.Sizeof(RID_DEVICE_INFO{})
	ret, _, _ = getRawInputDeviceListProc.Call(uintptr(unsafe.Pointer(&rawInputDeviceList[0])), uintptr(unsafe.Pointer(&deviceCount)), deviceSize)
	if ret == 0xFFFFFFFF {
		return nil, fmt.Errorf("GetRawInputDeviceList failed")
	}

	var mouseDevices []RID_DEVICE_INFO_MOUSE

	for _, device := range rawInputDeviceList {
		if device.Type == RID_DEVICE_INFO_MOUSE_TYPE {
			var deviceName [RID_DEVICE_NAME_SIZE]uint16
			deviceNameSize := RID_DEVICE_NAME_SIZE * 2
			getRawInputDeviceInfo.Call(device.ID, RID_DEVICE_INFO_MOUSE_TYPE, uintptr(unsafe.Pointer(&deviceName[0])), uintptr(unsafe.Pointer(&deviceNameSize)))

			mouseDevices = append(mouseDevices, RID_DEVICE_INFO_MOUSE{
				ID:   device.ID,
				Size: device.Type,
			})
		}
	}

	return mouseDevices, nil
}

func main() {
	// main1()
	mouseDevices, err := getMouseDevices()
	if err != nil {
		fmt.Println("Failed to get mouse devices:", err)
		return
	}

	for _, device := range mouseDevices {
		fmt.Printf("Mouse Device ID: %d\n", device.ID)
	}
}
