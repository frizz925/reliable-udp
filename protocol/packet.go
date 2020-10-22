package protocol

import (
	"io"
)

const (
	// Minimum MTU for IPv6 is 1280 bytes, the IPv6 header size is 40 bytes,
	// and UDP header size 8 bytes. From these numbers, the maximum packet
	// size this protocol can support is 1232 bytes.
	MaxPacketSize = 1232
	// u8 Type + u64 Connection ID + u32 Sequence
	PacketHeaderSize = 13
)

type Packet struct {
	ConnectionID
	Sequence uint32
	Type     PacketType
	Frame
}

func NewBuffer() []byte {
	return make([]byte, MaxPacketSize)
}

func EmptyBuffer() []byte {
	return make([]byte, 0, MaxPacketSize)
}

func Deserialize(r io.Reader) (p Packet, err error) {
	p.ConnectionID, err = ReadConnectionID(r)
	if err != nil {
		return
	}
	p.Sequence, err = ReadUint32(r)
	if err != nil {
		return
	}
	p.Type, err = ReadPacketType(r)
	if err != nil {
		return
	}
	p.Frame, err = ReadFrame(r)
	if err != nil {
		return
	}
	return p, nil
}

func (p Packet) Nonce() Nonce {
	return Uint32ToNonce(p.Sequence)
}

func (p Packet) Serialize(w io.Writer) error {
	if err := p.ConnectionID.Serialize(w); err != nil {
		return err
	}
	if err := WriteUint32(w, p.Sequence); err != nil {
		return err
	}
	if err := p.Type.Serialize(w); err != nil {
		return err
	}
	if p.Frame == nil {
		return nil
	}
	if err := p.Frame.Type().Serialize(w); err != nil {
		return err
	}
	return p.Frame.Serialize(w)
}
