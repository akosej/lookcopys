package main

import (

	//"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"time"
	"usbWatcher/models"
	"usbWatcher/routes"
	"usbWatcher/system"
)
type Client struct {
	name   string
	events chan *models.Logs
}
func main() {
	system.CreateDirectoryIfDoesntExist(system.Path+"/records")
	go system.MonitorUsb()
	// ------------------------------------------------------------------------------------------------------------
	// ------- RUN API-AGA-SENTINEL
	// ------------------------------------------------------------------------------------------------------------
	app := fiber.New(fiber.Config{
		AppName:       "OWL-MonitorUsb v1.0",
		CaseSensitive: true,
		//EnablePrintRoutes: true,
		//GETOnly: true,
	})

	app.Use(cors.New(cors.Config{
		AllowCredentials: false,
	}))

	app.Static("/", system.Path+"/frontend", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	//app.Get("/sse", adaptor.HTTPHandler(handler(dashboardHandler)))

	routes.RoutesApi(app)
	_ = app.Listen(":5323")

}
//
//func handler(f http.HandlerFunc) http.Handler {
//	return http.HandlerFunc(f)
//}
//func dashboardHandler(w http.ResponseWriter, r *http.Request) {
//	client := &Client{name: r.RemoteAddr, events: make(chan *DashBoard, 10)}
//	go updateDashboard(client)
//
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
//	w.Header().Set("Content-Type", "text/event-stream")
//	w.Header().Set("Cache-Control", "no-cache")
//	w.Header().Set("Connection", "keep-alive")
//
//	timeout := time.After(1 * time.Second)
//	select {
//	case ev := <-client.events:
//		var buf bytes.Buffer
//		enc := json.NewEncoder(&buf)
//		enc.Encode(ev)
//		fmt.Fprintf(w, "data: %v\n\n", buf.String())
//		fmt.Printf("data: %v\n", buf.String())
//	case <-timeout:
//		fmt.Fprintf(w, ": nothing to sent\n\n")
//	}
//
//	if f, ok := w.(http.Flusher); ok {
//		f.Flush()
//	}
//}
//
//func updateDashboard(client *Client) {
//	for {
//		db := &DashBoard{
//			User: uint(rand.Uint32()),
//		}
//		client.events <- db
//	}
//}