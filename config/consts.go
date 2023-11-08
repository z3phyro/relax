package config

import (
	"fmt"
	"os"
)

var rootPath string = "."

func GetRootPath() string {
	return rootPath
}

func InitConfig() {
	if len(os.Args) > 1 {
		fmt.Println(os.Args[1])
		rootPath = os.Args[1]
	}

}
