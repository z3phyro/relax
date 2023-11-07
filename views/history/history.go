package history

import (
	"fmt"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"stoicdynamics.com/relax/types"
)

var (
	wHistoryList    *widgets.List
	RequestsHistory []types.RequestLog = []types.RequestLog{}
	focused         bool               = false
)

func InitHistory() *widgets.List {
	wHistoryList = widgets.NewList()
	wHistoryList.Title = " History "
	wHistoryList.SelectedRowStyle = termui.NewStyle(termui.ColorGreen)

	return wHistoryList
}

func LogRequest(request types.Request, response string) {
	RequestsHistory = append(RequestsHistory, types.RequestLog{
		Request:  request,
		Response: response,
	})

	wHistoryList.Rows = append(wHistoryList.Rows, fmt.Sprintf("%s - %s At %s", request.Name, request.FileName, time.Now().Format("15:04:05")))
	wHistoryList.SelectedRow = len(wHistoryList.Rows) - 1
}

func ToggleFocus() {
	focused = !focused

	if focused {
		wHistoryList.Block.BorderStyle = termui.NewStyle(termui.ColorGreen)
		wHistoryList.TitleStyle = termui.NewStyle(termui.ColorGreen, termui.ColorClear, termui.ModifierBold)
	} else {
		wHistoryList.Block.BorderStyle = termui.NewStyle(termui.ColorWhite)
		wHistoryList.TitleStyle = termui.NewStyle(termui.ColorWhite)
	}
}
