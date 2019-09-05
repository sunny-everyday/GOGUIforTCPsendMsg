package gui

import (
	"encoding/xml"
	"fmt"
	"github.com/lxn/walk"
	"net"
	"strings"
)

type Condom struct {
	Index       int
	Name        string
	Type        string
	MessageInfo string
	Checked     bool
}


//message number
var Osmessagenumber int = 0
var Rsmessagenumber int = 0

var Onlysendmessagelist [100]string
var Readsendmessagelist [100]string

var Onlysendmessage [100]string
var Readsendmessage [100]string


type CondomMainWindow struct {
	*walk.MainWindow
	Model            *CondomModel
	Tv               *walk.TableView
	Messageforsocket *walk.TextEdit
}

func NewCondomMainWindow() *CondomMainWindow {
	mainWin := new(CondomMainWindow)
	//init model
	mainWin.Model = NewCondomModel()
	//init main gui

	return mainWin
}

func (m *CondomMainWindow) ResetRows() {
	items := []*Condom{}

	for i := 0; i < Osmessagenumber; i++ {
		x := &Condom{
			Index:       i + 1,
			Name:        Onlysendmessagelist[i],
			Type:        "send",
			MessageInfo: Onlysendmessage[i],
		}
		items = append(items, x)
	}

	for j := Osmessagenumber; j < Osmessagenumber+Rsmessagenumber; j++ {
		x := &Condom{
			Index:       j + 1,
			Name:        Readsendmessagelist[j-Osmessagenumber],
			Type:        "receive",
			MessageInfo: Readsendmessage[j-Osmessagenumber],
		}
		items = append(items, x)
	}
	m.Model.Items = items
	m.Model.PublishRowsReset()
	m.Tv.SetSelectedIndexes([]int{})
}


func (mw *CondomMainWindow) tv_ItemActivated() {
	msg := ``
	for _, i := range mw.Tv.SelectedIndexes() {
		msg = msg + "\n" + mw.Model.Items[i].Name + ":" + mw.Model.Items[i].MessageInfo
	}
	walk.MsgBox(mw, "title", msg, walk.MsgBoxIconInformation)
}

func (mw *CondomMainWindow) TcpClientReadandSend(ch chan bool, lnkclient net.Conn, lnkclientconnectFlag *bool) {
	var disconnectflag bool
	for {
		disconnectflag = <-ch
		if (disconnectflag) {
			lnkclient.Close()
			*lnkclientconnectFlag = false
			return
		}
		var buf [2000]byte
		n, err := lnkclient.Read(buf[:])

		if err != nil {
			fmt.Printf("read from connect failed, err: %v\n", err)
			break
		}
		str := string(buf[:n])
		fmt.Printf("receive from client, data: %v\n", str)
		mw.Messageforsocket.SetText("接收消息:" + "\r\n" + str + "\r\n")
		sendmessageName := mw.GetXMLanswer(str)
		if ("" != sendmessageName) {
			for _, x := range mw.Model.Items {
				if x.Name == sendmessageName {
					var err error
					_, err = lnkclient.Write([]byte(x.MessageInfo))
					if err != nil {
						fmt.Printf("write failed , err : %v\n", err)
						break
					}
				}
			}
		}
	}
}

func (mw *CondomMainWindow) GetXMLanswer(receiveXML string) string {
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
