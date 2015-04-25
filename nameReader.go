package main

import "io"

type NameReader struct {
	name   string
	reader io.Reader
}
