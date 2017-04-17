package main

import (
	"crypto/rand"
	"io"
	"unsafe"

	"github.com/scr34m/blake2b"
)

type Nonce struct {
	//	hash  hash.Hash
	nonce [24]byte
}

func NewNonce0() *Nonce {
	p := new(Nonce)
	io.ReadFull(rand.Reader, p.nonce[:])
	return p
}

func NewNonce1(b [24]byte) *Nonce {
	p := new(Nonce)
	for i := 0; i < 24; i++ {
		p.nonce[i] = b[i]
	}
	return p
}

func NewNonce2(client []byte, server []byte) *Nonce {
	p := new(Nonce)
	p.nonce = blake2b.Sum192(append(client[:], server[:]...))
	return p
}

func NewNonce3(nonce []byte, client []byte, server []byte) *Nonce {
	k := append(nonce[:], client[:]...)
	k = append(k, server[:]...)

	p := new(Nonce)
	p.nonce = blake2b.Sum192(k)
	return p
}

func increment(bytes []byte, value int) {
	for offset, carry := 0, false; carry == true || offset == 0; offset += 8 {
		ptr := (*int)(unsafe.Pointer((&bytes[offset])))
		old := *ptr
		*ptr += value
		if old > *ptr {
			// overflow, carry and continue
			value = 1
			carry = true
		} else {
			carry = false
		}
	}
}

func (nonce *Nonce) increment() [24]byte {
	increment(nonce.nonce[:], 2)
	return nonce.nonce
}
