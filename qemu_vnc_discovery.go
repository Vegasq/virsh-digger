package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type VNC struct {
	vmname  string
	vncport int
}

type Creds struct {
	username *string
	password *string
	host     *string
	port     *string
}

type QEMUVNCDiscovery struct {
	args      Creds
	vnc_nodes []VNC
}

func (vncd QEMUVNCDiscovery) get_dial() *ssh.Client {
	config := &ssh.ClientConfig{
		User: *vncd.args.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(*vncd.args.password),
		},
	}

	client, err := ssh.Dial("tcp", *vncd.args.host+":"+*vncd.args.port, config)
	if err != nil {
		log.Println("Failed to dial host: " + *vncd.args.host + ". Sleep for 10 and retry.")
		time.Sleep(100)
		return vncd.get_dial()
	}
	return client
}

func (vncd QEMUVNCDiscovery) connect() *ssh.Session {
	client := vncd.get_dial()
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
		os.Exit(1)
	}
	return session
}

func (vncd QEMUVNCDiscovery) run(command string) string {
	session := vncd.connect()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		log.Fatal("Failed to run: " + command + " : " + err.Error())
	}

	session.Close()
	return b.String()
}

func (vncd QEMUVNCDiscovery) get_nodes(session *ssh.Session) []string {
	nodes_names := make([]string, 100)
	virsh_list := vncd.run("virsh list")
	for _, line := range strings.Split(virsh_list, "\n") {
		if strings.Contains(line, "running") {
			items := strings.Split(line, " ")
			nodes_names = append(nodes_names, items[3])
		}
	}
	return nodes_names
}

func (vncd QEMUVNCDiscovery) get_vnc_port(session *ssh.Session, ch chan VNC, node_name string) {
	log.Println("Get vncdisplay for " + node_name)
	vncdisplay := vncd.run("virsh vncdisplay " + node_name)

	vncdisplay_s := strings.Replace(strings.Split(vncdisplay, ":")[1], "\n", "", -1)
	i, err := strconv.Atoi(vncdisplay_s)
	if err != nil {
		log.Fatalf("Can't convert port to int: %s", vncdisplay)
		os.Exit(1)
	}

	log.Println("Done with " + node_name + " port is " + strconv.Itoa(i))
	ch <- VNC{vmname: node_name, vncport: i}
}

func (vncd QEMUVNCDiscovery) start_vnc_viewer() {
	for _, vnc := range vncd.vnc_nodes {
		log.Println(vnc)
		vnc_addr := *vncd.args.host + ":" + strconv.Itoa(vnc.vncport)
		log.Println(vnc_addr)
		cmd := exec.Command("gvncviewer", vnc_addr)
		err := cmd.Start()

		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Waiting for command to finish...")
		err = cmd.Wait()
		log.Printf("Command finished with error: %v", err)
	}

}

func (vncd QEMUVNCDiscovery) init() {
	session := vncd.connect()

	defer session.Close()

	nodes := vncd.get_nodes(session)
	ch := make(chan VNC)

	for _, node_name := range nodes {
		if len(node_name) > 0 {
			go vncd.get_vnc_port(session, ch, node_name)
		}
	}
	for _, node_name := range nodes {
		if len(node_name) > 0 {
			log.Println("Await for response")
			vncd.vnc_nodes = append(vncd.vnc_nodes, <-ch)
			log.Println("Got response")
		}
	}

	vncd.start_vnc_viewer()
}

func args() Creds {
	host := flag.String("host", "", "SSH IP for compute")
	port := flag.String("port", "", "SSH port")
	user := flag.String("user", "", "SSH username")
	password := flag.String("password", "", "SSH password")
	flag.Parse()
	if (*host == "") || (*port == "") || (*user == "") || (*password == "") {
		log.Fatalln("define -host -port -user -password.")
		os.Exit(1)
	}
	return Creds{host: host, username: user, password: password, port: port}
}

func main() {
	vncd := QEMUVNCDiscovery{args: args()}
	vncd.init()
}
