package mp4

import (
	"encoding/binary"
	"io"
	"io/ioutil"
)

// SmhdBox - Sound Media Header Box (smhd - mandatory for sound tracks)
//
// Contained in : Media Information Box (minf)
//
// Status: decoded
type SmhdBox struct {
	Version byte
	Flags   [3]byte
	Balance uint16 // should be int16
}

// DecodeSmhd - box-specific decode
func DecodeSmhd(size uint64, startPos uint64, r io.Reader) (Box, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &SmhdBox{
		Version: data[0],
		Flags:   [3]byte{data[1], data[2], data[3]},
		Balance: binary.BigEndian.Uint16(data[4:6]),
	}, nil
}

// Type - box type
func (b *SmhdBox) Type() string {
	return "smhd"
}

// Size - calculated size of box
func (b *SmhdBox) Size() uint64 {
	return boxHeaderSize + 8
}

// Encode - write box to w
func (b *SmhdBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	binary.BigEndian.PutUint16(buf[4:], b.Balance)
	_, err = w.Write(buf)
	return err
}
