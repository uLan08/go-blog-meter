package blogmeter

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/grokify/html-strip-tags-go"
	"mvdan.cc/xurls"
)

type result struct {
	body    string
	success bool
	url     string
}

func hasPrice(body string) bool {
	baseExp := `\s?\d+\.?\,?\d+`
	for currency, symbol := range Currencies {
		currencyRe := regexp.MustCompile(currency + baseExp)
		var prefix string
		if symbol == `$` {
			prefix = `\` + symbol
		} else {
			prefix = symbol
		}
		symbolRe := regexp.MustCompile(prefix + baseExp)
		hasCurrency := currencyRe.MatchString(body)
		hasSymbol := symbolRe.MatchString(body)
		if hasCurrency || hasSymbol {
			return true
		}
	}
	return false
}

func uniqueSlice(slc []string) []string {
	copy := make([]string, 0, len(slc))
	set := make(map[string]struct{})

	for _, val := range slc {
		if _, ok := set[val]; !ok {
			set[val] = struct{}{}
			copy = append(copy, val)
		}
	}
	return copy
}

func cleanStr(str *string) {
	re := regexp.MustCompile(`\r?\n`)
	*str = strip.StripTags(*str)
	*str = re.ReplaceAllString(strings.Replace(*str, "&nbsp;", "", -1), "")
	*str = strings.Replace(*str, "\t", "", -1)
}

func getBody(url string) result {
	resp, err := http.Get(url)
	var res result
	if err != nil {
		fmt.Println(err)
		res.body = ""
		res.success = false
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		res.body = string(body)
		res.success = true
	}
	res.url = url
	return res
}

func extractUrls(str string) []string {
	return xurls.Relaxed().FindAllString(str, -1)
}

func resolveUrls(urls []string) int {
	var count = 0
	uniqueUrls := uniqueSlice(urls)
	var wg sync.WaitGroup
	resultsChan := make(chan result, len(urls))
	for _, url := range uniqueUrls {
		wg.Add(1)
		go func(url string) {
			res := getBody(url)
			resultsChan <- res
			wg.Done()
		}(url)
	}
	wg.Wait()
	close(resultsChan)
	for res := range resultsChan {
		if hasPrice(res.body) {
			fmt.Println(res.url)
			count++
		}
	}
	return count
}

// RateBlog rates the given url
func RateBlog(url string) (int, int) {
	res := getBody(url)
	if !res.success {
		log.Fatal("Seems like the url is not available, try another one")
	}
	str := res.body
	urls := extractUrls(str)
	count := resolveUrls(urls)
	cleanStr(&str)
	words := strings.Split(str, " ")
	return len(words), count
}
