package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"

	"github.com/augustoroman/hexdump"
)

var publicKey, privateKey *[32]byte
var serverPublicKey [32]byte
var sessionKey []byte
var encrypt_nonce *Nonce
var decrypt_nonce *Nonce
var conn net.Conn
var tickCounter = 0 // time in 1/60 seconds since login, sent every 10 seconds

func keepalive() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				cmd_10108(conn)
			}
		}
	}()
}

var homedata HomeData

func getFreeChest() int {
	for i, c := range homedata.chests {
		if c.id == -1 && c.removed == false {
			return i
		}
	}
	return -1
}

func getExpiredChest() int {
	for i, c := range homedata.chests {
		if c.status == 1 && c.removed == false {
			return i
		}
	}
	return -1
}

func getPendingChest() int {
	for i, c := range homedata.chests {
		if c.id != -2 && c.status == 8 {
			return i
		}
	}
	return -1
}

func getOneChest() int {
	for i, c := range homedata.chests {
		if c.id != -2 && c.status == 0 && c.removed == false {
			return i
		}
	}
	return -1
}

func ticker() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				tickCounter += 2
				p := Packet{}
				p.writeRRSInt(tickCounter * 100)
				p.writeByte(byte(tickCounter)) // checksum

				fmt.Printf("%+v\n", homedata.chests)
				if tickCounter >= 2 && len(homedata.chests) > 0 {
					if i := getFreeChest(); i != -1 {
						fmt.Printf("Claim a FREE chest %d\n", homedata.chests[i].id)
						homedata.chests[i].removed = true
						p.writeRRSInt(1)                       // command
						p.writeRRSInt(509)                     // id
						p.writeRRSInt((tickCounter * 100) - 1) // tickStart
						p.writeRRSInt((tickCounter * 100) - 1) // tickEnd
						p.writeRRSLong(accountId)              // accountID
					} else if i := getExpiredChest(); i != -1 {
						fmt.Printf("Claim a chest %d\n", homedata.chests[i].id)
						homedata.chests[i].removed = true
						p.writeRRSInt(1)                       // command
						p.writeRRSInt(503)                     // id
						p.writeRRSInt((tickCounter * 100) - 1) // tickStart
						p.writeRRSInt((tickCounter * 100) - 1) // tickEnd
						p.writeRRSLong(accountId)              // accountID
						p.writeRRSInt(homedata.chests[i].id)   // chest
					} else if i := getPendingChest(); i != -1 {
						fmt.Printf("Unlock in progress on chest %d\n", homedata.chests[i].id)
					} else if i := getOneChest(); i != -1 {
						fmt.Printf("Start unlock a chest %d\n", homedata.chests[i].id)
						homedata.chests[i].removed = true
						homedata.chests[i].status = 8
						p.writeRRSInt(1)                       // command
						p.writeRRSInt(502)                     // id
						p.writeRRSInt((tickCounter * 100) - 1) // tickStart
						p.writeRRSInt((tickCounter * 100) - 1) // tickEnd
						p.writeRRSLong(accountId)              // accountID
						p.writeRRSInt(homedata.chests[i].id)   // chest
					} else {
						p.writeByte(0) // command
					}
				} else {
					p.writeByte(0) // command
				}
				p.writeByte(0xff)
				p.writeByte(0xff)
				p.writeByte(0xff)
				p.writeByte(0xff)
				cmd_14102(conn, p)
			}
		}
	}()
}

func decrypt(msg Packet, dump bool) {
	msg.toByteArray()
	fmt.Printf("ID: %d, Len: %d, Version: %d\n", msg.ID, msg.Len, msg.Version)

	decrypt_nonce.increment()
	opened, ok := secretbox.Open(nil, msg.Buf, &decrypt_nonce.nonce, &serverPublicKey)
	if !ok {
		fmt.Printf("failed to open box\n")
	} else {
		if dump {
			fmt.Print(hexdump.Dump(opened))
		}
	}
}

var accountId uint64
var passToken string
var contentHash string
var androidId string
var live bool

// note, that variables are pointers
var accountIdFlag = flag.Int64("accountId", 0, "Unique account / game identifier")
var passTokenFlag = flag.String("passToken", "", "Used for pairing account / devices")
var contentHashFlag = flag.String("contentHash", "", "Game assets SHA1 hash value")
var androidIdFlag = flag.String("androidId", "", "Devide identifier")
var liveFlag = flag.Bool("live", false, "Connect to live server")

func main() {

	flag.Parse()

	if *accountIdFlag == 0 || *passTokenFlag == "" || *contentHashFlag == "" || *androidIdFlag == "" {
		fmt.Println("All of the arguments are required!")
		flag.PrintDefaults()
		os.Exit(1)
	}

	accountId = uint64(*accountIdFlag)
	passToken = *passTokenFlag
	contentHash = *contentHashFlag
	androidId = *androidIdFlag
	live = *liveFlag

	if 1 == 0 {
		p24101_file()
		//p10100_file()
		//p10101_file()
		//p14102()
		return
	}

	publicKey, privateKey, _ = box.GenerateKey(rand.Reader)

	var host string
	if live {
		host = "game.clashroyaleapp.com:9339"
		serverPublicKey = [32]byte{0x9e, 0x66, 0x57, 0xf2, 0xb4, 0x19, 0xc2, 0x37, 0xf6, 0xae, 0xef, 0x37, 0x08, 0x86, 0x90, 0xa6, 0x42, 0x01, 0x05, 0x86, 0xa7, 0xbd, 0x90, 0x18, 0xa1, 0x56, 0x52, 0xba, 0xb8, 0x37, 0x0f, 0x4f}
	} else {
		host = "127.0.0.1:9339"
		serverPublicKey = [32]byte{0x72, 0xf1, 0xa4, 0xa4, 0xc4, 0x8e, 0x44, 0xda, 0x0c, 0x42, 0x31, 0x0f, 0x80, 0x0e, 0x96, 0x62, 0x4e, 0x6d, 0xc6, 0xa6, 0x41, 0xa9, 0xd4, 0x1c, 0x3b, 0x50, 0x39, 0xd8, 0xdf, 0xad, 0xc2, 0x7e}
	}

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("Unable to connect to host!")
		os.Exit(1)
	}

	encrypt_nonce = NewNonce0()
	decrypt_nonce = NewNonce0()

	m1 := MsgClientHello{
		protocol:     1,
		keyVersion:   10,
		majorVersion: 3,
		minorVersion: 0,
		build:        193,
		contentHash:  contentHash,
		deviceType:   2,
		appStore:     2,
	}
	cmd_10100(conn, m1)

	var buf = make([]byte, 1024*1204)
	var buf_l uint32 = 0

	var read_buf = make([]byte, 2048)
	for {
		n, err := conn.Read(read_buf)
		if err != nil {
			fmt.Print(err)
			return
		}

		for i := 0; i < n; i++ {
			buf[buf_l] = read_buf[i]
			buf_l += 1
		}

		for {
			if buf_l < 5 {
				// not enough data id + len
				break
			}

			l := uint32(buf[4]) | uint32(buf[3])<<8 | uint32(buf[2])<<16
			l += 2 // ID
			l += 3 // length
			l += 2 // version
			if buf_l < l {
				// some part's of the packet is still missing
				break
			}

			msg := fromByteArray(buf)
			switch os := msg.ID; os {
			case 20100:
				cmd_20100(msg)

				m2 := MsgLogin{
					accountId:               accountId,
					passToken:               passToken,
					clientMajorVersion:      int32(m1.majorVersion),
					clientMinorVersion:      int32(m1.minorVersion),
					clientBuild:             int32(m1.build),
					resourceSha:             m1.contentHash,
					UDID:                    "",
					openUdid:                androidId,
					macAddress:              0,
					device:                  "GT-I9195",
					advertisingGuid:         "6b4d572b-da3f-4dc0-9c7f-bc653be81d7f",
					osVersion:               "4.4.2",
					isAndroid:               1,
					unknown1:                0,
					androidID:               androidId,
					preferredDeviceLanguage: "hu-HU",
					unknown2:                1,
					preferredLanguage:       0,
					facebookAttributionId:   "",
					advertisingEnabled:      1,
					appleIFV:                "",
					appStore:                byte(m1.appStore),
					kunlunSSO:               "",
					kunlunUID:               "",
					unknown3:                "",
					unknown4:                "",
					unknown5:                0,
				}
				cmd_10101(conn, m2)
			case 20104:
				cmd_20104(msg)
			case 20103:
				fmt.Printf("LoginError\n")
				return
			case 24101:
				cmd_24101(msg)
			case 24446:
				fmt.Printf("Unkown\n")
				decrypt(msg, false)
				cmd_14405(conn) // 24411 is reply for this, AvatarStream
				cmd_16103(conn) // 26108 is reply for this
				keepalive()
				m3 := MsgAskForBattleReplayStream{
					accountId: accountId,
				}
				cmd_14406(conn, m3) // 24413 is reply for this, BattleReportStream
			case 24311:
				fmt.Printf("Unkown\n")
				decrypt(msg, false)
			case 24411:
				fmt.Printf("Unkown\n")
				decrypt(msg, false)
			case 24413:
				fmt.Printf("Unkown\n")
				decrypt(msg, false)
				ticker()
			case 20108:
				fmt.Printf("KeepAliveOk\n")
				decrypt(msg, true)
				cmd_10107(conn)
			default:
				fmt.Printf("Unkown\n")
				decrypt(msg, true)
			}

			// copy left over data to the begginings
			for i := 0; i < int(buf_l)-int(l); i++ {
				buf[i] = buf[int(l)+i]
			}
			buf_l -= l
		}
	}

}
