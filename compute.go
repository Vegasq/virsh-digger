package main

import (
	"bytes"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type Compute struct {
	host     string
	port     string
	username string
	password string

	debug bool
}

func (compute *Compute) Exec(command string) string {
	session := compute.connect()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(command); err != nil {
		log.Fatal("Failed to run: " + command + " : " + err.Error())
	}

	session.Close()
	return b.String()
}

func (compute *Compute) getDial() *ssh.Client {
	config := &ssh.ClientConfig{
		User: compute.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(compute.password),
		},
	}

	client, err := ssh.Dial("tcp", compute.host+":"+compute.port, config)
	if err != nil {
		if compute.debug {
			log.Println("Failed to dial host: " + compute.host + ". Sleep for 1s and retry.")
		}
		time.Sleep(1000)
		return compute.getDial()
	}
	return client
}

func (compute *Compute) connect() *ssh.Session {
	client := compute.getDial()
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
		os.Exit(1)
	}
	return session
}
