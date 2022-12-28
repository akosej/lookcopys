package models

type InfoUsb struct {
	Path   string `json:"path"`
	Date   string `json:"date"`
	Serial string `json:"serial"`
	Model  string `json:"model"`
	Size   uint64 `json:"size"`
	Used   uint64 `json:"used"`
	Free   uint64 `json:"free"`
	Copy   uint64 `json:"copy"`
}

type Records struct {
	Date   string `json:"date"`
	Serial string `json:"serial"`
	Model  string `json:"model"`
	Copy   uint64 `json:"copy"`
}

type Logs struct {
	Date        string `json:"date"`
	Serial      string `json:"serial"`
	Model       string `json:"model"`
	Size        uint64 `json:"size"`
	Description string `json:"description"`
}

type State struct {
	Day       string
	Connected int
	CopiedMB  float64
	CopiedGB  float64
}
