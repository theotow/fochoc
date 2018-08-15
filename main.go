package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/olekukonko/tablewriter"
	survey "gopkg.in/AlecAivazis/survey.v1"
	cli "gopkg.in/urfave/cli.v1"

	"github.com/levigross/grequests"
)

const CoinList = "https://api.coinmarketcap.com/v2/ticker/?convert=BTC"
const AppName = "fochoc"
const AppVersion = "0.0.1"

var Providers = []Provider{
	Provider{id: "binance", factory: NewBinance()},
	Provider{id: "kraken", factory: NewKraken()},
	Provider{id: "poloniex", factory: NewPoloniex()},
	Provider{id: "bittrex", factory: NewBittrex()},
}

var ActiveProviders = []string{
	"binance",
	"kraken",
	"poloniex",
	"bittrex",
}

type ProviderInterface interface {
	Get(c ConfigInterface) ConfigProviderInterface
	ConfigKeys() []string
}

type Provider struct {
	instance ConfigProviderInterface
	factory  ProviderInterface
	id       string
}

func (p *Provider) isValid(c ConfigInterface) bool {
	keys := p.factory.ConfigKeys()
	var isValid = true
	for _, key := range keys {
		if c.GetKey(key) == "" {
			isValid = false
		}
	}
	return isValid
}
func (p *Provider) getCoinsOfProvider(coinMap map[string]Coin) []Balance {
	res := p.instance.GetAll(toArray(coinMap))
	var output []Balance
	for _, value := range res {
		if value.Balance > 0 {
			output = append(output, Balance{Provider: p, Coin: coinMap[value.Currency], Balance: value})
		}
	}
	return output
}

type CoinListResponse struct {
	Data map[string]Coin `json:"data"`
}

type Coin struct {
	Id       int                           `json:"id"`
	Name     string                        `json:"name"`
	Symbol   string                        `json:"symbol"`
	Quote    map[string]map[string]float64 `json:"quotes"`
	BtcPrice float64
	UsdPrice float64
}

type Balance struct {
	Coin     Coin
	Provider *Provider
	Balance  BalanceSimple
}

type BalanceSimple struct {
	Comment  string
	Address  string
	Currency string
	Balance  float64
}

func (b *Balance) getBtcBalance() float64 {
	return b.Coin.BtcPrice * b.Balance.Balance
}

func (b *Balance) getUsdBalance() float64 {
	return b.Coin.UsdPrice * b.Balance.Balance
}

func (b *Balance) getSymbolString() string {
	return b.Coin.Symbol
}

func (b *Balance) getBalanceString() string {
	return fmt.Sprintf("%f", b.Balance.Balance)
}

func (b *Balance) getProviderId() string {
	return b.Provider.id
}

func (b *Balance) getUsdBalanceString() string {
	return fmt.Sprintf("%f", b.getUsdBalance())
}

func (b *Balance) getBtcBalanceString() string {
	return fmt.Sprintf("%f", b.getBtcBalance())
}

type questions struct {
	answers map[string]interface{}
}

func startQuestions() *questions {
	return &questions{
		answers: make(map[string]interface{}),
	}
}

func (q *questions) Main() {
	survey.Ask([]*survey.Question{
		{
			Name:     "what",
			Validate: survey.Required,
			Prompt: &survey.Select{
				Message: "What you want to do:",
				Options: []string{"Add Exchange", "Reset Config"},
				Default: "Add Exchange",
			},
		},
	}, &q.answers)
	q.Logic()
}

func (q *questions) Exchange() {
	survey.Ask([]*survey.Question{
		{
			Name:     "exchange",
			Validate: survey.Required,
			Prompt: &survey.Select{
				Message: "Which:",
				Options: ActiveProviders,
				Default: ActiveProviders[0],
			},
		},
	}, &q.answers)
}

func (q *questions) ExchangeCreds(settings []string) {
	questions := []*survey.Question{}
	for _, settingName := range settings {
		questions = append(questions, &survey.Question{
			Name:     settingName,
			Prompt:   &survey.Input{Message: settingName},
			Validate: survey.Required,
		})
	}
	survey.Ask(questions, &q.answers)
}

func (q *questions) AreYouSure() {
	survey.Ask([]*survey.Question{
		{
			Name: "reset",
			Prompt: &survey.Select{
				Message: "Are you really sure?:",
				Options: []string{"yes", "no"},
				Default: "no",
			},
			Validate: survey.Required,
		},
	}, &q.answers)
	q.Logic()
}

func (q *questions) getKeySafe(key string) interface{} {
	if val, ok := q.answers[key]; ok {
		return val
	}
	return ""
}

func getConfigKeysOfProviderById(id string) []string {
	for _, provider := range Providers {
		if provider.id == id {
			return provider.factory.ConfigKeys()
		}
	}
	return []string{}
}

func (q *questions) Logic() {
	if q.getKeySafe("what") == "Add Exchange" {
		q.Exchange()
		for _, exchange := range ActiveProviders {
			if q.getKeySafe("exchange") == exchange {
				configKeys := getConfigKeysOfProviderById(exchange)
				q.ExchangeCreds(configKeys)
				config := NewFileConfig()
				configMap := config.Read()
				// write to config (merge)
				for _, name := range configKeys {
					configMap.Keys[name] = fmt.Sprint(q.getKeySafe(name))
				}
				config.Write(configMap)
				fmt.Println("Done!")
			}
		}
		return
	}
	if q.getKeySafe("what") == "Reset Config" {
		if q.getKeySafe("reset") == "" {
			q.AreYouSure()
		} else if q.getKeySafe("reset") == "yes" {
			config := NewFileConfig()
			config.Write(config.GetEmptyConfig())
			fmt.Println("Done!")
		}
		return
	}

}

func main() {
	app := cli.NewApp()
	app.Name = AppName
	app.Usage = "how to run this"
	app.Version = AppVersion
	app.Action = func(c *cli.Context) error {
		showOverview()
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Add exchanges and change settings",
			Action: func(c *cli.Context) error {
				startQuestions().Main()
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func showOverview() {
	coins := getCoins()
	config := NewFileConfig()
	providers := initProviders(ActiveProviders, config)
	res := getAllBalances(providers, coins)
	usdSum, btcSum := getAggSum(res)
	file, _ := ioutil.ReadFile("./assets/cookie.txt")
	fmt.Println(string(file[:]))
	renderTable(res, btcSum, usdSum)
}

func getAllBalances(providers []Provider, coins map[string]Coin) []Balance {
	var balances []Balance
	for _, provider := range providers {
		res := provider.getCoinsOfProvider(coins)
		balances = append(balances, res...)
	}
	return balances
}

func getAggSum(balances []Balance) (float64, float64) {
	usdSum := 0.00
	btcSum := 0.00
	for _, balance := range balances {
		usdSum = balance.getUsdBalance() + usdSum
		btcSum = balance.getBtcBalance() + btcSum
	}
	return usdSum, btcSum
}

func renderTable(data []Balance, sumBtc float64, sumUsd float64) {
	var tableData [][]string
	for _, balance := range data {
		tableData = append(tableData, []string{
			balance.getSymbolString(),
			balance.getBalanceString(),
			balance.getUsdBalanceString(),
			balance.getBtcBalanceString(),
			balance.getProviderId(),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Asset", "Amount", "Value USD", "Value BTC", "Sources"})
	table.SetFooter([]string{"", "", fmt.Sprintf("$%f", sumUsd), fmt.Sprintf("BTC %f", sumBtc), ""})
	table.SetBorder(false)
	table.AppendBulk(tableData)
	table.Render()
}

func toMap(array []string) map[string]bool {
	output := make(map[string]bool)
	for _, val := range array {
		output[val] = true
	}
	return output
}

func toArray(coinMap map[string]Coin) []string {
	var output []string
	for key, _ := range coinMap {
		output = append(output, key)
	}
	return output
}

func initProviders(neededProviders []string, config ConfigInterface) []Provider {
	var activeProvider []Provider
	neededProvidersIdMap := toMap(neededProviders)
	for _, provider := range Providers {
		if _, ok := neededProvidersIdMap[provider.id]; ok {
			// TODO: maybe refactor Provider to Provider / ProviderInited
			if provider.isValid(config) {
				activeProvider = append(activeProvider, Provider{
					id:       provider.id,
					factory:  provider.factory,
					instance: provider.factory.Get(config),
				})
			}
		}
	}
	return activeProvider
}

func getCoins() map[string]Coin {
	resp, err := grequests.Get(CoinList, nil)
	if err != nil {
		panic(errors.New("request failed"))
	}
	var coinRes CoinListResponse
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
