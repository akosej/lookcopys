package system

import (
	"fmt"
	"github.com/akosej/lookcopys/models"
	set "github.com/deckarep/golang-set"
	detector "github.com/deepakjois/gousbdrivedetector"
	"github.com/fsnotify/fsnotify"
	"github.com/ricochet2200/go-disk-usage/du"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func MonitorUsb() {
	//----------------------------------------------------------------------------------
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				//log.Println("event:", event)
				if event.Op.String() == "[no events]" {
					_ = watcher.Remove(event.Name)
				}
				//------------------------------
				if event.Has(fsnotify.Rename) {
					f, _ := os.Open(event.Name)
					files, _ := f.Readdir(0)
					for _, v := range files {
						if v.IsDir() && v.Name() == event.Name {
							_ = watcher.Remove(v.Name())
						}
					}
				}
				if event.Has(fsnotify.Chmod) {
					f, _ := os.Stat(event.Name)
					if f.Size() > 0 {
						for k, v := range ACTION {
							if k == event.Name {
								ACTION[event.Name] = models.Records{
									Path:    event.Name,
									Date:    CurrentTime(),
									Serial:  v.Serial,
									Model:   v.Model,
									Size:    Round(float64(f.Size()), 2),
									ModTime: f.ModTime(),
									Event:   fsnotify.Chmod,
								}
							}
						}
					}
				}
				//------------------------------
				if event.Has(fsnotify.Create) {
					f, _ := os.Open(event.Name)
					files, _ := f.Readdir(0)
					if len(files) > 0 {
						_ = filepath.Walk(event.Name, func(path string, info os.FileInfo, err error) error {
							if err != nil {
								return err
							}
							f, _ := os.Open(path)
							files, _ := f.Readdir(0)
							for _, v := range files {
								if v.IsDir() {
									_ = watcher.Add(path + "/" + v.Name())
								}
							}
							return nil
						})
					} else {
						f, _ := os.Stat(event.Name)
						if f.IsDir() {
							_ = watcher.Add(event.Name)
						} else {

							for _, v := range USB {
								if strings.Contains(event.Name, v.Path) {
									ACTION[event.Name] = models.Records{
										Path:    event.Name,
										Date:    CurrentTime(),
										Serial:  v.Serial,
										Model:   v.Model,
										Size:    Round(float64(f.Size()), 2),
										ModTime: f.ModTime(),
										Event:   fsnotify.Create,
									}
									break
								}
							}

						}
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	//----------------------------------------------------------------------------------
	reg := JsonReadInfoUsb()
	logs := JsonReadLogs()
	states := JsonReadState()
	records := JsonReadRecords()
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
						SendNotifyDesktop("Dispositivo conectado", model+" "+strconv.Itoa(int(usage.Size()/(GB)))+"GB")
						_ = watcher.Add(path)
						_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
							if err != nil {
								return err
							}
							f, _ := os.Open(path)
							files, _ := f.Readdir(0)
							for _, v := range files {
								if v.IsDir() {
									_ = watcher.Add(path + "/" + v.Name())
								}
							}
							return nil
						})
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
							Description: "Copiado " + strconv.FormatFloat(float64(v.Copy), 'g', 5, 64) + "MB, desconectado",
						})
						SendNotifyDesktop(
							"Dispositivo desconectado",
							v.Model+" "+strconv.Itoa(int(v.Size/(GB)))+"GB \n Copiado "+strconv.FormatFloat(float64(v.Copy), 'g', 5, 64)+"MB\n cobrar $"+strconv.FormatFloat(Round((float64(v.Copy)/1024)*5, 1), 'g', 5, 64))
						_ = watcher.Remove(v.Path)
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
		//fmt.Println(watcher.WatchList())
		for p, v := range ACTION {
			if v.Size > 0 {
				log.Println("COPY:", p, v.Serial)
				records = append(records, models.Records{
					Path:    v.Path,
					Date:    CurrentTime(),
					Serial:  v.Serial,
					Model:   v.Model,
					Size:    Round(v.Size, 2),
					ModTime: v.ModTime,
				})
				delete(ACTION, p)
			}
			//time.Sleep(5 * time.Second)
		}
		data = JsonMarshal(records)
		JsonWrite(data, Path+"/records/records.json")
	}
}
