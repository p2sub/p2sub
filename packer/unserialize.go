package packer

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

//Unserialize struct to binary data
type Unserialize struct {
	buffer *bytes.Reader
}

//New unserialize base on []byte
func New(d []byte) *Unserialize {
	return &Unserialize{buffer: bytes.NewReader(d)}
}

//Len of unserialized bytes
func (u *Unserialize) Len() int {
	return u.buffer.Len()
}

//Size of total size
func (u *Unserialize) Size() int {
	return int(u.buffer.Size())
}

//Pop new value to buffer
func (u *Unserialize) Pop(kind reflect.Kind, size ...int) (result interface{}) {
	resultLen := 0
	if len(size) == 0 {
		resultLen = u.buffer.Len()
	} else {
		resultLen = size[0]
	}
	switch kind {
	case reflect.Uint8:
		if b, err := u.buffer.ReadByte(); err == nil {
			result = b
		}
		break
	case reflect.Uint16:
		buf := make([]byte, 2)
		if n, err := u.buffer.Read(buf); err == nil && n == 2 {
			result = binary.BigEndian.Uint16(buf)
		}
		break
	case reflect.Uint32:
		buf := make([]byte, 4)
		if n, err := u.buffer.Read(buf); err == nil && n == 4 {
			result = binary.BigEndian.Uint32(buf)
		}
		break
	case reflect.Uint64:
		buf := make([]byte, 8)
		if n, err := u.buffer.Read(buf); err == nil && n == 8 {
			result = binary.BigEndian.Uint64(buf)
		}
		break
	//Read ulti null byte
	case reflect.String:
		tmp := ""
		for i := 0; u.buffer.Len() > 0 && i < resultLen; i++ {
			if b, err := u.buffer.ReadByte(); b != 0 && err == nil {
				tmp += string(b)
			} else {
				break
			}
		}
		result = tmp
		break
	case reflect.Slice:
		buf := make([]byte, resultLen)
		if n, err := u.buffer.Read(buf); err == nil && n == resultLen {
			result = buf
		}
		break
	default:
	}
	return
}
