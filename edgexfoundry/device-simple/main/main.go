package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	path := "README.md"

	Read(path)
}

func Read(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
