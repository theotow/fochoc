package main

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/levigross/grequests"
	"go.uber.org/zap"
)

var coinToResolverMapping = map[string]func(coinUppercase string, address string) (map[string]float64, error){
	"BTC":   getBalanceChainz,
	"LTC":   getBalanceChainz,
	"DASH":  getBalanceChainz,
	"STRAT": getBalanceChainz,
	"LUX":   getBalanceChainz,
	"DGB":   getBalanceChainz,
	"XZC":   getBalanceChainz,
	"VIA":   getBalanceChainz,
	"VTC":   getBalanceChainz,
	"ETH":   getER20Tokens,
	"LSK":   getBalanceLisk,
}

func getCoinMappingToArray() []string {
	var SupportedColdWalletCoins = []string{}
	for currency := range coinToResolverMapping {
		SupportedColdWalletCoins = append(SupportedColdWalletCoins, currency)
	}
	return SupportedColdWalletCoins
}

func resolveCoinHandler(input interface{}) (interface{}, error) {
	balanceStruct, ok := input.(BalanceSimple)
	if ok != true {
		return nil, errors.New("not BalanceSimple struct")
	}
	handler, ok := coinToResolverMapping[balanceStruct.Currency]
	if ok != true {
		return nil, errors.New("no handler found for: " + balanceStruct.Currency)
	}
	res, err := handler(balanceStruct.Currency, balanceStruct.Address)
	if err != nil {
		return nil, err
	}
	output := []BalanceSimple{}
	for currency, balance := range res {
		output = append(output, BalanceSimple{
			Address:  balanceStruct.Address,
			Comment:  balanceStruct.Comment,
			Currency: currency,
			Balance:  balance,
		})
	}

	return output, nil
}

func resolveCoins(coins []BalanceSimple) []BalanceSimple {
	// wrap
	workload := []AsyncTask{}
	for _, coin := range coins {
		workload = append(workload, AsyncTask{
			task: coin,
			exec: resolveCoinHandler,
		})
	}
	res := asyncRun(workload, 10)
	// unwrap
	output := []BalanceSimple{}
	for _, taskRes := range res {
		if t, ok := taskRes.GetResult().([]BalanceSimple); ok {
			output = append(output, t...)
		}
	}
	return output
}

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

func getBalanceEtc(coin string, address string) (map[string]float64, error) {
	output := make(map[string]float64)
	c := colly.NewCollector()

	c.OnHTML(".addr-details .dl-horizontal", func(e *colly.HTMLElement) {
		txt := e.DOM.Text()
		params := getParams(`(?P<Value>[0-9\.]{0,40})\sEther`, txt)
		if len(params) == 1 {
			output[coin], _ = strconv.ParseFloat(params["Value"], 64)
		}
	})
	c.Visit("https://gastracker.io/addr/" + address)
	return output, nil
}

func getBalanceNeo(coin string, address string) (map[string]float64, error) {
	output := make(map[string]float64)
	c := colly.NewCollector()

	c.OnHTML("#wallet-graph", func(e *colly.HTMLElement) {
		txt, _ := e.DOM.Find(".individual-wallet-box:nth-child(1)").Html()
		params := getParams(`\>(?P<Value>[0-9\.]{0,40})<\/p\>`, txt)
		if len(params) == 1 {
			output[coin], _ = strconv.ParseFloat(params["Value"], 64)
		}
	})
	c.Visit("https://neoscan.io/address/" + address)
	return output, nil
}

func getER20Tokens(coinUppercase string, address string) (map[string]float64, error) {
	output := make(map[string]float64)
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

	logger.Debug("fetched coin erc20", zap.Int("coins", len(output)))
	if len(output) == 0 {
		return output, errors.New("erc20-address: " + address + " didnt yield anything")
	}
	return output, nil
}

func getCoins(skip int, limit int) map[string]Coin {
	resp, err := grequests.Get("https://api.coinmarketcap.com/v2/ticker/?convert=BTC&sort=id&start="+strconv.Itoa(skip)+"&limit="+strconv.Itoa(limit), nil)
	if err != nil {
		panic(errors.New("request failed"))
	}
	var coinRes struct {
		Data map[string]Coin `json:"data"`
	}
	err = resp.JSON(&coinRes)
	if err != nil {
		panic(errors.New("request decode failed"))
	}
	output := make(map[string]Coin)
	for _, coin := range coinRes.Data {
		output[coin.Symbol] = Coin{
			Id:       coin.Id,
			Name:     coin.Name,
			Symbol:   coin.Symbol,
			Quote:    coin.Quote,
			BtcPrice: coin.Quote["BTC"]["price"], // TODO: find safer way
			UsdPrice: coin.Quote["USD"]["price"], // TODO: find safer way
		}
	}
	return output
}

func getBalanceChainz(coinUppercase string, address string) (map[string]float64, error) {
	output := make(map[string]float64)
	resp, err := grequests.Get("https://chainz.cryptoid.info/"+strings.ToLower(coinUppercase)+"/api.dws?q=getbalance&a="+address, nil)
	if err != nil {
		return output, errors.New("request failed to fetch balance of " + coinUppercase + " failed")
	}
	val, err := strconv.ParseFloat(resp.String(), 64)
	if err != nil {
		return output, errors.New("parse float blanace of " + coinUppercase + " failed")
	}
	output[coinUppercase] = val
	logger.Debug("fetched coin balance chainz", zap.String("coin", coinUppercase), zap.Float64("balance", val))
	return output, nil
}

func getBalanceLisk(coin string, address string) (map[string]float64, error) {
	output := make(map[string]float64)
	resp, err := grequests.Get("https://explorer.lisk.io/api/getAccount?address="+address, nil)
	if err != nil {
		return nil, errors.New("request failed to fetch balance of " + coin + " failed")
	}
	data := struct {
		Balance string `json:"balance"`
	}{}
	resp.JSON(&data)
	balanceFloat, err := strconv.ParseFloat(data.Balance, 64)
	output[coin] = balanceFloat / 100000000
	if err != nil {
		return nil, errors.New("int parse failed for " + coin)
	}
	return output, nil
}

type AsyncTask struct {
	error  error
	task   interface{}
	result interface{}
	exec   func(interface{}) (interface{}, error)
}

func (at *AsyncTask) Run() AsyncTask {
	res, err := at.exec(at.task)
	at.result = res
	at.error = err
	return *at
}
func (at *AsyncTask) GetResult() interface{} {
	return at.result
}

func asyncRun(workload []AsyncTask, workerCount int) []AsyncTask {
	inQueue := make(chan AsyncTask, len(workload))
	outQueue := make(chan AsyncTask, len(workload))
	done := []AsyncTask{}
	if len(workload) == 0 {
		return []AsyncTask{}
	}

	// workers
	for i := 0; i < workerCount; i++ {
		go func(workerId int) {
			for {
				select {
				case task, more := <-inQueue:
					if more {
						outQueue <- task.Run()
					} else {
						return
					}
				}
			}
		}(i)
	}
	// writer
	go func() {
		for _, task := range workload {
			inQueue <- task
		}
		close(inQueue)
	}()

	// chan to array
	for taskDone := range outQueue {
		done = append(done, taskDone)
		if len(done) == len(workload) {
			break
		}
	}
	return done
}

func getCoinsAsync() map[string]Coin {
	output := map[string]Coin{}
	limit := 100       // coins fetched per request
	totalLimit := 2000 // max coins fetched in total
	workload := []AsyncTask{}
	handler := func(input interface{}) (interface{}, error) {
		skip, ok := input.(int)
		if ok != true {
			return nil, errors.New("not int type")
		}
		res := getCoins(skip, limit)
		logger.Debug("fetched coin infos", zap.Int("skip", skip))
		return res, nil
	}
	// wrap
	for i := 0; i < totalLimit; {
		workload = append(workload, AsyncTask{
			task: i + 1,
			exec: handler,
		})
		i += limit
	}
	res := asyncRun(workload, 10)
	// unwrap
	for _, reqRes := range res {
		for currency, coin := range reqRes.GetResult().(map[string]Coin) {
			output[currency] = coin
		}
	}
	return output
}
