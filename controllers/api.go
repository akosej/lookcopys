package controllers

import (
	"github.com/gofiber/fiber/v2"
	"usbWatcher/models"
	"usbWatcher/system"
)

func ApiStatus(c *fiber.Ctx) error {
	states := system.JsonReadState()
	var mb, gb, cash []float64
	for _, v := range states {
		mb = append(mb, v.CopiedMB)
		gb = append(gb, system.Round(v.CopiedGB, 2))
		cash = append(cash, system.Round(v.CopiedGB*5, 2))
	}
	var series []fiber.Map
	series = append(series, fiber.Map{
		"name": "MB",
		"type": "line",
		"data": mb,
	})
	series = append(series, fiber.Map{
		"name": "GB",
		"type": "column",
		"data": gb,
	})
	series = append(series, fiber.Map{
		"name": "$",
		"type": "column",
		"data": cash,
	})
	online := system.JsonReadInfoUsb()
	var today models.State
	for _, v := range states {
		if v.Day == system.GetDay() {
			today = v
		}
	}

	return c.JSON(fiber.Map{
		"today": today,
		"online": online,
		"states": states,
		"gb":     gb,
		"series": series,
	})
}
