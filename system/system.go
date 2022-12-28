package system

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func CurrentTime() string {
	theTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond(), time.Local)
	date := theTime.Format("2006-01-02 03:04:05 pm")
	return date
}

func GetDay() string {
	date := strings.Split(CurrentTime(), " ")
	return date[0]
}

func Round(result float64, lugares int) float64 {
	var mostrar float64

	if lugares == 1 {
		mostrar = math.Round(result)
	} else if lugares == 2 {
		redondeo, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", result), 64)
		mostrar = redondeo
	} else if lugares == 3 {
		redondeo, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", result), 64)
		mostrar = redondeo
	} else {
		redondeo, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", result), 64)
		mostrar = redondeo
	}

	return mostrar
}

func CreateDirectoryIfDoesntExist(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.Mkdir(folder, 0755)
		if err != nil {
			// Aquí puedes manejar mejor el error, es un ejemplo
			fmt.Println(err)
			os.Exit(0)
		}
	}
}
