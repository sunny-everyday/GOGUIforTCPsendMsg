package gui

import (
	"github.com/lxn/walk"
	"strconv"
)

type CondomModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	Items      []*Condom
}

func NewCondomModel() *CondomModel {
	m := new(CondomModel)
	m.Items = make([]*Condom, 0)

	return m
}


func (m *CondomModel) RowCount() int {
	return len(m.Items)
}

func (m *CondomModel) Value(row, col int) interface{} {
	item := m.Items[row]

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
	return m.Items[row].Checked
}
func (m *CondomModel) GetCheckedItemlist() []int {
	var j int
	var checkedlist []int
	length := len(m.Items)
	for row := 0; row < length; row++ {
		if m.Checked(row) {
			checkedlist[j] = row
		}
	}
	return checkedlist
}

func (m *CondomModel) SetChecked(row int, checked bool) error {
	m.Items[row].Checked = checked
	return nil
}

func (m *CondomModel) Len() int {
	return len(m.Items)
}

func (m *CondomModel) Swap(i, j int) {
	m.Items[i], m.Items[j] = m.Items[j], m.Items[i]
}
func (m *CondomModel) FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(float64(input_num), 'f', 6, 64)
}

