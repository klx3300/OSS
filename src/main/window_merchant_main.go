package main

import (
	"candies"
	"database/sql"
	"dbconn"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

var mmwndRealnameBox *widget.Entry
var mmwndPhoneBox *widget.Entry
var mmwndAddressBox *widget.Entry
var mmwndUpdateBtn *widget.Button
var mmwndNewShopNameBox *widget.Entry
var mmwndNewShopBtn *widget.Button
var mmwndStatusLine *widget.Label
var mmwndShopGroup *widget.Box
var mmwndShopMap map[string]int

func setupWndMerchantMain() {
	mmwndShopMap = make(map[string]int)
	mmwndRealnameBox = widget.NewEntry()
	mmwndPhoneBox = widget.NewEntry()
	mmwndAddressBox = widget.NewMultiLineEntry()
	mmwndUpdateBtn = widget.NewButton("更新个人信息", func() { go mmwndOnUpdatePressed() })
	mmwndNewShopNameBox = widget.NewEntry()
	mmwndNewShopBtn = widget.NewButton("创建新商店", func() { go mmwndOnNewShopPressed() })
	mmwndShopGroup = widget.NewVBox()
	mmwndStatusLine = widget.NewLabelWithStyle("等待用户操作", fyne.TextAlignCenter, candies.NewTextStyle().SetBold().Fin())
	log.Debugln("Starting query information for merchant main window")
	uinfostmt := candies.FastSQLPrep("SELECT name, phone, address FROM user_info WHERE uid = ?", dbconn.Db, mainEventBus)
	defer uinfostmt.Close()
	var realnm, phone, addr string
	uinfoerr := uinfostmt.QueryRow(curruid).Scan(&realnm, &phone, &addr)
	if uinfoerr != nil {
		log.Warnln("User information missing:", uinfoerr.Error())
		mmwndStatusLine.SetText("用户数据损坏，请联系管理员")
	} else {
		mmwndRealnameBox.SetText(realnm)
		mmwndPhoneBox.SetText(phone)
		mmwndAddressBox.SetText(addr)
	}
	log.Debugln("Starting query shops for this merchant")
	shopstmt := candies.FastSQLPrep("SELECT shopid, name FROM shop_info WHERE belongs = ?", dbconn.Db, mainEventBus)
	defer shopstmt.Close()
	shopex, shopexerr := shopstmt.Query(curruid)
	if shopexerr != nil {
		log.Debugln("Shop query fail:", shopexerr.Error())
		mmwndStatusLine.SetText("数据库内部错误，请重启程序")
	}
	for true {
		var shopid int
		var shopnm string
		remain := shopex.Next()
		if !remain {
			log.Debugln("Shop iterate next false")
			break
		}
		shoperr := shopex.Scan(&shopid, &shopnm)
		if shoperr != nil {
			log.Debugln("ShopIterateFail:", shoperr.Error())
			break
		}
		log.Debugln("Shop:", shopnm, shopid)
		mmwndShopMap[shopnm] = shopid
		mmwndShopGroup.Append(widget.NewButton(shopnm, func() { go mmwndOnShopListClicked(shopid, shopnm) }))
	}
	currwnd = currapp.NewWindow("商户管理界面")
	currwnd.SetContent(
		widget.NewScrollContainer(
			widget.NewVBox(
				fyne.NewContainerWithLayout(
					layout.NewGridLayout(2),
					widget.NewLabel("真实姓名"),
					mmwndRealnameBox),
				fyne.NewContainerWithLayout(
					layout.NewGridLayout(2),
					widget.NewLabel("联系电话"),
					mmwndPhoneBox),
				fyne.NewContainerWithLayout(
					layout.NewGridLayout(2),
					widget.NewLabel("联系地址"),
					mmwndAddressBox),
				mmwndUpdateBtn,
				fyne.NewContainerWithLayout(
					layout.NewGridLayout(2),
					widget.NewLabel("新商店名称"),
					mmwndNewShopNameBox),
				mmwndNewShopBtn,
				mmwndStatusLine,
				widget.NewLabel("已创建商店列表"),
				mmwndShopGroup)))
	currwnd.Resize(fyne.NewSize(500, 600))
	currwnd.SetOnClosed(globalFastCloser)
	log.Debugln("Completed. returning to main window..")
}

func mmwndOnUpdatePressed() {
	mmwndOperable(false)
	defer mmwndOperable(true)
	mmwndStatusLine.SetText("请求中")
	infostmt := candies.FastSQLPrep("UPDATE user_info SET name = ?, phone = ?, address = ? WHERE uid = ? ", dbconn.Db, mainEventBus)
	_, execErr := infostmt.Exec(mmwndRealnameBox.Text, mmwndPhoneBox.Text, mmwndAddressBox.Text, curruid)
	if execErr != nil {
		log.Warnln("update uinfo db query fail:", execErr.Error())
		mmwndStatusLine.SetText("数据库内部错误，请重启程序")
		return
	}
	log.Succln("User detail updated.")
	mmwndStatusLine.SetText("用户信息更新成功")
}

func mmwndOnNewShopPressed() {
	mmwndOperable(false)
	mmwndOperable(true)
	mmwndStatusLine.SetText("请求中")
	sdupstmt := candies.FastSQLPrep("SELECT shopid FROM shop_info WHERE name = ?", dbconn.Db, mainEventBus)
	defer sdupstmt.Close()
	var dupid int
	sduperr := sdupstmt.QueryRow(mmwndNewShopNameBox.Text).Scan(&dupid)
	if sduperr != sql.ErrNoRows {
		if sduperr == nil {
			mmwndStatusLine.SetText("该商店名已被占用")
			log.Debugln("Duplicate shop name with shopid", dupid)
			return
		}
		mmwndStatusLine.SetText("数据库内部错误，请重启程序")
		log.Warnln("Shopname dup check query fail:", sduperr.Error())
		return
	}
	sappstmt := candies.FastSQLPrep("INSERT INTO shop_info (belongs, name) VALUES (?, ?)", dbconn.Db, mainEventBus)
	defer sappstmt.Close()
	_, sapperr := sappstmt.Exec(curruid, mmwndNewShopNameBox.Text)
	if sapperr != nil {
		mmwndStatusLine.SetText("数据库内部错误，请重启程序")
		log.Warnln("Shopname append query fail:", sduperr.Error())
		return
	}
	log.Debugln("newshop insert ok. start updating UI...")
	sreqstmt := candies.FastSQLPrep("SELECT shopid FROM shop_info WHERE name = ?", dbconn.Db, mainEventBus)
	defer sreqstmt.Close()
	var shopid int
	sreqerr := sdupstmt.QueryRow(mmwndNewShopNameBox.Text).Scan(&dupid)
	if sreqerr != nil {
		mmwndStatusLine.SetText("数据库内部错误，请重启程序")
		log.Warnln("Shopid requery fail:", sduperr.Error())
		return
	}
	shopnm := mmwndNewShopNameBox.Text
	mmwndShopMap[shopnm] = shopid
	mmwndShopGroup.Append(widget.NewButton(shopnm, func() { go mmwndOnShopListClicked(shopid, shopnm) }))
}

func mmwndOnShopListClicked(shopid int, shopnm string) {
	selshopid = shopid
	selshopnm = shopnm
	log.Debugln("Goto shop id", shopid)
	currwnd.Hide()
	chanTransist <- WND_MERCHANT_SHOP_DETAIL
}

func mmwndOperable(enabled bool) {
	if enabled {
		mmwndUpdateBtn.Show()
		mmwndNewShopBtn.Show()
		mmwndShopGroup.Show()
	} else {
		mmwndNewShopBtn.Hide()
		mmwndUpdateBtn.Hide()
		mmwndShopGroup.Hide()
	}
}
