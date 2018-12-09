package main

import (
	"github.com/rdavidnota/tiered-synchronous/source/commands/utils"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var pathFolder = "J:\\bkp\\Archivos PC\\Descargas\\Video\\show cams"
var needle = "WebCam de "

func main() {
	files, err := ioutil.ReadDir(pathFolder)
	utils.Check(err)

	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), needle) {
			name := strings.Split(file.Name(), needle)
			fileTarget, err := os.OpenFile(pathFolder+"\\"+name[1], os.O_WRONLY|os.O_CREATE, 0666)
			utils.Check(err)
			defer fileTarget.Close()

			fileOrigin, err := os.Open(pathFolder+"\\"+file.Name())
			utils.Check(err)
			defer fileOrigin.Close()

			io.Copy(fileTarget, fileOrigin)
		}
	}
}
