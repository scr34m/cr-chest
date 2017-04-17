package main

import (
	"fmt"

	"golang.org/x/crypto/nacl/box"

	"github.com/augustoroman/hexdump"
)

// LoginOk
func cmd_20104(msg Packet) {
	nonce := NewNonce3(encrypt_nonce.nonce[:], publicKey[:], serverPublicKey[:])
	opened, ok := box.Open(nil, msg.Buf, &nonce.nonce, &serverPublicKey, privateKey)
	if !ok {
		fmt.Printf("failed to open box")
	} else {
		for i := 0; i < 24; i++ {
			decrypt_nonce.nonce[i] = opened[i]
		}

		for i := 0; i < 32; i++ {
			serverPublicKey[i] = opened[24+i]
		}

		fmt.Print("Data:\n")
		fmt.Print(hexdump.Dump(opened[56:]))
	}
}
