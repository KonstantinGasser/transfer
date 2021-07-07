package main

import (
	"fmt"

	"github.com/KonstantinGasser/transfer/detect"
	"github.com/KonstantinGasser/transfer/shift"
)

func main() {
	// src, dst := ".test/src", ".test/dst"
	// files, _ := dirhash.DirFiles(".test/src", "")
	// fmt.Println(files)

	src := "USB-Sticks Anleger-, Neukundeninformation"
	signal := make(chan string)
	go detect.Listen(signal)

	done := make(chan string)

	for i := 0; i < 2; i++ {
		sig := <-signal
		go shift.Copy(src, sig, done)
	}
	close(signal)
	// for sig := range signal {
	// 	fmt.Println("SIG: ", sig)
	// 	go shift.Copy(src, sig, done)
	// }

	for dst := range done {
		fmt.Printf("[copy] copied data to USB-Stick: %q\n", dst)
	}
	// updatable, err := diff.Compare(src, dst).Updatable()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("updatable: ", updatable)
}
