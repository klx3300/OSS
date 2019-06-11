package main

import (
	"candies"
	"database/sql"
	"dbconn"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

var rwndUsernameBox *widget.Entry
var rwndPasswordBox *widget.Entry
var rwndRetypePassBox *widget.Entry
var rwndRealnameBox *widget.Entry
var rwndPhoneBox *widget.Entry
var rwndAddressBox *widget.Entry
var rwndMerchantCheck *widget.Check
var rwndCheckUsernameBtn *widget.Button
var rwndSubmitBtn *widget.Button
var rwndCancelBtn *widget.Button
var rwndStatusLine *widget.Label

func setupWndRegister() {
	rwndUsernameBox = widget.NewEntry()
	rwndPasswordBox = widget.NewPasswordEntry()
	rwndRetypePassBox = widget.NewPasswordEntry()
	rwndRealnameBox = widget.NewEntry()
	rwndPhoneBox = widget.NewEntry()
	rwndAddressBox = widget.NewMultiLineEntry()
	rwndMerchantCheck = widget.NewCheck("是否为商家", func(c bool) {})
	rwndCheckUsernameBtn = widget.NewButton("用户名检测", func() { go rwndOnCheckUsernamePressed() })
	rwndSubmitBtn = widget.NewButton("提交", func() { go rwndOnSubmitPressed() })
	rwndCancelBtn = widget.NewButton("取消", func() { go rwndOnCancelPressed() })
	rwndStatusLine = widget.NewLabel("等待用户操作")

	currwnd = currapp.NewWindow("新用户注册")
	currwnd.SetContent(
		widget.NewVBox(
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(2), widget.NewLabel("用户名"), rwndUsernameBox),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(2), widget.NewLabel("密码"), rwndPasswordBox),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(2), widget.NewLabel("确认密码"), rwndRetypePassBox),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(2), widget.NewLabel("真实姓名"), rwndRealnameBox),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(2), widget.NewLabel("联系电话"), rwndPhoneBox),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(2), widget.NewLabel("联系地址"), rwndAddressBox),
			rwndMerchantCheck,
			rwndStatusLine,
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(3), rwndCheckUsernameBtn, rwndSubmitBtn, rwndCancelBtn)))
	currwnd.Resize(fyne.NewSize(500, 400))
	currwnd.SetOnClosed(globalFastCloser)
	log.Debugln("register window settled")
}

func rwndOnCheckUsernamePressed() {
	rwndCheckUsernameUsed()
}

func rwndOnSubmitPressed() {
	rwndSetButtons(false)
	defer rwndSetButtons(true)
	rwndStatusLine.SetText("请求中")
	log.Debugln("Starting register process...")
	if rwndPasswordBox.Text != rwndRetypePassBox.Text {
		rwndStatusLine.SetText("两次输入密码不一致")
		return
	}
	if rwndRealnameBox.Text == "" {
		rwndStatusLine.SetText("真实姓名不能为空")
		return
	}
	if rwndPhoneBox.Text == "" {
		rwndStatusLine.SetText("联系电话不能为空")
		return
	}
	if rwndAddressBox.Text == "" {
		rwndStatusLine.SetText("联系地址不能为空")
		return
	}
	if rwndCheckUsernameUsed() {
		return
	}
	rwndSetButtons(false)
	rwndStatusLine.SetText("注册请求中")
	log.Debugln("register prereqs satisfy. starting db query...")
	stmt := candies.FastSQLPrep("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", dbconn.Db, mainEventBus)
	defer stmt.Close()
	role := "customer"
	if rwndMerchantCheck.Checked {
		role = "merchant"
	}
	_, execErr := stmt.Exec(rwndUsernameBox.Text, rwndPasswordBox.Text, role)
	if execErr != nil {
		log.Warnln("register db query fail:", execErr.Error())
		rwndStatusLine.SetText("数据库内部错误，请重启程序")
		return
	}
	log.Succln("Register succeeded. Starting insert user details")
	log.Debugln("Refetching user id...")
	uidstmt := candies.FastSQLPrep("SELECT uid FROM users where username=?", dbconn.Db, mainEventBus)
	defer uidstmt.Close()
	var uid int
	execErr = uidstmt.QueryRow(rwndUsernameBox.Text).Scan(&uid)
	if execErr != nil {
		log.Warnln("register db requery fail:", execErr.Error())
		rwndStatusLine.SetText("数据库内部错误，请重启程序")
		return
	}
	infostmt := candies.FastSQLPrep("INSERT INTO user_info (uid, name, phone, address) VALUES (?, ?, ?, ?) ", dbconn.Db, mainEventBus)
	_, execErr = infostmt.Exec(uid, rwndRealnameBox.Text, rwndPhoneBox.Text, rwndAddressBox.Text)
	if execErr != nil {
		log.Warnln("register db query fail:", execErr.Error())
		rwndStatusLine.SetText("数据库内部错误，请重启程序")
		return
	}
	rwndStatusLine.SetText("注册成功，请返回登录界面")
}

func rwndOnCancelPressed() {
	log.Debugln("Returning to login")
	currwnd.Hide()
	chanTransist <- WND_LOGIN
}

// false on unused.
func rwndCheckUsernameUsed() bool {
	rwndSetButtons(false)
	defer rwndSetButtons(true)
	rwndStatusLine.SetText("请求中")
	log.Debugln("Starting username used check query..")
	stmt := candies.FastSQLPrep("SELECT uid FROM users where username=?", dbconn.Db, mainEventBus)
	defer stmt.Close()
	var uid int
	queryerr := stmt.QueryRow(rwndUsernameBox.Text).Scan(&uid)
	if queryerr != nil {
		if queryerr == sql.ErrNoRows {
			rwndStatusLine.SetText("用户名可用")
			return false
		}
		rwndStatusLine.SetText("数据库异常，请重启程序")
		log.Warnln("Query failure:", queryerr.Error())
		return true
	}
	log.Debugln("Occupied by user id", uid)
	rwndStatusLine.SetText("用户名不可用")
	return true
}

func rwndSetButtons(enabled bool) {
	if !enabled {
		rwndCheckUsernameBtn.Hide()
		rwndSubmitBtn.Hide()
		rwndCancelBtn.Hide()
	} else {
		rwndCheckUsernameBtn.Show()
		rwndSubmitBtn.Show()
		rwndCancelBtn.Show()
	}
}
