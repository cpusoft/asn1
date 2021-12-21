package der

import (
	"bytes"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/cpusoft/asn1/der/coda"
)

var (
	ErrNodeIsConstructed    = errors.New("node is constructed")
	ErrNodeIsNotConstructed = errors.New("node is not constructed")
)

/*

golang asn1:

type RawValue struct {
	Class, Tag int
	IsCompound bool
	Bytes      []byte
	FullBytes  []byte // includes the tag and length
}

*/

type Node struct {
	class       int
	tag         int
	constructed bool // isCompound

	data  []byte      // Primitive:   (isCompound = false)
	value interface{} // Primitive:  int/bool/string/time... (isCompound = false)
	nodes []*Node     // Constructed: (isCompound = true)
}

func NewNode(class int, tag int) *Node {
	return &Node{
		class: class,
		tag:   tag,
	}
}

func CheckNode(n *Node, class int, tag int) error {
	if n.class != class {
		return fmt.Errorf("class: %d != %d", n.class, class)
	}
	if n.tag != tag {
		return fmt.Errorf("tag: %d != %d", n.tag, tag)
	}
	return nil
}

func (n *Node) GetTag() int {
	return n.tag
}

func (n *Node) getHeader() coda.Header {
	return coda.Header{
		Class:      n.class,
		Tag:        n.tag,
		IsCompound: n.constructed,
	}
}

func (n *Node) IsPrimitive() bool {
	return !(n.constructed)
}

func (n *Node) IsConstructed() bool {
	return (n.constructed)
}

func (n *Node) setHeader(h coda.Header) error {
	*n = Node{
		class:       h.Class,
		tag:         h.Tag,
		constructed: h.IsCompound,
	}
	return nil
}

func (n *Node) checkHeader(h coda.Header) error {
	k := n.getHeader()
	if !coda.EqualHeaders(k, h) {
		return errors.New("der: invalid header")
	}
	return nil
}

func EncodeNode(data []byte, n *Node) (rest []byte, err error) {

	header := n.getHeader()
	data, err = coda.EncodeHeader(data, &header)
	if err != nil {
		return nil, err
	}

	value, err := encodeValue(n)
	if err != nil {
		return nil, err
	}

	length := len(value)
	data, err = coda.EncodeLength(data, length)
	if err != nil {
		return nil, err
	}

	data = append(data, value...)
	return data, err
}

func DecodeNode(data []byte, n *Node) (rest []byte, err error) {

	var header coda.Header
	data, err = coda.DecodeHeader(data, &header)
	if err != nil {
		return nil, err
	}
	err = n.setHeader(header)
	if err != nil {
		return nil, err
	}

	var length int
	data, err = coda.DecodeLength(data, &length)
	if err != nil {
		return nil, err
	}
	if len(data) < length {
		return nil, errors.New("insufficient data length")
	}

	err = decodeValue(data[:length], n)
	if err != nil {
		return nil, err
	}

	rest = data[length:]

	return rest, nil
}

func encodeValue(n *Node) ([]byte, error) {
	if !n.constructed {
		return cloneBytes(n.data), nil
	}
	return encodeNodes(n.nodes)
}

func decodeValue(data []byte, n *Node) error {

	if !n.constructed {
		var err error
		n.data = cloneBytes(data)
		switch n.tag {
		case TAG_BOOLEAN:
			n.value, err = n.GetBool()
		case TAG_INTEGER:
			n.value, err = n.GetInt()
		case TAG_BIT_STRING:
			n.value, err = n.GetString()
		case TAG_OCTET_STRING:
			n.value, err = n.GetString()
		case TAG_NULL:
			n.value = nil
		case TAG_OID:
			n.value, err = n.GetOid()
		case TAG_REAL:
			n.value, err = n.GetReal()
		case TAG_ENUMERATED:
			n.value, err = n.GetEnumerated()
		case TAG_UTF8_STRING:
			n.value, err = n.GetString()
		case TAG_TIME:
			n.value, err = n.GetUTCTime()
		case TAG_IA5_STRING:
			n.value, err = n.GetString()
		case TAG_UTC_TIME:
			n.value, err = n.GetUTCTime()
		case TAG_GENERALIZED_TIME:
			n.value, err = n.GetGeneralizedTime()
		case TAG_BMP_STRING:
			n.value, err = n.GetString()
		default:
			err = errors.New("tag is not supported")
		}
		if err != nil {
			return err
		}
		return nil
	} else {

		ns, err := decodeNodes(data)
		if err != nil {
			return err
		}
		n.nodes = ns

		return nil
	}
}

//----------------------------------------------------------------------------

func (n *Node) SetNodes(ns []*Node) {
	n.constructed = true
	n.nodes = ns
}

func (n *Node) GetNodes() ([]*Node, error) {
	if !n.constructed {
		return nil, ErrNodeIsNotConstructed
	}
	return n.nodes, nil
}

func (n *Node) SetBool(b bool) {
	n.constructed = false
	n.data = boolEncode(b)
}

func (n *Node) GetBool() (bool, error) {
	if n.constructed {
		return false, ErrNodeIsConstructed
	}
	return boolDecode(n.data)
}

func (n *Node) SetInt(i int64) {
	n.constructed = false
	n.data = intEncode(i)
}

func (n *Node) GetInt() (int64, error) {
	if n.constructed {
		return 0, ErrNodeIsConstructed
	}
	return intDecode(n.data)
}

func (n *Node) SetUint(u uint64) {
	n.constructed = false
	n.data = uintEncode(u)
}

func (n *Node) GetUint() (uint64, error) {
	if n.constructed {
		return 0, ErrNodeIsConstructed
	}
	return uintDecode(n.data)
}

func (n *Node) SetBytes(bs []byte) {
	n.constructed = false
	n.data = bs
}

func (n *Node) GetBytes() ([]byte, error) {
	if n.constructed {
		return nil, ErrNodeIsConstructed
	}
	return n.data, nil
}

func (n *Node) SetString(s string) {
	n.constructed = false
	n.data = []byte(s)
}

func (n *Node) GetString() (string, error) {
	if n.constructed {
		return "", ErrNodeIsConstructed
	}
	if !utf8.Valid(n.data) {
		return "", errors.New("invalid utf8 string")
		//return "", errors.New("data is not utf-8 string")
	}
	return string(n.data), nil
}

func (n *Node) SetUTCTime(t time.Time) error {
	data, err := encodeUTCTime(t)
	if err != nil {
		return err
	}
	n.constructed = false
	n.data = data
	return nil
}

func (n *Node) GetUTCTime() (time.Time, error) {
	if n.constructed {
		return time.Time{}, ErrNodeIsConstructed
	}
	return decodeUTCTime(n.data)
}

func (n *Node) GetOid() (string, error) {
	if n.constructed {
		return "", ErrNodeIsConstructed
	}
	oids := make([]uint32, len(n.data)+2)
	//the first byte using: first_arc* 40+second_arc
	//the later , when highest bit is 1, will add to next to calc
	// https://msdn.microsoft.com/en-us/library/windows/desktop/bb540809(v=vs.85).aspx
	f := uint32(n.data[0])
	if f < 80 {
		oids[0] = f / 40
		oids[1] = f % 40
	} else {
		oids[0] = 2
		oids[1] = f - 80
	}
	var tmp uint32
	for i := 2; i <= len(n.data); i++ {
		f = uint32(n.data[i-1])
		//	fmt.Printf("f:0x%x\r\n", f)
		if f >= 0x80 {
			//		fmt.Printf("tmp<<8:0x%x +   (f&0x7f)0x%x\r\n", tmp<<8, (f & 0x7f))
			tmp = tmp<<7 + (f & 0x7f)
			//		fmt.Printf("tmp:0x%x\r\n", tmp)
		} else {
			oids[i] = tmp<<7 + (f & 0x7f)
			//		fmt.Printf("oids[i]:0x%x\r\n", oids[i])
			tmp = 0
		}
	}
	var buffer bytes.Buffer
	for i := 0; i < len(oids); i++ {
		if oids[i] == 0 {
			continue
		}
		buffer.WriteString(fmt.Sprint(oids[i]) + ".")
	}
	return buffer.String()[0 : len(buffer.String())-1], nil
}

func (n *Node) GetReal() (float64, error) {
	//https://github.com/guidoreina/asn1/blob/58d422657c0378218587c89647b771348e2f7d07/asn1/ber/common.cpp
	return 0, errors.New("not supported")
}

func (n *Node) GetEnumerated() (int64, error) {
	return n.GetInt()
}

func (n *Node) GetGeneralizedTime() (time.Time, error) {
	return n.GetUTCTime()
}
