package main

import (
	"fmt"

	"github.com/KonstantinGasser/transfer/detect"
)

func main() {
	// src, dst := ".test/src", ".test/dst"
	// files, _ := dirhash.DirFiles(".test/src", "")
	// fmt.Println(files)
	signal := make(chan string)
	go detect.Listen(signal)

	for sig := range signal {
		fmt.Println("SIG: ", sig)
	}
	// updatable, err := diff.Compare(src, dst).Updatable()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("updatable: ", updatable)
}
