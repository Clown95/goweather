package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Rss struct {
	Channel []Channel `xml:"channel,version"`
}
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	PubDate     string `xml:"pubDate"`
	Category    string `xml:"category"`
	Item        []Item `xml:"item"`
}
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid,isPermaLink"`
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "天气预报, version 0.1\n")
		fmt.Fprintf(os.Stderr, "输入城市")
		flag.PrintDefaults()
	}
	cityname := os.Args[1]
	cityid := getcityid(cityname)
	weather := getweather(cityid)
	fmt.Println(weather)
}

func getcityid(cityname string) (cityid string) {
	data := url.Values{}
	data.Set("search", cityname)
	loc := []string{}
	redirect := func(req *http.Request, via []*http.Request) error {
		loc = append(loc, req.URL.Path)
		return fmt.Errorf("重定向取消")
	}
	tr := &http.Transport{}
	client := &http.Client{
		Transport:     tr,
		CheckRedirect: redirect,
	}
	reqest, _ := http.NewRequest("POST", "http://weather.raychou.com/?/search", strings.NewReader(data.Encode()))
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqest.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/28.0.1500.52 Chrome/28.0.1500.52 Safari/537.36")
	reqest.Header.Set("Connection", "close")
	reqest.Header.Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
	reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Set("Cache-Control", "max-age=0")
	reqest.Header.Set("Connection", "keep-alive")
	response, _ := client.Do(reqest)
	if response.StatusCode == 302 {
		//		fmt.Println(response.Header)
		//		fmt.Println("\n")
		//		fmt.Println(response.Request.Header)
		j := response.Header.Get("Location")
		id := Substr(j, 11, 5)
		return string(id)
	}
	return
}

func getweather(city string) (t string) {
	url := "http://weather.raychou.com/?/detail/" + city + "/count_1/rss"
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", url, nil)
	response, _ := client.Do(reqest)
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		var f Rss
		err := xml.Unmarshal(body, &f)
		if err != nil {
			panic(err)
		}
		for _, channel := range f.Channel {
			for _, item := range channel.Item {
				t := item.Description
				return t
			}
		}
	}
	return
}

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
