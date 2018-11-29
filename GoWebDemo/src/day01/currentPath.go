package main


import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//func main() {
//	fmt.Println("GoLang 获取程序运行绝对路径")
//	fmt.Println(GetCurrPath())
//}

func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return ret
}
