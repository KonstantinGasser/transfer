package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	usbdrivedetector "github.com/deepakjois/gousbdrivedetector"
	"github.com/otiai10/copy"
)

const (
	defaultPath = "Zeichnungsdokumente"
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
				var wg sync.WaitGroup

				for _, device := range devices {
					path := filepath.Join(device, defaultPath)
					fmt.Printf("[copy] coping data from %q to %q\n", defaultPath, path)
					wg.Add(1)
					go func(p string, d string) {
						defer wg.Done()

						start := time.Now()
						if err := copy.Copy(defaultPath, p); err != nil {
							fmt.Printf("[copy] ERROR: could not copy data to device: %v\n", err)
							return
						}
						fmt.Printf("[copy] SUCCESS: copied data to device: %q, (in %v seconds)\n", d, time.Since(start).Seconds())
					}(path, device)

				}
				wg.Wait()
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
