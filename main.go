package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/olekukonko/tablewriter"
	cli "gopkg.in/urfave/cli.v1"
)

// AppName is the name if the binary on the end users machine
const AppName = "fochoc"

// AppVersion is the version of the current app
const AppVersion = "0.0.6"

// Providers is an array of available providers
var Providers = []Provider{
	{id: "binance", factory: NewBinance()},
	{id: "kraken", factory: NewKraken()},
	{id: "poloniex", factory: NewPoloniex()},
	{id: "bittrex", factory: NewBittrex()},
	{id: "coldwallet", factory: NewColdWallet()},
}

// ActiveProviders is a list of all active provider ids
var ActiveProviders = []string{
	"binance",
	"kraken",
	"poloniex",
	"bittrex",
	"coldwallet",
}

// ActiveExchangeProviders is a list of all active exchange provider ids
var ActiveExchangeProviders = []string{
	"binance",
	"kraken",
	"poloniex",
	"bittrex",
}

// ProviderInterface the interface of a provider factory
type ProviderInterface interface {
	Get(c ConfigInterface) (ConfigProviderInterface, error)
	ConfigKeys() []string
}

// Provider is a struct which describes a provider with its id, maybe instance, factory
type Provider struct {
	instance ConfigProviderInterface
	factory  ProviderInterface
	id       string
}

// ProviderError contains meta data and error message to return to enduser
type ProviderError struct {
	errorMsg   string
	providerID string
}

func (e *ProviderError) Error() string {
	return "provider(" + e.providerID + ") could not be initialized, maybe the key/secret is invalid, to edit config visit ~/.fochocconfig.json, error: " + e.errorMsg
}

// NewProviderError returns a Provider Error which contains meta data and error message to return to enduser
func NewProviderError(err string, id string) *ProviderError {
	return &ProviderError{errorMsg: err, providerID: id}
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
	listOfCoinsToFetch := toArray(coinMap)
	res := p.instance.GetAll(listOfCoinsToFetch)
	output := []Balance{}
	for _, value := range res {
		if value.Balance > 0 {
			output = append(output, Balance{Provider: *p, Coin: coinMap[value.Currency], Balance: value})
		}
	}
	return output
}

// Coin contains meta data like price usd... for a certain coin
type Coin struct {
	Id       int                           `json:"id"`
	Name     string                        `json:"name"`
	Symbol   string                        `json:"symbol"`
	Quote    map[string]map[string]float64 `json:"quotes"`
	BtcPrice float64
	UsdPrice float64
}

// Balance is the main struct which contains the links to provider,
// the coin and its current value in usd/btc and the actual address and comment
type Balance struct {
	Coin     Coin
	Provider Provider
	Balance  BalanceSimple
}

// NewBalance allows to create a Balance struct with all related structs
// mainly for testing
func NewBalance(
	coinName string,
	balance float64,
	id int,
	btcPrice float64,
	usdPrice float64,
	providerID string,
	providerFactory ProviderInterface,
) Balance {
	return Balance{
		Coin: Coin{
			Id:       id,
			Symbol:   coinName,
			Name:     "name",
			Quote:    make(map[string]map[string]float64),
			BtcPrice: btcPrice,
			UsdPrice: usdPrice,
		},
		Balance: BalanceSimple{
			Comment:  "comment",
			Address:  "none",
			Currency: coinName,
			Balance:  balance,
		},
		Provider: Provider{
			id:      providerID,
			factory: providerFactory,
		},
	}
}

// BalanceSimple stores the address related balance
// if crypto is stored on exchange the balance is tied to the currency
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

func (b *Balance) getAddressString() string {
	return b.Balance.Address
}

func (b *Balance) getCommentString() string {
	return b.Balance.Comment
}

func (b *Balance) getProviderID() string {
	return b.Provider.id
}

func (b *Balance) getUsdBalanceString() string {
	return fmt.Sprintf("%f", b.getUsdBalance())
}

func (b *Balance) getBtcBalanceString() string {
	return fmt.Sprintf("%f", b.getBtcBalance())
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("I crashed, im so so sorry")
		}
	}()
	app := cli.NewApp()
	app.Name = AppName
	app.Usage = "how to run this"
	app.Version = AppVersion
	app.Action = func(c *cli.Context) error {
		return showOverview()
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
		fmt.Println(err)
	}
}

type tableData struct {
	balances []Balance
	usd      float64
	btc      float64
	errorMsg error
}

func getTableData() tableData {
	coins := getCoinsAsync()
	config := NewFileConfig()
	providers, err := initProviders(ActiveProviders, config)
	if err != nil {
		return tableData{errorMsg: err}
	}
	res := getAllBalances(providers, coins)
	usdSum, btcSum := getAggSum(res)
	return tableData{balances: res, usd: usdSum, btc: btcSum}
}

func printLogo() {
	file, _ := ioutil.ReadFile("./assets/cookie.txt")
	fmt.Println(string(file[:]))
}

func showOverview() error {
	ch := make(chan tableData)
	go func(c chan tableData) {
		c <- getTableData()
	}(ch)
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	res := <-ch
	close(ch)
	s.Stop()
	if res.errorMsg != nil {
		return res.errorMsg
	}
	printLogo()
	renderTable(res.balances, res.btc, res.usd)
	return nil
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
			balance.getProviderID(),
			balance.getAddressString(),
			balance.getCommentString(),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Asset", "Amount", "Value USD", "Value BTC", "Sources", "Address", "Comment"})
	table.SetFooter([]string{"", "", fmt.Sprintf("$%f", sumUsd), fmt.Sprintf("BTC %f", sumBtc), "", "", ""})
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
	for key := range coinMap {
		output = append(output, key)
	}
	return output
}

func initProviders(neededProviders []string, config ConfigInterface) ([]Provider, error) {
	var activeProvider []Provider
	neededProvidersIDMap := toMap(neededProviders)
	for _, provider := range Providers {
		if _, ok := neededProvidersIDMap[provider.id]; ok {
			// TODO: maybe refactor Provider to Provider / ProviderInited
			if provider.isValid(config) {
				instance, err := provider.factory.Get(config)
				if err != nil {
					return nil, NewProviderError(err.Error(), provider.id)
				}
				activeProvider = append(activeProvider, Provider{
					id:       provider.id,
					factory:  provider.factory,
					instance: instance,
				})
			}
		}
	}
	return activeProvider, nil
}
