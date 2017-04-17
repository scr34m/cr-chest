package main

import (
	"fmt"
	"io/ioutil"
	"net"
)

type MsgClientHello struct {
	protocol     uint32 // 1
	keyVersion   uint32 // 10
	majorVersion uint32 // 3
	minorVersion uint32 // 0
	build        uint32 // 193
	contentHash  string // 18e66ca14c9086d3ff2b900783beeee8844c18d8"
	deviceType   uint32 // 2
	appStore     uint32 // 2
}

func cmd_10100(conn net.Conn, m MsgClientHello) {
	msg := Packet{
		ID: 10100,
	}
	msg.writeInt(int(m.protocol))
	msg.writeInt(int(m.keyVersion))
	msg.writeInt(int(m.majorVersion))
	msg.writeInt(int(m.minorVersion))
	msg.writeInt(int(m.build))
	msg.writeString(m.contentHash)
	msg.writeInt(int(m.deviceType))
	msg.writeInt(int(m.appStore))
	_, err := conn.Write(msg.toByteArray())
	if err != nil {
		fmt.Print(err)
		return
	}
}

/*
ls -al *-10100.bin
*/
func p10100_file() {
	b, _ := ioutil.ReadFile("cr-proxy/replay/4901-10100.bin")
	p := Packet{
		Buf: b,
		Pos: 7,
	}
	p10100(&p)
}

func p10100(p *Packet) {
	//dump(*p)
	m := MsgClientHello{
		protocol:     p.readInt(),
		keyVersion:   p.readInt(),
		majorVersion: p.readInt(),
		minorVersion: p.readInt(),
		build:        p.readInt(),
		contentHash:  p.readString(),
		deviceType:   p.readInt(),
		appStore:     p.readInt(),
	}
	fmt.Printf("%+v\n", m)
}
