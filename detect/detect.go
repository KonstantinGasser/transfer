package detect

import (
	"fmt"
	"time"

	usbdrivedetector "github.com/deepakjois/gousbdrivedetector"
)

func Listen(signal chan string) {
	var ticker = time.NewTicker(time.Second * 2)

	// remember known devices
	var known = make(map[string]bool)
	startDevices, err := Lookup()
	if err != nil {
		panic(fmt.Errorf("could not check on initial devices: %v", err))
	}
	for _, d := range startDevices {
		known[d] = true
	}
	fmt.Println("Init devices: ", known)
	// listen for new devices
	for {
		select {
		case <-signal: // return if chan closed
			return
		case <-ticker.C:
			devices, err := Lookup()
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, d := range devices {
				if found, ok := known[d]; !found && !ok {
					fmt.Printf("[detect] a new USB device was detected. Device: %q\n", d)
					known[d] = true
					signal <- d
				}
			}
		}
	}
}

func Lookup() ([]string, error) {
	return usbdrivedetector.Detect()
}
