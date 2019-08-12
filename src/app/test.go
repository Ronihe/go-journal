package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	//3rd party lib
	"github.com/urfave/cli"
)

func init() {
	fmt.Println(" 😎  😎  😎  \n")
}

var app = cli.NewApp()

var pizza = []string{"Enjoy your pizza with some delicious"}

func info() {
	app.Name = "Simple Test cli"
	app.Usage = "an example cli to write your working journal"
	app.Author = "Roni"
	app.Version = "1.0.0"
}

var currentTime = time.Now()
var fileName = currentTime.Format("2006-01-02")
var layout = "2016-01-02 15:04:05"
var timeStampStr = currentTime.Format(layout)
var path, err = filepath.Abs("code/go-journal/src/app/journals/" + fileName + ".md")
var fileTile = currentTime.Format("2006-01-02 Mon")
var content []string

func commands() {
	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "What did you accomplish",
			Action: func(c *cli.Context) {

				shouldWrite := startJournal()
				if shouldWrite {
					fmt.Println("Hey, How was your day? Tell me everything!\n")
					fmt.Println("what did you finish today? use t 'blah blah' to let me know\n")
					writeJournal(content)
				}

			},
		},
		{
			Name:    "task",
			Aliases: []string{"t"},
			Usage:   "task and how did you finish, how do you feel\n",
			Action: func(c *cli.Context) {
				fmt.Println("what else did you do? I want to know everything!\n")
				timeStamp, err := time.Parse(layout, timeStampStr)
				if err != nil {
					return
				}
				hr, min, sec := timeStamp.Clock()
				hrStr := strconv.Itoa(hr)
				minStr := strconv.Itoa(min)
				secStr := strconv.Itoa(sec)
				achieved := hrStr + ":" + minStr + ":" + secStr + ": " + c.Args().First()
				content = append(content, achieved)
				writeJournal(content)

			},
		},
		{
			Name:    "done",
			Aliases: []string{"d"},
			Usage:   "write to the file",
			Action: func(c *cli.Context) {
				fmt.Println("Thanks for telling me everything!\n")
				readJournal()

			},
		},
	}
}

// create a file today's date

func startJournal() bool {
	//put in today's date

	//detect if file exsit

	var _, err = os.Stat(path)

	// create file if not exist

	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if isError(err) {
			return false
		}

		defer file.Close()
	} else {
		fmt.Println("😜  already created! Please go ahead and write your journal\n")
		return false
	}
	content = append(content, fileTile)

	fmt.Println(" 😎  ==>  created your  journal ", path)
	return true

}

// write the array of string to the file

func writeJournal(content []string) {
	// open your journal using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if isError(err) {
		return
	}
	defer file.Close()

	// read the content one by one and write to the file one by one
	for index, element := range content {
		if index == 0 {
			_, err := file.WriteString(element + "\n")
			if isError(err) {
				return
			}

		} else {
			idxStr := strconv.Itoa(index)
			_, err := file.WriteString(idxStr + element + "\n")
			if isError(err) {
				return
			}
		}
	}
	// save changes
	err = file.Sync()
	if isError(err) {
		return
	}

}

func readJournal() {
	// re-open file
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		return
	}
	defer file.Close()

	// read file, line by line
	var text = make([]byte, 1024)
	for {
		_, err = file.Read(text)

		// break if finally arrived at end of file
		if err == io.EOF {
			break
		}

		// break if error occured
		if err != nil && err != io.EOF {
			isError(err)
			break
		}
	}

	fmt.Println(" 😎  ==> This is your journal today\n")
	fmt.Println(string(text))
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

func main() {
	info()

	commands()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
