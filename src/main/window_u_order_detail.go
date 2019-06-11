package main

import (
	"candies"
	"dbconn"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

var uodwndStatusLine *widget.Label
var uodwndOrderIdBox *widget.Label
var uodwndCustNameBox *widget.Label
var uodwndShopNameBox *widget.Label
var uodwndStockNameBox *widget.Label
var uodwndAmountBox *widget.Label
var uodwndPriceBox *widget.Label
var uodwndTotalBox *widget.Label
var uodwndPaymentBox *widget.Entry
var uodwndDeliveryBox *widget.Entry
var uodwndOrderStatBox *widget.Label
var uodwndSetPaidBtn *widget.Button
var uodwndSetFinishedBtn *widget.Button
var uodwndSetCancelBtn *widget.Button
var uodwndReturnBtn *widget.Button

var uodwndOrdStat string

func setupWndUserOrderDetail() {
	// setup widgets
	uodwndStatusLine = widget.NewLabelWithStyle("等待用户操作", fyne.TextAlignCenter, candies.NewTextStyle().Fin())
	uodwndOrderIdBox = widget.NewLabel("ID")
	uodwndCustNameBox = widget.NewLabel("CustName")
	uodwndShopNameBox = widget.NewLabel("ShopName")
	uodwndStockNameBox = widget.NewLabel("StockName")
	uodwndAmountBox = widget.NewLabel("Amnt")
	uodwndPriceBox = widget.NewLabel("Price")
	uodwndTotalBox = widget.NewLabel("Total")
	uodwndPaymentBox = widget.NewMultiLineEntry()
	uodwndPaymentBox.SetReadOnly(true)
	uodwndDeliveryBox = widget.NewMultiLineEntry()
	uodwndDeliveryBox.SetReadOnly(true)
	uodwndOrderStatBox = widget.NewLabel("Stat")
	uodwndSetPaidBtn = widget.NewButton("设置为已付款", func() { go uodwndOnPaidPressed() })
	uodwndSetFinishedBtn = widget.NewButton("设置为已完成", func() { go uodwndOnFinishPressed() })
	uodwndSetCancelBtn = widget.NewButton("取消该订单", func() { go uodwndOnCancelPressed() })
	uodwndReturnBtn = widget.NewButton("返回上级", func() { go uodwndOnReturnPressed() })
	// setup layout to ensure labels works
	currwnd = currapp.NewWindow("订单详情")
	currwnd.SetContent(widget.NewVBox(
		uodwndOrderIdBox,
		widget.NewLabel("客户姓名"), uodwndCustNameBox,
		widget.NewLabel("商店名称"), uodwndShopNameBox,
		widget.NewLabel("商品名称"), uodwndStockNameBox,
		widget.NewLabel("购买数量"), uodwndAmountBox,
		widget.NewLabel("单个价格"), uodwndPriceBox,
		widget.NewLabel("总计金额"), uodwndTotalBox,
		widget.NewLabel("付款方式"), uodwndPaymentBox,
		widget.NewLabel("配送信息"), uodwndDeliveryBox,
		uodwndOrderStatBox,
		uodwndStatusLine,
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			uodwndSetPaidBtn, uodwndSetFinishedBtn),
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			uodwndSetCancelBtn, uodwndReturnBtn)))
	// start query information for these labels
	log.Debugln("start querying order detail")
	orderstmt := candies.FastSQLPrep("SELECT cust, stock, shop, amount, inst_price, price_sum, payment_detail, delivery_detail, stat FROM order_info WHERE orderid = ?", dbconn.Db, mainEventBus)
	var cust, stock, shop, amount int
	var inst_price, price_sum float64
	var payment_detail, delivery_detail, custname, stockname, shopname string
	ordererr := orderstmt.QueryRow(selorderid).Scan(&cust, &stock, &shop, &amount, &inst_price, &price_sum, &payment_detail, &delivery_detail, &uodwndOrdStat)
	if ordererr != nil {
		uodwndStatusLine.SetText(FAILMSG_DBINTERNAL)
	} else {
		custname = uodwndFastGetCustName(cust)
		stockname = uodwndFastGetStockName(stock)
		shopname = uodwndFastGetShopName(shop)
		// setup label informations
		uodwndOrderIdBox.SetText(strconv.Itoa(selorderid))
		uodwndCustNameBox.SetText(custname)
		uodwndShopNameBox.SetText(shopname)
		uodwndStockNameBox.SetText(stockname)
		uodwndAmountBox.SetText(strconv.Itoa(amount))
		uodwndPriceBox.SetText(candies.Ftoa(inst_price))
		uodwndTotalBox.SetText(candies.Ftoa(price_sum))
		uodwndPaymentBox.SetText(payment_detail)
		uodwndDeliveryBox.SetText(delivery_detail)
		uodwndOrderStatBox.SetText(fancyOrderStat(uodwndOrdStat))
		// complete.
	}
	// resize it
	currwnd.Resize(fyne.NewSize(400, 600))
}

func uodwndOnPaidPressed() {
	// allowed transist: issued -> paid
	if uodwndOrdStat != "issued" {
		uodwndStatusLine.SetText("无效操作")
		return
	}
	updquery := candies.FastSQLPrep("UPDATE order_info SET stat = ? WHERE orderid = ?",
		dbconn.Db, mainEventBus)
	_, upderr := updquery.Exec("paid", selorderid)
	if upderr != nil {
		uodwndStatusLine.SetText(FAILMSG_DBINTERNAL)
		return
	}
	uodwndOrdStat = "paid"
	uodwndOrderStatBox.SetText(fancyOrderStat(uodwndOrdStat))
	uodwndStatusLine.SetText("操作成功")
}

func uodwndOnFinishPressed() {
	// allowed transist: delivered -> finished
	if uodwndOrdStat != "delivered" {
		uodwndStatusLine.SetText("无效操作")
		return
	}
	updquery := candies.FastSQLPrep("UPDATE order_info SET stat = ? WHERE orderid = ?",
		dbconn.Db, mainEventBus)
	_, upderr := updquery.Exec("finished", selorderid)
	if upderr != nil {
		uodwndStatusLine.SetText(FAILMSG_DBINTERNAL)
		return
	}
	uodwndOrdStat = "finished"
	uodwndOrderStatBox.SetText(fancyOrderStat(uodwndOrdStat))
	uodwndStatusLine.SetText("操作成功")
}

func uodwndOnCancelPressed() {
	// allowed transist: issued -> cancelled
	if uodwndOrdStat != "issued" {
		uodwndStatusLine.SetText("无效操作")
		return
	}
	updquery := candies.FastSQLPrep("UPDATE order_info SET stat = ? WHERE orderid = ?",
		dbconn.Db, mainEventBus)
	_, upderr := updquery.Exec("cancelled", selorderid)
	if upderr != nil {
		uodwndStatusLine.SetText(FAILMSG_DBINTERNAL)
		return
	}
	uodwndOrdStat = "cancelled"
	uodwndOrderStatBox.SetText(fancyOrderStat(uodwndOrdStat))
	uodwndStatusLine.SetText("操作成功")
}

func uodwndOnReturnPressed() {
	currwnd.Hide()
	chanTransist <- WND_USER_ORDER_LIST
}

func uodwndFastGetStockName(stockid int) string {
	snstmt := candies.FastSQLPrep("SELECT name FROM stock_info WHERE stockid = ?",
		dbconn.Db, mainEventBus)
	snstr := ""
	snerr := snstmt.QueryRow(stockid).Scan(&snstr)
	if snerr != nil {
		snstr = "!!数据库错误!!"
	}
	return snstr
}

func uodwndFastGetCustName(custid int) string {
	snstmt := candies.FastSQLPrep("SELECT name FROM user_info WHERE uid = ?",
		dbconn.Db, mainEventBus)
	snstr := ""
	snerr := snstmt.QueryRow(custid).Scan(&snstr)
	if snerr != nil {
		snstr = "!!数据库错误!!"
	}
	return snstr
}

func uodwndFastGetShopName(shopid int) string {
	snstmt := candies.FastSQLPrep("SELECT name FROM shop_info WHERE shopid = ?",
		dbconn.Db, mainEventBus)
	snstr := ""
	snerr := snstmt.QueryRow(shopid).Scan(&snstr)
	if snerr != nil {
		snstr = "!!数据库错误!!"
	}
	return snstr
}
