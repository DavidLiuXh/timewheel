package main

import (
	"bufio"
	"fmt"
	"os"
	"timewheel"
)

var runTime = 1
var timeWheel *timewheel.TimeWheel

func showMsg(datas []interface{}) {
	fmt.Printf("runTime:%d\n", runTime)
	runTime++
	for _, v := range datas {
		fmt.Printf("data %v | %T\n", v, v)

		dd := v.([]interface{})
		for _, d := range dd {
			fmt.Printf("data %v | %T\n", d, d)
			if d == "c" {
				fmt.Printf("stop%v\n", d)
				timeWheel.Stop()
			}
		}
	}
}

func main() {
	timeWheel = timewheel.NewTimeWhell(1, 30, showMsg)
	timeWheel.Start()
	timeWheel.Add(3, "a")
	timeWheel.Add(9, "b")
	timeWheel.Add(15, "c")
	bio := bufio.NewReader(os.Stdin)
	bio.ReadLine()
}
