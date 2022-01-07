package der

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"strconv"
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

	Data  []byte      `json:"-"`               // Primitive:   (isCompound = false)
	Value interface{} `json:"value,omitempty"` // Primitive:  int/bool/string/time... (isCompound = false)
	Nodes []*Node     `json:"nodes,omitempty"` // Constructed: (isCompound = true)
}

/*
func (n *Node) MarshalJSON() ([]byte, error) {
	if !n.constructed {
		jsonKey, jsonValue, err := n.toJsonKeyValue()
		if err != nil {
			fmt.Println("MarshalJSON fail:", err)
			return nil, err
		}
		ret := `{"` + jsonKey + `":"` + jsonValue + `"},`
		return []byte(ret), nil
	} else {
		for i, child := range n.Nodes {
			if b, err := child.MarshalJSON(); err != nil {
				fmt.Println("MarshalJSON child fail:", i, err)
				return nil, err
			}
		}
	}
	return []byte(strconv.FormatFloat(float64(f), 'f', -1, 64)), nil
}
*/
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

func (n *Node) GetClass() int {
	return n.class
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

func (n *Node) toJsonKeyValue() (jsonKey, jsonValue string, err error) {

	jsonKey = tagName(n.tag)
	value := n.Value
	switch n.tag {
	case TAG_BOOLEAN:
		if b, ok := value.(bool); ok {
			jsonValue = strconv.FormatBool(b)
		} else {
			err = errors.New("data is not bool")
		}

	case TAG_ENUMERATED:
		fallthrough
	case TAG_INTEGER:
		// some interger is too big, so use string
		if s, p := value.(big.Int); p {
			jsonValue = s.Text(16)
		} else {
			err = errors.New("data is not int")
		}
	case TAG_REAL:
		if f, p := value.(float32); p {
			jsonValue = strconv.FormatFloat(float64(f), 'f', -1, 32)
		} else {
			err = errors.New("data is not real")
		}

		if f, p := value.(float64); p {
			jsonValue = strconv.FormatFloat(f, 'f', -1, 32)
		} else {
			err = errors.New("data is not real")
		}

	case TAG_BIT_STRING:
		fallthrough
	case TAG_OCTET_STRING:
		if f, p := value.([]byte); p {
			jsonValue = printBytes(f)
		} else {
			err = errors.New("data is not bytes")
		}
	case TAG_BMP_STRING:
		fallthrough
	case TAG_OID:
		fallthrough
	case TAG_UTF8_STRING:
		fallthrough
	case TAG_NUMBERIC_STRING:
		fallthrough
	case TAG_PRINTABLE_STRING:
		fallthrough
	case TAG_T61_STRING:
		fallthrough
	case TAG_VIDEOTEX_STRING:
		fallthrough
	case TAG_IA5_STRING:
		if s, ok := value.(string); ok {
			jsonValue = s
		} else {
			err = errors.New("data is not string")
		}

	case TAG_TIME:
		fallthrough
	case TAG_UTC_TIME:
		fallthrough
	case TAG_GENERALIZED_TIME:
		if t, ok := value.(time.Time); ok {
			jsonValue = t.Local().Format("2006-01-02 15:04:05 UTC")
		} else {
			err = errors.New("data is not time")
		}
	case TAG_END_OF_CONTENT:
		jsonValue = ""
	case TAG_NULL:
		jsonValue = ""
	default:
		err = errors.New("tag is not supported")
	}
	if err != nil {
		return "", "", err
	}
	//fmt.Println("k:", k, "   v:", v)
	return jsonKey, jsonValue, nil
}

func valueToString(tag int, value interface{}) (string, error) {
	var k, v string
	var err error
	switch tag {
	case TAG_BOOLEAN:
		k = "bool"
		if b, ok := value.(bool); ok {
			v = strconv.FormatBool(b)
		} else {
			err = errors.New("data is not bool")
		}

	case TAG_ENUMERATED:
		fallthrough
	case TAG_INTEGER:
		// some interger is too big, so use string
		k = "int"
		if s, p := value.(big.Int); p {
			v = s.Text(16)
		} else {
			err = errors.New("data is not int")
		}
	case TAG_REAL:
		k = "float"
		if f, p := value.(float32); p {
			v = strconv.FormatFloat(float64(f), 'f', -1, 32)
		} else {
			err = errors.New("data is not real")
		}

		if f, p := value.(float64); p {
			v = strconv.FormatFloat(f, 'f', -1, 32)
		} else {
			err = errors.New("data is not real")
		}

	case TAG_BIT_STRING:
		fallthrough
	case TAG_OCTET_STRING:
		k = "bytes"
		if f, p := value.([]byte); p {
			v = printBytes(f)
		} else {
			err = errors.New("data is not bytes")
		}
	case TAG_BMP_STRING:
		fallthrough
	case TAG_OID:
		fallthrough
	case TAG_UTF8_STRING:
		fallthrough
	case TAG_NUMBERIC_STRING:
		fallthrough
	case TAG_PRINTABLE_STRING:
		fallthrough
	case TAG_T61_STRING:
		fallthrough
	case TAG_VIDEOTEX_STRING:
		fallthrough
	case TAG_IA5_STRING:
		k = "string"
		if s, ok := value.(string); ok {
			v = s
		} else {
			err = errors.New("data is not string")
		}

	case TAG_TIME:
		fallthrough
	case TAG_UTC_TIME:
		fallthrough
	case TAG_GENERALIZED_TIME:
		k = "time"
		if t, ok := value.(time.Time); ok {
			v = t.Local().Format("2006-01-02 15:04:05 UTC")
		} else {
			err = errors.New("data is not time")
		}
	case TAG_END_OF_CONTENT:
		k = "end"
		v = ""
	case TAG_NULL:
		k = "nil"
		v = ""
	default:
		err = errors.New("tag is not supported")
	}
	if err != nil {
		return "", err
	}
	//fmt.Println("k:", k, "   v:", v)
	return "{" + k + ":" + v + "}", nil
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
		fmt.Println("DecodeNode(): DecodeHeader fail:", data, err)
		return nil, err
	}
	err = n.setHeader(header)
	if err != nil {
		fmt.Println("DecodeNode(): setHeader fail:", err)
		return nil, err
	}

	var length int
	data, err = coda.DecodeLength(data, &length)
	if err != nil {
		fmt.Println("DecodeNode(): DecodeLength fail:", err)
		return nil, err
	}
	if len(data) < length {
		fmt.Println("DecodeNode():len(data) < length fail:")
		return nil, errors.New("insufficient data length")
	}

	err = decodeValue(data[:length], n)
	if err != nil {
		fmt.Println("DecodeNode(): decodeValue fail:", err)
		return nil, err
	}

	rest = data[length:]
	return rest, nil
}

func encodeValue(n *Node) ([]byte, error) {
	if !n.constructed {
		return cloneBytes(n.Data), nil
	}
	return encodeNodes(n.Nodes)
}

func decodeValue(data []byte, n *Node) error {

	if !n.constructed {
		var err error
		n.Data = cloneBytes(data)
		switch n.tag {
		case TAG_END_OF_CONTENT:
			n.Value = nil
			fmt.Println("decodeValue(): TAG_END_OF_CONTENT ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_BOOLEAN:
			n.Value, err = n.GetBool()
			fmt.Println("decodeValue(): TAG_BOOLEAN ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_INTEGER:

			/*
				n.Value, err = n.GetInt()
				var ret uint64
				for _, i := range data {
					ret = ret * 256
					ret = ret + uint64(i)
				}
				n.Value = ret
			*/
			// bigint

			b := new(big.Int).SetBytes(data)
			fmt.Println("big.Int", b)
			n.Value = *b
			fmt.Println("decodeValue(): TAG_INTEGER ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)

		case TAG_BIT_STRING:
			n.Value = n.Data
			fmt.Println("decodeValue(): TAG_BIT_STRING ", "  tag:", n.tag, "   len(data):", len(n.Data))
		case TAG_OCTET_STRING:
			n.Value = n.Data
			fmt.Println("decodeValue(): TAG_OCTET_STRING ", "  tag:", n.tag, "   len(data):", len(n.Data))
		case TAG_NULL:
			n.Value = nil
			fmt.Println("decodeValue(): TAG_NULL ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_OID:
			n.Value, err = n.GetOid()
			fmt.Println("decodeValue(): TAG_OID ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_REAL:
			n.Value, err = n.GetReal()
			fmt.Println("decodeValue(): TAG_REAL ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_ENUMERATED:
			n.Value, err = n.GetEnumerated()
			fmt.Println("decodeValue(): TAG_ENUMERATED ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_UTF8_STRING:
			n.Value, err = n.GetString()
			fmt.Println("decodeValue(): TAG_UTF8_STRING ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_TIME:
			n.Value, err = n.GetUTCTime()
			fmt.Println("decodeValue(): TAG_TIME ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_NUMBERIC_STRING:
			n.Value, err = n.GetString()
			fmt.Println("decodeValue(): TAG_NUMBERIC_STRING ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_PRINTABLE_STRING:
			n.Value, err = n.GetString()
			fmt.Println("decodeValue(): TAG_PRINTABLE_STRING ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_T61_STRING:
			n.Value, err = n.GetString()
			fmt.Println("decodeValue(): TAG_T61_STRING ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_VIDEOTEX_STRING:
			n.Value, err = n.GetString()
			fmt.Println("decodeValue(): TAG_VIDEOTEX_STRING ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_IA5_STRING:
			n.Value, err = n.GetString()
			fmt.Println("decodeValue(): TAG_IA5_STRING ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_UTC_TIME:
			n.Value, err = n.GetUTCTime()
			fmt.Println("decodeValue(): TAG_UTC_TIME ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_GENERALIZED_TIME:
			n.Value, err = n.GetGeneralizedTime()
			fmt.Println("decodeValue(): TAG_GENERALIZED_TIME ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		case TAG_BMP_STRING:
			n.Value, err = n.GetString()
			fmt.Println("decodeValue(): TAG_BMP_STRING ", "  tag:", n.tag, "   data:", printBytes(n.Data), n.Value)
		default:
			err = errors.New("tag is not supported")
		}
		if err != nil {
			fmt.Println("decodeValue(): fail, tag:", n.tag, "  data:", printBytes(n.Data), err)
			return err
		}
		return nil
	} else {

		ns, err := decodeNodes(data)
		if err != nil {
			return err
		}
		n.Nodes = ns

		return nil
	}
}

//----------------------------------------------------------------------------

func (n *Node) SetNodes(ns []*Node) {
	n.constructed = true
	n.Nodes = ns
}

func (n *Node) GetNodes() ([]*Node, error) {
	if !n.constructed {
		return nil, ErrNodeIsNotConstructed
	}
	return n.Nodes, nil
}

func (n *Node) SetBool(b bool) {
	n.constructed = false
	n.Data = boolEncode(b)
}

func (n *Node) GetBool() (bool, error) {
	if n.constructed {
		return false, ErrNodeIsConstructed
	}
	return boolDecode(n.Data)
}

func (n *Node) SetInt(i int64) {
	n.constructed = false
	n.Data = intEncode(i)
}

func (n *Node) GetInt() (int64, error) {
	if n.constructed {
		return 0, ErrNodeIsConstructed
	}
	return intDecode(n.Data)
}

func (n *Node) SetUint(u uint64) {
	n.constructed = false
	n.Data = uintEncode(u)
}

func (n *Node) GetUint() (uint64, error) {
	if n.constructed {
		return 0, ErrNodeIsConstructed
	}
	return uintDecode(n.Data)
}

func (n *Node) SetBytes(bs []byte) {
	n.constructed = false
	n.Data = bs
}

func (n *Node) GetBytes() ([]byte, error) {
	if n.constructed {
		return nil, ErrNodeIsConstructed
	}
	return n.Data, nil
}

func (n *Node) SetString(s string) {
	n.constructed = false
	n.Data = []byte(s)
}

func (n *Node) GetString() (string, error) {
	if n.constructed {
		return "", ErrNodeIsConstructed
	}
	if !utf8.Valid(n.Data) {
		return "", errors.New("invalid utf8 string")
		//return "", errors.New("data is not utf-8 string")
	}
	return string(n.Data), nil
}

func (n *Node) SetUTCTime(t time.Time) error {
	data, err := encodeUTCTime(t)
	if err != nil {
		return err
	}
	n.constructed = false
	n.Data = data
	return nil
}

func (n *Node) GetUTCTime() (time.Time, error) {
	if n.constructed {
		return time.Time{}, ErrNodeIsConstructed
	}
	return decodeUTCTime(n.Data)
}

func (n *Node) GetOid() (string, error) {
	if n.constructed {
		return "", ErrNodeIsConstructed
	}
	oids := make([]uint32, len(n.Data)+2)
	//the first byte using: first_arc* 40+second_arc
	//the later , when highest bit is 1, will add to next to calc
	// https://msdn.microsoft.com/en-us/library/windows/desktop/bb540809(v=vs.85).aspx
	f := uint32(n.Data[0])
	if f < 80 {
		oids[0] = f / 40
		oids[1] = f % 40
	} else {
		oids[0] = 2
		oids[1] = f - 80
	}
	var tmp uint32
	for i := 2; i <= len(n.Data); i++ {
		f = uint32(n.Data[i-1])
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
	return parseGeneralizedTime(n.Data)
}
