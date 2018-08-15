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

func GetER20Tokens(address string) map[string]float64 {
	var output = make(map[string]float64)
	c := colly.NewCollector()

	c.OnHTML("#balancelist a", func(e *colly.HTMLElement) {
		html, _ := e.DOM.Html()
		result := getParams(`\"\>(?P<Name>[A-Z]{0,3})\<\/i\>(.*)\<br\/\>(?P<Value>[0-9,]{0,10})\s[A-Z]{3}`, html)
		if len(result) == 3 {
			output[result["Name"]], _ = strconv.ParseFloat(strings.Replace(result["Value"], ",", ".", -1), 64)
		}
	})

	c.OnHTML("#ContentPlaceHolder1_divSummary table", func(e *colly.HTMLElement) {
		html := e.DOM.Find("tbody tr:nth-child(1)").Text()
		result := getParams(`(?P<Value>[0-9\.]{0,20})\sEther`, html)
		if reflect.TypeOf(result["Value"]).String() == "string" && len(result["Value"]) > 0 {
			output["ETH"], _ = strconv.ParseFloat(result["Value"], 64)
		}
	})

	c.Visit("https://etherscan.io/address/" + address)

	return output
}
