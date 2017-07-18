package packet

import (
	"encoding/binary"
	"errors"
)

// ============================================================================

/*
	packet format:
		length 4 bytes
		op     4 bytes
		body   []byte
		sid    8 bytes (optional)
*/

// ============================================================================

const (
	CMAXPACKETLEN = 1024 * 1024
)

// ============================================================================

var (
	ErrPacketLength = errors.New("invalid packet length")
)

// ============================================================================

type Reader struct {
	data []byte
	ptr  []byte
	l    uint32
}

func NewReader() *Reader {
	t := &Reader{}
	t.reset()
	return t
}

func (self *Reader) reset() {
	self.data = make([]byte, 4)
	self.ptr = self.data
	self.l = 0
}

func (self *Reader) Read(buf []byte) (p Packet, bufPtr []byte, err error) {

	bufPtr = buf

	// read length
	if self.l == 0 {
		n := copy(self.ptr, bufPtr)
		self.ptr, bufPtr = self.ptr[n:], bufPtr[n:]

		// partial
		if len(self.ptr) > 0 {
			return
		}

		// get length
		self.l = binary.BigEndian.Uint32(self.data)

		// check length
		if self.l < 8 || self.l > CMAXPACKETLEN {
			self.l = 0 // for safety
			return nil, nil, ErrPacketLength
		}

		// alloc packet buffer
		//	+8: provide more room for appending sid without buffer copy
		newData := make([]byte, self.l+8)
		copy(newData, self.data)
		self.data = newData
		self.ptr = self.data[4:self.l]
	}

	// read rest
	n := copy(self.ptr, bufPtr)
	self.ptr, bufPtr = self.ptr[n:], bufPtr[n:]

	// partial
	if len(self.ptr) > 0 {
		return
	}

	// full packet
	p = self.data[:self.l]
	self.reset()

	return
}

// ============================================================================

type Writer struct {
}

func NewWriter() *Writer {
	return &Writer{}
}

func (self *Writer) Write(p Packet) (buf []byte) {
	return p
}

// ============================================================================

type Packet []byte

func (self Packet) Op() uint32 {
	return binary.BigEndian.Uint32(self[4:8])
}

func (self Packet) Body() []byte {
	return self[8:]
}

func (self *Packet) AddSid(sid uint64) {
	l := len(*self)

	*self = (*self)[:l+8]
	binary.BigEndian.PutUint64((*self)[l:l+8], sid)
	binary.BigEndian.PutUint32((*self)[:4], uint32(l+8))
}

func (self *Packet) RemoveSid() uint64 {
	l := len(*self)

	sid := binary.BigEndian.Uint64((*self)[l-8:])

	*self = (*self)[:l-8]
	binary.BigEndian.PutUint32((*self)[:4], uint32(l-8))

	return sid
}

// ============================================================================

func Assemble(op uint32, body []byte) Packet {
	l := 4 + 4 + len(body)

	buf := make([]byte, l, l+8)

	binary.BigEndian.PutUint32(buf[:4], uint32(l))
	binary.BigEndian.PutUint32(buf[4:8], op)
	copy(buf[8:], body)

	return Packet(buf)
}
