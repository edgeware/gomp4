package mp4

import "io"

// MinfBox -  Media Information Box (minf - mandatory)
//
// Contained in : Media Box (mdia)
//
// Status: partially decoded (hmhd - hint tracks - and nmhd - null media - are ignored)
type MinfBox struct {
	Vmhd  *VmhdBox
	Smhd  *SmhdBox
	Stbl  *StblBox
	Dinf  *DinfBox
	Hdlr  *HdlrBox
	boxes []Box
}

// DecodeMinf - box-specific decode
func DecodeMinf(size uint64, startPos uint64, r io.Reader) (Box, error) {
	l, err := DecodeContainer(size, startPos, r)
	if err != nil {
		return nil, err
	}
	m := &MinfBox{}
	m.boxes = l
	for _, b := range l {
		switch b.Type() {
		case "vmhd":
			m.Vmhd = b.(*VmhdBox)
		case "smhd":
			m.Smhd = b.(*SmhdBox)
		case "stbl":
			m.Stbl = b.(*StblBox)
		case "dinf":
			m.Dinf = b.(*DinfBox)
		case "hdlr":
			m.Hdlr = b.(*HdlrBox)
		}
	}
	return m, nil
}

// Type - box type
func (b *MinfBox) Type() string {
	return "minf"
}

// Size - calculated size of box
func (b *MinfBox) Size() uint64 {
	return containerSize(b.boxes)
}

// Dump - print box info
func (b *MinfBox) Dump() {
	b.Stbl.Dump()
}

// Encode - write box to w
func (b *MinfBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	if b.Vmhd != nil {
		err = b.Vmhd.Encode(w)
		if err != nil {
			return err
		}
	}
	if b.Smhd != nil {
		err = b.Smhd.Encode(w)
		if err != nil {
			return err
		}
	}
	err = b.Dinf.Encode(w)
	if err != nil {
		return err
	}
	err = b.Stbl.Encode(w)
	if err != nil {
		return err
	}
	if b.Hdlr != nil {
		return b.Hdlr.Encode(w)
	}
	return nil
}
