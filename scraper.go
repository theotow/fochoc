package main

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func getParams(regEx, url string) (paramsMap map[string]string) {

	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}

// GetER20Tokens returns eth balance + erc20 balances
func GetER20Tokens(address string) map[string]float64 {
	var output = make(map[string]float64)
	c := colly.NewCollector()

	c.OnHTML("#balancelist a", func(e *colly.HTMLElement) {
		html, _ := e.DOM.Html()
		result := getParams(`(\>|\s)(?P<Name>[A-Z]{0,3})\<\/i\>(.*)\<br\/\>(?P<Value>[0-9,\.]{0,20})\s[A-Z]{0,3}`, html)
		if len(result) == 3 {
			output[result["Name"]], _ = strconv.ParseFloat(strings.Replace(result["Value"], ",", "", -1), 64)
		}
	})

	c.OnHTML("#ContentPlaceHolder1_divSummary table", func(e *colly.HTMLElement) {
		html, _ := e.DOM.Find("tbody tr:nth-child(1)").Html()
		html = strings.Replace(html, "<b>.</b>", ".", -1)
		html = strings.Replace(html, ",", "", -1)
		result := getParams(`(?P<Value>[0-9\.]{0,40})\sEther`, html)
		if reflect.TypeOf(result["Value"]).String() == "string" && len(result["Value"]) > 0 {
			res, err := strconv.ParseFloat(result["Value"], 64)
			if err == nil && res > 0 {
				output["ETH"] = res
			}
		}
	})

	c.Visit("https://etherscan.io/address/" + address)

	return output
}
