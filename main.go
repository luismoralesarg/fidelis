package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type KeyValueStore struct {
	store map[string]string
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{make(map[string]string)}
}

func (kvs *KeyValueStore) Set(key, value string) {
	kvs.store[key] = value
}

func (kvs *KeyValueStore) Get(key string) (string, bool) {
	value, found := kvs.store[key]
	return value, found
}

func (kvs *KeyValueStore) Delete(key string) {
	delete(kvs.store, key)
}

func main() {
	kvs := NewKeyValueStore()

	ln, err := net.Listen("tcp", ":11211")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()
	fmt.Println("Listening on :11211")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		go handleConnection(conn, kvs)
	}
}

func handleConnection(conn net.Conn, kvs *KeyValueStore) {
	defer conn.Close()
	fmt.Println("New client connected:", conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) < 2 {
			conn.Write([]byte("CLIENT_ERROR invalid command\r\n"))
			continue
		}

		cmd, key := parts[0], parts[1]
		switch cmd {
		case "get":
			value, found := kvs.Get(key)
			if found {
				response := fmt.Sprintf("VALUE %s 0 %d\r\n%s\r\n", key, len(value), value)
				conn.Write([]byte(response))
			}
			conn.Write([]byte("END\r\n"))
		case "set":
			if len(parts) != 4 {
				conn.Write([]byte("CLIENT_ERROR invalid command line format\r\n"))
				continue
			}
			kvs.Set(key, parts[3])
			conn.Write([]byte("STORED\r\n"))
		case "delete":
			kvs.Delete(key)
			conn.Write([]byte("DELETED\r\n"))
		case "quit":
			conn.Write([]byte("QUITTING\r\n"))
			return
		default:
			conn.Write([]byte("ERROR\r\n"))
		}
	}
}
