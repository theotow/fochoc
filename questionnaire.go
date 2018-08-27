package main

import (
	"fmt"

	survey "gopkg.in/AlecAivazis/survey.v1"
)

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
				Options: []string{"Add Exchange", "Add Coldwallet", "Reset Config"},
				Default: "Add Exchange",
			},
		},
	}, &q.answers)
	q.Logic()
}

func (q *questions) ColdWallets() {
	coins := getCoinMappingToArray()
	survey.Ask([]*survey.Question{
		{
			Name:     "coldwallet",
			Validate: survey.Required,
			Prompt: &survey.Select{
				Message: "Which one:",
				Options: coins,
				Default: coins[0],
			},
		},
	}, &q.answers)
}

func (q *questions) Exchange() {
	survey.Ask([]*survey.Question{
		{

			Name:     "exchange",
			Validate: survey.Required,
			Prompt: &survey.Select{
				Message: "Which:",
				Options: ActiveExchangeProviders,
				Default: ActiveExchangeProviders[0],
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

func (q *questions) AskForAddressComment() {
	survey.Ask([]*survey.Question{
		{
			Name:     "address",
			Prompt:   &survey.Input{Message: "Address"},
			Validate: survey.Required,
		},
		{
			Name:   "comment",
			Prompt: &survey.Input{Message: "Comment"},
		},
	}, &q.answers)
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

func getConfigKeysOfProviderByID(id string) []string {
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
		for _, exchange := range ActiveExchangeProviders {
			if q.getKeySafe("exchange") == exchange {
				configKeys := getConfigKeysOfProviderByID(exchange)
				q.ExchangeCreds(configKeys)
				config := NewFileConfig()
				configMap := config.Read()
				// write to config (merge)
				for _, name := range configKeys {
					configMap.addKey(name, fmt.Sprint(q.getKeySafe(name)))
				}
				config.Write(configMap)
				fmt.Println("Done!")
			}
		}
		return
	}
	if q.getKeySafe("what") == "Add Coldwallet" {
		q.ColdWallets()
		q.AskForAddressComment()
		config := NewFileConfig()
		configMap := config.Read()
		configMap.addColdWalletCoin(ColdWalletCoin{
			Currency: fmt.Sprint(q.getKeySafe("coldwallet")),
			Address:  fmt.Sprint(q.getKeySafe("address")),
			Comment:  fmt.Sprint(q.getKeySafe("comment")),
		})
		config.Write(configMap)
		fmt.Println("Done!")
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
