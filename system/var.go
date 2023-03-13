package system

import (
	"github.com/akosej/lookcopys/models"
	"os"
)

var (
	Path, _ = os.Getwd()
	KB      = uint64(1024)
	MB      = KB * KB
	GB      = MB * KB
	USB     = map[string]models.InfoUsb{}
	ACTION  = map[string]models.Records{}
)
