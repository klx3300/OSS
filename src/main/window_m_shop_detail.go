package main

import (
	"candies"
	"dbconn"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

var msdwndShopNameBox *widget.Entry
var msdwndShopNameUpdateBtn *widget.Button
var msdwndReturnBtn *widget.Button
var msdwndStockListBox *widget.Box
var msdwndAddStockBtn *widget.Button
var msdwndOrderListBox *widget.Box
var msdwndRefreshBtn *widget.Button
var msdwndStatusLine *widget.Label

var msdwndStockMap map[int]string

const FAILMSG_DBINTERNAL = "数据库内部错误，请重启程序"

func setupWndMerchantShopDetail() {
	msdwndStockMap = make(map[int]string)
	log.Debugln("Start merchant shop detail window initialization")
	msdwndShopNameBox = widget.NewEntry()
	msdwndShopNameBox.SetText(selshopnm)
	msdwndShopNameUpdateBtn = widget.NewButton("商店更名", func() { go msdwndOnRenamePressed() })
	msdwndReturnBtn = widget.NewButton("返回上级", func() { go msdwndOnReturnPressed() })
	msdwndStockListBox = widget.NewVBox()
	msdwndAddStockBtn = widget.NewButton("上新商品", func() { go msdwndOnAddStockPressed() })
	msdwndOrderListBox = widget.NewVBox()
	msdwndRefreshBtn = widget.NewButton("刷新", func() { go msdwndOnRefreshPressed() })
	msdwndStatusLine = widget.NewLabelWithStyle("等待用户操作", fyne.TextAlignCenter, candies.NewTextStyle().Fin())
	log.Debugln("Start loading information for shopid", selshopid)
	// stock list first. display: stockid, name, price, promo, stock
	allstockstmt := candies.FastSQLPrep("SELECT stockid, name, price, promotion, stock_count FROM stock_info WHERE sold_at = ?", dbconn.Db, mainEventBus)
	allstockreslt, allstockerr := allstockstmt.Query(selshopid)
	if allstockerr != nil {
		log.Warnln("Query stock failure:", allstockerr.Error())
		msdwndStatusLine.SetText(FAILMSG_DBINTERNAL)
	} else {
		for true {
			var stockid int
			var stocknm string
			var stockprice float64
			var stockpromo float64
			var stockcnt int
			hasnext := allstockreslt.Next()
			if !hasnext {
				log.Debugln("stock not next")
				break
			}
			err := allstockreslt.Scan(&stockid, &stocknm, &stockprice, &stockpromo, &stockcnt)
			if err != nil {
				log.Warnln("Scan failure:", err.Error())
				msdwndStatusLine.SetText(FAILMSG_DBINTERNAL)
				break
			}
			log.Debugln("Stock", stockid, stocknm, stockprice, stockpromo, stockcnt)
			msdwndStockMap[stockid] = stocknm
			msdwndStockListBox.Append(
				widget.NewButton("No."+strconv.Itoa(stockid)+", "+stocknm+", 价格: "+
					candies.Ftoa(stockprice)+", 折扣: "+candies.Ftoa(stockpromo)+", 余货: "+strconv.Itoa(stockcnt),
					func() { go msdwndOnStockSelect(stockid) }))
		}
	}
	// start parsing for orders. much alike let's copy
	allorderstmt := candies.FastSQLPrep("SELECT orderid, stock, amount, stat FROM order_info WHERE shop = ?", dbconn.Db, mainEventBus)
	allorderreslt, allordererr := allorderstmt.Query(selshopid)
	if allordererr != nil {
		log.Warnln("Query order failure:", allordererr.Error())
		msdwndStatusLine.SetText(FAILMSG_DBINTERNAL)
	} else {
		for true {
			var orderid, stockid, amount int
			var orderstat string
			hasnext := allorderreslt.Next()
			if !hasnext {
				log.Debugln("order not next")
				break
			}
			err := allorderreslt.Scan(&orderid, &stockid, &amount, &orderstat)
			if err != nil {
				log.Warnln("Scan failure:", err.Error())
				msdwndStatusLine.SetText(FAILMSG_DBINTERNAL)
				break
			}
			log.Debugln("Order", orderid, stockid, amount, orderstat)
			msdwndOrderListBox.Append(
				widget.NewButton("No."+strconv.Itoa(orderid)+", "+msdwndStockMap[stockid]+", 数量: "+
					strconv.Itoa(amount)+", 状态: "+fancyOrderStat(orderstat),
					func() { go msdwndOnOrderSelect(orderid) }))
		}
	}
	log.Debugln("setting up shop detail layout")
	currwnd = currapp.NewWindow("商店详情页面")
	currwnd.SetContent(
		widget.NewScrollContainer(
			widget.NewVBox(
				fyne.NewContainerWithLayout(layout.NewGridLayout(2),
					widget.NewLabel("商店名称更新"), msdwndShopNameBox),
				msdwndShopNameUpdateBtn,
				msdwndStatusLine,
				fyne.NewContainerWithLayout(layout.NewGridLayout(3),
					msdwndReturnBtn, msdwndAddStockBtn, msdwndRefreshBtn),
				widget.NewLabelWithStyle("商店货物列表", fyne.TextAlignCenter, candies.NewTextStyle().Fin()),
				msdwndStockListBox,
				widget.NewLabelWithStyle("商店订单列表", fyne.TextAlignCenter, candies.NewTextStyle().Fin()),
				msdwndOrderListBox)))
	currwnd.Resize(fyne.NewSize(400, 600))
	currwnd.SetOnClosed(globalFastCloser)
	log.Debugln("Shop detail succeeded.")
}

func msdwndOnStockSelect(stockid int) {
	selstockid = stockid
	currwnd.Hide()
	chanTransist <- WND_MERCHANT_STOCK_DETAIL
}

func msdwndOnOrderSelect(orderid int) {
	selorderid = orderid
	currwnd.Hide()
	chanTransist <- WND_MERCHANT_ORDER_DETAIL
}

func msdwndOnRenamePressed() {
	msdwndOperable(false)
	defer msdwndOperable(true)
	msdwndStatusLine.SetText("请求中")
	updstmt := candies.FastSQLPrep("UPDATE shop_info SET name = ? WHERE shopid = ?", dbconn.Db, mainEventBus)
	_, execErr := updstmt.Exec(msdwndShopNameBox.Text, selshopid)
	if execErr != nil {
		log.Warnln("update shop name fail:", execErr.Error())
		mmwndStatusLine.SetText(FAILMSG_DBINTERNAL)
		return
	}
	log.Succln("shop name update to", msdwndShopNameBox.Text)
	msdwndStatusLine.SetText("更新成功")
}

func msdwndOnReturnPressed() {
	log.Debugln("return to merchant main")
	currwnd.Hide()
	chanTransist <- WND_MERCHANT_MAIN
}

func msdwndOnAddStockPressed() {
	currwnd.Hide()
	chanTransist <- WND_MERCHANT_STOCK_ADD
}

func msdwndOnRefreshPressed() {
	log.Debugln("force reload shop detail")
	currwnd.Hide()
	chanTransist <- WND_MERCHANT_SHOP_DETAIL
}

func msdwndOperable(ena bool) {
	if ena {
		msdwndShopNameUpdateBtn.Show()
		msdwndReturnBtn.Show()
		msdwndAddStockBtn.Show()
		msdwndRefreshBtn.Show()
		msdwndOrderListBox.Show()
		msdwndStockListBox.Show()
	} else {
		msdwndShopNameUpdateBtn.Hide()
		msdwndReturnBtn.Hide()
		msdwndAddStockBtn.Hide()
		msdwndRefreshBtn.Hide()
		msdwndStockListBox.Hide()
		msdwndOrderListBox.Hide()
	}
}

func fancyOrderStat(stat string) string {
	if stat == "issued" {
		return "已下单"
	}
	if stat == "paid" {
		return "已付款"
	}
	if stat == "delivered" {
		return "已发货"
	}
	if stat == "finished" {
		return "已完成"
	}
	if stat == "cancelled" {
		return "已取消"
	}
	return "未知错误"
}
