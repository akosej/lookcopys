package models

import (
	"github.com/fsnotify/fsnotify"
	"time"
)

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

type Operation struct {
	Size    float64
	ModTime time.Time
	Event   fsnotify.Op
}

type Records struct {
	Path    string `json:"path"`
	Date    string `json:"date"`
	Serial  string `json:"serial"`
	Model   string `json:"model"`
	Size    float64
	ModTime time.Time
	Event   fsnotify.Op
}

// -- Ouput
type InfoUsbOut struct {
	Id      int       `json:"id"`
	Path    string    `json:"path"`
	Date    string    `json:"date"`
	Serial  string    `json:"serial"`
	Model   string    `json:"model"`
	Size    uint64    `json:"size"`
	Used    uint64    `json:"used"`
	Free    uint64    `json:"free"`
	Copy    uint64    `json:"copy"`
	Records []Records `json:"records"`
}
