package main

import (
	"weqiniu/conf"
	"weqiniu/uploader"
)

func main() {
	c := conf.NewConfig()
	err := c.LoadJSON("config.json")
	if err != nil {
		panic(err)
	}
	u := uploader.NewUploader(c)
	u.Serve()
}
