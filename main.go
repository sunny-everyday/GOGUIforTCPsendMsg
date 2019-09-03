// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by  a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
)

import (
	
	"encoding/xml"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"net"
	"path/filepath"

	"strconv"
	"strings"
)

type Condom struct {
	Index       int
	Name        string
	Type        string
	MessageInfo string
	checked     bool
}

type CondomModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*Condom
}

func (m *CondomModel) RowCount() int {
	return len(m.items)
}

func (m *CondomModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index
	case 1:
		return item.Name
	case 2:
		return item.Type
	case 3:
		return item.MessageInfo

	}
	panic("unexpected col")
}

func (m *CondomModel) Checked(row int) bool {
	return m.items[row].checked
}
func (m *CondomModel) GetCheckedItemlist() []int {
	var j int
	var checkedlist []int
	length := len(m.items)
	for row := 0; row < length; row++ {
		if m.Checked(row) {
			checkedlist[j] = row
		}
	}
	return checkedlist
}

func (m *CondomModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked
	return nil
}

func (m *CondomModel) Len() int {
	return len(m.items)
}

func (m *CondomModel) Swap(i, j int) {
	m.items[i], m.items[j] = m.items[j], m.items[i]
}
func (m *CondomModel) FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(float64(input_num), 'f', 6, 64)
}

//message number
var osmessagenumber int = 0
var rsmessagenumber int = 0

var onlysendmessagelist [100]string
var readsendmessagelist [100]string

var onlysendmessage [100]string
var readsendmessage [100]string

func NewCondomModel() *CondomModel {
	m := new(CondomModel)
	m.items = make([]*Condom, 0)

	return m
}

type CondomMainWindow struct {
	*walk.MainWindow
	model *CondomModel
	tv    *walk.TableView
	messageforsocket *walk.TextEdit
}

func (m *CondomMainWindow) ResetRows() {
	items := []*Condom{}


	for i := 0; i < osmessagenumber; i++ {
		x := &Condom{
			Index:       i + 1,
			Name:        onlysendmessagelist[i],
			Type:        "send",
			MessageInfo: onlysendmessage[i],
		}
		items = append(items, x)
	}

	for j := osmessagenumber; j < osmessagenumber+rsmessagenumber; j++ {
		x := &Condom{
			Index:       j + 1,
			Name:        readsendmessagelist[j-osmessagenumber],
			Type:        "receive",
			MessageInfo: readsendmessage[j-osmessagenumber],
		}
		items = append(items, x)
	}
	m.model.items = items
	m.model.PublishRowsReset()
	m.tv.SetSelectedIndexes([]int{})
}
func main() {
	mw := &CondomMainWindow{model: NewCondomModel()}
	var IP, Port,connectstatus,messageinfo *walk.TextEdit
	var lnkclient net.Conn
	var lnkclientconnectFlag bool = false
	var Hexflag  bool = true
	tcpdisconnect := make(chan bool)
	var db *walk.DataBinder
	MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Robot Simulator",
		Size:     Size{800, 500},
		Layout:   VBox{},
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "animal",
			DataSource:     Hexflag,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		Children: []Widget{
			TabWidget{
				Pages: []TabPage{
					{
						Title:  "Robot",
						Layout: VBox{},
						Children: []Widget{
							GroupBox{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									PushButton{
										Text: "Reload",
										OnClicked: func() {
											fpt, err := filepath.Abs("go_build_minfang827.exe")
											if err != nil {
												panic(err)
											}
											fmt.Println(fpt)
											way1 := strings.Replace(fpt, "go_build_minfang827.exe", `Config\onlysend\`, 1)
											way2 := strings.Replace(fpt, "go_build_minfang827.exe", `Config\readsend\`, 1)

											//onlysend
											files, _ := ioutil.ReadDir(way1)
											fmt.Println(way1)
											osmessagenumber = 0
											for _, f := range files {
												fmt.Println(f.Name())
												onlysendmessagelist[osmessagenumber] = f.Name()
												data, _ := ioutil.ReadFile(way1 + f.Name())
												fmt.Println(data)
												onlysendmessage[osmessagenumber] = string(data)
												osmessagenumber ++
											}
											//readsend
											files, _ = ioutil.ReadDir(way2)
											fmt.Println(way2)
											rsmessagenumber = 0
											for _, f := range files {
												fmt.Println(f.Name())
												readsendmessagelist[rsmessagenumber] = f.Name()
												data, _ := ioutil.ReadFile(way2 + f.Name())
												fmt.Println(data)
												readsendmessage[rsmessagenumber] = string(data)
												rsmessagenumber ++
											}
											fmt.Println(rsmessagenumber + osmessagenumber)
											mw.ResetRows()
										},
									},
									PushButton{
										Text: "Connect",
										OnClicked: func() {
											serveraddr := IP.Text() + ":" + Port.Text()
											conn, err := net.Dial("tcp", serveraddr)

											if err != nil {
												fmt.Printf("connect failed, err : %v\n", err.Error())
												return
											}
											lnkclient = conn
											lnkclientconnectFlag = true
											connectstatus.SetText("connect for " + serveraddr)
											go mw.TcpClientReadandSend(tcpdisconnect, lnkclient, &lnkclientconnectFlag)
										},
									},
									PushButton{
										Text: "Send",
										OnClicked: func() {
											buf := bytes.NewBufferString("发送消息:")
											fmt.Println(buf.String())
											for _, x := range mw.model.items {
												if x.checked {
													fmt.Printf("checked: %v\n", x)
													var err error
													//判断十六进制
													if(Hexflag){
														result, buffer := stringtoASCII(x.MessageInfo)
														if(result){
															var delim byte = 0x23 //在stringtoASCII处理中增加的结束符
															//fmt.Printf("before read buffer: %v,len: %v \n", buffer.String(),buffer.Len())
															line, _ := buffer.ReadString(delim)
															fmt.Printf("after read buffer: %v,len: %v\n", buffer.String(),buffer.Len())
															_, err = lnkclient.Write([]byte(strings.Trim(line,"#")))
														}

													}else{
														_, err = lnkclient.Write([]byte(x.MessageInfo))
													}
													if err != nil {
														fmt.Printf("write failed , err : %v\n", err)
														break
													}
													//将newString这个string写到buf的尾部
													buf.WriteString(x.MessageInfo)
												}
											}
											mw.messageforsocket.SetText(buf.String() + "\r\n")
											fmt.Println()

										},
									},
									PushButton{
										Text: "disconnect",
										OnClicked: func() {
											if (lnkclientconnectFlag) {
												go disconnect(tcpdisconnect, lnkclient)
											}
											connectstatus.SetText("disconnect")

										},
									},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									HSplitter{
										Children: []Widget{
											Label{
												Text: "IP",
											},
											TextEdit{AssignTo: &IP},
											Label{
												Text: "PORT",
											},
											TextEdit{AssignTo: &Port},
											Label{
												Text: "status for connection",
											},
											TextEdit{AssignTo: &connectstatus},
										},
									},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									HSplitter{
										Children: []Widget{
											Label{
												Text: "message for socket",
											},
											TextEdit{
												AssignTo: &mw.messageforsocket,
												ReadOnly: false,
												HScroll:  true,
												VScroll:  true,
												Text:     "waiting for send or receive",
											},
										},
									},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								ContextMenuItems: []MenuItem{

									Action{
										Text: "E&xit",
										OnTriggered: func() {
											lnkclient.Close()
											mw.Close()
										},
									},
								},
								Children: []Widget{
									TableView{
										AssignTo:         &mw.tv,
										CheckBoxes:       true,
										ColumnsOrderable: true,
										MultiSelection:   true,
										Columns: []TableViewColumn{
											{Title: "编号"},
											{Title: "消息名称"},
											{Title: "消息类型"},
											//{Title: "消息内容"},
										},
										Model: mw.model,
										OnCurrentIndexChanged: func() {
											i := mw.tv.CurrentIndex()
											if 0 <= i {
												fmt.Printf("OnCurrentIndexChanged: %v\n", mw.model.items[i].Name)
												messageinfo.SetText(mw.model.items[i].MessageInfo)

											}
										},
										OnSelectedIndexesChanged: func() {
											fmt.Printf("SelectedIndexes: %v\n", mw.tv.SelectedIndexes())
										},
									},
									//TextEdit{AssignTo: &messageinfo},
									TextEdit{
										AssignTo: &messageinfo,
										ReadOnly: false,
										HScroll:  true,
										VScroll:  true,
										Text:     "",
									},
								},
							},
							Composite{
								Layout: Grid{Columns: 2},
								Children: []Widget{
									Label{
										Text: "十六进制",
									},
									CheckBox{
										Checked: Bind("Hexflag"),
									},
								},
							},
						},
					},
					{
						Title:  "104",
						Layout: VBox{},
					},
				},
			},
		},
	}.Run()
}
func disconnect(ch chan bool, lnkclient net.Conn) {
	//lnkclient.Close()
	//lnkclientconnectFlag = false
	ch <- true
}
func stringtoASCII(loadinfo string) (bool, bytes.Buffer)  {
	var reinfo bytes.Buffer
	var outinfo []byte = make([]byte, 1)
	loadinfoByte := strings.Split(loadinfo, " ")
	if(loadinfoByte[0] == ""){
		fmt.Printf("文件是空的")
		reinfo.Write([]byte(""))
		return false, reinfo
	}
		
	lenforload := len(loadinfoByte)
	fmt.Printf("byte number is: %v\n", lenforload)
	for  i:=0; i < lenforload; i +=1{
		fmt.Printf(loadinfoByte[i])
		single := []byte(loadinfoByte[i])

		single[1] = ASCIItoBi(single[1])
		single[0] = ASCIItoBi(single[0]) 
		outinfo[0] = single[0] * 16 + single[1]
		fmt.Printf("append byte  is: %v\n", outinfo)
		reinfo.Write(outinfo)
	}
	
	reinfo.Write([]byte("#"))
	//fmt.Println(reinfo.String())
	return true, reinfo

}
func ASCIItoBi(IN byte)byte {
	if (IN >=48 && IN <= 57){
		return (IN - 48)
	}else if(IN >=97 && IN <= 102){
		return (IN - 97 + 10)
	}else if(IN >=65 && IN <= 70){
		return (IN - 65 + 10)
	}else{
		return 0xFF
	}
		
}
func (mw *CondomMainWindow) tv_ItemActivated() {
	msg := ``
	for _, i := range mw.tv.SelectedIndexes() {
		msg = msg + "\n" + mw.model.items[i].Name + ":" + mw.model.items[i].MessageInfo
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
		mw.messageforsocket.SetText("接收消息:"+ "\r\n" + str + "\r\n")
		sendmessageName := mw.GetXMLanswer(str)
		if ("" != sendmessageName) {
			for _, x := range mw.model.items {
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
