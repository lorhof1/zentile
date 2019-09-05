package state

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
	log "github.com/sirupsen/logrus"
)

var (
	X           *xgbutil.XUtil
	DeskCount   uint
	ActiveWin   xproto.Window
	CurrentDesk uint
	Stacking    []xproto.Window
	workArea    []ewmh.Workarea
)

func init() {
	populateState()

	win := xwindow.New(X, X.RootWin())
	win.Listen(xproto.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(stateUpdate).Connect(X, X.RootWin())
}

func populateState() {
	var err error
	X, err = xgbutil.NewConn()
	checkErr(err)

	DeskCount, err = ewmh.NumberOfDesktopsGet(X)
	checkErr(err)

	ActiveWin, err = ewmh.ActiveWindowGet(X)
	checkErr(err)

	CurrentDesk, err = ewmh.CurrentDesktopGet(X)
	checkErr(err)

	Stacking, err = ewmh.ClientListGet(X)
	checkErr(err)

	workArea, err = ewmh.WorkareaGet(X)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal("Error populating state: ", err)
	}
}

func stateUpdate(X *xgbutil.XUtil, e xevent.PropertyNotifyEvent) {
	var err error
	if aname, _ := xprop.AtomName(X, e.Atom); aname == "_NET_ACTIVE_WINDOW" {
		ActiveWin, err = ewmh.ActiveWindowGet(X)
	} else if aname, _ := xprop.AtomName(X, e.Atom); aname == "_NET_CURRENT_DESKTOP" {
		CurrentDesk, err = ewmh.CurrentDesktopGet(X)
	} else if aname, _ := xprop.AtomName(X, e.Atom); aname == "_NET_NUMBER_OF_DESKTOPS" {
		DeskCount, err = ewmh.NumberOfDesktopsGet(X)
	} else if aname, _ := xprop.AtomName(X, e.Atom); aname == "_NET_CLIENT_LIST_STACKING" {
		Stacking, err = ewmh.ClientListStackingGet(X)
	} else if aname, _ := xprop.AtomName(X, e.Atom); aname == "_NET_WORKAREA" {
		workArea, err = ewmh.WorkareaGet(X)
	}

	if err != nil {
		log.Warn("Error updating state: ", err)
	}
}

func WorkAreaDimensions(num uint) (x, y, width, height int) {
	w := workArea[num]
	x = w.X
	y = w.Y
	width = int(w.Width)
	height = int(w.Height)
	return
}
