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

var msawndNameBox *widget.Entry
var msawndTypeBox *widget.Entry
var msawndDescBox *widget.Entry
var msawndPriceBox *widget.Entry
var msawndPicPathBox *widget.Entry
var msawndPromoBox *widget.Entry
var msawndCountBox *widget.Entry
var msawndPrecheckBtn *widget.Button
var msawndSubmitBtn *widget.Button
var msawndReturnBtn *widget.Button
var msawndStatusLine *widget.Label
var msawndPicturePreview *canvas.Image
var msawndContainBox *widget.Box

func setupWndMerchantStockAdd() {
	// setup parts
	msawndNameBox = widget.NewEntry()
	msawndTypeBox = widget.NewEntry()
	msawndDescBox = widget.NewMultiLineEntry()
	msawndPriceBox = widget.NewEntry()
	msawndPicPathBox = widget.NewEntry()
	msawndPromoBox = widget.NewEntry()
	msawndCountBox = widget.NewEntry()
	msawndStatusLine = widget.NewLabel("等待操作")
	msawndPrecheckBtn = widget.NewButton("预览", func() { go msawndOnPrecheckPressed() })
	msawndReturnBtn = widget.NewButton("返回上级", func() { go msawndOnReturnPressed() })
	msawndSubmitBtn = widget.NewButton("上架新商品", func() { go msawndOnSubmitPressed() })
	msawndPicturePreview = &canvas.Image{FillMode: canvas.ImageFillOriginal}
	msawndPicturePreview.SetMinSize(fyne.NewSize(250, 200))
	msawndContainBox = widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("商品名称"), msawndNameBox),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("商品类型"), msawndTypeBox),
		widget.NewLabel("商品描述"), msawndDescBox,
		widget.NewLabel("图片路径"), msawndPicPathBox,
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("商品价格"), msawndPriceBox),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("折扣因数"), msawndPromoBox),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("商品数量"), msawndCountBox),
		fyne.NewContainerWithLayout(layout.NewGridLayout(3),
			msawndPrecheckBtn, msawndReturnBtn, msawndSubmitBtn),
		msawndStatusLine)

	// we dont need much preload in this scene. load layout now.
	currwnd = currapp.NewWindow("新增商品")
	currwnd.SetContent(
		fyne.NewContainerWithLayout(layout.NewBorderLayout(msawndContainBox, nil, nil, nil),
			msawndContainBox,
			msawndPicturePreview))
	currwnd.Resize(fyne.NewSize(500, 700))
	currwnd.SetOnClosed(globalFastCloser)
}

func msawndOnPrecheckPressed() {
	msawndPrecheck()
}

func msawndOnReturnPressed() {
	log.Debugln("back to shop detail")
	currwnd.Hide()
	chanTransist <- WND_MERCHANT_SHOP_DETAIL
}

func msawndOnSubmitPressed() {
	msawndOperable(false)
	defer msawndOperable(true)
	if !msawndPrecheck() {
		return
	}
	log.Debugln("precheck pass, start submitting")
	submitstmt := candies.FastSQLPrep(
		"INSERT INTO stock_info (sold_at, name, type, description, price, picture, promotion, stock_count, enabled) VALUES (?,?,?,?,?,?,?,?,?)",
		dbconn.Db, mainEventBus)
	sold_at := selshopid
	name := msawndNameBox.Text
	stype := msawndTypeBox.Text
	description := msawndDescBox.Text
	price, _ := candies.Atof(msawndPriceBox.Text)
	picture := msawndPicPathBox.Text
	promotion, _ := candies.Atof(msawndPromoBox.Text)
	stock_count, _ := strconv.Atoi(msawndCountBox.Text)
	enabled := true
	_, execerr := submitstmt.Exec(sold_at, name, stype, description, price, picture, promotion, stock_count, enabled)
	if execerr != nil {
		log.Warnln("stock add fail:", execerr.Error())
		msawndStatus(FAILMSG_DBINTERNAL)
	} else {
		log.Succln("new stock succ")
		msawndStatus("上新成功")
	}
}

func msawndPrecheck() bool {
	// check the price & promo & count: they have to be numbers
	_, convsucc := strconv.Atoi(msawndCountBox.Text)
	if convsucc != nil {
		msawndStatus("商品数量无效")
		return false
	}
	_, convsucc = candies.Atof(msawndPriceBox.Text)
	if convsucc != nil {
		msawndStatus("商品价格无效")
		return false
	}
	_, convsucc = candies.Atof(msawndPromoBox.Text)
	if convsucc != nil {
		msawndStatus("折扣因数无效")
		return false
	}
	// do the picture job
	msawndPicturePreview.File = msawndPicPathBox.Text
	canvas.Refresh(msawndPicturePreview)
	currwnd.Resize(fyne.NewSize(500, 700))
	msawndStatus("预览完毕")
	return true
}

func msawndStatus(s string) {
	msawndStatusLine.SetText(s)
}

func msawndOperable(ena bool) {
	if ena {
		msawndPrecheckBtn.Show()
		msawndReturnBtn.Show()
		msawndSubmitBtn.Show()
	} else {
		msawndPrecheckBtn.Hide()
		msawndReturnBtn.Hide()
		msawndSubmitBtn.Hide()
	}
}
