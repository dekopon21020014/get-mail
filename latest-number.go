package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	filename string = "latest-mail.txt"
)

func GetLatestNumber() (int, error){	
	// open the file
	file, err := os.Open(filename)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	// read first line
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	if err := scanner.Err(); err != nil {		
		return -1, err
	}

	// convert to integer
	number, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return -1, err
	}

	return number, nil
}

// update latest-mail.txt
func UpdateLatestNum(number uint32) error {
	err := os.WriteFile(filename, []byte(fmt.Sprintf("%d\n", number)), 0644)
	if err != nil {
		return err
	}
	return nil
}
