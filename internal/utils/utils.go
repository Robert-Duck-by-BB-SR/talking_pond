package utils

import (
	"fmt"
	"os"
)

var DebugFile *os.File

func init() {
	DebugFile, _ = os.Create("debug.log")
}

func FileDebug(content any) {
	DebugFile.Write([]byte(fmt.Sprintf("%+v\n", content)))
}
