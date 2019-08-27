package main

import (
    "github.com/andlabs/ui"
	_"github.com/andlabs/ui/winmanifest"
	"io/ioutil"
)

var mainwin *ui.Window
var rs *ui.MultilineEntry
var filepath []byte
var result string

func makeBasicControlsPage() ui.Control{
	//整体模块
    var protocol, test bool
	var address,filename string
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
	port :=ui.NewEntry()
	lbox.Append(ip, false)
	lbox.Append(ui.NewLabel("Server Port"),false)
	lbox.Append(port, false)
    //左边第二部分
	lbox.Append(ui.NewLabel("Please choose the protocol "),false)
	cbox := ui.NewEditableCombobox()
	lbox.Append(cbox, false)
	cbox.Append("QUIC Protocol")
	cbox.Append("TCP  Protocol")
	//左边第三部分
	sbutton :=ui.NewButton("select file")
	grid := ui.NewGrid()
	grid.SetPadded(true)
	lbox.Append(grid, false)
	grid.Append(sbutton,0, 0, 1, 1,true, ui.AlignFill, false, ui.AlignFill)

	entry := ui.NewNonWrappingMultilineEntry()
	entry.Handle()
	lbox.Append(entry,true)

	sbutton.OnClicked(func(*ui.Button){
		files:=""
		filess, _ := ioutil.ReadDir("./")
		for _, f := range filess {
			files+=f.Name()+"\n"
		}
		filepath= [] byte(files)
		entry.Append(string(filepath)+"\n")
	})
	//右边部分
	rgroup:=ui.NewGroup("results")
	rgroup.SetMargined(true)
	box.Append(rgroup, true)
	rbox := ui.NewVerticalBox()
	rbox.SetPadded(true)
	rgroup.SetChild(rbox)
	//右边第一部分
	grid = ui.NewGrid()
	grid.SetPadded(true)
	rbox.Append(grid, false)
	ubutton := ui.NewButton("Upload File")
	uentry := ui.NewEntry()
	uentry.SetReadOnly(true)
	test=false
	ubutton.OnClicked(func(*ui.Button) {
		filename = ui.OpenFile(mainwin)
		if filename == "" {
			filename = "you have not choose file"
		}
		uentry.SetText(filename)
		test=true
	})
	grid.Append(ubutton,0, 0, 1, 1,false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(uentry,1, 0, 1, 1,true, ui.AlignFill, false, ui.AlignFill)
	dbutton := ui.NewButton("Download File")
	dentry := ui.NewEntry()
	dentry.SetText("please input the filename ")
	dbutton.OnClicked(func (*ui.Button){
		ui.MsgBox(mainwin,
			"you should fill in the filename ",
			"you can press the list file button to see the filename")
	})
	grid.Append(dbutton,0, 1, 1, 1,false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(dentry,1, 1, 1, 1,true, ui.AlignFill, false, ui.AlignFill)
    //右边第二部分
	msggrid := ui.NewGrid()
	msggrid.SetPadded(true)
	grid.Append(msggrid,0, 2, 2, 1,true, ui.AlignCenter, false, ui.AlignStart)
	button := ui.NewButton("  confirm  ")
	msggrid.Append(button,1, 0, 1, 1,true, ui.AlignFill, false, ui.AlignFill)
	cbutton := ui.NewButton("  close  ")
	msggrid.Append(cbutton,2, 0, 1, 1,true, ui.AlignFill, false, ui.AlignFill)
	button.OnClicked(func(*ui.Button) {
		address =ip.Text()+":"+port.Text()
		if cbox.Text()=="QUIC Protocol"{
			protocol=true
		}else{
			protocol=false
		}
		if !test{
			filename=dentry.Text()
		}
		if protocol{
			go Client(address,filename,test)
		} else {
			go TCPClient(address,filename,test)
		}
	})
	cbutton.OnClicked(func(*ui.Button) {
		ui.Quit()
	})
	//有右边第三部分
	rs=ui.NewNonWrappingMultilineEntry()
	rs.Handle()
	rs.SetText(result)
	rbox.Append(ui.NewLabel("the file result" ),false)
	rbox.Append(rs, true)
    return box
}
func setupUI(){
	mainwin=ui.NewWindow("this is a quic transmission for client",640,480,true)
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

	tab.Append("client controls", makeBasicControlsPage())
	tab.SetMargined(0, true)

	mainwin.Show()
}

func main(){
	ui.Main(setupUI)
}
