// Copyright 2017 The Walk Authors. All rights reserved.
// Use of this source code is governed by  a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"GOGUIforTCPsendMsg/common"
	"bytes"
)

import (
	gui "GOGUIforTCPsendMsg/gui"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"net"
	"path/filepath"
	"strings"
)

func main() {
	mw := gui.NewCondomMainWindow()
	var IP, Port, connectstatus *walk.LineEdit
	IP.SetText("127.0.0.1")
	Port.SetText("60000")
	var messageinfo *walk.TextEdit
	var lnkclient net.Conn
	var lnkclientconnectFlag bool = false
	var Hexflag bool = true
	tcpdisconnect := make(chan bool)
	var db *walk.DataBinder
	MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Simulator",
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
											gui.Osmessagenumber = 0
											for _, f := range files {
												fmt.Println(f.Name())
												gui.Onlysendmessagelist[gui.Osmessagenumber] = f.Name()
												data, _ := ioutil.ReadFile(way1 + f.Name())
												fmt.Println(data)
												gui.Onlysendmessage[gui.Osmessagenumber] = string(data)
												gui.Osmessagenumber++
											}
											//readsend
											files, _ = ioutil.ReadDir(way2)
											fmt.Println(way2)
											gui.Rsmessagenumber = 0
											for _, f := range files {
												fmt.Println(f.Name())
												gui.Readsendmessagelist[gui.Rsmessagenumber] = f.Name()
												data, _ := ioutil.ReadFile(way2 + f.Name())
												fmt.Println(data)
												gui.Readsendmessage[gui.Rsmessagenumber] = string(data)
												gui.Rsmessagenumber++
											}
											fmt.Println(gui.Rsmessagenumber + gui.Osmessagenumber)
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
											for _, x := range mw.Model.Items {
												if x.Checked {
													fmt.Printf("checked: %v\n", x)
													var err error
													//判断十六进制
													if Hexflag {
														result, buffer := common.StringtoASCII(x.MessageInfo)
														if result {
															var delim byte = 0x23 //在stringtoASCII处理中增加的结束符
															//fmt.Printf("before read buffer: %v,len: %v \n", buffer.String(),buffer.Len())
															line, _ := buffer.ReadString(delim)
															fmt.Printf("after read buffer: %v,len: %v\n", buffer.String(), buffer.Len())
															_, err = lnkclient.Write([]byte(strings.Trim(line, "#")))
														}

													} else {
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
											mw.Messageforsocket.SetText(buf.String() + "\r\n")
											fmt.Println()

										},
									},
									PushButton{
										Text: "disconnect",
										OnClicked: func() {
											if lnkclientconnectFlag {
												go disconnect(tcpdisconnect, lnkclient)
											}
											connectstatus.SetText("disconnect")

										},
									},
								},
							},
							GroupBox{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{
										Text: "IP",
									},
									LineEdit{AssignTo: &IP},
									Label{
										Text: "PORT",
									},
									LineEdit{AssignTo: &Port},
									Label{
										Text: "status for connection",
									},
									LineEdit{AssignTo: &connectstatus},
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
												AssignTo: &mw.Messageforsocket,
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
										AssignTo:         &mw.Tv,
										CheckBoxes:       true,
										ColumnsOrderable: true,
										MultiSelection:   true,
										Columns: []TableViewColumn{
											{Title: "编号"},
											{Title: "消息名称"},
											{Title: "消息类型"},
											//{Title: "消息内容"},
										},
										Model: mw.Model,
										OnCurrentIndexChanged: func() {
											i := mw.Tv.CurrentIndex()
											if 0 <= i {
												fmt.Printf("OnCurrentIndexChanged: %v\n", mw.Model.Items[i].Name)
												messageinfo.SetText(mw.Model.Items[i].MessageInfo)

											}
										},
										OnSelectedIndexesChanged: func() {
											fmt.Printf("SelectedIndexes: %v\n", mw.Tv.SelectedIndexes())
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
