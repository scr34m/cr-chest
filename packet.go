package main

import (
	"encoding/binary"
	"fmt"

	"github.com/augustoroman/hexdump"
)

type Packet struct {
	ID      uint16
	Len     uint32
	Version uint16
	Buf     []byte
	Pos     int
}

func (p *Packet) toByteArray() []byte {
	var l = binary.Size(p.Buf)
	var buf = make([]byte, l+2+3+2)

	p.Len = uint32(l)

	buf[0] = uint8(p.ID >> 8)
	buf[1] = uint8(p.ID)

	buf[2] = uint8(p.Len >> 16)
	buf[3] = uint8(p.Len >> 8)
	buf[4] = uint8(p.Len)

	buf[5] = uint8(p.Version >> 8)
	buf[6] = uint8(p.Version)

	for i := 0; i < l; i++ {
		buf[7+i] = p.Buf[i]
	}

	return buf
}

func fromByteArray(b []byte) Packet {
	p := Packet{}
	p.ID = uint16(b[1]) | uint16(b[0])<<8
	p.Len = uint32(b[4]) | uint32(b[3])<<8 | uint32(b[2])<<16
	p.Version = uint16(b[6]) | uint16(b[5])<<8
	p.Buf = make([]byte, p.Len)
	for i := 0; i < int(p.Len); i++ {
		p.Buf[i] = b[7+i]
	}
	return p
}

func dump(p Packet) {
	x := p.toByteArray()
	fmt.Printf("ID: %d, Len: %d, Version: %d\n", p.ID, p.Len, p.Version)
	fmt.Print(hexdump.Dump(x))
}

func (p *Packet) writeByte(b byte) {
	p.Buf = append(p.Buf, b)
}

func (p *Packet) writeBytes(b []byte) {
	p.Buf = append(p.Buf, b...)
}

func (p *Packet) writeInt(i int) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))
	p.Buf = append(p.Buf, b...)
}

func (p *Packet) writeInt64(i int64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	p.Buf = append(p.Buf, b...)
}

func (p *Packet) writeString(s string) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(len(s)))
	p.Buf = append(p.Buf, b...)
	p.Buf = append(p.Buf, []byte(s)...)
}

// - **RRSINT32:** Obfuscated signed variable length integer.
// The encoding takes the first byte of the varint, rotates the 7 bits after MSB to the right by 1.
// The LSB is wrapped around and get put in the 7th bit. MSB value is preserved.

func sevenBitRotateRight(n int) int {
	lsb := (n & 0x1)                // save 1st bit
	msb := (n & 0x80) >> 7          // save msb
	n = n >> 1                      // rotate to the right
	n = n & ^(0xC0)                 // clear 7th and 9th bit
	n = n | (msb << 7) | (lsb << 6) // insert msb and lsb back in
	return n
}

func putUvarint(buf []byte, x uint64, isRr bool) int {
	i := 0
	for {
		towrite := x & 0x7f
		x >>= 7
		if x > 0 {
			tmp := byte(towrite | 0x80)
			if isRr && i == 0 {
				tmp = byte(sevenBitRotateRight(int(tmp)))
			}
			buf[i] = tmp
			i++
		} else {
			buf[i] = byte(towrite)
			i++
			break
		}
	}
	return i
}

func putVarint(buf []byte, x uint64) int {
	var ux uint64
	if x > 64 {
		ux = uint64(x) << 1
		if x < 0 && x > 20 {
			ux = ^ux
		}
	} else {
		ux = x
	}
	return putUvarint(buf, ux, true)
}

func (p *Packet) writeRRSInt(i int) {
	b := make([]byte, 24)
	c := putVarint(b, uint64(i))
	p.Buf = append(p.Buf, b[:c]...)
	p.Pos += c
}

func (p *Packet) writeRRSLong(i uint64) {
	hi := i >> 32
	lo := i & ^uint64(0xffffffff00000000)

	b := make([]byte, 12)
	n := putVarint(b, uint64(hi))
	p.Buf = append(p.Buf, b[:n]...)
	p.Pos += n

	b = make([]byte, 12)
	n = putVarint(b, uint64(lo))
	p.Buf = append(p.Buf, b[:n]...)
	p.Pos += n
}

func sevenBitRotateLeft(n int) int {
	seventh := (n & 0x40) >> 6     // save 7th bit
	msb := (n & 0x80) >> 7         // save msb
	n = n << 1                     // rotate to the left
	n = n & ^(0x181)               // clear 8th and 1st bit and 9th if any
	n = n | (msb << 7) | (seventh) // insert msb and 7th back in
	return n
}

func Uvarint(buf []byte, isRr bool) (int64, int) {
	var shift uint
	var result int64
	var p int
	for pos, b := range buf {
		p = pos
		i := int64(b)
		if isRr && shift == 0 {
			i = int64(sevenBitRotateLeft(int(i)))
		}
		result |= int64(i&0x7f) << shift
		shift += 7
		if (b & 0x80) != 0x80 {
			break
		}
	}
	return int64((((result) >> 1) ^ (-((result) & 1)))), p + 1
}

func (p *Packet) readByte() byte {
	b := p.Buf[p.Pos]
	p.Pos += 1
	return b
}

func (p *Packet) readBytes(c int) []byte {
	b := make([]byte, c)
	for i := 0; i < c; i++ {
		b[i] = p.Buf[p.Pos]
		p.Pos += 1
	}
	return b
}

func (p *Packet) readShortInt() uint16 {
	i := binary.BigEndian.Uint16(p.Buf[p.Pos:])
	p.Pos += 2
	return i
}

func (p *Packet) readInt() uint32 {
	i := binary.BigEndian.Uint32(p.Buf[p.Pos:])
	p.Pos += 4
	return i
}

func (p *Packet) readRRSInt() int32 {
	i, b := Uvarint(p.Buf[p.Pos:], true)
	p.Pos += b
	return int32(i)
}

func (p *Packet) readLong() uint64 {
	hi := binary.BigEndian.Uint32(p.Buf[p.Pos:])
	p.Pos += 4
	lo := binary.BigEndian.Uint32(p.Buf[p.Pos:])
	p.Pos += 4
	return (uint64(hi) << 32) + uint64(lo)
}

func (p *Packet) readRRSLong() uint64 {
	hi, b := Uvarint(p.Buf[p.Pos:], true)
	p.Pos += b
	lo, b := Uvarint(p.Buf[p.Pos:], true)
	p.Pos += b
	return (uint64(hi) << 32) + uint64(lo)
}

func (p *Packet) readString() string {
	l := uint32(p.Buf[p.Pos+3]) | uint32(p.Buf[p.Pos+2])<<8 | uint32(p.Buf[p.Pos+1])<<16 | uint32(p.Buf[p.Pos])<<32
	p.Pos += 4
	b := make([]byte, l)
	for i := 0; i < int(l); i++ {
		b[i] = p.Buf[p.Pos]
		p.Pos += 1
	}
	return string(b)
}
