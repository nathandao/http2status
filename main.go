package main

import (
	"fmt"

	. "github.com/nathandao/http2status/http2status"
)

var url = "golang.org"

func main() {
	s, res, err := Http2Status(url)
	if err != nil {
		fmt.Println("Opps, error: ", err)
	}

	if !s {
		fmt.Println("Not http2")
	} else {
		fmt.Println("HTTP2 detected!")
		fmt.Println(res)
	}
}
