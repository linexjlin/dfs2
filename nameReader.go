package main

import "io"

type NameReader struct {
	Name  string
	Rchan chan io.Reader
}
