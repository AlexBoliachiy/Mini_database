// Lab2_file_database project main.go
package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "23"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		fmt.Println("accept")

		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleConnection(conn)
	}
}

// Handles incoming requests.
func handleConnection(conn net.Conn) {
	// Make a buffer to hold incoming data.
	var buf [512]byte
	for {
		_, err := conn.Read(buf[0:])
		if err != nil {
			return

		}
	}
}

func handleRequest(command string) {

	params := strings.Split(command, " ")
	// remove empty entries and remove whitespaces
	cnt := 0
	for _, element := range params {
		element = strings.Replace(element, " ", "", -1)
		if element == "" {
			cnt++
		}
	}
	// Если длина массива не равна кнт + 1 значит в массиве есть пустые элементы, следовательно убираем их
	if len(params) != cnt+1 {
		replace_arr := make([]string, cnt+1)
		cnt = 0
		for _, element := range params {
			if element == "" {
				replace_arr[cnt] = element
				cnt++
			}
		}
		params = replace_arr
	}
	// remove empty entries and remove whitespaces
	switch params[0] {
	case "insert":
		fmt.Println("insert")
	case "select":
		fmt.Println("select")
	case "delete":
		fmt.Println("delete")
		
	case "update":
		fmt.Println("update")
	}
}

func insert_(key string, value string, table_name string) {

}

func select_(key string, table_name string) []string {

}

func delete_(key string, table_name string) {

}

func update_(key string, table_name string) {

}
