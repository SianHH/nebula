package main

import (
	"fmt"
	"github.com/slackhq/nebula/package/common"
)

func main() {
	fmt.Println(common.GenerateCa("example-ca", "./example-ca.key", "./example-ca.crt"))
	fmt.Println(common.GenerateSign(
		"example-a",
		"10.11.11.100/24",
		"normal",
		"./example-ca.key",
		"./example-ca.crt",
		"./example-a.crt",
		"./example-a.key",
	))
}
