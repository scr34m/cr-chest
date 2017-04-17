package main

import (
	"fmt"
	"net"

	"golang.org/x/crypto/nacl/secretbox"
)

// AskForAvatarStream
func cmd_14405(conn net.Conn) {
	encrypt_nonce.increment()
	data := []byte{}
	tosend := secretbox.Seal(nil, data, &encrypt_nonce.nonce, &serverPublicKey)
	msg := Packet{
		ID:  14405,
		Buf: tosend,
	}
	// dump(msg)
	_, err := conn.Write(msg.toByteArray())
	if err != nil {
		fmt.Print(err)
		return
	}
}
