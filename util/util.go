package util

import (
	"io"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

func GetEnvFilePath() string {
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "../.env")
}

func PanicOnError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func HandleClose(closer io.Closer) {
	if closer != nil {
		err := closer.Close()
		if err != nil {
			log.Printf("error on close %T: %s\n", closer, err.Error())
		}
	}
}
