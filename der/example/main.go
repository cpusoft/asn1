package main

import (
	"bytes"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cpusoft/asn1/der"
	"github.com/cpusoft/goutil/fileutil"
	"github.com/cpusoft/goutil/jsonutil"
)

func main() {
	//fn := derFile
	fn := derHex
	//fn := testUint64
	//fn := testIntDER
	//fn := testIntJSON
	//fn := testPersone
	//fn := testFloat

	if err := fn(); err != nil {
		fmt.Println(err)
	}
}
func derFile() error {
	//file := `apnic-rpki-root-iana-origin.cer`
	file := os.Args[1]
	data1, err := fileutil.ReadFileToBytes(file)
	fmt.Println("file:", file, "  len(data1):", len(data1), err)
	n := new(der.Node)

	_, err = der.DecodeNode(data1, n)
	if err != nil {
		fmt.Println("derFile():  DecodeNode fail:", err)
		return err
	}

	s, err := der.ConvertToString(n)
	if err != nil {
		fmt.Println("derFile():  ConvertToString fail:", err)
		return err
	}

	fmt.Println(s)
	return nil
}

func derHex() error {

	hexDump := `30-2E-A0-03-02-01-01-A1 03-02-01-01-A2-03-02-01
01-A3-08-0C-06-31-32-33 34-35-36-A4-13-17-11-31
35-31-32-31-37-31-37-34 38-34-34-2B-30-33-30-30`

	// checklist .sig
	hexDump = `30 81 9c 30 14 a1 12 30 10 30 0e 04 01 02 30 09 03 07 00 20 01 06 7c 20 8c 30 0b 06 09 60 86 48 01 65 03 04 02 01 30 77 30 34 16 10 62 34 32 5f 69 70 76 36 5f 6c 6f 61 2e 70 6e 67 04 20 95 16 dd 64 be 7c 17 25 b9 fc a1 17 12 0e 58 e8 d8 42 a5 20 68 73 39 9b 3d df fc 91 c4 b6 ac f0 30 3f 16 1b 62 34 32 5f 73 65 72 76 69 63 65 5f 64 65 66 69 6e 69 74 69 6f 6e 2e 6a 73 6f 6e 04 20 0a e1 39 47 22 00 5c d9 2f 4c 6a a0 24 d5 d6 b3 e2 e6 7d 62 9f 11 72 0d 94 78 a6 33 a1 17 a1 c7`

	// mft
	//hexDump = `30 82 01 33 02 14 01 0d 0c 9f 43 28 57 6d 51 cc 73 c0 42 cf c1 73 e8 5c ae a2 18 0f 32 30 32 30 30 38 31 33 31 35 30 37 33 34 5a 18 0f 32 30 32 30 30 38 31 36 31 37 30 30 30 30 5a 06 09 60 86 48 01 65 03 04 02 01 30 81 ed 30 4d 16 28 34 31 31 30 35 31 64 61 2d 63 30 32 65 2d 33 39 32 34 2d 62 35 31 32 2d 38 62 66 35 36 63 39 30  36 66 36 38 2e 72 6f 61 03 21 00 4c ed d8 ef 52 25 65 fc d1 2b 94 34 1a 79 ee 50 22 a7 dd 39 c7 9a 37 ff 97 f5 75 ce b3 05 21 68 30 4d 16 28 37 62 61 35 39 33 61 66 2d 36 66 39 35 2d 33 37 30 66 2d 38 31 63 38 2d 34 62 35 62 61 64 36 61 30 34 36 34 2e 72 6f 61 03 21 00 ea 36 17 ae 1a a7 80 38 d7 5a  21 20 19 18 0e 66 71 ee d8 18 ce c5 33 4b e5 87 a4 f0 76 b5 27 d7 30 4d 16 28 64 65 35 31 66 64 37 64 2d 33 37 64 37 2d 34 64 61 36 2d 38 64 65 38 2d 37 64 33 34 36 65 30 30 30 36 61 61 2e 63 72 6c 03 21 00 48 78 9d e1 96 e6 10 d7 9b 74 69 c1 1f 8b d1 fb 64 5c 42 46 c6 99 fd 1d 31 14 41 46 05 b1 06 34`
	s := onlyHex(hexDump)

	data1, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	n := new(der.Node)

	_, err = der.DecodeNode(data1, n)
	if err != nil {
		return err
	}

	s = jsonutil.MarshalJson(n)
	fmt.Println("Json:\n" + s)

	n, _ = changeNode(n)
	s = jsonutil.MarshalJson(n)
	fmt.Println("changeNode Json:\n" + s)
	/*
	   {
	       "nodes": [
	           {
	               "nodes": [
	                   {
	                       "nodes": [
	                           {
	                               "nodes": [
	                                   {
	                                       "nodes": [
	                                           {
	                                               "value": "Ag=="
	                                           },
	                                           {
	                                               "nodes": [
	                                                   {
	                                                       "value": "ACABBnwgjA=="
	                                                   }
	                                               ]
	                                           }
	                                       ]
	                                   }
	                               ]
	                           }
	                       ]
	                   }
	               ]
	           },
	           {
	               "nodes": [
	                   {
	                       "value": "2.16.840.1.101.3.4.2.1"
	                   }
	               ]
	           },
	           {
	               "nodes": [
	                   {
	                       "nodes": [
	                           {
	                               "value": "b42_ipv6_loa.png"
	                           },
	                           {
	                               "value": "lRbdZL58FyW5/KEXEg5Y6NhCpSBoczmbPd/8kcS2rPA="
	                           }
	                       ]
	                   },
	                   {
	                       "nodes": [
	                           {
	                               "value": "b42_service_definition.json"
	                           },
	                           {
	                               "value": "CuE5RyIAXNkvTGqgJNXWs+LmfWKfEXINlHimM6EXocc="
	                           }
	                       ]
	                   }
	               ]
	           }
	       ]
	   }
	*/
	/*
		checkList := CheckList{}
		err = jsonutil.UnmarshalJson(s, &checkList)
		fmt.Println(checkList, err)

			s, err = der.ConvertToString(n)
			if err != nil {
				return err
			}

			fmt.Println(s)
			s = strings.TrimSpace(s)
			s0 := strings.Replace(s, "],],],]", "]]]]", -1)
			s1 := strings.Replace(s0, "],],]", "]]]", -1)
			s2 := strings.Replace(s1, "],]", "]]", -1)
			fmt.Println("s2:\n" + s2)
			s3 := strings.Replace(s2, `"},[`, `"}],[`, -1)
			fmt.Println("s3:\n" + s3)
			checkList := CheckList{}
			err = jsonutil.UnmarshalJson(s3, &checkList)
			fmt.Println(checkList, err)
	*/
	/*
		checklist
		[
			[
				[
					[
										[
											{
												"tagOctetString": "02"
											},
											[
												{
													"tagBitString": "00 20 01 06 7c 20 8c"
												}
											]
										]
									]
								]
							],
							[
								{
									"tagOid": "2.16.840.1.101.3.4.2.1"
								}
							],
							[
								[
									{
										"tagIa5String": "b42_ipv6_loa.png"
									},
									{
										"tagOctetString": "95 16 dd 64 be 7c 17 25 b9 fc a1 17 12 0e 58 e8 d8 42 a5 20 68 73 39 9b 3d df fc 91 c4 b6 ac f0"
									}
								],
								[
									{
										"tagIa5String": "b42_service_definition.json"
									},
									{
										"tagOctetString": "0a e1 39 47 22 00 5c d9 2f 4c 6a a0 24 d5 d6 b3 e2 e6 7d 62 9f 11 72 0d 94 78 a6 33 a1 17 a1 c7"
									}
								]
							]
						]

					mft
					[
					{
						"tagInteger": "10d0c9f4328576d51cc73c042cfc173e85caea2"
					},
					{
						"tagGeneralizedTime": "2020-08-13 23:07:34 UTC"
					},
					{
						"tagGeneralizedTime": "2020-08-17 01:00:00 UTC"
					},
					{
						"tagOid": "2.16.840.1.101.3.4.2.1"
					},
					[
						[
							{
								"tagIa5String": "411051da-c02e-3924-b512-8bf56c906f68.roa"
							},
							{
								"tagBitString": "00 4c ed d8 ef 52 25 65 fc d1 2b 94 34 1a 79 ee 50 22 a7 dd 39 c7 9a 37 ff 97 f5 75 ce b3 05 21 68"
							}
						],
						[
							{
								"tagIa5String": "7ba593af-6f95-370f-81c8-4b5bad6a0464.roa"
							},
							{
								"tagBitString": "00 ea 36 17 ae 1a a7 80 38 d7 5a 21 20 19 18 0e 66 71 ee d8 18 ce c5 33 4b e5 87 a4 f0 76 b5 27 d7"
							}
						],
						[
							{
								"tagIa5String": "de51fd7d-37d7-4da6-8de8-7d346e0006aa.crl"
							},
							{
								"tagBitString": "00 48 78 9d e1 96 e6 10 d7 9b 74 69 c1 1f 8b d1 fb 64 5c 42 46 c6 99 fd  1d 31 14 41 46 05 b1 06 34"
							}
						]
					]
				]


	*/

	data2, err := der.EncodeNode(nil, n)
	if err != nil {
		return err
	}

	fmt.Printf("equal: %t\n", bytes.Equal(data1, data2))

	return nil
}

func byteIsHex(b byte) bool {

	if (b >= '0') && (b <= '9') {
		return true
	}

	if (b >= 'a') && (b <= 'f') {
		return true
	}

	if (b >= 'A') && (b <= 'F') {
		return true
	}

	return false
}

func onlyHex(s string) string {

	data := []byte(s)

	var res []byte
	for _, b := range data {
		if byteIsHex(b) {
			res = append(res, b)
		}
	}

	return string(res)
}

type uint64Sample struct {
	val  uint64
	data []byte
}

func newUint64Sample(v uint64, s string) *uint64Sample {
	s = onlyHex(s)
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err.Error())
	}
	return &uint64Sample{v, data}
}

func testUint64() error {

	var as = []uint64{0, 1, 2}

	var x uint64 = 4
	for i := 0; i < 64; i++ {

		as = append(as, x-1)
		as = append(as, x)
		as = append(as, x+1)

		x *= 2
	}

	for _, a := range as {

		data, err := der.Marshal(a)
		if err != nil {
			return err
		}

		fmt.Printf("newUint64Sample(%d, \"%X\"),\n", a, data)
	}

	return nil
}

type int64Sample struct {
	val  int64
	data []byte
}

func newInt64Sample(v int64, s string) *int64Sample {
	s = onlyHex(s)
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err.Error())
	}
	return &int64Sample{v, data}
}

func testInt64() error {

	var as = []int64{0, 1, -1, 2, -2}

	var x int64 = 4
	for i := 0; i < 64; i++ {

		as = append(as, x-1)
		as = append(as, -(x - 1))
		as = append(as, x)
		as = append(as, -x)
		as = append(as, x+1)
		as = append(as, -(x + 1))

		x *= 2
	}

	for _, a := range as {

		data, err := der.Marshal(a)
		if err != nil {
			return err
		}

		fmt.Printf("newInt64Sample(%d, \"%X\"),\n", a, data)
	}

	return nil
}

func testIntDER() error {

	var a int64 = -100000

	data, err := der.Marshal(a)
	if err != nil {
		return err
	}

	fmt.Printf("%X\n", data)

	var b int32

	err = der.Unmarshal(data, &b)
	if err != nil {
		return err
	}

	fmt.Println(b)

	return nil
}

func testIntJSON() error {

	var a int = -108987

	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	var b int

	err = json.Unmarshal(data, &b)
	if err != nil {
		return err
	}

	fmt.Println(b)

	return nil
}

type Persone struct {
	Name string `asn1:"tag:0" der:"tag:0"`
	Age  int    `asn1:"tag:1" der:"tag:1"`
	Desc string `asn1:"tag:2" der:"tag:2,optional"`
}

func newString(s string) *string {
	return &s
}

func testPersone() error {

	a := Persone{
		Name: "John",
		Age:  -97,
		Desc: "Ароза упала на лапу Азора",
	}

	data, err := der.Marshal(&a)
	if err != nil {
		return err
	}

	fmt.Printf("%X\n", data)

	var b Persone

	err = der.Unmarshal(data, &b)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", b)
	fmt.Println("desc:", b.Desc)

	data, err = asn1.Marshal(a)
	if err != nil {
		return err
	}
	fmt.Printf("%X\n", data)

	return nil
}

func testFloat() error {

	var a float64 = 3.14

	data, err := der.Marshal(a)
	//data, err := asn1.Marshal(a)
	if err != nil {
		return err
	}
	fmt.Printf("%X\n", data)

	return nil
}
func printBytes(data []byte) (ret string) {
	for _, b := range data {
		ret += fmt.Sprintf("%02x ", b)
	}
	return ret
}

type CheckList struct {
	Node1s []Node1 `json:"nodes"`
}

type Node1 struct {
	Node2s         []Node2          `json:"nodes"`
	Oid            Oid              `json:"-"`
	CheckListBlock []CheckListBlock `json:"-"`
}
type Node2 struct {
	Node3s []Node3 `json:"nodes"`
}
type Node3 struct {
	Node4s []Node4 `json:"nodes"`
}
type Node4 struct {
	Node5s []Node5 `json:"nodes"`
}
type Node5 struct {
	IpFamliyBytes []byte `json:"value"`
	Node6s        Node6  `json:"nodes"`
}
type Node6 struct {
	IpAddresses []byte `json:"value"` //[]IpAddress
}

type IpAddress struct {
	IpAddressPrefix []byte `json:"value"`
}
type Oid struct {
	Oid string `json:"value"`
}

type NameAndHash struct {
	Name string `json:"value"`
}

func changeNode(n *der.Node) (newNode *der.Node, err error) {
	ipNode := n.Nodes[0].Nodes[0].Nodes[0].Nodes[0]
	ipFamliy := ipNode.Nodes[0]
	ipAddress := ipNode.Nodes[1].Nodes[0]
	oidNode := n.Nodes[1].Nodes[0]
	fileHashNode := n.Nodes[2]
	fmt.Println(jsonutil.MarshalJson(ipFamliy))
	fmt.Println(jsonutil.MarshalJson(ipAddress))

	fmt.Println(jsonutil.MarshalJson(oidNode))
	for _, child := range fileHashNode.Nodes {
		fmt.Println(jsonutil.MarshalJson(child.Nodes[0]))
		fmt.Println(jsonutil.MarshalJson(child.Nodes[1]))
		nameAndHashStr := jsonutil.MarshalJson(child.Nodes)
		fmt.Println(nameAndHashStr)
		nameAndHash := make([]NameAndHash, 0)
		err = jsonutil.UnmarshalJson(nameAndHashStr, &nameAndHash)
		fmt.Println(nameAndHash, err)
	}

	return n, nil
}
