package logs

import (
	"log"
	"os"
	"time"

	"crypto/rand"
	"encoding/base64"
)

const LogTimeFormat = "2006-01-02"

func GetLogFile() (*os.File, error) {
	dir := "./logs"
	b, mkdirErr := DirExists(dir)
	if !b && mkdirErr == nil {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Println("logs mkdir err~")
			return nil, err
		} else {
			log.Println("logs mkdir success")
		}
	}
	//randomString, _ := generateRandomString(5)
	file := "./logs/" + time.Now().Format(LogTimeFormat) + "_info" + ".log"
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		defer f.Close()
		return nil, err
	}
	return f, err
}

func DirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func generateRandomString(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	randomString := base64.RawURLEncoding.EncodeToString(buffer)
	return randomString[:length], nil
}
