package gui

import (
	"GOGUIforTCPsendMsg/common"
	"fmt"
	"github.com/lxn/walk"
	"net"
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
	for {
		select {
		case disconnectflag := <-ch:
			if disconnectflag {
				lnkclient.Close()
				*lnkclientconnectFlag = false
				fmt.Println("Exit!")
				return
			}
		default:
		}

		var buf [2000]byte
		n, err := lnkclient.Read(buf[:])

		if err != nil {
			fmt.Printf("read from connect failed, err: %v\n", err)
			break
		}
		str := string(buf[:n])
		fmt.Printf("receive from client, data: %s\n", str)
		mw.Messageforsocket.SetText("接收消息:" + "\r\n" + str + "\r\n")
		sendmessageName := common.GetXMLanswer(str)
		if "" != sendmessageName {
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
