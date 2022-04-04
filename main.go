/**
* @Author : jiahongming
* @Description :
* @Time : 2022/4/3 5:33 PM
* @File : main
* @Software: GoLand
**/
package main

import (
	"ThreeKingdoms/config"
	"fmt"
)

func main() {
	http := config.File.MustValue("login_server", "host", "127.0.0.1")
	fmt.Println(http)
}
