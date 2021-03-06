package driver

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func ReadFile(pathToFile string) string {
	file, err := os.Open(pathToFile)
	if err != nil {
		fmt.Print("FATAL: Error on reading the file!!")
		log.Fatal(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	out := string(data)
	//fmt.Print(out)
	return out
}

func WriteFile(pathToFile string, toWrite string) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(pathToFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(toWrite)); err != nil {
		f.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func OverWriteFile(pathToFile string, toWrite string) {
	err := ioutil.WriteFile(pathToFile, []byte(toWrite), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateValueInFile(pathToFile string, change int) {
	fileValue := ReadFile(pathToFile)
	oldValue, err := strconv.Atoi(fileValue)
	if err != nil {
		log.Println(errors.New("Cannot read from value File: " + pathToFile))
	}

	newValue := oldValue + change
	OverWriteFile(pathToFile, strconv.Itoa(newValue))
}

func DeleteFile(filePath string) {
	// delete file
	var err = os.Remove(filePath)
	if err != nil {
		fmt.Println("error in deleting file")
	} else {
		fmt.Println("==> done deleting file")
	}

}
