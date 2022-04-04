/**
* @Author : jiahongming
* @Description :
* @Time : 2022/4/3 5:00 PM
* @File : config.go
* @Software: GoLand
**/
package config

import (
	"errors"
	"fmt"
	"github.com/Unknwon/goconfig"
	"log"
	"os"
)

const configFile = "/conf/conf.ini"

var File *goconfig.ConfigFile

//init
func init() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + configFile

	//参数
	len := len(os.Args)
	if len > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + configFile

		}
	}
	if !fileExist(configPath) {
		panic(errors.New("配置文件不存在"))
	}

	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		log.Fatal("读取配置文件出错:", err)
	}
	//fmt.Print(File)
}

func fileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}
func A() {
	fmt.Println("AAA")
}
