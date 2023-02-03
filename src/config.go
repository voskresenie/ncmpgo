package main

type MetadataFormat string

const (
	NowPlayingLine1Fmt MetadataFormat = "{$5%a$9 $1|$9 %t}|{$5%a$9}|{%t}"
	NowPlayingLine2Fmt MetadataFormat = "{$2%b$9 $1|$9 $6(%y)$9}|{$2%b$9}|{$6(%y)$9}"
)
