package blogmeter

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"mvdan.cc/xurls"
)

func main() {
	resp, err := http.Get("http://lauftechnik.de/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	str := string(body)
	urls := xurls.Strict().FindAllString(str, -1)
	for _, url := range urls {
		fmt.Println(url)
	}
}
