package main

import (
	"candies"
	"dbconn"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

var uolwndContainer *widget.Box
var uolwndReturnBtn *widget.Button

func setupWndUserOrderList() {
	uolwndReturnBtn = widget.NewButton("返回上级", func() { go uolwndOnReturnPressed() })
	uolwndContainer = widget.NewVBox(uolwndReturnBtn)
	currwnd = currapp.NewWindow("订单列表")
	// start query
	log.Debugln("user order list query start")
	allorderstmt := candies.FastSQLPrep("SELECT orderid, stock, amount, stat FROM order_info WHERE cust = ?", dbconn.Db, mainEventBus)
	allorderreslt, allordererr := allorderstmt.Query(curruid)
	if allordererr != nil {
		uolwndContainer.Prepend(widget.NewLabel(FAILMSG_DBINTERNAL))
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
				uolwndContainer.Prepend(widget.NewLabel(FAILMSG_DBINTERNAL))
				break
			}
			log.Debugln("Order", orderid, stockid, amount, orderstat)
			uolwndContainer.Prepend(
				widget.NewButton("No."+strconv.Itoa(orderid)+", "+uolwndFastGetStockName(stockid)+" x"+
					strconv.Itoa(amount)+", 状态: "+fancyOrderStat(orderstat),
					func() { go uolwndOnOrderSelect(orderid) }))
		}
	}
	currwnd.SetContent(widget.NewScrollContainer(uolwndContainer))
	currwnd.SetOnClosed(globalFastCloser)
	currwnd.Resize(fyne.NewSize(400, 500))
}

func uolwndOnReturnPressed() {
	currwnd.Hide()
	chanTransist <- WND_USER_MAIN
}

func uolwndFastGetStockName(stockid int) string {
	snstmt := candies.FastSQLPrep("SELECT name FROM stock_info WHERE stockid = ?",
		dbconn.Db, mainEventBus)
	snstr := ""
	snerr := snstmt.QueryRow(stockid).Scan(&snstr)
	if snerr != nil {
		snstr = "!!数据库错误!!"
	}
	return snstr
}

func uolwndOnOrderSelect(orderid int) {
	selorderid = orderid
	currwnd.Hide()
	chanTransist <- WND_USER_ORDER_DETAIL
}
