package main

import (
	"candies"
	"dbconn"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type OrderedStockStat struct {
	remain int
	price  float64
	promo  float64
}

var unowndStockStat OrderedStockStat
var unowndStockIdBox *widget.Entry
var unowndStockNameBox *widget.Entry
var unowndStockRemainBox *widget.Entry
var unowndStockPriceBox *widget.Entry
var unowndStockPromoBox *widget.Entry
var unowndStockCostBox *widget.Entry
var unowndAmountBox *widget.Entry
var unowndTotalBox *widget.Entry
var unowndPaymentInfoBox *widget.Entry
var unowndDeliveryInfoBox *widget.Entry
var unowndPreviewBtn *widget.Button
var unowndConfirmBtn *widget.Button
var unowndReturnBtn *widget.Button
var unowndStatusLine *widget.Label

const MSG_PLEASEWAIT = "请等待"

func setupWndUserNewOrder() {
	unowndStockIdBox = widget.NewEntry()
	unowndStockIdBox.SetText(strconv.Itoa(selstockid))
	unowndStockIdBox.SetReadOnly(true)
	unowndStockNameBox = widget.NewEntry()
	unowndStockNameBox.SetText(selstocknm)
	unowndStockNameBox.SetReadOnly(true)
	unowndStockRemainBox = widget.NewEntry()
	unowndStockRemainBox.SetText(MSG_PLEASEWAIT)
	unowndStockRemainBox.SetReadOnly(true)
	unowndStockPriceBox = widget.NewEntry()
	unowndStockPriceBox.SetReadOnly(true)
	unowndStockPromoBox = widget.NewEntry()
	unowndStockPromoBox.SetReadOnly(true)
	unowndStockCostBox = widget.NewEntry()
	unowndStockCostBox.SetReadOnly(true)
	unowndAmountBox = widget.NewEntry()
	unowndTotalBox = widget.NewEntry()
	unowndTotalBox.SetReadOnly(true)
	unowndPaymentInfoBox = widget.NewMultiLineEntry()
	unowndDeliveryInfoBox = widget.NewMultiLineEntry()
	unowndStatusLine = widget.NewLabel("")
	unowndPreviewBtn = widget.NewButton("计算并检测", func() { go unowndOnPreviewPressed() })
	unowndConfirmBtn = widget.NewButton("立即下单", func() { go unowndOnConfirmPressed() })
	unowndReturnBtn = widget.NewButton("返回上级", func() { go unowndOnReturnPressed() })
	currwnd = currapp.NewWindow("新增订单")
	currwnd.SetContent(
		widget.NewVBox(
			widget.NewLabel("商品ID"),
			unowndStockIdBox,
			widget.NewLabel("商品名称"),
			unowndStockNameBox,
			widget.NewLabel("剩余库存"),
			unowndStockRemainBox,
			widget.NewLabel("当前价格"),
			unowndStockPriceBox,
			widget.NewLabel("当前折扣"),
			unowndStockPromoBox,
			widget.NewLabel("购买价格"),
			unowndStockCostBox,
			widget.NewLabel("购买数量"),
			unowndAmountBox,
			widget.NewLabel("总共支付"),
			unowndTotalBox,
			widget.NewLabel("支付信息"),
			unowndPaymentInfoBox,
			widget.NewLabel("配送信息"),
			unowndDeliveryInfoBox,
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				unowndPreviewBtn, unowndReturnBtn, unowndConfirmBtn),
			unowndStatusLine))
	unowndRequeryStockInfo()
	unowndRefreshCalculations()
	unowndStatusLine.SetText("等待用户操作")
	// setup layout
	currwnd.Resize(fyne.NewSize(300, 600))
	currwnd.SetOnClosed(globalFastCloser)
}

func unowndRequeryStockInfo() {
	log.Debugln("user new order refreshing..")
	unowndStatusLine.SetText("正在刷新商品当前数据")
	var price, promotion float64
	var stock_count int
	querystmt := candies.FastSQLPrep("SELECT price, promotion, stock_count FROM stock_info WHERE stockid = ?",
		dbconn.Db, mainEventBus)
	queryerr := querystmt.QueryRow(selstockid).Scan(&price, &promotion, &stock_count)
	if queryerr != nil {
		log.Warnln("stock detail query fail:", queryerr.Error())
		unowndStatusLine.SetText(FAILMSG_DBINTERNAL)
	} else {
		// setups
		log.Debugln("refreshed Stock info", price, promotion, stock_count)
		unowndStockStat.price = price
		unowndStockStat.promo = promotion
		unowndStockStat.remain = stock_count
	}
}

func unowndRefreshCalculations() {
	unowndStockPriceBox.SetText(candies.Ftoa(unowndStockStat.price))
	unowndStockPromoBox.SetText(candies.Ftoa(unowndStockStat.promo))
	unowndStockRemainBox.SetText(strconv.Itoa(unowndStockStat.remain))
	unowndStockCostBox.SetText(candies.Ftoa(unowndStockStat.price * unowndStockStat.promo))
	amnt, err := strconv.Atoi(unowndAmountBox.Text)
	if err == nil {
		unowndTotalBox.SetText(candies.Ftoa(float64(amnt) * unowndStockStat.price * unowndStockStat.promo))
	}
}

func unowndPrecheck() bool {
	// check input correction
	unowndRequeryStockInfo()
	defer unowndRefreshCalculations()
	amnt, err := strconv.Atoi(unowndAmountBox.Text)
	if err != nil {
		unowndStatusLine.SetText("数量无效")
		return false
	}
	if amnt > unowndStockStat.remain {
		unowndStatusLine.SetText("数量多于余量")
		return false
	}
	if len(unowndPaymentInfoBox.Text) < 10 {
		unowndStatusLine.SetText("请提供有效的付款信息")
		return false
	}
	if len(unowndDeliveryInfoBox.Text) < 10 {
		unowndStatusLine.SetText("请提供有效的配送信息")
		return false
	}
	unowndStatusLine.SetText("刷新完成")
	return true
}

func unowndOnPreviewPressed() {
	unowndOperable(false)
	defer unowndOperable(true)
	unowndPrecheck()
}

func unowndOnConfirmPressed() {
	unowndOperable(false)
	defer unowndOperable(true)
	infocorr := unowndPrecheck()
	if !infocorr {
		return
	}
	buyAmount, _ := strconv.Atoi(unowndAmountBox.Text)
	// decrease amount first
	amntstmt := candies.FastSQLPrep(
		"UPDATE stock_info SET stock_count = ? WHERE stockid = ?",
		dbconn.Db, mainEventBus)
	_, amnterr := amntstmt.Exec(unowndStockStat.remain-buyAmount, selstockid)
	if amnterr != nil {
		log.Warnln("update amount fail:", amnterr.Error())
		unowndStatusLine.SetText(FAILMSG_DBINTERNAL)
		return
	}
	orderstmt := candies.FastSQLPrep(
		"INSERT INTO order_info (cust, stock, shop, amount, inst_price, price_sum, payment_detail, delivery_detail, stat) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		dbconn.Db, mainEventBus)
	_, ordererr := orderstmt.Exec(curruid, selstockid, selordershop, buyAmount, unowndStockStat.price*unowndStockStat.promo,
		float64(buyAmount)*unowndStockStat.price*unowndStockStat.promo,
		unowndPaymentInfoBox.Text, unowndDeliveryInfoBox.Text, "issued")
	if ordererr != nil {
		log.Warnln("Failed to order:", ordererr)
		unowndStatusLine.SetText(FAILMSG_DBINTERNAL)
		return
	}
	unowndStatusLine.SetText("下单成功 冻结5秒以防误下单...")
	<-time.After(5 * time.Millisecond)
	unowndRequeryStockInfo()
	unowndRefreshCalculations()
}

func unowndOnReturnPressed() {
	// back to stock detail
	currwnd.Hide()
	chanTransist <- WND_USER_STOCK_DETAIL
}

func unowndOperable(ena bool) {
	if ena {
		unowndReturnBtn.Show()
		unowndConfirmBtn.Show()
		unowndPreviewBtn.Show()
	} else {
		unowndReturnBtn.Hide()
		unowndConfirmBtn.Hide()
		unowndPreviewBtn.Hide()
	}
}
