package main

import (
	"flag"
	"fmt"
)

type CLIArgs struct {
	host     string
	port     string
	user     string
	password string

	// Filters
	mac  bool
	node string

	// debug
	debug bool

	vncaddr bool
}

func PrintNodes(nodes Nodes) {
	for _, node := range nodes.nodes {
		fmt.Println(node.name, node.port, node.macs)
	}
}

func PrintVNCAddr(compute Compute, nodes Nodes) {
	for _, node := range nodes.nodes {
		fmt.Println(compute.host + ":" + node.port)
	}
}

func args() CLIArgs {
	host := flag.String("host", "", "SSH IP for compute")
	port := flag.String("port", "", "SSH port")
	user := flag.String("user", "", "SSH username")
	password := flag.String("password", "", "SSH password")

	node := flag.String("node", "", "Node name in virsh")
	mac := flag.Bool("mac", false, "Collect macs")

	debug := flag.Bool("debug", false, "Node mac in virsh")

	vncaddr := flag.Bool("vncaddr", false, "Show vnc camatible addrs")

	flag.Parse()

	if *host == "" {
		*host = "127.0.0.1"
	}

	if *port == "" {
		*port = "22"
	}

	if *user == "" {
		*user = "root"
	}

	if *password == "" {
		*password = ""
	}

	return CLIArgs{
		host: *host, port: *port,
		user: *user, password: *password,
		node: *node, mac: *mac,
		debug:   *debug,
		vncaddr: *vncaddr}
}

func main() {
	cliArgs := args()

	compute := Compute{
		host: cliArgs.host, port: cliArgs.port,
		username: cliArgs.user, password: cliArgs.password,
		debug: cliArgs.debug}

	nodes := Nodes{compute: compute, debug: cliArgs.debug}
	// First param part of node name to filter
	nodes.load(cliArgs.node)

	nodes.request_vnc_ports()

	if cliArgs.mac {
		nodes.request_mac_addresses()
	}

	if cliArgs.vncaddr {
		PrintVNCAddr(compute, nodes)
	} else {
		PrintNodes(nodes)
	}
}
