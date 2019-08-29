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
	"encoding/xml"
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
	m.items = make([]*Condom, 0)

	return m
}

type CondomMainWindow struct {
	*walk.MainWindow
	model *CondomModel
	tv    *walk.TableView
}

func (m *CondomMainWindow) ResetRows() {
	items := []*Condom{}
	

	for i := 0; i<osmessagenumber ; i++ {
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
			Index: j+1,
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
	var lnkclientconnectFlag bool = false
	tcpdisconnect :=make(chan bool)
	
	MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Robot Simulator",
		Size:     Size{800, 500},
		Layout:   VBox{},

		Children: []Widget{
			Composite{
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
							lnkclientconnectFlag = true
							go mw.TcpClientReadandSend(tcpdisconnect,lnkclient,&lnkclientconnectFlag)
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

						},
					},
					PushButton{
						Text: "disconnect",
						OnClicked: func() {
							if(lnkclientconnectFlag){
								go disconnect(tcpdisconnect,lnkclient)
							}
							
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
func disconnect(ch chan bool,lnkclient net.Conn){
	//lnkclient.Close()
	//lnkclientconnectFlag = false
	ch <- true
}

func (mw *CondomMainWindow) tv_ItemActivated() {
	msg := ``
	for _, i := range mw.tv.SelectedIndexes() {
		msg = msg + "\n" + mw.model.items[i].Name + ":" + mw.model.items[i].MessageInfo
	}
	walk.MsgBox(mw, "title", msg, walk.MsgBoxIconInformation)
}

func (mw *CondomMainWindow) TcpClientReadandSend(ch chan bool,lnkclient net.Conn,lnkclientconnectFlag *bool){
	var disconnectflag bool
	for {
		disconnectflag = <-ch
		if(disconnectflag){
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
		sendmessageName := mw.GetXMLanswer(str)
		if("" != sendmessageName){
			for _, x := range mw.model.items {
				if x.Name == sendmessageName {
					var err error
					_,err = lnkclient.Write([]byte(x.MessageInfo))
					if err != nil {
						fmt.Printf("write failed , err : %v\n", err)
						break
					}
				}
			}
		}	
	}
}
func (mw *CondomMainWindow) GetXMLanswer(receiveXML string)string {
	var t xml.Token
    var err error
    inputReader := strings.NewReader(receiveXML)

    // 从文件读取，如可以如下：
    // content, err := ioutil.ReadFile("studygolang.xml")
    // decoder := xml.NewDecoder(bytes.NewBuffer(content))

    decoder := xml.NewDecoder(inputReader)
    var TYPEflag bool
    for t,err = decoder.Token(); err == nil; t,err = decoder.Token() {
        switch token := t.(type) {
        // 处理元素开始（标签）
        case xml.StartElement:
            name := token.Name.Local
            fmt.Printf("Token name: %s\n", name)
            if("Type" == name) {
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
            if(TYPEflag){
				content := string([]byte(token))
				fmt.Printf("This is the content: %v\n", content)
				if("1" == content){
					return "control58Res"
				}
			}
			TYPEflag=false
        default:
            // ...
        }
    }
	return ""
}