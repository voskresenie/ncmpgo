package main

import (
	"fmt"
	"log"

	"github.com/fhs/gompd/v2/mpd"
)

type Client interface {
	Close()
}

type client struct {
	mpc *mpd.Client
}

func (c *client) Close() {
	c.mpc.Close()
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
