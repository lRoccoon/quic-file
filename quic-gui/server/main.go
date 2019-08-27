package main

import (
    "github.com/andlabs/ui"
	_"github.com/andlabs/ui/winmanifest"
	"io/ioutil"
)

var mainwin *ui.Window
var rs *ui.MultilineEntry
var filepaths []byte
var result string

func makeBasicControlsPage() ui.Control{
	//整体模块
    var protocol, test bool
	//var address string
	box := ui.NewHorizontalBox()
	box.SetPadded(true)
    //左边部分
	group := ui.NewGroup("choices")
	group.SetMargined(true)
	box.Append(group, true)
	lbox :=ui.NewVerticalBox()
	lbox.SetPadded(true)
	group.SetChild(lbox)
    //左边第一部分
	lbox.Append(ui.NewLabel("Server Ip"),false)
	ip :=ui.NewEntry()
	ip.SetReadOnly(false)
	port :=ui.NewEntry()
	lbox.Append(ip, false)
	lbox.Append(ui.NewLabel("Server Port"),false)
	lbox.Append(port, false)
    //左边第二部分
	lbox.Append(ui.NewLabel("Please choose the protocol "),false)
	cbox := ui.NewEditableCombobox()
	lbox.Append(cbox, false)
	cbox.Append("QUIC Protocol")
	cbox.Append("TCP Protocol")
	//左边第三部分
	rs=ui.NewNonWrappingMultilineEntry()
	rs.Handle()
	lbox.Append(ui.NewLabel("the file result" ),false)
	lbox.Append(rs, true)
    //右边部分
	rgroup:=ui.NewGroup("results")
	rgroup.SetMargined(true)
	box.Append(rgroup, true)
	rbox := ui.NewVerticalBox()
	rbox.SetPadded(true)
	rgroup.SetChild(rbox)
	//第一部分
	rbox.Append(ui.NewLabel("Please choose the ways "),false)
	dbox := ui.NewEditableCombobox()
	rbox.Append(dbox, false)
	dbox.Append("upload file")
	dbox.Append("download file")
	//第二部分
	button :=ui.NewButton("start server")
	button2 :=ui.NewButton("select file")
	grid := ui.NewGrid()
	grid.SetPadded(true)
	rbox.Append(grid, false)
	grid.Append(button,0, 0, 1, 1,true, ui.AlignFill, false, ui.AlignFill)
	grid.Append(button2,4, 0, 1, 1,true, ui.AlignFill, false, ui.AlignFill)
	entry := ui.NewNonWrappingMultilineEntry()
	entry.Handle()
	rbox.Append(entry,true)
	button.OnClicked(func(*ui.Button){
		//判断协议
		if (cbox.Text()=="QUIC Protocol"){
			protocol=true
		} else {
			protocol=false
		}
		//判断上传下载
		if (dbox.Text()=="upload file"){
			test=true
		} else{
			test=false
		}
		if protocol{
		   go Server(ip.Text()+":"+port.Text(),test)
		} else {
		   go TCPServer(ip.Text(),port.Text(),test)
		}
	})
	filename:=""
	files, _ := ioutil.ReadDir("./")
    for _, f := range files {
		filename+=f.Name()+"\n"
	}
    filepaths= [] byte(filename)
	button2.OnClicked(func(*ui.Button) {
		entry.Append(string(filepaths)+"\n")
	})
    return box
}

func setupUI(){
	mainwin=ui.NewWindow("this is a quic transmission for server",640,480,true)
	mainwin.OnClosing(func(*ui.Window) bool{
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool{
		mainwin.Destroy()
		return true
	})

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	tab.Append("server controls", makeBasicControlsPage())
	tab.SetMargined(0, true)

	mainwin.Show()
}

func main(){
	ui.Main(setupUI)
}
