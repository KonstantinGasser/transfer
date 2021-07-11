package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	usbdrivedetector "github.com/deepakjois/gousbdrivedetector"
	"github.com/otiai10/copy"
)

const (
	defaultPath = "USB-Sticks Anleger-, Neukundeninformation"
)

func main() {
	concurrent := flag.Int64("concurrent", 2, "concurrent devices contecet to the device")
	flag.Parse()

	fmt.Printf("[set-up] using %q as source path\n", defaultPath)
	fmt.Printf("[set-up] processing %d USB-Sticks concurrently\n", *concurrent)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	for {
		select {
		case sig := <-sigs:
			fmt.Printf("[main] STOPPING process due to os.Signal: %v\n", sig)
			os.Exit(0)
		default:
			devices, err := usbdrivedetector.Detect()
			if err != nil {
				panic(err)
			}
			if len(devices) == int(*concurrent) {
				fmt.Printf("[main] device(s) detected - starting processing...\n")
				for _, device := range devices {
					path := filepath.Join(device, defaultPath)
					fmt.Printf("[copy] coping data from %q to %q\n", defaultPath, path)
					if err := copy.Copy(defaultPath, path); err != nil {
						fmt.Printf("[copy] ERROR: could not copy data to device: %v\n", err)
						continue
					}
					fmt.Printf("[copy] SUCCESS: copied data to device: %q\n", device)
				}
				wait(sigs)
			}
			time.Sleep(time.Second * 3)
		}
	}
}

func wait(sigs chan os.Signal) {
	fmt.Println("[wait] unplugg all connected USB devices...")
	t := time.NewTicker(time.Second * 3)
	for {
		select {
		case sig := <-sigs:
			fmt.Printf("[main] process interrupted - terminating program with SIG: %v\n", sig)
			os.Exit(0)
		default:
			<-t.C
			ds, _ := usbdrivedetector.Detect()
			fmt.Println("[wait] devices: ", ds)
			if len(ds) == 0 {
				fmt.Println("[wait] all devices disconnected - starting next iteration")
				return
			}
		}
	}
}
