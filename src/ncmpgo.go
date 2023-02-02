package main

import (
	"fmt"
)

func main() {
	c := NewFmtClient("192.168.1.100", "6600", "tcp")
	defer c.Close()

	lines := c.NowPlaying()
	fmt.Printf("Now playing:\n\t%s\n\t%s", lines[0], lines[1])
}
