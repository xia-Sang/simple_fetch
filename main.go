package main

import (
	"fmt"
	"neofetch/helper"
)

func main() {
	info, err := helper.NewPCInfo()
	fmt.Println(info, err)
}
