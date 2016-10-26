package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Node struct {
	name string
	port string
	macs []string

	debug bool
}

func (node *Node) get_vnc_port(compute Compute, ch chan Node2Port) {
	if node.debug {
		log.Println("Get vncdisplay for " + node.name)
	}

	vncdisplay := compute.Exec("virsh vncdisplay " + node.name)

	vncdisplay_s := strings.Replace(strings.Split(vncdisplay, ":")[1], "\n", "", -1)
	i, err := strconv.Atoi(vncdisplay_s)
	s := strconv.Itoa(i)

	if err != nil {
		log.Fatalf("Can't convert port to int: %s", vncdisplay)
		os.Exit(1)
	}
	if node.debug {
		log.Println("Done with " + node.name + " port is " + strconv.Itoa(i))
	}

	ch <- Node2Port{node_name: node.name, port: s}
}

func (node *Node) get_mac_addresses(compute Compute, ch chan Node2Macs) {
	r, err := regexp.Compile(`\S{2}:\S{2}:\S{2}:\S{2}:\S{2}:\S{2}`)
	if err != nil {
		log.Fatal("Can't compile regex.")
		log.Fatal(err)
		os.Exit(1)
	}

	if node.debug {
		log.Println("Get network interfaces for " + node.name)
	}

	domiflist := compute.Exec("virsh domiflist " + node.name)
	macs := r.FindAllString(domiflist, -1)
	ch <- Node2Macs{macs: macs, node_name: node.name}
}
