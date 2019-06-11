package main

import (
	"candies"
	"dbconn"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

var uswndSearchCondBox *widget.Entry
var uswndSearchBtn *widget.Button
var uswndBackBtn *widget.Button
var uswndDisplayBox *widget.Box
var uswndStatusLine *widget.Label

var uswndCondStr string

func setupWndUserSearch() {
	uswndSearchCondBox = widget.NewEntry()
	uswndSearchBtn = widget.NewButton("搜索", func() { go uswndOnSearchPressed() })
	uswndBackBtn = widget.NewButton("返回", func() { go uswndOnBackPressed() })
	uswndStatusLine = widget.NewLabelWithStyle("等待用户操作", fyne.TextAlignCenter, candies.NewTextStyle().Fin())
	uswndDisplayBox = widget.NewVBox()
	uswndSearchCondBox.SetText(uswndCondStr)
	uswndDisplayBox.Append(widget.NewLabel("搜索条件"))
	uswndDisplayBox.Append(uswndSearchCondBox)
	uswndDisplayBox.Append(fyne.NewContainerWithLayout(
		layout.NewGridLayout(2),
		uswndSearchBtn, uswndBackBtn))
	uswndDisplayBox.Append(uswndStatusLine)
	if uswndCondStr != "" {
		searchstmt := candies.FastSQLPrep("SELECT stockid, name, price, promotion FROM stock_info WHERE name LIKE ? AND stock_count > 0",
			dbconn.Db, mainEventBus)
		// Execute it
		searchreslt, searcherr := searchstmt.Query("%" + uswndCondStr + "%")
		if searcherr != nil {
			log.Warnln("Search fail:", searcherr.Error())
			uswndStatusLine.SetText(FAILMSG_DBINTERNAL)
			return
		}
		for true {
			hasnext := searchreslt.Next()
			if !hasnext {
				log.Debugln("search no next")
				break
			}
			var stockid int
			var name string
			var price, promo float64
			searchreslt.Scan(&stockid, &name, &price, &promo)
			log.Debugln("Search reslt:", stockid, name, price, promo)
			uswndDisplayBox.Append(widget.NewButton(
				name+"| ￥"+candies.Ftoa(price)+", 折扣为"+candies.Ftoa((1.00-promo)*100.0)+"%",
				func() { go uswndOnStockSelected(stockid) }))
		}
		uswndStatusLine.SetText("搜索完成")
	}
	currwnd.Hide()
	currwnd = currapp.NewWindow("搜索")
	currwnd.SetContent(
		widget.NewScrollContainer(
			uswndDisplayBox))
	currwnd.SetOnClosed(globalFastCloser)
	currwnd.Resize(fyne.NewSize(300, 600))
}

func uswndOnSearchPressed() {
	uswndStatusLine.SetText("开始检索")
	uswndCondStr = uswndSearchCondBox.Text
	chanTransist <- WND_USER_SEARCH
}

func uswndOnBackPressed() {
	currwnd.Hide()
	uswndCondStr = ""
	chanTransist <- WND_USER_MAIN
}

func uswndOnStockSelected(stockid int) {
	selstockid = stockid
	currwnd.Hide()
	chanTransist <- WND_USER_STOCK_DETAIL
}
