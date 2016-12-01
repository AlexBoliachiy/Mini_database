// Lab2_file_database project main.go
package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
)

type cache struct {
	tables []table
}

func (c *cache) get(name string) *table {
	//Проверяем есть ли таблицца на оперативке
	fmt.Println("Try to get table")
	for _, t := range c.tables {
		if t.name == name {
			fmt.Println("table find in cache")
			return &t
		}
	}
	//Пытаемся загрузить базу из файлововой системы
	f, err := os.Open(name)
	if err == nil {
		data := make([][]string, 0, 0)
		r := csv.NewReader(bufio.NewReader(f))

		for {
			record, err := r.Read()
			// Stop at EOF.
			if err == io.EOF {
				break
			}

			row_data := []string{record[0], record[1]}
			data = append(data, row_data)
		}
		t := table{name: name, data: &data, isLocked: false}
		c.tables = append(c.tables, t)
		return &t
	}
	//Ну и если ни черта не вышло создаем пустую базу
	t := table{name: name, data: new([][]string), isLocked: false}
	c.tables = append(c.tables, t)
	fmt.Println("table returned")
	return &t
}

type table struct {
	isLocked bool
	name     string
	data     *[][]string
}

func (t *table) save() {
	for t.isLocked {

	}

	t.isLocked = true
	file, _ := os.Create(t.name)
	defer file.Close()
	for _, row := range *t.data {
		file.WriteString(row[0] + "," + row[1] + "\n")
	}
	t.isLocked = false
}

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8888"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.
	c := cache{}
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
		go handleConnection(conn, &c)
	}
}

// Handles incoming requests.
func handleConnection(conn net.Conn, c *cache) {
	// Make a buffer to hold incoming data.
	var buf [512]byte
	empty_string := ""
	for {
		copy(buf[:], empty_string) //make  buffer empty
		_, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		n := bytes.Index(buf[:], []byte{0})
		s := string(buf[:n])
		go handleRequest(s, c, conn)
	}
}

func handleRequest(command string, c *cache, conn net.Conn) {
	fmt.Println("handle request: " + command)

	params := strings.Fields(command)
	// remove empty entries and remove whitespaces
	fmt.Println(params)
	fmt.Println(params[0])
	// remove empty entries and remove whitespaces
	switch strings.ToLower(params[0]) {
	case "exit":
		os.Exit(0)
	case "insert":
		fmt.Println("switching")

		fmt.Println("insert")
		insert_regex := regexp.MustCompile(`^insert\s+\w+\s+\w+\s+into\s+[A-z]+[A-z_0-9]*`)
		if insert_regex.MatchString(command) {
			t := c.get(params[4])
			t.insert_(params[1], params[2])

		} else {
			conn.Write([]byte(string("Your command didn't match the pattern\n")))
		}
	case "select":
		fmt.Println("try to select")
		select_regex := regexp.MustCompile(`^select\s+\w+\s+from\s+[A-z]+[A-z_0-9]*`)
		if select_regex.MatchString(command) {

			t := c.get(params[3])
			data := ""
			for _, elem := range t.select_(params[1]) {
				data += elem + " "
				conn.Write([]byte(string(elem) + " "))

			}
			conn.Write([]byte("\n"))

			fmt.Println(data + "\n newline was sended")
		} else {
			conn.Write([]byte(string("Your command didn't match the pattern\n")))
		}
		fmt.Println("select")
	case "delete":
		delete_regex := regexp.MustCompile(`^delete\s+\w+\s+in\s+[A-z]+[A-z_0-9]*`)
		if delete_regex.MatchString(command) {
			t := c.get(params[3])
			t.delete_(params[1])
		} else {
			conn.Write([]byte(string("Your command didn't match the pattern\n")))
		}
		fmt.Println("delete")

	case "update":
		update_regex := regexp.MustCompile(`^update\s+\w+\s+to\s+\w+\s+in\s+[A-z]+[A-z_0-9]*`)
		if update_regex.MatchString(command) {
			t := c.get(params[5])
			t.update_(params[1], params[3])
		} else {
			conn.Write([]byte(string("Your command didn't match the pattern\n")))
		}
		fmt.Println("update")
	default:
		conn.Write([]byte(string("Your command didn't match the pattern\n")))
	}
}

func (t *table) insert_(key string, value string) {
	row := []string{key, value}
	//Возможное место для мютекса
	fmt.Println(t.data)
	t.additem(row)
	fmt.Println(t.data)
	go t.save()
}
func (t *table) additem(value []string) {
	*t.data = append(*t.data, value)
}
func (t *table) select_(key string) []string {
	selecting_data := []string{}
	fmt.Println(*t.data)
	for _, row := range *t.data {
		if row[0] == key {
			selecting_data = append(selecting_data, row[1])
		}
	}

	return selecting_data
}

func (t *table) delete_(key string) {
	for index := 0; index < len(*t.data); index++ {
		if (*t.data)[index][0] == key {
			*t.data = append((*t.data)[:index], (*t.data)[index+1:]...)
		}
	}
	go t.save()
}

func (t *table) update_(key string, value string) {
	for i, row := range *t.data {
		if row[0] == key {
			(*t.data)[i][1] = value
		}

	}
	go t.save()
}
