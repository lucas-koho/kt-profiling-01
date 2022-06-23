package main

import (
	"fmt"
	"github.com/blackfireio/go-blackfire"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	blackfire.Disable()
	// memory
	http.HandleFunc("/strings", func(w http.ResponseWriter, r *http.Request) {
		makeStrings()
	})

	// callstack
	http.HandleFunc("/fib", func(w http.ResponseWriter, r *http.Request) {
		fibonacci()
	})

	// io http
	http.HandleFunc("/requests", func(w http.ResponseWriter, r *http.Request) {
		requests()
	})

	// io disk
	http.HandleFunc("/diskio", func(w http.ResponseWriter, r *http.Request) {
		diskio()
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	//ender.End()
}

func leakyFunction(wg sync.WaitGroup) {
	defer wg.Done()
	s := make([]string, 3)
	for i := 0; i < 1000000; i++ {
		s = append(s, "bananas")
		if (i % 100000) == 0 {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func makeStrings() {

	list := []string{}

	for i := 0; i < 100000; i++ {
		list = append(list, strings.Repeat("a", i))
	}
}

func fibonacci() {
	fib(30)
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func requests() {
	http.Get("http://www.koho.ca")
}

func diskio() {
	for i := 0; i < 100; i++ {
		f, err := os.OpenFile("./diskio", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Println(err)
		}
		f.WriteString("Here is the data I am writing")
		f.Close()
	}
}

func diskiobatch() {

	s := []string{}

	for i := 0; i < 1000; i++ {
		s = append(s, "Here is the data I am writing")
	}

	f, err := os.OpenFile("./diskiobatch", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err)
	}
	f.WriteString(strings.Join(s, "\\n"))
	f.Close()
}
