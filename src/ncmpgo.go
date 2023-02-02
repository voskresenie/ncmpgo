package main

import (
	"fmt"
)

func main() {
	c := NewClient("192.168.1.100", "6600", "tcp")
	defer c.Close()

	fmt.Printf("Now playing: %v", c.NowPlaying())
}
