package main

type Chest struct {
	unknown_1 byte
	unknown_2 byte
	status    int // 0 closed, 1 opened, 8 countdown
	time1     int
	time2     int
	timestamp int
	id        int
	unknown_3 uint32
	typ       byte
	removed   bool
}

type Card struct {
	id       int32
	level    byte
	isNew    int32
	cards    int32
	b1       int32
	b2       int32
	b3       int32
	selected bool
}

type Deck struct {
	index int
	cards []int32
}

type News struct {
	id        int32
	title     string
	unknown_1 []byte
	title2    string
	json      string
}

type HomeData struct {
	unknown_1   uint32
	unknown_2   int32
	unknown_3   int32
	unknown_4   []byte
	unknown_5   int32
	unknown_6   int32
	timestamp_1 int32
	unknown_7   byte
	decks       []Deck
	unknown_8   byte
	cards       []Card
	unknown_9   []byte
	unknown_10  []byte
	news        []News
	unknown_11  []int32
	unknown_12  []int32
	unknown_13  string
	unknown_14  uint16
	chests      []Chest
	unknown_15  []int32
	unknown_16  []int32
	unknown_17  []int32
	unknown_18  []int32
	unknown_19  []int32
	unknown_20  []int32
	unknown_21  [][]int32
	unknown_22  []int32
	unknown_23  []int32
	unknown_24  []int32
	unknown_25  []int32
	unknown_26  [][]int32
	unknown_27  []int32
	unknown_28  []int32
	unknown_29  []int32
	unknown_30  []int32
	unknown_31  []int32
	accountId   uint64
	accountId2  uint64
	accountId3  uint64
	nick        string
	unknown_32  uint16
	thropyCount int32
	unknown_33  []int32
	unknown_34  []byte
	unknown_35  [][]int32
	unknown_36  []byte
	unknown_37  [][]int32
	unknown_38  [][]int32
	unknown_39  [][]int32
	unknown_40  [][]int32
	gem1        int32
	gem2        int32
	xp          int32
	level       int32
	unknown_41  byte
	rank        byte // 0 demo, 1 not demo, 9 elder
	cardsFound  int32
	guildId     int32
	guild       string
	unknown_42  int32
	unknown_43  int32
	unknown_44  int32
	unknown_45  int32
	unknown_46  int32
	matchWin    int32
	matchList   int32
	unknown_47  byte // 7F 1st loose, 7E 2nd loose, 1 1st win
	progress    int32
	unknown_48  []byte
	unknown_49  []byte
}
