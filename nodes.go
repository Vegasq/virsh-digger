package main

import (
	"errors"
	"log"
	"os"
	"strings"
)

type Node2Port struct {
	node_name string
	port      string
}

type Node2Macs struct {
	node_name string
	macs      []string
}

type Nodes struct {
	nodes   []*Node
	compute Compute

	debug bool
}

func (nodes *Nodes) GetNodeByName(node_name string) (*Node, error) {
	for _, node := range nodes.nodes {
		if node.name == node_name {
			return node, nil
		}
	}
	return &Node{}, errors.New("Can't find node " + node_name)
}

func (nodes *Nodes) load(node_name_filter string) {
	if nodes.debug {
		log.Println("Loading nodes.")
	}

	session := nodes.compute.connect()
	defer session.Close()

	nodes_names := make([]string, 100)

	virsh_list := nodes.compute.Exec("virsh list")

	for _, line := range strings.Split(virsh_list, "\n") {
		if strings.Contains(line, "running") {
			items := strings.Split(line, " ")
			nodes_names = append(nodes_names, items[3])

			node := Node{name: items[3], debug: nodes.debug}

			if node_name_filter != "" && strings.Contains(node.name, node_name_filter) {
				nodes.nodes = append(nodes.nodes, &node)
			} else if node_name_filter == "" {
				nodes.nodes = append(nodes.nodes, &node)
			}
		}
	}
}

func (nodes *Nodes) request_vnc_ports() {
	n2p_ch := make(chan Node2Port)

	for _, node := range nodes.nodes {
		go node.get_vnc_port(nodes.compute, n2p_ch)
	}

	for _, node := range nodes.nodes {
		n2p := <-n2p_ch

		node_obj, err := nodes.GetNodeByName(n2p.node_name)
		if err != nil {
			log.Fatalln("Node not found: " + node.name)
			os.Exit(2)
		}

		node_obj.port = n2p.port
	}
}

func (nodes *Nodes) request_mac_addresses() {
	n2m_ch := make(chan Node2Macs)

	for _, node := range nodes.nodes {
		go node.get_mac_addresses(nodes.compute, n2m_ch)
	}

	for _, node := range nodes.nodes {
		n2m := <-n2m_ch

		node_obj, err := nodes.GetNodeByName(n2m.node_name)
		if err != nil {
			log.Fatalln("Node not found: " + node.name)
			os.Exit(2)
		}

		node_obj.macs = n2m.macs
	}
}
