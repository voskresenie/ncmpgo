package main

import (
	"fmt"

	"github.com/fhs/gompd/v2/mpd"
)

type Track interface {
}

type track struct {
	attrs map[string]string
}

func NewTrack(a mpd.Attrs) Track {
	return &track{attrs: a}
}

func (t *track) String() string {
	return fmt.Sprintf("%s - %s", t.attrs["Artist"], t.attrs["Title"])
}
