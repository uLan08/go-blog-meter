package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/grokify/html-strip-tags-go"
)

func cleanStr(str *string) {
	re := regexp.MustCompile(`\r?\n`)
	*str = strip.StripTags(*str)
	*str = re.ReplaceAllString(strings.Replace(*str, "&nbsp;", "", -1), "")
	*str = strings.Replace(*str, "\t", "", -1)
}

func getBody(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func main() {
	str := getBody("http://lauftechnik.de/")
	cleanStr(&str)
	words := strings.Split(str, " ")
	for _, word := range words {
		fmt.Println(word)
	}

	// fmt.Println(words, len(words))
	// urls := xurls.Strict().FindAllString(str, -1)
	// for _, url := range urls {
	// 	fmt.Println(url)
	// }
}
