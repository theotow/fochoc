package main


type erc20 struct {
	ResultRaw []BalanceSimple
}

type methodsEcr20 struct{}

func initEcr20(c ConfigInterface) *erc20 {
	output := []BalanceSimple{}
	// TODO: token is a misleading word here
	tokens := c.GetTokens()
	for _, token := range tokens {
		currencyMap := GetER20Tokens(token.Address)
		for currency, balance := range currencyMap {
			output = append(output, BalanceSimple{
				Address:  token.Address,
				Comment:  token.Comment,
				Currency: currency,
				Balance:  balance,
			})
		}
	}
	return &erc20{ResultRaw: output}
}

func (m *methodsEcr20) Get(c ConfigInterface) ConfigProviderInterface {
	return initEcr20(c)
}

func (m *methodsEcr20) ConfigKeys() []string {
	return []string{}
}

func NewEcr20() *methodsEcr20 {
	return &methodsEcr20{}
}

func (k *erc20) GetCurrencyValue(name string) float64 {
	return 0.00
}

func (k *erc20) GetAll(keys []string) []BalanceSimple {
	needleMap := toMap(keys)
	output := []BalanceSimple{}
	for _, result := range k.ResultRaw {
		// only let whitelisted pass
		if needleMap[result.Currency] == true {
			output = append(output, result)
		}
	}
	return output
}

func (k *erc20) AddTestBalance(name string, value float64) {

}
