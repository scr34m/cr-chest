package main

// ServerHello
func cmd_20100(msg Packet) {
	sessionKey = msg.Buf[4:]
}
