package main

type Encoder interface {
	Encode(rec Record) ([]byte, error)
}
