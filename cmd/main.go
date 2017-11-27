package main

import (
	"fmt"

	"github.com/uLan08/go-blog-meter"
)

func main() {
	words, links := blogmeter.RateBlog("http://lauftechnik.de/")
	fmt.Println("Links:", links)
	fmt.Println("Words:", words)
}
