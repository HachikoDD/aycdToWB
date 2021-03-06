package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"moul.io/banner"
)

type session struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Id      string `json:"_id"`
	Proxy   string `json:"proxy"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Checkboxes(label string, opts []string) string {
	res := ""
	prompt := &survey.Select{
		Message: label,
		Options: opts,
	}
	survey.AskOne(prompt, &res)

	return res
}

func StringPrompt(label string) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	var options []string
	for _, file := range files {
		if strings.Contains(file.Name(), ".csv") {
			options = append(options, file.Name())
		}
	}
	options = append(options, "exit")

	answer := Checkboxes(
		"Which file?", options)

	// exit program
	if answer == "exit" {
		os.Exit(0)
	}
	records := readCsvFile(answer)
	var sessions [][]byte
	for _, record := range records[1:] {
		s := &session{
			Name:    record[1],
			Enabled: true,
			Id:      RandStringRunes(16),
			Proxy:   record[5],
		}
		sjson, _ := json.Marshal(s)
		sessions = append(sessions, sjson)
	}
	f, err := os.Create("captchaSessions.db")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, s := range sessions {
		_, err2 := f.WriteString(string(s) + "\n")
		if err2 != nil {
			log.Fatal(err2)
		}
	}
	fmt.Printf("captchaSessions.db generated by %d gmail!\n", len(sessions))

}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func main() {
	fmt.Println(banner.Inline("DingD"))
	for {
		StringPrompt("What is your name?")
	}
}
