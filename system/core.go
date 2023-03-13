package system

import (
	"github.com/akosej/lookcopys/models"
	"github.com/jaypipes/ghw"
)

func IfSerial(reg []interface{}, serial string) bool {
	for _, v := range reg {
		if v == serial {
			return true
		}
	}
	return false
}
func IsTheSerial(reg []models.InfoUsb, serial string) bool {
	for _, v := range reg {
		if v.Serial == serial {
			return true
		}
	}
	return false
}
func isInsertSerial(serial string) (bool, string) {
	for _, v := range USB {
		if v.Serial == serial {
			return true, v.Model
		}
	}
	return false, ""
}

func IsDayState(states []models.State) bool {
	day := GetDay()
	for _, v := range states {
		if v.Day == day {
			return true
		}
	}
	return false
}

func GetCopiedDay(regs []models.InfoUsb) float64 {
	var copied float64
	for _, reg := range regs {
		if reg.Date == GetDay() {
			copied += float64(reg.Copy)
		}
	}
	return copied
}

func GetInfoDay(regs []models.InfoUsb) []models.InfoUsb {
	for i, reg := range regs {
		if reg.Date != GetDay() {
			if i+1 <= len(regs) {
				regs = append(regs[:i], regs[i+1:]...)
			}
		}
	}
	return regs
}

func GetModelUsd(size uint64) (string, string) {
	block, _ := ghw.Block()
	model := ""
	serial := ""

	for _, disk := range block.Disks {
		result, mod := isInsertSerial(disk.SerialNumber)
		if !result {
			if mod != disk.Model {
				if size == disk.SizeBytes/(KB*KB*KB) {
					model = disk.Model
					serial = disk.SerialNumber
					break
				}
			}
		}

	}
	return model, serial
}
