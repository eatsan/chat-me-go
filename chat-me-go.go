package main

import (
	"flag"
	"log"
	"sync"
)

type User struct {
	id      int
	name    string
	address string
	role    string
}

type UserID struct {
	sync.Mutex // ensures ID is go-routine safe
	id         int
}

func main() {

	serverFlag := flag.Bool("server", true, "-server flag enables the server functionality of chat-me-go.")
	portFlag := flag.Int("port", 8090, "-port flag defines the server's listening TCP port")

	flag.Parse()

	if *serverFlag {

		s, err := NewServer("127.0.0.1", *portFlag)
		if err != nil {
			log.Fatalf("Failed to create a new server. err: %s", err)
		}

		defer (*s.listener).Close()

		for {
			conn, err := (*s.listener).Accept()
			if err != nil {
				log.Fatalf("Failed to accept connection request. err:%s", err)
			}

			go s.handleConnection(conn)

		}
	} else {
		//Client code
	}

}

// setUsername sets the username field in user type with the given name
// TODO: check existing active usernames and return error if given name cannot be selected.
func (u *User) setUsername(name string) error {
	log.Printf("Setting username of user %s to %s (userId: %d)", u.name, name, u.id)
	u.name = name
	return nil
}

func (i *UserID) ID() (id int) {
	i.Lock()
	defer i.Unlock()

	id = i.id
	i.id++
	return
}
