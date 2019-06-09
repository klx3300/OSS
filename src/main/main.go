package main

import (
	"candies"
	"dbconn"
	"evbus"
	"logger"
	"os"
	"osstheme"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

var mainEventBus = evbus.NewEventBus()
var log = logger.Logger{
	Name:     "Main:",
	LogLevel: 0,
}

var currapp fyne.App
var currwnd fyne.Window

const (
	WND_LOADING = iota
	WND_LOGIN
	WND_REGISTER
	WND_MERCHANT_MAIN
	WND_MERCHANT_SHOP_DETAIL
	WND_MERCHANT_STOCK_ADD
	WND_MERCHANT_STOCK_DETAIL
	WND_MERCHANT_ORDER_LIST
	WND_MERCHANT_ORDER_DETAIL
	WND_USER_MAIN
	WND_USER_ORDER_LIST
	WND_USER_ORDER_DETAIL
	WND_USER_SHOP_LIST
	WND_USER_SHOP_DETAIL
	WND_USER_STOCK_DETAIL
	WND_USER_NEW_ORDER
)

var chanTransist = make(chan int)
var nextTransist = 0

func main() {
	log.Debugln("Startup..")
	candies.Initialization(mainEventBus)
	log.Infoln("GUI Initializing...")
	osstheme.Initialization(mainEventBus)
	currapp = app.New()
	currapp.Settings().SetTheme(osstheme.CNFont{})
	currwnd = currapp.NewWindow("ZHWK Shopping System")
	currwnd.SetContent(
		widget.NewVBox(
			widget.NewLabelWithStyle(
				"< ZHWK 在线购物系统 >",
				fyne.TextAlignCenter, candies.NewTextStyle().SetMono().Fin()),
			widget.NewLabel("你的全新在线购物体验，从这里开始"),
			widget.NewLabelWithStyle("正在完成初始化...",
				fyne.TextAlignCenter, candies.NewTextStyle().Fin()),
		))
	currwnd.SetOnClosed(globalFastCloser)
	handleUnrecoverable()
	defer dbconn.Finalization(mainEventBus)
	// from this point on, we cannot use os.Exit() to forcefully quit the application.
	// publish UnrecoverableFailureEvent to  do this
	currwnd.Show()
	go transistionHandler()
	go preloginInitializer()
	currapp.Run()
}

func transistionHandler() {
	for true {
		nextTransist = <-chanTransist
		log.Debugln("New transist occurred.")
		currwnd.Hide()
		switch nextTransist {
		case WND_LOGIN:
			setupWndLogin()
			log.Debugln("Starting display..")
			currwnd.Show()
		case WND_REGISTER:
			setupWndRegister()
			log.Debugln("Starting display..")
			currwnd.Show()
		case WND_MERCHANT_MAIN:
		case WND_MERCHANT_SHOP_DETAIL:
		case WND_MERCHANT_STOCK_DETAIL:
		case WND_MERCHANT_STOCK_ADD:
		case WND_MERCHANT_ORDER_LIST:
		case WND_MERCHANT_ORDER_DETAIL:
		case WND_USER_MAIN:
		case WND_USER_SHOP_LIST:
		case WND_USER_SHOP_DETAIL:
		case WND_USER_NEW_ORDER:
		case WND_USER_STOCK_DETAIL:
		case WND_USER_ORDER_LIST:
		case WND_USER_ORDER_DETAIL:
		default:
			log.FailLn("transisting to undefined window", nextTransist)
			return
		}
	}
}

func globalFastCloser() {
	candies.FastGG("User Close", mainEventBus)
}

func handleUnrecoverable() {
	log := logger.Logger{
		Name:     "Finalizer:",
		LogLevel: 0,
	}
	var subscriber evbus.Subscriber
	subscriber.SubscribedTypes = []int{candies.EventTypeUnrecoverableFail}
	subscriber.SubscriberPriority = -999
	subscriber.MessageFunctor = func(s evbus.Subscriber, e evbus.Event) bool {
		if e.Payload.(string) == "User Close" {
			log.Debugln("User close detected.")
			log.Debugln("Closing database connection...")
			dbconn.Finalization(mainEventBus)
			os.Exit(0)
		} else {
			log.FailLn("Unrecoverable Failure caught:", e.Payload.(string))
			log.Debugln("Closing database connection...")
			dbconn.Finalization(mainEventBus)
			os.Exit(-1)
		}
		return true
	}
	mainEventBus.Subscribe(subscriber)
}
