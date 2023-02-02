package main

import (
	"fmt"

	"github.com/fhs/gompd/v2/mpd"
)

type Track interface {
	// TODO: consider removing bool second param
	Attr(string) (string, bool)
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

func (t *track) Attr(a string) (string, bool) {
	val, ok := t.attrs[a]
	return val, ok
}
