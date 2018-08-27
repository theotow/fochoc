package main

type coldwallet struct {
	ResultRaw []BalanceSimple
}

type methodsColdWallet struct{}

func initColdWallet(c ConfigInterface) (*coldwallet, error) {
	workload := []BalanceSimple{}
	coins := c.GetColdWalletCoins()
	for _, coin := range coins {
		workload = append(workload, BalanceSimple{
			Address:  coin.Address,
			Comment:  coin.Comment,
			Currency: coin.Currency,
			Balance:  0.00,
		})
	}
	return &coldwallet{ResultRaw: resolveCoins(workload)}, nil
}

func (m *methodsColdWallet) Get(c ConfigInterface) (ConfigProviderInterface, error) {
	return initColdWallet(c)
}

func (m *methodsColdWallet) ConfigKeys() []string {
	return []string{}
}

// NewColdWallet returns a new coldwallet provider
func NewColdWallet() *methodsColdWallet {
	return &methodsColdWallet{}
}

func (k *coldwallet) GetCurrencyValue(name string) float64 {
	return 0.00
}

func (k *coldwallet) GetAll(keys []string) []BalanceSimple {
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

func (k *coldwallet) AddTestBalance(name string, value float64) {

}
