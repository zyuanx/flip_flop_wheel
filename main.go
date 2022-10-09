package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"golang.org/x/sys/windows/registry"
)

const (
	BASE_PATH = `SYSTEM\CurrentControlSet\Enum\HID`
	MOUSE_RE  = "@msmouse.inf,%hid.mousedevice%;"
)

var re = regexp.MustCompile(MOUSE_RE)

func getMouseDevice() map[string]bool {
	HIDKey, err := registry.OpenKey(registry.LOCAL_MACHINE, BASE_PATH, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer HIDKey.Close()

	names, err := HIDKey.ReadSubKeyNames(-1)
	if err != nil {
		log.Println(err)
		return nil
	}

	mouseMap := make(map[string]bool)
	for _, name := range names {
		subKey, err := registry.OpenKey(HIDKey, name, registry.ENUMERATE_SUB_KEYS)
		if err != nil {
			log.Println(err)
			return nil
		}
		defer subKey.Close()
		children, err := subKey.ReadSubKeyNames(-1)
		if err != nil {
			log.Println(err)
			return nil
		}

		for _, child := range children {
			childKey, err := registry.OpenKey(subKey, child, registry.QUERY_VALUE)
			if err != nil {
				log.Println(err)
				return nil
			}
			defer childKey.Close()
			s, _, err := childKey.GetStringValue("DeviceDesc")
			if err != nil {
				log.Println(err)
				return nil
			}
			if re.MatchString(s) {
				mouseMap[name+`\`+child] = true
			}
		}

	}
	return mouseMap
}

func setMouseDevice(mouseMap map[string]bool) {
	for path := range mouseMap {
		a, _ := registry.OpenKey(registry.LOCAL_MACHINE, BASE_PATH+`\`+path+`\Device Parameters`, registry.QUERY_VALUE)
		r, _, err := registry.CreateKey(a, "", registry.SET_VALUE)
		if err != nil {
			log.Println("Reopen the program in superuser mode.")
			log.Println(err)
			return
		}
		defer r.Close()
		err = r.SetDWordValue("FlipFlopWheel", 1)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	log.Println("Finish")
}

func main() {
	mouseMap := getMouseDevice()
	log.Println("Found mouse device:", len(mouseMap))
	log.Println("Press input y or Y to set mouse device...")
	yes := ""
	fmt.Scanln(&yes)
	if yes == "y" || yes == "Y" {
		setMouseDevice(mouseMap)
	}
	log.Printf("Press any key to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
