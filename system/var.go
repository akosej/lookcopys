package system

import (
	"os"
	"usbWatcher/models"
)

var (
	Path, _   = os.Getwd()
	KB        = uint64(1024)
	MB        = KB * KB
	GB        = MB * KB
	USB       = map[string]models.InfoUsb{}
)
