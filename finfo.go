package main

import (
	"bytes"
	"encoding/gob"
	"io"
	"time"
)

type Finfo struct {
	IsDir   bool
	ModTime time.Time
	Name    string
	Size    int64
}

func (i *Finfo) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	en := gob.NewEncoder(w)
	err := en.Encode(i.IsDir)
	if err != nil {
		return nil, err
	}
	err = en.Encode(i.ModTime)
	if err != nil {
		return nil, err
	}
	err = en.Encode(i.Name)
	if err != nil {
		return nil, err
	}
	err = en.Encode(i.Size)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (i *Finfo) GobEncodeToWriter(w io.Writer) error {
	en := gob.NewEncoder(w)
	err := en.Encode(i.IsDir)
	if err != nil {
		return err
	}
	err = en.Encode(i.ModTime)
	if err != nil {
		return err
	}
	err = en.Encode(i.Name)
	if err != nil {
		return err
	}
	err = en.Encode(i.Size)
	if err != nil {
		return err
	}
	return nil
}

func (i *Finfo) GobDecode(buf []byte) error {
	r := bytes.NewReader(buf)
	de := gob.NewDecoder(r)
	err := de.Decode(&i.IsDir)
	if err != nil {
		return err
	}
	err = de.Decode(&i.ModTime)
	if err != nil {
		return err
	}
	err = de.Decode(&i.Name)
	if err != nil {
		return err
	}
	err = de.Decode(&i.Size)
	if err != nil {
		return err
	}
	return nil
}

func (i *Finfo) GobDecodeFromReader(r io.Reader) error {
	de := gob.NewDecoder(r)
	err := de.Decode(&i.IsDir)
	if err != nil {
		return err
	}
	err = de.Decode(&i.ModTime)
	if err != nil {
		return err
	}
	err = de.Decode(&i.Name)
	if err != nil {
		return err
	}
	err = de.Decode(&i.Size)
	if err != nil {
		return err
	}
	return nil
}
