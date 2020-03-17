package packer_test

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/p2sub/p2sub/packer"
)

func TestPush(t *testing.T) {
	serialized := packer.NewSerialize()
	serialized.Push(uint64(0x1122334455667788))
	got := hex.EncodeToString(serialized.Bytes())
	if got != "1122334455667788" {
		t.Errorf("serialized.Push(uint64(0x1122334455667788)) = %s; want 1122334455667788", got)
	}
}

func TestWrite(t *testing.T) {
	serialized := packer.NewSerialize()
	serialized.Write(uint64(0x1122334455667788),
		uint32(0xaabbccdd),
		uint8(0xff),
		uint16(0xffaa),
		[]byte{255, 255, 255},
		"Hello! I'm Mary")
	got := hex.EncodeToString(serialized.Bytes())
	if got != "1122334455667788aabbccddffffaaffffff48656c6c6f212049276d204d617279" {
		t.Errorf("serialized.Write(...) = %s; want 1122334455667788aabbccddffffaaffffff48656c6c6f212049276d204d617279", got)
	}
}

func BenchmarkPush(b *testing.B) {
	serialized := packer.NewSerialize()
	for i := 0; i < b.N; i++ {
		serialized.Push(uint64(0x1122334455667788))
	}
}

func BenchmarkWrite(b *testing.B) {
	serialized := packer.NewSerialize()
	for i := 0; i < b.N; i++ {
		serialized.Write(uint64(0x1122334455667788), uint32(0xffaabbcc), 0xdd)
	}
}

func TestReadString(t *testing.T) {
	uString, err := packer.NewUnserialize([]byte("Hello!\nI'm chiro\x00This is another string.\x00")).ReadString(6)
	if uString != "Hello!" && err != nil {
		t.Errorf("unserialized.Readstring(6) = %s want Hello!", uString)
	}
	uString, err = packer.NewUnserialize([]byte("Hello!\nI'm chiro\x00This is another string.\x00")).ReadString()
	if uString != "Hello!\nI'm chiro" && err != nil {
		t.Errorf("unserialized.Readstring() = %s want Hello!\\nI'm chiro", uString)
	}
}

func BenchmarkReadUint8(b *testing.B) {
	b.StopTimer()
	s := b.N
	buf := randomBytes(s)
	unserialized := packer.NewUnserialize(buf)
	b.StartTimer()
	for i := 0; i < s; i++ {
		unserialized.ReadUint8()
	}
}

func randomBytes(s int) []byte {
	buf := make([]byte, s)
	_, _ = rand.Read(buf)
	return buf
}

func TestDataReader(t *testing.T) {
	testMap := map[int]int{
		1: 254,
		2: 252,
		3: 248,
		4: 240,
		5: 220,
		6: 0,
	}
	result := true
	unserialized := packer.NewUnserialize(randomBytes(255))
	unserialized.ReadUint8()
	result = result && (unserialized.Len() == testMap[1])
	unserialized.ReadUint16()
	result = result && (unserialized.Len() == testMap[2])
	unserialized.ReadUint32()
	result = result && (unserialized.Len() == testMap[3])
	unserialized.ReadUint64()
	result = result && (unserialized.Len() == testMap[4])
	unserialized.ReadBytes(20)
	result = result && (unserialized.Len() == testMap[5])
	unserialized.ReadBytes()
	result = result && (unserialized.Len() == testMap[6])
	result = result && unserialized.Size() == 255
	_, err := unserialized.ReadUint8()
	result = result && err != nil
	_, err = unserialized.ReadUint16()
	result = result && err != nil
	_, err = unserialized.ReadUint32()
	result = result && err != nil
	_, err = unserialized.ReadUint64()
	result = result && err != nil
	_, err = unserialized.ReadBytes()
	result = result && err != nil
	_, err = unserialized.ReadString()
	result = result && err != nil
	if !result {
		t.Error("Data remaining was wrong")
	}
}
