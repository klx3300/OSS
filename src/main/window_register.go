package main

import (
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
}

func rwndOnCheckUsernamePressed() {

}

func rwndOnSubmitPressed() {

}

func rwndOnCancelPressed() {

}
