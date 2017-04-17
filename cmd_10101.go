package main

import (
	"fmt"
	"io/ioutil"
	"net"

	"golang.org/x/crypto/nacl/box"
)

type MsgLogin struct {
	accountId               uint64
	passToken               string
	clientMajorVersion      int32
	clientMinorVersion      int32
	clientBuild             int32
	resourceSha             string
	UDID                    string
	openUdid                string
	macAddress              uint32
	device                  string
	advertisingGuid         string
	osVersion               string
	isAndroid               byte
	unknown1                uint32
	androidID               string
	preferredDeviceLanguage string
	unknown2                byte
	preferredLanguage       byte
	facebookAttributionId   string
	advertisingEnabled      byte
	appleIFV                string
	appStore                byte
	kunlunSSO               string
	kunlunUID               string
	unknown3                string
	unknown4                string
	unknown5                byte
}

//Login
func cmd_10101(conn net.Conn, m MsgLogin) {
	nonce := NewNonce2(publicKey[:], serverPublicKey[:])

	msg := Packet{
		ID: 10101,
	}
	msg.writeInt64(int64(m.accountId))
	msg.writeString(m.passToken)
	msg.writeRRSInt(int(m.clientMajorVersion))
	msg.writeRRSInt(int(m.clientMinorVersion))
	msg.writeRRSInt(int(m.clientBuild))
	msg.writeString(m.resourceSha)
	msg.writeString(m.UDID)
	msg.writeString(m.openUdid)
	msg.writeInt(int(m.macAddress))
	msg.writeString(m.device)
	msg.writeString(m.advertisingGuid)
	msg.writeString(m.osVersion)
	msg.writeByte(m.isAndroid)
	msg.writeInt(int(m.unknown1))
	msg.writeString(m.androidID)
	msg.writeString(m.preferredDeviceLanguage)
	msg.writeByte(m.unknown2)
	msg.writeByte(m.preferredLanguage)
	msg.writeString(m.facebookAttributionId)
	msg.writeByte(m.advertisingEnabled)
	msg.writeString(m.appleIFV)
	msg.writeByte(m.appStore)
	msg.writeString(m.kunlunSSO)
	msg.writeString(m.kunlunUID)
	msg.writeString(m.unknown3)
	msg.writeString(m.unknown4)
	msg.writeByte(m.unknown5)

	message := append(sessionKey[:], encrypt_nonce.nonce[:]...)
	message = append(message, msg.Buf[:]...)
	ciphertext := box.Seal(nil, message, &nonce.nonce, &serverPublicKey, privateKey)
	msg.Buf = append(publicKey[:], ciphertext[:]...)
	// dump(msg)
	_, err := conn.Write(msg.toByteArray())
	if err != nil {
		fmt.Print(err)
		return
	}
}

// ls -al *-10101.bin
func p10101_file() {
	b, _ := ioutil.ReadFile("cr-proxy/replay/4903-10101.bin")
	p := Packet{
		Buf: b,
		Pos: 7,
	}
	p10101(&p)
}

func p10101(p *Packet) {
	dump(*p)
	m := MsgLogin{
		accountId:               p.readLong(),
		passToken:               p.readString(),
		clientMajorVersion:      p.readRRSInt(),
		clientMinorVersion:      p.readRRSInt(),
		clientBuild:             p.readRRSInt(),
		resourceSha:             p.readString(),
		UDID:                    p.readString(),
		openUdid:                p.readString(),
		macAddress:              p.readInt(),
		device:                  p.readString(),
		advertisingGuid:         p.readString(),
		osVersion:               p.readString(),
		isAndroid:               p.readByte(),
		unknown1:                p.readInt(),
		androidID:               p.readString(),
		preferredDeviceLanguage: p.readString(),
		unknown2:                p.readByte(),
		preferredLanguage:       p.readByte(),
		facebookAttributionId:   p.readString(),
		advertisingEnabled:      p.readByte(),
		appleIFV:                p.readString(),
		appStore:                p.readByte(),
		kunlunSSO:               p.readString(),
		kunlunUID:               p.readString(),
		unknown3:                p.readString(),
		unknown4:                p.readString(),
		unknown5:                p.readByte(),
	}
	fmt.Printf("%+v\n", m)
}
