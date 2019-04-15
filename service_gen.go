package members

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Info) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zb0001}
		return
	}
	z.ID, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err, "Services")
		return
	}
	if cap(z.Services) >= int(zb0002) {
		z.Services = (z.Services)[:zb0002]
	} else {
		z.Services = make([]Service, zb0002)
	}
	for za0001 := range z.Services {
		var zb0003 uint32
		zb0003, err = dc.ReadArrayHeader()
		if err != nil {
			err = msgp.WrapError(err, "Services", za0001)
			return
		}
		if zb0003 != 2 {
			err = msgp.ArrayError{Wanted: 2, Got: zb0003}
			return
		}
		z.Services[za0001].Name, err = dc.ReadString()
		if err != nil {
			err = msgp.WrapError(err, "Services", za0001, "Name")
			return
		}
		z.Services[za0001].Port, err = dc.ReadUint16()
		if err != nil {
			err = msgp.WrapError(err, "Services", za0001, "Port")
			return
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Info) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 2
	err = en.Append(0x92)
	if err != nil {
		return
	}
	err = en.WriteString(z.ID)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Services)))
	if err != nil {
		err = msgp.WrapError(err, "Services")
		return
	}
	for za0001 := range z.Services {
		// array header, size 2
		err = en.Append(0x92)
		if err != nil {
			return
		}
		err = en.WriteString(z.Services[za0001].Name)
		if err != nil {
			err = msgp.WrapError(err, "Services", za0001, "Name")
			return
		}
		err = en.WriteUint16(z.Services[za0001].Port)
		if err != nil {
			err = msgp.WrapError(err, "Services", za0001, "Port")
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Info) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 2
	o = append(o, 0x92)
	o = msgp.AppendString(o, z.ID)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Services)))
	for za0001 := range z.Services {
		// array header, size 2
		o = append(o, 0x92)
		o = msgp.AppendString(o, z.Services[za0001].Name)
		o = msgp.AppendUint16(o, z.Services[za0001].Port)
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Info) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zb0001}
		return
	}
	z.ID, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "ID")
		return
	}
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Services")
		return
	}
	if cap(z.Services) >= int(zb0002) {
		z.Services = (z.Services)[:zb0002]
	} else {
		z.Services = make([]Service, zb0002)
	}
	for za0001 := range z.Services {
		var zb0003 uint32
		zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, "Services", za0001)
			return
		}
		if zb0003 != 2 {
			err = msgp.ArrayError{Wanted: 2, Got: zb0003}
			return
		}
		z.Services[za0001].Name, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			err = msgp.WrapError(err, "Services", za0001, "Name")
			return
		}
		z.Services[za0001].Port, bts, err = msgp.ReadUint16Bytes(bts)
		if err != nil {
			err = msgp.WrapError(err, "Services", za0001, "Port")
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Info) Msgsize() (s int) {
	s = 1 + msgp.StringPrefixSize + len(z.ID) + msgp.ArrayHeaderSize
	for za0001 := range z.Services {
		s += 1 + msgp.StringPrefixSize + len(z.Services[za0001].Name) + msgp.Uint16Size
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Service) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zb0001}
		return
	}
	z.Name, err = dc.ReadString()
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	z.Port, err = dc.ReadUint16()
	if err != nil {
		err = msgp.WrapError(err, "Port")
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Service) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 2
	err = en.Append(0x92)
	if err != nil {
		return
	}
	err = en.WriteString(z.Name)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	err = en.WriteUint16(z.Port)
	if err != nil {
		err = msgp.WrapError(err, "Port")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Service) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 2
	o = append(o, 0x92)
	o = msgp.AppendString(o, z.Name)
	o = msgp.AppendUint16(o, z.Port)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Service) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 2 {
		err = msgp.ArrayError{Wanted: 2, Got: zb0001}
		return
	}
	z.Name, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	z.Port, bts, err = msgp.ReadUint16Bytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Port")
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Service) Msgsize() (s int) {
	s = 1 + msgp.StringPrefixSize + len(z.Name) + msgp.Uint16Size
	return
}
