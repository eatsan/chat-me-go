package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	address         string
	connectionCount int
	activeUsers     map[int]*User
	listener        *net.Listener
	ids             UserID
}

//NewServer returns a pointer to a new Server
func NewServer(hostname string, port int) (*Server, error) {

	s := Server{
		address:         fmt.Sprintf("%s:%d", hostname, port),
		connectionCount: 0,
	}

	l, err := net.Listen("tcp4", s.address)
	if err != nil {
		return nil, err
	}
	s.listener = &l
	s.activeUsers = make(map[int]*User)
	return &s, nil
}

// handleConnection handles required actions from chat client messages
func (s *Server) handleConnection(c net.Conn) {

	defer c.Close()

	s.increaseConnectionCount()
	defer s.decreaseConnectionCount()
	log.Printf("Received a new connection to the server from %s (Total active connections: %d)", c.RemoteAddr().String(), s.connectionCount)

	r := bufio.NewReader(c)

	// userId generation is go-routines safe so should be unique on a server
	currentUserID := s.ids.ID()

	s.activeUsers[currentUserID] = &User{
		id:      currentUserID,
		name:    fmt.Sprintf("user%d", currentUserID),
		address: c.RemoteAddr().String(),
		role:    s.setDefaultUserRole(),
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
			log.Printf("[%s] %s", s.activeUsers[currentUserID].name, line)
		case "\\exit":
			log.Printf("Closing chat connection to the client %s", c.RemoteAddr().String())
			return
		case "\\name":
			(s.activeUsers[currentUserID]).setUsername(commandArg)
		default:
			log.Printf("Unknown command. failing command: %s\n", command)
		}

	}
}

// setDefaultUserRole returns the default user role configured in the server
// This function can be useful in the future when users might have different roles in the chat room
func (s *Server) setDefaultUserRole() string {
	return "user"
}

func (s *Server) increaseConnectionCount() {
	s.connectionCount++
}
func (s *Server) decreaseConnectionCount() {
	s.connectionCount--
}
