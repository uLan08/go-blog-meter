package blogmeter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/grokify/html-strip-tags-go"
	"mvdan.cc/xurls"
)

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

func getBody(url string, result chan string, wg *sync.WaitGroup) {
	resp, err := http.Get(url)
	var output string
	if err != nil {
		output = ""
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		output = string(body)
	}
	result <- output
	wg.Done()
}

func extractUrls(str string) []string {
	return xurls.Relaxed().FindAllString(str, -1)
}

func resolveUrls(urls []string) int {
	var count = 0
	uniqueUrls := uniqueSlice(urls)
	var wg sync.WaitGroup
	resultsChan := make(chan string, len(urls))
	for _, url := range uniqueUrls {
		fmt.Println(url)
		wg.Add(1)
		go getBody(url, resultsChan, &wg)
	}
	wg.Wait()
	close(resultsChan)
	for body := range resultsChan {
		if hasPrice(body) {
			count++
		}
	}
	return count
}

func RateBlog(url string) (int, int) {
	result := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go getBody(url, result, &wg)
	wg.Wait()
	str := <-result
	urls := extractUrls(str)
	count := resolveUrls(urls)
	cleanStr(&str)
	words := strings.Split(str, " ")
	return len(words), count
}
