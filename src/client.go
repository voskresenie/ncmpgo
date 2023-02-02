package main

import (
	"fmt"
	"log"

	"github.com/fhs/gompd/v2/mpd"
)

type Client interface {
	Close()

	NowPlaying() Track
}

type client struct {
	mpc *mpd.Client
}

func (c *client) Close() {
	c.mpc.Close()
}

func (c *client) NowPlaying() Track {
	attrs, err := c.mpc.CurrentSong()
	if err != nil {
		log.Fatalln(err)
	}
	return NewTrack(attrs)
}

func NewClient(host, port, protocol string) Client {
	mpc, err := mpd.Dial(protocol, fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalln(err)
	}
	return &client{
		mpc: mpc,
	}
}
