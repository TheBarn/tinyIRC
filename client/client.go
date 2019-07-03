package main

import (
	"fmt"

	"github.com/TheBarn/tinyIRC/utils"
)

func main() {
	port := utils.ChoosePort()
	fmt.Println("Hello", port)
}
