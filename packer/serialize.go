package packer

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

//Serialize struct to binary data
type Serialize struct {
	buffer bytes.Buffer
}

//New serialize
func New() *Serialize {
	return &Serialize{}
}

//Push new value to buffer
func (s *Serialize) Push(value interface{}) {
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Uint8:
		s.buffer.Write([]byte{byte(v.Uint())})
		break
	case reflect.Uint16:
		buf := make([]byte, 2)
		binary.BigEndian.PutUint16(buf, uint16(v.Uint()))
		s.buffer.Write(buf)
		break
	case reflect.Uint32:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v.Uint()))
		s.buffer.Write(buf)
		break
	case reflect.Uint64:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(v.Uint()))
		s.buffer.Write(buf)
		break
	case reflect.String:
		s.buffer.Write([]byte(v.String()))
		break
	case reflect.Slice:
		s.buffer.Write(v.Bytes())
		break
	default:
		//Do nothing
	}
}

//Write many thing at once instead of push on by one
func (s *Serialize) Write(v ...interface{}) {
	for i := 0; i < len(v); i++ {
		s.Push(v[i])
	}
}

//Bytes return buffer bytes array
func (s *Serialize) Bytes() []byte {
	return s.buffer.Bytes()
}
