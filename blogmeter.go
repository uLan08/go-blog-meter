package main

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

var currencies = map[string]string{
	"AED": "د.إ",
	"AFN": "؋",
	"ALL": "L",
	"AMD": "֏",
	"ANG": "ƒ",
	"AOA": "Kz",
	"ARS": "$",
	"AUD": "$",
	"AWG": "ƒ",
	"AZN": "₼",
	"BAM": "KM",
	"BBD": "$",
	"BDT": "৳",
	"BGN": "лв",
	"BHD": ".د.ب",
	"BIF": "FBu",
	"BMD": "$",
	"BND": "$",
	"BOB": "$b",
	"BRL": "R$",
	"BSD": "$",
	"BTC": "฿",
	"BTN": "Nu.",
	"BWP": "P",
	"BYR": "p.",
	"BZD": "BZ$",
	"CAD": "$",
	"CDF": "FC",
	"CHF": "CHF",
	"CLP": "$",
	"CNY": "¥",
	"COP": "$",
	"CRC": "₡",
	"CUC": "$",
	"CUP": "₱",
	"CVE": "$",
	"CZK": "Kč",
	"DJF": "Fdj",
	"DKK": "kr",
	"DOP": "RD$",
	"DZD": "دج",
	"EEK": "kr",
	"EGP": "£",
	"ERN": "Nfk",
	"ETB": "Br",
	"ETH": "Ξ",
	"EUR": "€",
	"FJD": "$",
	"FKP": "£",
	"GBP": "£",
	"GEL": "₾",
	"GGP": "£",
	"GHC": "₵",
	"GHS": "GH₵",
	"GIP": "£",
	"GMD": "D",
	"GNF": "FG",
	"GTQ": "Q",
	"GYD": "$",
	"HKD": "$",
	"HNL": "L",
	"HRK": "kn",
	"HTG": "G",
	"HUF": "Ft",
	"IDR": "Rp",
	"ILS": "₪",
	"IMP": "£",
	"INR": "₹",
	"IQD": "ع.د",
	"IRR": "﷼",
	"ISK": "kr",
	"JEP": "£",
	"JMD": "J$",
	"JOD": "JD",
	"JPY": "¥",
	"KES": "KSh",
	"KGS": "лв",
	"KHR": "៛",
	"KMF": "CF",
	"KPW": "₩",
	"KRW": "₩",
	"KWD": "KD",
	"KYD": "$",
	"KZT": "лв",
	"LAK": "₭",
	"LBP": "£",
	"LKR": "₨",
	"LRD": "$",
	"LSL": "M",
	"LTC": "Ł",
	"LTL": "Lt",
	"LVL": "Ls",
	"LYD": "LD",
	"MAD": "MAD",
	"MDL": "lei",
	"MGA": "Ar",
	"MKD": "ден",
	"MMK": "K",
	"MNT": "₮",
	"MOP": "MOP$",
	"MRO": "UM",
	"MUR": "₨",
	"MVR": "Rf",
	"MWK": "MK",
	"MXN": "$",
	"MYR": "RM",
	"MZN": "MT",
	"NAD": "$",
	"NGN": "₦",
	"NIO": "C$",
	"NOK": "kr",
	"NPR": "₨",
	"NZD": "$",
	"OMR": "﷼",
	"PAB": "B/.",
	"PEN": "S/.",
	"PGK": "K",
	"PHP": "₱",
	"PKR": "₨",
	"PLN": "zł",
	"PYG": "Gs",
	"QAR": "﷼",
	"RMB": "￥",
	"RON": "lei",
	"RSD": "Дин.",
	"RUB": "₽",
	"RWF": "R₣",
	"SAR": "﷼",
	"SBD": "$",
	"SCR": "₨",
	"SDG": "ج.س.",
	"SEK": "kr",
	"SGD": "$",
	"SHP": "£",
	"SLL": "Le",
	"SOS": "S",
	"SRD": "$",
	"SSP": "£",
	"STD": "Db",
	"SVC": "$",
	"SYP": "£",
	"SZL": "E",
	"THB": "฿",
	"TJS": "SM",
	"TMT": "T",
	"TND": "د.ت",
	"TOP": "T$",
	"TRL": "₤",
	"TRY": "₺",
	"TTD": "TT$",
	"TVD": "$",
	"TWD": "NT$",
	"TZS": "TSh",
	"UAH": "₴",
	"UGX": "USh",
	"USD": "$",
	"UYU": "$U",
	"UZS": "лв",
	"VEF": "Bs",
	"VND": "₫",
	"VUV": "VT",
	"WST": "WS$",
	"XAF": "FCFA",
	"XBT": "Ƀ",
	"XCD": "$",
	"XOF": "CFA",
	"XPF": "₣",
	"YER": "﷼",
	"ZAR": "R",
	"ZWD": "Z$",
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
		output = " "
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		output = string(body)
	}
	result <- output
	wg.Done()
}

func hasPrice(body string) bool {
	baseExp := `\s?\d+\.?\,?\d+`
	for currency, symbol := range currencies {
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

func resolveUrls(urls []string) int {
	var count = 0
	var wg sync.WaitGroup
	resultsChan := make(chan string, len(urls))
	for _, url := range urls {
		fmt.Println(url)
		wg.Add(1)
		go getBody(url, resultsChan, &wg)
	}
	wg.Wait()
	for body := range resultsChan {
		fmt.Println(len(resultsChan), cap(resultsChan))
		fmt.Println(count)
		if hasPrice(body) {
			count++
		}
	}
	return count
}

func main() {
	result := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go getBody("http://lauftechnik.de/", result, &wg)
	wg.Wait()
	str := <-result
	urls := xurls.Relaxed().FindAllString(str, -1)
	count := resolveUrls(urls)
	fmt.Println("count", count)
	cleanStr(&str)
	words := strings.Split(str, " ")
	fmt.Println(len(words))
	// for _, word := range words {
	// 	fmt.Println(word)
	// }

	// fmt.Println(words, len(words))
	// for _, url := range urls {
	// 	fmt.Println(url)
	// }
}
