package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func handler(w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	fmt.Println(r.RequestURI)
	matched := regexp.MustCompile("Android").MatchString(userAgent)
	if matched {
		fmt.Println("Android")
		http.Redirect(w, r, "http://www.dfordsoft.com", 302)
		return
	}
	matched = regexp.MustCompile("iPod").MatchString(userAgent)
	if matched {
		fmt.Println("iPod")
		return
	}
	matched = regexp.MustCompile("iPhone").MatchString(userAgent)
	if matched {
		fmt.Println("iPhone")
		return
	}
	matched = regexp.MustCompile("iPad").MatchString(userAgent)
	if matched {
		fmt.Println("iPad")
		return
	}
	matched = regexp.MustCompile("BlackBerry").MatchString(userAgent)
	if matched {
		fmt.Println("BlackBerry")
		return
	}
	matched = regexp.MustCompile("IEMobile").MatchString(userAgent)
	if matched {
		fmt.Println("Windows")
		return
	}
	fmt.Println("Unknown")
}

func main() {
	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
