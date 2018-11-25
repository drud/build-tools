package main

import (
	"fmt"
"github.com/lextoumbourou/goodhosts"

)

func main() {
	fmt.Println("This is just a silly main.")
	_, err := goodhosts.NewHosts()
	if err != nil {
		return
	}

}
