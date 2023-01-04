package system

import (
	"github.com/gen2brain/beeep"
	"reflect"
	"tawesoft.co.uk/go/dialog"
)

func SendNotifyDesktop(title, msg string) {
	_ = beeep.Beep(600.0, 1000)
	_ = beeep.Notify(title, msg, "./frontend/look.png")
}
func SendAlertDesktop(arg interface{}) {
	dialog.Alert(reflect.ValueOf(arg).String())
}
