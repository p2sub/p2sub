// Copyright 2019 Trần Anh Dũng <chiro@fkguru.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package packer

import (
	"bytes"
	"encoding/binary"
	"io"
)

//Unserialize struct to binary data
type Unserialize struct {
	buffer *bytes.Reader
}

//NewUnserialize new unserialize instance based on []byte
func NewUnserialize(d []byte) *Unserialize {
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

//ReadUint8 from buffer
func (u *Unserialize) ReadUint8() (uint8, error) {
	b, err := u.buffer.ReadByte()
	if err == nil {
		return uint8(b), err
	}
	return 0, err
}

//ReadUint16 from buffer
func (u *Unserialize) ReadUint16() (uint16, error) {
	buf := make([]byte, 2)
	n, err := u.buffer.Read(buf)
	if err == nil && n == 2 {
		return binary.BigEndian.Uint16(buf), nil
	}
	return 0, err
}

//ReadUint32 from buffer
func (u *Unserialize) ReadUint32() (uint32, error) {
	buf := make([]byte, 4)
	n, err := u.buffer.Read(buf)
	if err == nil && n == 4 {
		return binary.BigEndian.Uint32(buf), nil
	}
	return 0, err
}

//ReadUint64 from buffer
func (u *Unserialize) ReadUint64() (uint64, error) {
	buf := make([]byte, 8)
	n, err := u.buffer.Read(buf)
	if err == nil && n == 8 {
		return binary.BigEndian.Uint64(buf), nil
	}
	return 0, err
}

//ReadString from buffer
func (u *Unserialize) ReadString(size ...int) (string, error) {
	if u.buffer.Len() == 0 {
		return "", io.EOF
	}
	resultLen := 0
	if len(size) == 0 {
		resultLen = u.buffer.Len()
	} else {
		resultLen = size[0]
	}
	tmp := ""
	for i := 0; u.buffer.Len() > 0 && i < resultLen; i++ {
		if b, err := u.buffer.ReadByte(); b != 0 && err == nil {
			tmp += string(b)
		} else {
			return tmp, err
		}
	}
	return tmp, nil
}

//ReadBytes from buffer
func (u *Unserialize) ReadBytes(size ...int) ([]byte, error) {
	resultLen := 0
	if len(size) == 0 {
		resultLen = u.buffer.Len()
	} else {
		resultLen = size[0]
	}
	buf := make([]byte, resultLen)
	n, err := u.buffer.Read(buf)
	if err == nil && n == resultLen {
		return buf, nil
	}
	return nil, err
}
