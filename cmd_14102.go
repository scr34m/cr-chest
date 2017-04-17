package main

import (
	"fmt"
	"io/ioutil"
	"net"

	"golang.org/x/crypto/nacl/secretbox"
)

// EndClientTurn
func cmd_14102(conn net.Conn, msg Packet) {
	encrypt_nonce.increment()
	msg.ID = 14102
	dump(msg)
	msg.Buf = secretbox.Seal(nil, msg.Buf, &encrypt_nonce.nonce, &serverPublicKey)
	_, err := conn.Write(msg.toByteArray())
	if err != nil {
		fmt.Print(err)
		return
	}
}

/*
./find.sh "*-14102.bin" b607
502 (b607) - Start Unlocking -> 1032 1283 1467 2514
509 - Collect Free Chest -> 2654
503 - Start Reward Claim -> 4041
511 - Collect Crown Chest
526 - Collected Reward -> 2655, 2658, 2670, 2681, 2693, 2708, 2709, 2710, 4042, 4045, 4046, 4047, 4048, 4052
210 - Claim Reward -> 2657, 4044
*/
// collect opened chest order 503 -> 526 -> 210
func p14102() {
	b, _ := ioutil.ReadFile("cr-proxy/replay/4052-14102.bin")
	p := Packet{
		Buf: b,
		Pos: 7,
	}

	fmt.Printf("tick: %d\n", p.readRRSInt())
	fmt.Printf("checksum: %d\n", p.readRRSInt())
	commands := p.readRRSInt()
	fmt.Printf("commands: %d\n", commands)
	for i := 0; i < int(commands); i++ {
		id := p.readRRSInt()
		fmt.Printf("id: %d\n", id)
		if id == 502 {
			fmt.Printf("tickStart: %d\n", p.readRRSInt())
			fmt.Printf("tickEnd: %d\n", p.readRRSInt())
			fmt.Printf("accountId: %d\n", p.readRRSLong())
			fmt.Printf("chestId: %d\n", p.readRRSInt())
		} else if id == 509 {
			fmt.Printf("tickStart: %d\n", p.readRRSInt())
			fmt.Printf("tickEnd: %d\n", p.readRRSInt())
			fmt.Printf("accountId: %d\n", p.readRRSLong())
		} else if id == 503 {
			fmt.Printf("tickStart: %d\n", p.readRRSInt())
			fmt.Printf("tickEnd: %d\n", p.readRRSInt())
			fmt.Printf("accountId: %d\n", p.readRRSLong())
			fmt.Printf("chestId: %d\n", p.readRRSInt())
		} else {
			panic("ismeretlen id")
		}
	}
	fmt.Printf("Byte(4): %02x %02x %02x %02x\n", p.readByte(), p.readByte(), p.readByte(), p.readByte())
}
