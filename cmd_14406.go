package main

import (
	"fmt"
	"net"

	"golang.org/x/crypto/nacl/secretbox"
)

// AskForBattleReplayStream
type MsgAskForBattleReplayStream struct {
	accountId uint64
}

func cmd_14406(conn net.Conn, m MsgAskForBattleReplayStream) {
	encrypt_nonce.increment()
	msg := Packet{
		ID: 14406,
	}
	msg.writeInt64(int64(m.accountId))
	msg.Buf = secretbox.Seal(nil, msg.Buf, &encrypt_nonce.nonce, &serverPublicKey)
	// dump(msg)
	_, err := conn.Write(msg.toByteArray())
	if err != nil {
		fmt.Print(err)
		return
	}
}
