package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/uLan08/go-blog-meter"
)

func main() {
	result := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go blogmeter.GetBody("http://lauftechnik.de/", result, &wg)
	wg.Wait()
	str := <-result
	urls := blogmeter.ExtractUrls(str)
	count := blogmeter.ResolveUrls(urls)
	blogmeter.CleanStr(&str)
	words := strings.Split(str, " ")
	fmt.Println("Links:", count)
	fmt.Println("Words:", len(words))
}
