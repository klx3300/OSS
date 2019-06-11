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

var ustdwndNameBox *widget.Entry
var ustdwndTypeBox *widget.Entry
var ustdwndDescBox *widget.Entry
var ustdwndPriceBox *widget.Entry
var ustdwndPromoBox *widget.Entry
var ustdwndCountBox *widget.Entry
var ustdwndOrderBtn *widget.Button
var ustdwndReturnBtn *widget.Button
var ustdwndStatusLine *widget.Label
var ustdwndPicturePreview *canvas.Image
var ustdwndContainBox *widget.Box

func setupWndUserStockDetail() {
	ustdwndNameBox = widget.NewEntry()
	ustdwndTypeBox = widget.NewEntry()
	ustdwndDescBox = widget.NewMultiLineEntry()
	ustdwndPriceBox = widget.NewEntry()
	ustdwndPromoBox = widget.NewEntry()
	ustdwndCountBox = widget.NewEntry()
	ustdwndStatusLine = widget.NewLabel("等待操作")
	ustdwndReturnBtn = widget.NewButton("返回上级", func() { go ustdwndOnReturnPressed() })
	ustdwndOrderBtn = widget.NewButton("立即下单", func() { go ustdwndOnOrderPressed() })
	ustdwndPicturePreview = &canvas.Image{FillMode: canvas.ImageFillOriginal}
	ustdwndPicturePreview.SetMinSize(fyne.NewSize(250, 200))
	// start querying from
	var name, stype, description, picture string
	var price, promotion float64
	var stock_count int
	querystmt := candies.FastSQLPrep("SELECT name, type, description, price, picture, promotion, stock_count, sold_at FROM stock_info WHERE stockid = ?",
		dbconn.Db, mainEventBus)
	queryerr := querystmt.QueryRow(selstockid).Scan(&name, &stype, &description, &price, &picture, &promotion, &stock_count, &selordershop)
	if queryerr != nil {
		log.Warnln("stock detail query fail:", queryerr.Error())
		ustdwndStatus(FAILMSG_DBINTERNAL)
	} else {
		// setups
		log.Debugln("Stock info", name, stype, description, price, promotion, stock_count, picture, "@", selordershop)
		selstocknm = name // the fast path..
		ustdwndNameBox.SetText(name)
		ustdwndTypeBox.SetText(stype)
		ustdwndDescBox.SetText(description)
		ustdwndPriceBox.SetText(candies.Ftoa(price))
		ustdwndPromoBox.SetText(candies.Ftoa(promotion))
		ustdwndCountBox.SetText(strconv.Itoa(stock_count))
		ustdwndPicturePreview.File = picture
		canvas.Refresh(ustdwndPicturePreview)
	}
	// we dont need much preload in this scene. load layout now.
	ustdwndContainBox = widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("商品名称"), ustdwndNameBox),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("商品类型"), ustdwndTypeBox),
		widget.NewLabel("商品描述"), ustdwndDescBox,
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("商品价格"), ustdwndPriceBox),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("折扣因数"), ustdwndPromoBox),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabel("商品数量"), ustdwndCountBox),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			ustdwndReturnBtn, ustdwndOrderBtn),
		ustdwndStatusLine)
	currwnd = currapp.NewWindow("查看商品详情")
	currwnd.SetContent(
		fyne.NewContainerWithLayout(layout.NewBorderLayout(ustdwndContainBox, nil, nil, nil),
			ustdwndContainBox,
			ustdwndPicturePreview))
	currwnd.Resize(fyne.NewSize(500, 700))
	currwnd.SetOnClosed(globalFastCloser)
}

func ustdwndOnReturnPressed() {
	log.Debugln("back to search")
	currwnd.Hide()
	chanTransist <- WND_USER_SEARCH
}

func ustdwndOnOrderPressed() {
	currwnd.Hide()
	chanTransist <- WND_USER_NEW_ORDER
}

func ustdwndStatus(s string) {
	ustdwndStatusLine.SetText(s)
}

func ustdwndOperable(ena bool) {
	if ena {
		ustdwndReturnBtn.Show()
		ustdwndOrderBtn.Show()
	} else {
		ustdwndReturnBtn.Hide()
		ustdwndOrderBtn.Hide()
	}
}
