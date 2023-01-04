package system

import (
	"encoding/json"
	"log"
	"os"
	"usbWatcher/models"
)

func JsonReadInfoUsb() []models.InfoUsb {
	byteValue, _ := os.ReadFile(Path + "/records/data.json")
	var result []models.InfoUsb
	_ = json.Unmarshal(byteValue, &result)
	return result
}
func JsonReadLogs() []models.Logs {
	byteValue, _ := os.ReadFile(Path + "/records/logs.json")
	var result []models.Logs
	_ = json.Unmarshal(byteValue, &result)
	return result
}
func JsonReadState() []models.State {
	byteValue, _ := os.ReadFile(Path + "/records/state.json")
	var result []models.State
	_ = json.Unmarshal(byteValue, &result)
	return result
}
func JsonReadRecords() []models.Records {
	byteValue, _ := os.ReadFile(Path + "/records/records.json")
	var result []models.Records
	_ = json.Unmarshal(byteValue, &result)
	return result
}

func JsonMarshal(reg interface{}) []byte {
	data, err := json.Marshal(reg)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func JsonWrite(data []byte, path string) {
	fp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	_, err = fp.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}
