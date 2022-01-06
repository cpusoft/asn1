package der

import (
	"bytes"
	"fmt"
	"strings"
)

func ConvertToString(n *Node) (string, error) {
	var buf bytes.Buffer
	//err := nodeToString(n, &buf, 0)
	err := nodeToJson(n, &buf)
	if err != nil {
		return "", err
	}
	str := strings.Replace(buf.String(), "},]", "}]", -1)
	return str, nil
}
func nodeToJson(n *Node, buf *bytes.Buffer) (err error) {

	if !n.constructed {

		jsonKey, jsonValue, err := n.toJsonKeyValue()
		if err != nil {
			return err
		}
		fmt.Println(jsonKey, jsonValue)
		//s = hex.EncodeToString(n.data)
		if _, err = buf.WriteString(`{"` + jsonKey + `":"` + jsonValue + `"},`); err != nil {
			return err
		}

	} else {

		buf.WriteString("[")

		for _, child := range n.nodes {
			if err = nodeToJson(child, buf); err != nil {
				return err
			}
		}

		buf.WriteString("]")
	}

	return nil
}

func nodeToString(n *Node, buf *bytes.Buffer, indent int) error {

	indentBuff := make([]byte, indent)
	for i := 0; i < indent; i++ {
		indentBuff[i] = '\t'
	}

	_, err := buf.Write(indentBuff)
	if err != nil {
		return err
	}

	className := classShortName(n.class)
	s := fmt.Sprintf("%s(%d):", className, n.GetTag())
	if _, err = buf.WriteString(s); err != nil {
		return err
	}

	if !n.constructed {

		buf.WriteByte(' ')

		s, err = valueToString(n.tag, n.value)
		if err != nil {
			return err
		}
		//s = hex.EncodeToString(n.data)
		if _, err = buf.WriteString(s); err != nil {
			return err
		}

		buf.WriteByte('\n')

	} else {

		buf.WriteString(" {\n")

		for _, child := range n.nodes {
			if err = nodeToString(child, buf, indent+1); err != nil {
				return err
			}
		}

		buf.Write(indentBuff)
		buf.WriteString("}\n")
	}

	return nil
}
