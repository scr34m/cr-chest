package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"

	"golang.org/x/crypto/nacl/secretbox"
)

// OwnHomeData
func cmd_24101(msg Packet) {
	decrypt_nonce.increment()
	opened, ok := secretbox.Open(nil, msg.Buf, &decrypt_nonce.nonce, &serverPublicKey)
	if !ok {
		fmt.Printf("failed to open box")
	} else {
		ioutil.WriteFile("5203.dat", opened, 0644)
		// fmt.Print(hexdump.Dump(opened))
		p := Packet{
			Buf: opened,
			Pos: 0,
		}
		p24101(&p, 24101)
	}
}

/*
ls -al *-24101.bin
*/
func p24101_file() {
	nr := 5342
	b, _ := ioutil.ReadFile(fmt.Sprintf("cr-proxy/replay/%d-24101.bin", nr))
	p := Packet{
		Buf: b,
		Pos: 7,
	}
	/*
		b, _ := ioutil.ReadFile("24101.dat")
		p := Packet{
			Buf: b,
			Pos: 0,
		}
	*/
	p24101(&p, nr)
}

func p24101(p *Packet, nr int) {
	hd := HomeData{}

	hd.unknown_1 = p.readInt()
	hd.unknown_2 = p.readRRSInt()
	hd.unknown_3 = p.readRRSInt()

	// innen
	// 4879 00 00 39 20 25 09
	// 4905 00 13 E3 8A B2 05 80 02
	if hd.unknown_3 == 0 {
		hd.unknown_4 = p.readBytes(4)
	} else {
		hd.unknown_4 = p.readBytes(6)
	}
	// eddig

	hd.unknown_5 = p.readRRSInt()
	hd.unknown_6 = p.readRRSInt()
	hd.timestamp_1 = p.readRRSInt()
	hd.unknown_7 = p.readByte()

	c1 := int(p.readByte())
	for i := 0; i < c1; i++ {
		deck := Deck{}
		deck.index = i
		c2 := int(p.readByte())
		for j := 0; j < c2; j++ {
			deck.cards = append(deck.cards, p.readRRSInt())
		}
		hd.decks = append(hd.decks, deck)
	}

	hd.unknown_8 = p.readByte()

	// aktuális kiválasztott deck
	getCard(p, &hd, true)
	getCard(p, &hd, true)
	getCard(p, &hd, true)
	getCard(p, &hd, true)
	getCard(p, &hd, true)
	getCard(p, &hd, true)
	if hd.unknown_8 >= 0x7f {
		getCard(p, &hd, true)
	}
	if hd.unknown_8 == 0xff {
		getCard(p, &hd, true)
	}

	// többi kártya
	if hd.unknown_8 == 0xff {
		c1 = int(p.readRRSInt())
		for i := 0; i < c1; i++ {
			getCard(p, &hd, false)
		}
		hd.unknown_9 = p.readBytes(1)
	} else {
		hd.unknown_9 = p.readBytes(2)
	}

	hd.unknown_10 = p.readBytes(10)

	c1 = int(p.readShortInt())
	for i := 0; i < c1; i++ {
		news := News{}
		news.id = p.readRRSInt()
		news.title = p.readString()
		news.unknown_1 = p.readBytes(24)
		news.title2 = p.readString()
		news.json = p.readString()

		news.title = ""
		news.title2 = ""
		news.json = ""

		hd.news = append(hd.news, news)
	}

	v1 := p.readRRSInt()
	hd.unknown_11 = append(hd.unknown_11, v1)
	if v1 != 0 {
		panic("after news 1\n")
	}
	v1 = p.readRRSInt()
	hd.unknown_11 = append(hd.unknown_11, v1)
	if v1 != 0 {
		hd.unknown_11 = append(hd.unknown_11, p.readRRSInt())
	}
	v1 = p.readRRSInt()
	hd.unknown_11 = append(hd.unknown_11, v1)
	if v1 != 0 {
		panic("after news 3\n")
	}
	v1 = p.readRRSInt()
	hd.unknown_11 = append(hd.unknown_11, v1)
	if v1 != 0 {
		panic("after news 4\n")
	}
	v1 = p.readRRSInt()
	hd.unknown_11 = append(hd.unknown_11, v1)
	if v1 != 0 {
		panic("after news 5\n")
	}
	v1 = p.readRRSInt()
	hd.unknown_11 = append(hd.unknown_11, v1)
	if v1 != 0 {
		panic("after news 6\n")
	}
	hd.unknown_11 = append(hd.unknown_11, p.readRRSInt())
	hd.unknown_11 = append(hd.unknown_11, p.readRRSInt())
	v1 = p.readRRSInt()
	hd.unknown_11 = append(hd.unknown_11, v1)
	if v1 != 0 {
		hd.unknown_11 = append(hd.unknown_11, p.readRRSInt())
	}
	v1 = p.readRRSInt()
	hd.unknown_11 = append(hd.unknown_11, v1)
	if v1 != 0 {
		panic("after news 10\n")
	}

	if len(hd.news) > 0 {
		for i := 0; i < len(hd.news); i++ {
			hd.unknown_12 = append(hd.unknown_12, p.readRRSInt())
			hd.unknown_12 = append(hd.unknown_12, p.readRRSInt())
		}
	} else {
		hd.unknown_12 = append(hd.unknown_12, p.readRRSInt())
	}
	hd.unknown_12 = append(hd.unknown_12, p.readRRSInt())

	hd.unknown_13 = p.readString()

	hd.unknown_14 = p.readShortInt()
	for {
		typ := p.readByte()
		if typ == 0 {
			break
		}
		getChest(p, typ, &hd)
	}

	freeChestT1 := int(p.readRRSInt())
	freeChestT2 := int(p.readRRSInt())
	freeChestTS := int(p.readRRSInt())

	crownChestT1 := int(p.readRRSInt())
	crownChestT2 := int(p.readRRSInt())
	crownChestTS := int(p.readRRSInt())

	hd.chests = append(hd.chests, Chest{id: -2, time1: crownChestT1, time2: crownChestT2, timestamp: crownChestTS})

	freechest := p.readByte()
	if freechest != 0 {
		if freeChestT1 == 0 {
			hd.chests = append(hd.chests, Chest{id: -1, time1: freeChestT1, time2: freeChestT2, timestamp: freeChestTS})
			hd.chests = append(hd.chests, Chest{id: -1, time1: freeChestT1, time2: freeChestT2, timestamp: freeChestTS})
		}
		if freeChestT1 < freeChestT2 {
			hd.chests = append(hd.chests, Chest{id: -1, time1: freeChestT1, time2: freeChestT2, timestamp: freeChestTS})
		}
		for i := 0; i < 9; i++ {
			hd.unknown_15 = append(hd.unknown_15, p.readRRSInt())
		}
	}

	// 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1704240, 1728000, 1492060855, 734000, 1701740, 1492012343
	for i := 0; i < 17; i++ {
		hd.unknown_16 = append(hd.unknown_16, p.readRRSInt())
	}

	if p.readRRSInt() != 0 { // 0 vagy 1
		// 19, 38, 1, 392, 0, -64, 0, 0, 0
		for i := 0; i < 9; i++ {
			hd.unknown_17 = append(hd.unknown_17, p.readRRSInt())
		}
	}

	if p.readRRSInt() != 0 {
		panic("0 kellene, hogy legyen 1")
	}

	if p.readRRSInt() != 0 {
		panic("0 kellene, hogy legyen 2")
	}

	// -64, 3, 0, 0, 0, 0, 0, 0, 0
	for i := 0; i < 9; i++ {
		hd.unknown_18 = append(hd.unknown_18, p.readRRSInt())
	}

	// 2, 2, 54
	for i := 0; i < 3; i++ {
		hd.unknown_19 = append(hd.unknown_19, p.readRRSInt())
	}

	// 6, 2571181, 4, 3, 1201260, 1201260, 1492034399
	for i := 0; i < 7; i++ {
		hd.unknown_20 = append(hd.unknown_20, p.readRRSInt())
	}

	/*
		1 66 2 0 0 0 0 0 0 26 2 0
		1 66 2 0 0 0 0 0 0 26 11 1
		1 66 2 0 0 0 0 0 0 27 8 2
	*/
	c1 = int(p.readRRSInt())
	for i := 0; i < c1; i++ {
		u := []int32{}
		for i := 0; i < 12; i++ {
			u = append(u, p.readRRSInt())
		}
		hd.unknown_21 = append(hd.unknown_21, u)
	}

	c1 = int(p.readRRSInt())
	if c1 != 0 {
		hd.unknown_22 = append(hd.unknown_22, int32(c1))
		c2 := p.readRRSInt()
		var r int
		if c2 == 2 {
			r = 10
		} else if c2 == 3 {
			r = 8
		} else {
			panic("unknown_22 hossz gond \n")
		}
		for i := 0; i < r; i++ {
			hd.unknown_22 = append(hd.unknown_22, p.readRRSInt())
		}
	}
	// fmt.Printf("%+v\n", hd.unknown_22)

	// 0, 0, -64, 0, 0, -64, 0, 0, -64
	for i := 0; i < 9; i++ {
		hd.unknown_23 = append(hd.unknown_23, p.readRRSInt())
	}

	hd.unknown_24 = append(hd.unknown_24, p.readRRSInt())
	hd.unknown_24 = append(hd.unknown_24, p.readRRSInt())
	hd.unknown_24 = append(hd.unknown_24, p.readRRSInt())
	hd.unknown_24 = append(hd.unknown_24, p.readRRSInt())
	hd.unknown_24 = append(hd.unknown_24, p.readRRSInt())
	c2 := p.readRRSInt()
	hd.unknown_24 = append(hd.unknown_24, c2)
	if c2 != 0 {
		hd.unknown_24 = append(hd.unknown_24, p.readRRSInt())
	}
	hd.unknown_24 = append(hd.unknown_24, p.readRRSInt())
	c2 = p.readRRSInt()
	hd.unknown_24 = append(hd.unknown_24, c2)
	for i := 0; i < int(c2); i++ {
		hd.unknown_24 = append(hd.unknown_24, p.readRRSInt())
	}
	c2 = p.readRRSInt()
	hd.unknown_24 = append(hd.unknown_24, c2)
	if c2 != 0 {
		hd.unknown_24 = append(hd.unknown_24, p.readRRSInt())
	}
	hd.unknown_24 = append(hd.unknown_24, p.readRRSInt())
	hd.unknown_24 = append(hd.unknown_24, p.readRRSInt())

	// 9, 0, 0, 0, 0, -505, 1
	for i := 0; i < 7; i++ {
		hd.unknown_25 = append(hd.unknown_25, p.readRRSInt())
	}

	for i := 0; i < 8; i++ {
		u := []int32{}
		for i := 0; i < 7; i++ {
			u = append(u, p.readRRSInt())
		}
		hd.unknown_26 = append(hd.unknown_26, u)
	}

	// 26000046, 1, 26000046, 0
	for i := 0; i < 4; i++ {
		hd.unknown_27 = append(hd.unknown_27, p.readRRSInt())
	}

	if p.readRRSInt() != 0 {
		panic("0 kellene, hogy legyen\n")
	}

	// 4, 66000009, 66000010, 66000011, 66000012
	for i := 0; i < 5; i++ {
		hd.unknown_28 = append(hd.unknown_28, p.readRRSInt())
	}

	// 5229 [0 1 1608786000 0 1 1 1 54000010]
	// 5289 [0 1 1608786000 0 1 1 1 54000010]
	for i := 0; i < 8; i++ {
		hd.unknown_29 = append(hd.unknown_29, p.readRRSInt())
	}

	c1 = int(p.readRRSInt())
	if c1 != 0 {
		// 5229 [2 0 -1543952872 1491552036 1491811236 2861408 1 0 0]
		// 5289 [2 1 -1885200892 1492156836 1492416036 3120975 0 1 0]
		hd.unknown_30 = append(hd.unknown_30, int32(c1))
		hd.unknown_30 = append(hd.unknown_30, p.readRRSInt())
		hd.unknown_30 = append(hd.unknown_30, int32(p.readInt()))
		hd.unknown_30 = append(hd.unknown_30, int32(p.readInt()))
		hd.unknown_30 = append(hd.unknown_30, int32(p.readInt()))
		hd.unknown_30 = append(hd.unknown_30, p.readRRSInt())
		hd.unknown_30 = append(hd.unknown_30, p.readRRSInt())
		hd.unknown_30 = append(hd.unknown_30, p.readRRSInt())
		hd.unknown_30 = append(hd.unknown_30, p.readRRSInt())
	}

	// 5229 [0 0 0 0 0 0]
	// 5289 [0 0 0 0 0 1]
	for i := 0; i < 6; i++ {
		hd.unknown_31 = append(hd.unknown_31, p.readRRSInt())
	}

	hd.accountId = p.readRRSLong()
	hd.accountId2 = p.readRRSLong()
	hd.accountId3 = p.readRRSLong()
	hd.nick = p.readString()
	hd.unknown_32 = p.readShortInt()
	hd.thropyCount = p.readRRSInt()

	// 5289 [334 0 0 0 0]
	for i := 0; i < 5; i++ {
		hd.unknown_33 = append(hd.unknown_33, p.readRRSInt())
	}

	// 5289 [0 0 30 0 0 0 0 0 7]
	for i := 0; i < 9; i++ {
		hd.unknown_34 = append(hd.unknown_34, p.readByte())
	}

	c1 = int(p.readByte())
	for i := 0; i < c1; i++ {
		u := []int32{}
		for i := 0; i < 3; i++ {
			u = append(u, p.readRRSInt())
		}
		hd.unknown_35 = append(hd.unknown_35, u)
	}

	c1 = int(p.readByte())
	for i := 0; i < c1; i++ {
		hd.unknown_36 = append(hd.unknown_36, p.readByte())
	}

	if c1 != 0 && p.readByte() != 0 {
		panic("0 kellene, hogy legyen")
	}

	c1 = int(p.readByte())
	for i := 0; i < c1; i++ {
		u := []int32{}
		for i := 0; i < 3; i++ {
			u = append(u, p.readRRSInt())
		}
		hd.unknown_37 = append(hd.unknown_37, u)
	}

	c1 = int(p.readByte())
	for i := 0; i < c1; i++ {
		u := []int32{}
		for i := 0; i < 3; i++ {
			u = append(u, p.readRRSInt())
		}
		hd.unknown_38 = append(hd.unknown_38, u)
	}

	c1 = int(p.readByte())
	for i := 0; i < c1; i++ {
		u := []int32{}
		for i := 0; i < 3; i++ {
			u = append(u, p.readRRSInt())
		}
		hd.unknown_39 = append(hd.unknown_39, u)
	}

	c1 = int(p.readByte())
	for i := 0; i < c1; i++ {
		u := []int32{}
		for i := 0; i < 3; i++ {
			u = append(u, p.readRRSInt())
		}
		hd.unknown_40 = append(hd.unknown_40, u)
	}

	if p.readByte() != 0 {
		panic("0 kellene, hogy legyen")
	}
	hd.gem1 = p.readRRSInt()
	hd.gem2 = p.readRRSInt()
	hd.xp = p.readRRSInt()
	hd.level = p.readRRSInt()
	hd.unknown_41 = p.readByte()
	hd.rank = p.readByte()
	hd.cardsFound = p.readRRSInt()
	hd.guildId = p.readRRSInt()
	if hd.guildId != 0 {
		hd.guild = p.readString()
		hd.unknown_42 = p.readRRSInt()
		hd.unknown_43 = p.readRRSInt()
		hd.unknown_44 = p.readRRSInt()
		hd.unknown_45 = p.readRRSInt()
	}
	hd.unknown_46 = p.readRRSInt()

	hd.matchWin = p.readRRSInt()
	hd.matchList = p.readRRSInt()
	hd.unknown_47 = p.readByte()
	hd.progress = p.readRRSInt()
	if hd.progress != 7 {
		panic("7 kellene, ha már nem demo acc")
	}

	for i := 0; i < 13; i++ {
		hd.unknown_48 = append(hd.unknown_48, p.readByte())
	}

	for i := 0; i < 4; i++ {
		hd.unknown_49 = append(hd.unknown_49, p.readByte())
	}

	fmt.Printf("\n---> POS: %d 0x%x\n", p.Pos, p.Pos)

	cs := spew.NewDefaultConfig()
	cs.Indent = " "
	cs.MaxDepth = 3
	cs.DisablePointerAddresses = true
	cs.ContinueOnMethod = true
	cs.DisableCapacities = true

	f, _ := os.Create(fmt.Sprintf("%d.txt", nr))
	cs.Fdump(f, hd)

	homedata = hd
}

func getChest(p *Packet, typ byte, hd *HomeData) {
	chest := Chest{}
	chest.unknown_1 = p.readByte()
	chest.unknown_2 = p.readByte()
	chest.status = int(p.readByte())
	if chest.status == 8 {
		chest.time1 = int(p.readRRSInt())
		chest.time2 = int(p.readRRSInt())
		chest.timestamp = int(p.readRRSInt())
	}
	chest.id = int(p.readRRSInt())
	chest.unknown_3 = p.readInt()
	chest.typ = typ
	hd.chests = append(hd.chests, chest)
}

func getCard(p *Packet, hd *HomeData, selected bool) {
	card := Card{}
	card.id = p.readRRSInt()
	card.level = p.readByte() + 1
	card.isNew = p.readRRSInt()
	card.cards = p.readRRSInt()
	card.b1 = p.readRRSInt()
	card.b2 = p.readRRSInt()
	card.b3 = p.readRRSInt()
	card.selected = selected
	hd.cards = append(hd.cards, card)
}
