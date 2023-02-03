package main

type FmtClient interface {
	Close()

	NowPlaying() []string
}

type fmtclient struct {
	client Client
}

func (c *fmtclient) Close() {
	c.client.Close()
}

var nowPlayingLine1Text, _ = NewText(NowPlayingLine1Fmt)
var nowPlayingLine2Text, _ = NewText(NowPlayingLine2Fmt)

func (c *fmtclient) NowPlaying() []string {
	t := c.client.NowPlaying()
	return []string{
		nowPlayingLine1Text.Format(t, false),
		nowPlayingLine2Text.Format(t, false),
	}
}

func NewFmtClient(host, port, protocol string) FmtClient {
	client := NewClient(host, port, protocol)
	return &fmtclient{
		client: client,
	}
}
