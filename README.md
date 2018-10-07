# Fochoc - Fortune Choc [![Build Status](https://travis-ci.com/theotow/fochoc.svg?branch=master)](https://travis-ci.com/theotow/fochoc) [![Build status](https://ci.appveyor.com/api/projects/status/w0s545dhqeqmqhna?svg=true)](https://ci.appveyor.com/project/theotow/fochoc) [![Go Report Card](https://goreportcard.com/badge/github.com/theotow/fochoc)](https://goreportcard.com/report/github.com/theotow/fochoc)

<img src="https://github.com/theotow/fochoc/blob/master/assets/preview.png " alt="" width="600" />

## Installation

1. ``` make build ```
2. ```mv ./fochoc_(win|mac|linux) /usr/local/bin/fochoc```
3. use ```$ fochoc```

## Supported Exchanges
- binance
- kraken
- poloniex
- bittrex

## Supported Coins / Tokens
- https://api.coinmarketcap.com/v2/listings/ (kraken, poloniex, binance, bittrex)
- ERC-20 coin names with 3 digits listed on coinmarketcap
- Coldwallet Coins (BTC, LTC, DASH, STRAT, LUX, DGB, XZC, VIA, VTC, ETH, ETC, NEO, LSK), more on request

## Platform Support

- Mac
- Windows
- Linux

## Planned Features

- [x] Improve error reporting
- [x] Add ERC-20 token support
- [x] Add Bittrex exchange
- [x] Improve test coverage
- [x] Coldwallet coins

## Ideas

- CSV export
- Group By exchange / address / etc.
- Order by
- Support for more currencies
- Faster

