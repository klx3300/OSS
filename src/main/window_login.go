package main

import (
	"candies"
	"database/sql"
	"dbconn"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

var usernameBox *widget.Entry
var passwordBox *widget.Entry
var loginStatusLine *widget.Label
var loginBtn *widget.Button
var registerBtn *widget.Button

func setupWndLogin() {
	log.Debugln("Setup login window..")
	usernameBox = widget.NewEntry()
	passwordBox = widget.NewPasswordEntry()
	loginStatusLine = widget.NewLabelWithStyle("等待用户操作中", fyne.TextAlignCenter, candies.NewTextStyle().SetBold().Fin())
	currwnd = currapp.NewWindow("用户登录")
	loginBtn = widget.NewButton("登录", func() { go onLoginPressed() })
	registerBtn = widget.NewButton("注册", func() { go onRegisterPressed() })

	currwnd.SetContent(
		widget.NewVBox(
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				widget.NewLabel("用户名"), usernameBox),
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				widget.NewLabel("密码"), passwordBox),
			loginStatusLine,
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				loginBtn, registerBtn)))
	currwnd.SetOnClosed(globalFastCloser)
	currwnd.Resize(fyne.NewSize(400, 200))

	log.Debugln("Login window settled. Waiting to display...")
}

func onLoginPressed() {
	loginBtn.Hide()
	registerBtn.Hide()
	loginStatusLine.SetText("请求中")
	log.Debugln("Start querying db...")
	stmt := candies.FastSQLPrep("SELECT uid, role FROM users where username=? and password=?", dbconn.Db, mainEventBus)
	defer stmt.Close()
	var uid int
	var role string
	queryerr := stmt.QueryRow(usernameBox.Text, passwordBox.Text).Scan(&uid, &role)
	if queryerr != nil {
		if queryerr == sql.ErrNoRows {
			loginStatusLine.SetText("不正确的用户名/密码")
		} else {
			loginStatusLine.SetText("数据库异常，请重启程序")
		}
		loginBtn.Show()
		registerBtn.Show()
		return
	}
	log.Succln("Login success:", uid, role)
	curruid = uid
	// TODO: spawn corresponding interface according to user role: admin/merchant/customer
	if role == "merchant" {
		log.Debugln("go to merchant main")
		currwnd.Hide()
		chanTransist <- WND_MERCHANT_MAIN
		return
	}
	if role == "customer" {
		log.Debugln("go to customer main")
		currwnd.Hide()
		chanTransist <- WND_USER_MAIN
		return
	}
	candies.FastGG("Unknown USER TYPE: "+role, mainEventBus)
}

func onRegisterPressed() {
	log.Debugln("Preparing register..")
	currwnd.Hide()
	chanTransist <- WND_REGISTER
}
