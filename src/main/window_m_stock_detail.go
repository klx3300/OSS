package main

import (
	"candies"
	"dbconn"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

var mstdwndNameBox *widget.Entry
var mstdwndTypeBox *widget.Entry
var mstdwndDescBox *widget.Entry
var mstdwndPriceBox *widget.Entry
var mstdwndPicPathBox *widget.Entry
var mstdwndPromoBox *widget.Entry
var mstdwndCountBox *widget.Entry
var mstdwndPrecheckBtn *widget.Button
var mstdwndSubmitBtn *widget.Button
var mstdwndReturnBtn *widget.Button
var mstdwndStatusLine *widget.Label
var mstdwndPicturePreview *canvas.Image
var mstdwndContainBox *widget.Box

func setupWndMerchantStockDetail() {
	// setup parts
	mstdwndNameBox = widget.NewEntry()
	mstdwndTypeBox = widget.NewEntry()
	mstdwndDescBox = widget.NewMultiLineEntry()
	mstdwndPriceBox = widget.NewEntry()
	mstdwndPicPathBox = widget.NewEntry()
	mstdwndPromoBox = widget.NewEntry()
	mstdwndCountBox = widget.NewEntry()
	mstdwndStatusLine = widget.NewLabel("等待操作")
	mstdwndPrecheckBtn = widget.NewButton("预览", func() { go mstdwndOnPrecheckPressed() })
	mstdwndReturnBtn = widget.NewButton("返回上级", func() { go mstdwndOnReturnPressed() })
	mstdwndSubmitBtn = widget.NewButton("更新商品", func() { go mstdwndOnSubmitPressed() })
	mstdwndPicturePreview = &canvas.Image{FillMode: canvas.ImageFillOriginal}
	mstdwndPicturePreview.SetMinSize(fyne.NewSize(250, 200))
	// start querying from
	var name, stype, description, picture string
	var price, promotion float64
	var stock_count int
	querystmt := candies.FastSQLPrep("SELECT name, type, description, price, picture, promotion, stock_count FROM stock_info WHERE stockid = ?",
		dbconn.Db, mainEventBus)
	queryerr := querystmt.QueryRow(selstockid).Scan(&name, &stype, &description, &price, &picture, &promotion, &stock_count)
	if queryerr != nil {
		log.Warnln("stock detail query fail:", queryerr.Error())
		mstdwndStatus(FAILMSG_DBINTERNAL)
	} else {
		// setups
		mstdwndNameBox.SetText(name)
		mstdwndTypeBox.SetText(stype)
		mstdwndDescBox.SetText(description)
		mstdwndPriceBox.SetText(candies.Ftoa(price))
		mstdwndPromoBox.SetText(candies.Ftoa(promotion))
		mstdwndCountBox.SetText(strconv.Itoa(stock_count))
		mstdwndPicPathBox.SetText(picture)
		mstdwndPicturePreview.File = mstdwndPicPathBox.Text
		canvas.Refresh(mstdwndPicturePreview)
	}
	// we dont need much preload in this scene. load layout now.
	mstdwndContainBox = widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("商品名称"), mstdwndNameBox),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("商品类型"), mstdwndTypeBox),
		widget.NewLabel("商品描述"), mstdwndDescBox,
		widget.NewLabel("图片路径"), mstdwndPicPathBox,
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("商品价格"), mstdwndPriceBox),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("折扣因数"), mstdwndPromoBox),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("商品数量"), mstdwndCountBox),
		fyne.NewContainerWithLayout(layout.NewGridLayout(3),
			mstdwndPrecheckBtn, mstdwndReturnBtn, mstdwndSubmitBtn),
		mstdwndStatusLine)
	currwnd = currapp.NewWindow("修改商品详情")
	currwnd.SetContent(
		fyne.NewContainerWithLayout(layout.NewBorderLayout(mstdwndContainBox, nil, nil, nil),
			mstdwndContainBox,
			mstdwndPicturePreview))
	currwnd.Resize(fyne.NewSize(500, 700))
	currwnd.SetOnClosed(globalFastCloser)
}

func mstdwndOnPrecheckPressed() {
	mstdwndPrecheck()
}

func mstdwndOnReturnPressed() {
	log.Debugln("back to shop detail")
	currwnd.Hide()
	chanTransist <- WND_MERCHANT_SHOP_DETAIL
}

func mstdwndOnSubmitPressed() {
	mstdwndOperable(false)
	defer mstdwndOperable(true)
	if !mstdwndPrecheck() {
		return
	}
	log.Debugln("precheck pass, start submitting")
	submitstmt := candies.FastSQLPrep(
		"UPDATE stock_info SET sold_at = ?, name = ?, type = ?, description = ?, price = ?, picture = ?, promotion = ?, stock_count = ?, enabled = ? WHERE stockid = ?",
		dbconn.Db, mainEventBus)
	sold_at := selshopid
	name := mstdwndNameBox.Text
	stype := mstdwndTypeBox.Text
	description := mstdwndDescBox.Text
	price, _ := candies.Atof(mstdwndPriceBox.Text)
	picture := mstdwndPicPathBox.Text
	promotion, _ := candies.Atof(mstdwndPromoBox.Text)
	stock_count, _ := strconv.Atoi(mstdwndCountBox.Text)
	enabled := true
	_, execerr := submitstmt.Exec(sold_at, name, stype, description, price, picture, promotion, stock_count, enabled, selstockid)
	if execerr != nil {
		log.Warnln("stock add fail:", execerr.Error())
		mstdwndStatus(FAILMSG_DBINTERNAL)
	} else {
		log.Succln("new stock succ")
		mstdwndStatus("修改成功")
	}
}

func mstdwndPrecheck() bool {
	// check the price & promo & count: they have to be numbers
	_, convsucc := strconv.Atoi(mstdwndCountBox.Text)
	if convsucc != nil {
		mstdwndStatus("商品数量无效")
		return false
	}
	_, convsucc = candies.Atof(mstdwndPriceBox.Text)
	if convsucc != nil {
		mstdwndStatus("商品价格无效")
		return false
	}
	_, convsucc = candies.Atof(mstdwndPromoBox.Text)
	if convsucc != nil {
		mstdwndStatus("折扣因数无效")
		return false
	}
	if mstdwndPicPathBox.Text == "" {
		mstdwndStatus("请重新设定图片")
	}
	// do the picture job
	mstdwndPicturePreview.File = mstdwndPicPathBox.Text
	canvas.Refresh(mstdwndPicturePreview)
	currwnd.Resize(fyne.NewSize(500, 700))
	mstdwndStatus("预览完毕")
	return true
}

func mstdwndStatus(s string) {
	mstdwndStatusLine.SetText(s)
}

func mstdwndOperable(ena bool) {
	if ena {
		mstdwndPrecheckBtn.Show()
		mstdwndReturnBtn.Show()
		mstdwndSubmitBtn.Show()
	} else {
		mstdwndPrecheckBtn.Hide()
		mstdwndReturnBtn.Hide()
		mstdwndSubmitBtn.Hide()
	}
}
