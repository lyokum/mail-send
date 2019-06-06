package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
	"os/exec"
)

var (
	ErrNoRecipient = errors.New("No recipient specified")
	ErrTooManyArgs = errors.New("Too many arguments given")

	Providers = map[string]string{"att": "txt.att.net", "tmobile": "tmomail.net", "sprint": "messaging.sprintpcs.com", "verizon": "vtext.com"}
)

func sendMail(to, subject, body string) error {
	// fill in mail message to send
	message := fmt.Sprintf("Subject: %s\n%s", subject, body)

	// create command
	cmd := exec.Command("sendmail", to)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, message)
	}()

	err = cmd.Run()

	return err
}

// main function of app
func appHandler(c *cli.Context) (err error) {
	// check n args
	if c.NArg() == 0 {
		cli.ShowAppHelpAndExit(c, 1)
	}

	if c.NArg() > 2 {
		return ErrTooManyArgs
	}

	// set subject and body
	var subject, body string
	if c.NArg() == 1 {
		subject = "GENERAL"
		body = c.Args()[0]
	} else {
		subject = c.Args()[0]
		body = c.Args()[1]
	}

	// choose recipient
	to := c.String("recipient")
	if to == "" {
		return ErrNoRecipient
	}

	// send mail to recipient
	err = sendMail(to, subject, body)
	if err != nil {
		return
	}

	return nil
}

func main() {
	// Create app
	app := cli.NewApp()

	// Fill fields
	app.Name = "mail-send"
	app.Usage = "mail-send BODY, mail-send SUBJECT BODY"
	app.Description = "wrapper for sendmail command"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Logan Yokum",
			Email: "lyokum@nd.edu",
		},
	}
	app.Action = appHandler

	// Fill flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "recipient, r",
			Usage: "specify `RECIPIENT` email address of message (required)",
			Value: "",
		},
	}

	// Remove commands
	app.Setup()
	app.Commands = []cli.Command{}

	// Run app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
