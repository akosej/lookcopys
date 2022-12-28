package system

import (
	"fmt"
	set "github.com/deckarep/golang-set"
	detector "github.com/deepakjois/gousbdrivedetector"
	"github.com/ricochet2200/go-disk-usage/du"
	"os"
	"strconv"
	"time"
	"usbWatcher/models"
)

func MonitorUsb() {
	reg := JsonReadInfoUsb()
	logs := JsonReadLogs()
	states := JsonReadState()
	for {
		var UsbAll, UsbConnect []interface{}
		if len(reg) > 0 {
			for i, v := range reg {
				UsbAll = append(UsbAll, v.Serial)
				if v.Date != GetDay() {
					if i+1 <= len(reg) {
						reg = append(reg[:i], reg[i+1:]...)
					}
				}
			}
		} else {
			reg = []models.InfoUsb{}
		}
		if !IsDayState(states) {
			states = append(states, models.State{
				Day:       GetDay(),
				Connected: len(GetInfoDay(reg)),
				CopiedMB:  0,
				CopiedGB:  0,
			})
		}
		if drives, err := detector.Detect(); err == nil {
			if len(drives) > 0 {
				for _, path := range drives {
					usage := du.NewDiskUsage(path)
					value, exists := USB[path]
					//	--------------------------------------------------------------------
					if !exists {
						model, serial := GetModelUsd(usage.Size() / (GB))
						var copy uint64
						if !IsTheSerial(reg, serial) {
							for i, state := range states {
								if state.Day == GetDay() {
									states = append(states[:i], states[i+1:]...)
									// -----------------------
									states = append(states, models.State{
										Day:       GetDay(),
										Connected: state.Connected + 1,
										CopiedMB:  state.CopiedMB,
										CopiedGB:  state.CopiedGB,
									})
								}
							}
							UsbConnect = append(UsbConnect, serial)
							copy = 0
						} else {
							UsbConnect = append(UsbConnect, serial)
							for _, v := range reg {
								if v.Serial == serial {
									copy = v.Copy
								}
							}
						}
						logs = append(logs, models.Logs{
							Date:        CurrentTime(),
							Model:       model,
							Serial:      serial,
							Size:        usage.Size() / (GB),
							Description: "Connected",
						})

						USB[path] = models.InfoUsb{
							Path:   path,
							Date:   GetDay(),
							Serial: serial,
							Model:  model,
							Size:   usage.Size() / (MB),
							Used:   usage.Used() / (MB),
							Free:   usage.Free() / (MB),
							Copy:   copy,
						}
						//	--------------------------------------------------------------------
					} else {
						UsbConnect = append(UsbConnect, value.Serial)
						if value.Free > usage.Free()/(MB) {
							cop := value.Copy + (value.Free - usage.Free()/(MB))
							USB[path] = models.InfoUsb{
								Path:   path,
								Date:   GetDay(),
								Serial: value.Serial,
								Model:  value.Model,
								Size:   usage.Size() / (MB),
								Used:   usage.Used() / (MB),
								Free:   usage.Free() / (MB),
								Copy:   cop,
							}
						}
						//-- Si se borra
						if value.Free < usage.Free()/(MB) {
							USB[path] = models.InfoUsb{
								Path:   path,
								Serial: value.Serial,
								Model:  value.Model,
								Size:   usage.Size() / (MB),
								Used:   usage.Used() / (MB),
								Free:   usage.Free() / (MB),
								Copy:   value.Copy,
							}
						}
					}
				}
			}
		} else {
			fmt.Println(err)
		}
		//-------------------------
		SliceUsbAll := set.NewSetFromSlice(UsbAll)
		SliceUsbConnect := set.NewSetFromSlice(UsbConnect)
		Difference := SliceUsbAll.Difference(SliceUsbConnect)

		reg = []models.InfoUsb{}
		for k, v := range USB {
			if IfSerial(Difference.ToSlice(), v.Serial) {
				for i, state := range states {
					if state.Day == GetDay() {
						states = append(states[:i], states[i+1:]...)
						// -----------------------
						states = append(states, models.State{
							Day:       GetDay(),
							Connected: state.Connected + 1,
							CopiedMB:  state.CopiedMB + float64(v.Copy),
							CopiedGB:  state.CopiedGB + float64(v.Copy)/1024,
						})
						logs = append(logs, models.Logs{
							Date:        CurrentTime(),
							Model:       v.Model,
							Serial:      v.Serial,
							Size:        v.Size / (GB),
							Description: "Copiado " + strconv.FormatFloat(float64(v.Copy), 'g', 5, 64) + "MB, desconectado" ,
						})
					}
				}
				delete(USB, k)
			} else {
				reg = append(reg, v)
			}
		}
		//------------------------
		var data []byte
		_ = os.RemoveAll(Path + "/records/data.json")
		_ = os.RemoveAll(Path + "/records/state.json")
		data = JsonMarshal(reg)
		JsonWrite(data, Path+"/records/data.json")
		data = JsonMarshal(states)
		JsonWrite(data, Path+"/records/state.json")
		data = JsonMarshal(logs)
		JsonWrite(data, Path+"/records/logs.json")
		//-------------------------
		time.Sleep(5 * time.Second)
	}
}
