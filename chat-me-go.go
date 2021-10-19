package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

type user struct {
	name    string
	address string
	role    string
}

var connectionCount int

func main() {

	serverFlag := flag.Bool("server", true, "-server flag enables the server functionality of chat-me-go.")
	portFlag := flag.Int("port", 8090, "-port flag defines the server's listening TCP port")

	flag.Parse()

	if *serverFlag {
		//Server code
		connectionCount = 0
		address := fmt.Sprintf("127.0.0.1:%d", *portFlag)
		l, err := net.Listen("tcp4", address)
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()

		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go handleConnection(conn)

		}
	} else {
		//Client code
	}

}

// handleConnection handles required actions from chat client messages
func handleConnection(c net.Conn) {

	defer c.Close()
	//TODO: do something with the connection object.
	fmt.Printf("Received a new connection to the server from %s\n", c.RemoteAddr().String())

	r := bufio.NewReader(c)

	increaseConnectionCount()

	user := user{
		name:    generateUsername(),
		address: c.RemoteAddr().String(),
		role:    setDefaultUserRole(),
	}

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			log.Printf("Error reading line from the connection. Err: %s\n", err)
			return
		}

		line = strings.Trim(line, " \n")

		var command string
		//TODO: command arguments should be a slice of strings
		var commandArg string

		// Check if line starts with a chat room command (e.g. starts with \ rune)
		if []rune(line)[0] == '\\' {
			ss := strings.Split(line, " ")
			command = ss[0]
			if len(ss) > 1 {
				commandArg = ss[1]
			}
		}

		switch command {
		case "":
			log.Printf("[%s] %s", user.name, line)
		case "\\exit":
			log.Printf("Closing chat connection to the client %s\n", c.RemoteAddr().String())
			return
		case "\\name":
			user.setUsername(commandArg)
		default:
			log.Printf("Unknown command. failing command: %s\n", command)
		}

	}
}

func increaseConnectionCount() {
	connectionCount++
}

// generateUsername returns a generated username in the /userSequenceID/ format.
func generateUsername() string {
	return fmt.Sprintf("user%d", connectionCount)
}

// setUsername sets the username field in user type with the given name
// TODO: check existing active usernames and return error if given name cannot be selected.
func (u *user) setUsername(name string) error {
	log.Printf("Setting username of user %s to %s\n", u.name, name)
	u.name = name
	return nil
}

// setDefaultUserRole returns the default user role configured in the server
// This function can be useful in the future when users might have different roles in the chat room
func setDefaultUserRole() string {
	return "user"
}
