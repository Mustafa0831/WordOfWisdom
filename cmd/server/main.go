package main

import (
	"context"
	"flag"
)

var add = flag.String("addr", "0.0.0.0:1024", "")

func main(){
	flag.Parse()
	// ctx,cancel:=context.WithCancel(context.Background())

}