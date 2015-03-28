package main

import "crypto/rand"

var (
	data64B  = make([]byte, 64)
	data256B = make([]byte, 256)
	data1KB  = make([]byte, 1024)
)

func init() {
	rand.Read(data64B)
	rand.Read(data256B)
	rand.Read(data1KB)
}
