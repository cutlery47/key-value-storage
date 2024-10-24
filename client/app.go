package client

import "github.com/cutlery47/key-value-storage/client/internal/client"

func Run() {
	cl := client.NewHTTP()
	app := client.New(cl)
	app.Run()
}
