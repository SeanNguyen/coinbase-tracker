package bitfinex

import (
	"github.com/seannguyen/coin-tracker/internal/services/cryto_exchanges"
	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"github.com/spf13/viper"
	"strconv"
	"github.com/rhymond/go-money"
	"strings"
)

const exchangeName = "bitfinex"

type Exchange struct{}

func (*Exchange) GetBalances() ([]*cryto_exchanges.BalanceData, error) {
	client := bitfinex.NewClient().Auth(viper.GetString("BITFINEX_API_KEY"), viper.GetString("BITFINEX_API_SECRET"))
	balances, err := client.Balances.All()

	if err != nil {
		return nil, err
	}

	var balancesData []*cryto_exchanges.BalanceData
	for _, balance := range balances {
		amount, err := strconv.ParseFloat(balance.Amount, 64)
		if err != nil {
			return nil, err
		}
		if amount <= 0 {
			continue
		}
		currency := strings.ToUpper(balance.Currency)
		balancesData = append(balancesData, &cryto_exchanges.BalanceData{
			Type: getCurrencyType(currency),
			Amount: amount,
			ExchangeName: exchangeName,
			Currency: currency,
		})
	}
	return balancesData, nil
}

func getCurrencyType(currency string) int {
	currencyRecord := money.GetCurrency(currency)
	if currencyRecord != nil {
		return cryto_exchanges.Fiat
	}
	return cryto_exchanges.Crypto
}