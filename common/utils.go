package common

import (
	"bytes"
	"fmt"
	"strings"
	"encoding/xml"
)

func StringtoASCII(loadinfo string) (bool, bytes.Buffer) {
	var reinfo bytes.Buffer
	var outinfo []byte = make([]byte, 1)
	loadinfoByte := strings.Split(loadinfo, " ")
	if (loadinfoByte[0] == "") {
		fmt.Printf("文件是空的")
		reinfo.Write([]byte(""))
		return false, reinfo
	}

	lenforload := len(loadinfoByte)
	fmt.Printf("byte number is: %v\n", lenforload)
	for i := 0; i < lenforload; i += 1 {
		fmt.Printf(loadinfoByte[i])
		single := []byte(loadinfoByte[i])
		if len(single) <= 0 {
			continue
		}
		single[1] = ASCIItoBi(single[1])
		single[0] = ASCIItoBi(single[0])
		outinfo[0] = single[0]*16 + single[1]
		fmt.Printf("append byte  is: %v\n", outinfo)
		reinfo.Write(outinfo)
	}

	reinfo.Write([]byte("#"))
	//fmt.Println(reinfo.String())
	return true, reinfo

}

func ASCIItoBi(IN byte) byte {
	if (IN >= 48 && IN <= 57) {
		return (IN - 48)
	} else if (IN >= 97 && IN <= 102) {
		return (IN - 97 + 10)
	} else if (IN >= 65 && IN <= 70) {
		return (IN - 65 + 10)
	} else {
		return 0xFF
	}
}

func GetXMLanswer(receiveXML string) string {
	var t xml.Token
	var err error
	inputReader := strings.NewReader(receiveXML)

	// 从文件读取，如可以如下：
	// content, err := ioutil.ReadFile("studygolang.xml")
	// decoder := xml.NewDecoder(bytes.NewBuffer(content))

	decoder := xml.NewDecoder(inputReader)
	var TYPEflag bool
	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
	switch token := t.(type) {
	// 处理元素开始（标签）
	case xml.StartElement:
		name := token.Name.Local
		fmt.Printf("Token name: %s\n", name)
		if ("Type" == name) {
		TYPEflag = true
	}
		for _, attr := range token.Attr {
		attrName := attr.Name.Local
		attrValue := attr.Value
		fmt.Printf("An attribute is: %s %s\n", attrName, attrValue)
	}

		// 处理元素结束（标签）
	case xml.EndElement:
		fmt.Printf("Token of '%s' end\n", token.Name.Local)
		// 处理字符数据（这里就是元素的文本）
	case xml.CharData:
		if (TYPEflag) {
		content := string([]byte(token))
		fmt.Printf("This is the content: %v\n", content)
		if ("1" == content) {
		return "control58Res"
	}
	}

		TYPEflag = false
	default:
		// ...
	}
	}
		return ""
	}
