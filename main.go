package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func getQueryField(query url.Values, field string) (string, bool) {
	if q, ok := query[field]; ok && len(q) >= 1 {
		return q[0], true
	}
	return "", false
}

func hasCriteria(query url.Values, criteria string) bool {
	if c, ok := getQueryField(query, criteria); ok {
		b, _ := strconv.ParseBool(c)
		return b
	}
	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	var (
		upper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		symbols = "!$%&/()=?-_{}[]*^<>.:,;\\|@#"
		numbers = "0123456789"
		length  = 16
	)

	chars := "abcdefghijklmnopqrstuvwxyz"
	query := r.URL.Query()

	if hasCriteria(query, "upper") {
		chars += upper
	}
	if hasCriteria(query, "symbols") {
		chars += symbols
	}
	if hasCriteria(query, "numbers") {
		chars += numbers
	}

	if l, ok := getQueryField(query, "length"); ok {
		if li, err := strconv.Atoi(l); err == nil && li > 0 && li <= 256 {
			length = li
		}
	}

	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = chars[rand.Intn(len(chars))]
	}

	for i := len(buf) - 1; i > 0; i-- {
		// Fisher–Yates shuffle
		j := rand.Intn(i + 1)
		buf[i], buf[j] = buf[j], buf[i]
	}
	w.Write(buf)
}

func main() {
	http.HandleFunc("/", handler)

	errs := make(chan error, 2)
	go func() {
		log.Printf("server started on port %s\n", ":8080")
		errs <- http.ListenAndServe(":8080", nil)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Printf("service terminated: %s\n", <-errs)
}
