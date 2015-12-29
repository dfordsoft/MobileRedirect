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
	WeiXin       string `json:"wx"`
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

	pairs := []struct {
		Regexp string
		To     string
		Desc   string
	}{
		{Regexp: "MQQBrowser", To: redirect.WeiXin, Desc: "wx"},
		{Regexp: "TBS", To: redirect.WeiXin, Desc: "wx"},
		{Regexp: "Android", To: redirect.Android, Desc: "android"},
		{Regexp: "iPod", To: redirect.IPod, Desc: "ipod"},
		{Regexp: "iPhone", To: redirect.IPhone, Desc: "iphone"},
		{Regexp: "iPad", To: redirect.IPad, Desc: "ipad"},
		{Regexp: "BlackBerry", To: redirect.BlackBerry, Desc: "bb"},
		{Regexp: "IEMobile", To: redirect.WindowsPhone, Desc: "wp"},
	}

	if redirect != nil {
		userAgent := r.UserAgent()
		for _, p := range pairs {
			matched := regexp.MustCompile(p.Regexp).MatchString(userAgent)
			if matched {
				log.Println(p.Desc, p.To)
				http.Redirect(w, r, p.To, 302)
				return
			}
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
