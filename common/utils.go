package common

import (
	"bytes"
	"fmt"
	"strings"
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

