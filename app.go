package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var proxy, username string

func init() {
	tsh := readConf()

	ext := func() {
		log.Println("Error for tsh config")
		os.Exit(1)
	}

	t := strings.Split(tsh, "#")
	if len(t) != 2 {
		ext()
	}

	x := strings.Split(t[0], "=")
	if x[0] != "proxy" {
		ext()
	}

	y := strings.Split(t[1], "=")
	if y[0] != "username" {
		ext()
	}
	proxy, username = x[1], y[1]
}

func main() {
	// Get ARGS
	var tshFlag, cmd string
	flag.StringVar(&tshFlag, "s", "", "flag for service name")
	flag.StringVar(&cmd, "c", "", "flag for command")
	flag.Parse()

	// Validate
	if tshFlag == "" {
		log.Println("Please define -s=[SERVICE NAME]")
		return
	} else if cmd == "" {
		log.Println("Please define -c=[COMMAND]")
		return
	}

	servSlice := lsTsh(tshFlag)
	var wg sync.WaitGroup
	for i, v := range servSlice {
		wg.Add(1)
		go func(i int, v string, wg *sync.WaitGroup) {
			defer wg.Done()
			command := fmt.Sprintf("tsh --proxy=%s --user=%s ssh root@%s %s", proxy, username, v, cmd)
			x, err := exec.Command("/bin/sh", "-c", command).Output()
			if err != nil {
				log.Println("Error Command for /bin/sh", string(x), err)
			}
			log.Println("===============", v, "===============")
			log.Println(string(x))
		}(i, v, &wg)
	}
	wg.Wait()
	log.Println("DONE!")
}

func readConf() string {
	f, _ := os.Open("tsh_conf.txt")
	fileScanner := bufio.NewScanner(f)

	tsh_conf := ""
	for fileScanner.Scan() {
		tsh_conf += fileScanner.Text()
	}

	return tsh_conf
}

func lsTsh(servName string) []string {
	cmd := "tsh --proxy=%s --user=%s ls | grep %s | awk '{print $1}'"

	x, err := exec.Command("/bin/sh", "-c", fmt.Sprintf(cmd, proxy, username, servName)).Output()
	if err != nil {
		log.Println("[lsTsh] Error Command for /bin/sh", string(x), err)
	}
	s := strings.TrimSpace(string(x))
	return strings.Split(s, "\n")
}
