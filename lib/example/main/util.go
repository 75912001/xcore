package main

import (
	"fmt"
	xutil "xcore/lib/util"
)

func exampleUtil() {
	if false {
		outStr, errStr, err := xutil.Command("ls -l")
		fmt.Println(outStr, errStr, err)
	}
	return
}
