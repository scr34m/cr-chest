package main

import (
	"fmt"
	"math/rand"
	"net"

	"golang.org/x/crypto/nacl/secretbox"
)

// ClientCapabilities
func cmd_10107(conn net.Conn) {
	encrypt_nonce.increment()

	msg := Packet{
		ID: 10107,
	}
	msg.writeRRSInt(rand.Intn(10) * rand.Intn(30))
	// string length: 4
	msg.writeByte(0x00)
	msg.writeByte(0x00)
	msg.writeByte(0x00)
	msg.writeByte(0x04)
	// string: eth1
	msg.writeByte(0x65)
	msg.writeByte(0x74)
	msg.writeByte(0x68)
	msg.writeByte(0x31)
	msg.Buf = secretbox.Seal(nil, msg.Buf, &encrypt_nonce.nonce, &serverPublicKey)
	// dump(msg)
	_, err := conn.Write(msg.toByteArray())
	if err != nil {
		fmt.Print(err)
		return
	}
}
