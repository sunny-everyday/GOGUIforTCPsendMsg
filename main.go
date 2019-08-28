package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"net"
	"path/filepath"
	"io/ioutil"
	"strings"
	"strconv"
)

type Condom struct {
	Index   int
	Name    string
	Type    string
	MessageInfo string
	checked bool
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
	length :=len(m.items)
	for row := 0; row < length; row++ {
		if m.Checked(row){
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
func (m *CondomModel)FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(float64(input_num), 'f', 6, 64)
}
//message number
var osmessagenumber int  = 0
var rsmessagenumber int  = 0

var onlysendmessagelist [100]string
var readsendmessagelist [100]string

var onlysendmessage [100]string
var readsendmessage [100]string


func NewCondomModel() *CondomModel {
	m := new(CondomModel)
	m.items = make([]*Condom, 5)

	m.items[0] = &Condom{
		Index: 1,
		Name:  "消息1",
		Type: "2接口",
		MessageInfo:"aaa",
		
	}

	m.items[1] = &Condom{
		Index: 2,
		Name:  "消息2",
		Type: "1接口",
		MessageInfo:"bbb",
	}

	m.items[2] = &Condom{
		Index: 3,
		Name:  "消息3",
		Type: "1接口",
		MessageInfo:"ccc",
	}
	m.items[3] = &Condom{
		Index: 3,
		Name:  "消息3",
		Type: "2接口",
		MessageInfo:"ddd",
	}

	m.items[4] = &Condom{
		Index: 3,
		Name:  "消息3",
		Type: "2接口",
		MessageInfo:`MESSAGE sip: 前端设备地址编码@前端设备所属系统域名或IP地址 SIP/2.0
		From: <sip: 用户地址编码@用户所属系统域名或IP地址> ;tag=BK32B1U8DKDrB
		To: <sip: 前端设备地址编码@前端设备所属系统域名或IP地址>
		Contact: <sip: 用户地址编码@用户所属系统域名或IP地址>
		Via: SIP/2.0/UDP 用户所属系统IP地址;branch=z9hG4bK
		Call-ID: c47k42
		CSeq:1 MESSAGE
		Content-type: application/xml
		Content-Length: 消息体的长度
		
		<?xml version="1.0" encoding="UTF-8"?>
		<SIP_XML EventType=" Alg_Ability_Query">
		<!-- 对象地址编码最小对象为前端系统 -->
		<Item Code="对象地址编码"/>
		</SIP_XML>`,
	}
	return m
}

type CondomMainWindow struct {
	*walk.MainWindow
	model *CondomModel
	tv    *walk.TableView
}

func (m *CondomMainWindow) ResetRows() {
	items := []*Condom{}
	

	for i := 0; i<osmessagenumber -1; i++ {
		x:= &Condom{
			Index: i+1,
			Name:   onlysendmessagelist[i],
			Type:   "send",
			MessageInfo:  onlysendmessage[i],
		}
		items = append(items, x)
	}

	for j := osmessagenumber; j < osmessagenumber +rsmessagenumber; j++ {
		x:= &Condom{
			Index: j,
			Name:   readsendmessagelist[j-osmessagenumber],
			Type:   "receive",
			MessageInfo:  readsendmessage[j-osmessagenumber],
		}
		items = append(items, x)
	}
	m.model.items = items
	m.model.PublishRowsReset()
	m.tv.SetSelectedIndexes([]int{})
}
func main() {
	mw := &CondomMainWindow{model: NewCondomModel()}
	var IP, Port *walk.TextEdit
	var lnkclient net.Conn
	MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Robot Simulator",
		Size:     Size{800, 500},
		Layout:   VBox{},
		Children: []Widget{
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
						},
					},
					//PushButton{
					//	Text:      "Reset Rows",
					//	OnClicked: mw.ResetRows,
					//},
					PushButton{
						Text: "Reload",
						OnClicked: func() {
							fpt, err := filepath.Abs("___go_build_minfang827.exe")
							if err != nil {
								panic(err)
							}
							fmt.Println(fpt)
							way1 := strings.Replace(fpt, "___go_build_minfang827.exe", `minfang827\Config\onlysend\`, 1)
							way2 := strings.Replace(fpt, "___go_build_minfang827.exe", `minfang827\Config\readsend\`, 1)
							

							//onlysend
							files, _ := ioutil.ReadDir(way1)
							fmt.Println(way1)
							osmessagenumber = 0
							for _, f := range files {
								fmt.Println(f.Name())
								onlysendmessagelist[osmessagenumber] = f.Name()
								data, _ := ioutil.ReadFile(way1+f.Name())
								fmt.Println(data)
								onlysendmessage[osmessagenumber]  = string(data)
								osmessagenumber ++ 
							}
							//readsend
							files, _ = ioutil.ReadDir(way2)
							fmt.Println(way2)
							rsmessagenumber = 0
							for _, f := range files {
								fmt.Println(f.Name())
								readsendmessagelist[rsmessagenumber] = f.Name()
								data, _ := ioutil.ReadFile(way2+f.Name())
								fmt.Println(data)
								readsendmessage[rsmessagenumber] = string(data)
								rsmessagenumber ++
							}
							fmt.Println(rsmessagenumber +osmessagenumber)
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
						},
					},
					PushButton{
						Text: "Send",
						OnClicked: func() {
							for _, x := range mw.model.items {
								if x.checked {
									fmt.Printf("checked: %v\n", x)
									var err error
									_,err = lnkclient.Write([]byte(x.MessageInfo))
									
									if err != nil {
										fmt.Printf("write failed , err : %v\n", err)
										break
									}
								}
							}
							fmt.Println()
							lnkclient.Close()
						},
					},
				},
			},
			Composite{
				Layout: VBox{},
				ContextMenuItems: []MenuItem{
					Action{
						Text:        "I&nfo",
						OnTriggered: mw.tv_ItemActivated,
					},
					Action{
						Text: "E&xit",
						OnTriggered: func() {
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
							{Title: "消息内容"},
						},
						Model: mw.model,
						OnCurrentIndexChanged: func() {
							i := mw.tv.CurrentIndex()
							if 0 <= i {
								fmt.Printf("OnCurrentIndexChanged: %v\n", mw.model.items[i].Name)

							}
						},
						OnSelectedIndexesChanged: func() {
							fmt.Printf("SelectedIndexes: %v\n", mw.tv.SelectedIndexes())
						},
						OnItemActivated: mw.tv_ItemActivated,
					},

				},
			},
		},
	}.Run()
}

func (mw *CondomMainWindow) tv_ItemActivated() {
	msg := ``
	for _, i := range mw.tv.SelectedIndexes() {
		msg = msg + "\n" + mw.model.items[i].Name + ":" + mw.model.items[i].MessageInfo
	}
	walk.MsgBox(mw, "title", msg, walk.MsgBoxIconInformation)
}
