package main

import (
	"candies"
	"dbconn"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

var umwndRealnameBox *widget.Entry
var umwndPhoneBox *widget.Entry
var umwndAddressBox *widget.Entry
var umwndUpdateBtn *widget.Button
var umwndSearchBtn *widget.Button
var umwndOrderListBtn *widget.Button
var umwndStatusLine *widget.Label

func setupWndUserMain() {
	umwndRealnameBox = widget.NewEntry()
	umwndPhoneBox = widget.NewEntry()
	umwndAddressBox = widget.NewMultiLineEntry()
	umwndUpdateBtn = widget.NewButton("更新个人信息", func() { go umwndOnUpdatePressed() })
	umwndSearchBtn = widget.NewButton("商品搜索界面", func() { go umwndOnSearchPressed() })
	umwndOrderListBtn = widget.NewButton("我的订单", func() { go umwndOrderListPressed() })
	umwndStatusLine = widget.NewLabelWithStyle("等待用户操作", fyne.TextAlignCenter, candies.NewTextStyle().SetBold().Fin())
	log.Debugln("Starting query information for merchant main window")
	uinfostmt := candies.FastSQLPrep("SELECT name, phone, address FROM user_info WHERE uid = ?", dbconn.Db, mainEventBus)
	defer uinfostmt.Close()
	var realnm, phone, addr string
	uinfoerr := uinfostmt.QueryRow(curruid).Scan(&realnm, &phone, &addr)
	if uinfoerr != nil {
		log.Warnln("User information missing:", uinfoerr.Error())
		umwndStatusLine.SetText("用户数据损坏，请联系管理员")
	} else {
		umwndRealnameBox.SetText(realnm)
		umwndPhoneBox.SetText(phone)
		umwndAddressBox.SetText(addr)
	}
	currwnd = currapp.NewWindow("用户中心")
	currwnd.SetContent(
		widget.NewScrollContainer(
			widget.NewVBox(
				fyne.NewContainerWithLayout(
					layout.NewGridLayout(2),
					widget.NewLabel("真实姓名"),
					umwndRealnameBox),
				fyne.NewContainerWithLayout(
					layout.NewGridLayout(2),
					widget.NewLabel("联系电话"),
					umwndPhoneBox),
				fyne.NewContainerWithLayout(
					layout.NewGridLayout(2),
					widget.NewLabel("联系地址"),
					umwndAddressBox),
				umwndUpdateBtn,
				umwndSearchBtn,
				umwndOrderListBtn,
				umwndStatusLine)))
	currwnd.Resize(fyne.NewSize(400, 400))
	currwnd.SetOnClosed(globalFastCloser)
	log.Debugln("Completed. returning to main window..")
}

func umwndOnUpdatePressed() {
	umwndOperable(false)
	defer umwndOperable(true)
	umwndStatusLine.SetText("请求中")
	infostmt := candies.FastSQLPrep("UPDATE user_info SET name = ?, phone = ?, address = ? WHERE uid = ? ", dbconn.Db, mainEventBus)
	_, execErr := infostmt.Exec(umwndRealnameBox.Text, umwndPhoneBox.Text, umwndAddressBox.Text, curruid)
	if execErr != nil {
		log.Warnln("update uinfo db query fail:", execErr.Error())
		umwndStatusLine.SetText("数据库内部错误，请重启程序")
		return
	}
	log.Succln("User detail updated.")
	umwndStatusLine.SetText("用户信息更新成功")
}

func umwndOnSearchPressed() {
	currwnd.Hide()
	chanTransist <- WND_USER_SEARCH
}

func umwndOrderListPressed() {
	currwnd.Hide()
	chanTransist <- WND_USER_ORDER_LIST
}

func umwndOperable(enabled bool) {
	if enabled {
		umwndUpdateBtn.Show()
		umwndSearchBtn.Show()
		umwndOrderListBtn.Show()
	} else {
		umwndUpdateBtn.Hide()
		umwndSearchBtn.Hide()
		umwndOrderListBtn.Hide()
	}
}
