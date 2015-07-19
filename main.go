package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Redirect struct {
	URI          string `json:"uri"`
	Android      string `json:"android"`
	IPhone       string `json:"iphone"`
	IPod         string `json:"ipod"`
	IPad         string `json:"ipad"`
	BlackBerry   string `json:"bb"`
	WindowsPhone string `json:"wp"`
	WinXin       string `json:"wx"`
	Unknown      string `json:"unknown"`
}

type Config struct {
	Address   string      `json:"address"`
	Port      uint16      `json:"port"`
	Redirects []*Redirect `json:"redirect"`
}

var (
	config *Config = new(Config)
)

func handler(w http.ResponseWriter, r *http.Request) {
	var redirect *Redirect = nil
	for _, v := range config.Redirects {
		if v.URI == r.RequestURI {
			redirect = v
			break
		}
	}
	if redirect != nil {
		userAgent := r.UserAgent()
		matched := regexp.MustCompile("Android").MatchString(userAgent)
		if matched {
			log.Println("android", redirect.Android)
			http.Redirect(w, r, redirect.Android, 302)
			return
		}
		matched = regexp.MustCompile("iPod").MatchString(userAgent)
		if matched {
			log.Println("ipod", redirect.IPod)
			http.Redirect(w, r, redirect.IPod, 302)
			return
		}
		matched = regexp.MustCompile("iPhone").MatchString(userAgent)
		if matched {
			log.Println("iphone", redirect.IPhone)
			http.Redirect(w, r, redirect.IPhone, 302)
			return
		}
		matched = regexp.MustCompile("iPad").MatchString(userAgent)
		if matched {
			log.Println("ipad", redirect.IPad)
			http.Redirect(w, r, redirect.IPad, 302)
			return
		}
		matched = regexp.MustCompile("BlackBerry").MatchString(userAgent)
		if matched {
			log.Println("bb", redirect.BlackBerry)
			http.Redirect(w, r, redirect.BlackBerry, 302)
			return
		}
		matched = regexp.MustCompile("IEMobile").MatchString(userAgent)
		if matched {
			log.Println("wp", redirect.WindowsPhone)
			http.Redirect(w, r, redirect.WindowsPhone, 302)
			return
		}
		matched = regexp.MustCompile("MQQBrowser").MatchString(userAgent)
		if matched {
			log.Println("wx", redirect.WinXin)
			http.Redirect(w, r, redirect.WinXin, 302)
			return
		}
		matched = regexp.MustCompile("TBS").MatchString(userAgent)
		if matched {
			log.Println("wx", redirect.WinXin)
			http.Redirect(w, r, redirect.WinXin, 302)
			return
		}

		log.Println("unknown", redirect.Unknown)
		http.Redirect(w, r, redirect.Unknown, 302)
		return
	}
	log.Println("unmatched", config.Redirects[0].Unknown)
	http.Redirect(w, r, config.Redirects[0].Unknown, 302)
}

func main() {
	fh, err := os.Open("config.json")
	if err != nil {
		log.Fatal("opening config.json failed")
		return
	}
	defer fh.Close()
	configcontent, err := ioutil.ReadAll(fh)
	if err != nil {
		log.Fatal("reading config.json failed")
		return
	}

	if err = json.Unmarshal(configcontent, config); err != nil {
		log.Fatal("parsing config.json failed")
		return
	}

	http.HandleFunc("/", handler)
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf("%s:%d", config.Address, config.Port),
			nil))
}
