package main

import (
	"fmt"
	"net"

	"golang.org/x/crypto/nacl/secretbox"
)

func cmd_16103(conn net.Conn) {
	encrypt_nonce.increment()
	data := []byte{0x00, 0x00, 0x00}
	tosend := secretbox.Seal(nil, data, &encrypt_nonce.nonce, &serverPublicKey)
	msg := Packet{
		ID:  16103,
		Buf: tosend,
	}
	//dump(msg)
	_, err := conn.Write(msg.toByteArray())
	if err != nil {
		fmt.Print(err)
		return
	}
}
