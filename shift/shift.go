package shift

import (
	"fmt"
)

func Copy(src, dst string, done chan string) {
	defer func() {
		done <- dst
	}()
	fmt.Printf("copy data from: %q to: %q\n", src, dst)
	// err := copy.Copy(src, dst)
	// if err != nil {
	// 	fmt.Printf("[copy] could not copy data from %q to %q\n", src, dst)
	// 	return
	// }
}
