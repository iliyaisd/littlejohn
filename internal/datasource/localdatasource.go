package datasource

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/iliyaisd/littlejohn/ljlib"
	"github.com/shopspring/decimal"
)

const mockDailyPriceIncrement = 0.5

var mockRoughTickerPrices = map[string]float64{
	"AAPL": 150,
	"MSFT": 300,
	"GOOG": 100,
	"AMZN": 100,
	"META": 200,
	"TSLA": 200,
	"NVDA": 300,
	"JPM":  130,
	"BABA": 80,
	"JNJ":  160,
	"WMT":  140,
	"PG":   130,
	"PYPL": 70,
	"DIS":  90,
	"ADBE": 340,
	"PFE":  45,
	"V":    200,
	"MA":   350,
	"CRM":  150,
	"NFLX": 280,
}

var mockUsers = []ljlib.User{
	{ID: uuid.MustParse("8a8d28aa-6c15-43be-8363-eb9862466063"), Username: "johndoe"},
	{ID: uuid.MustParse("f2f208c8-16a4-4ef6-80e3-88103f6471a2"), Username: "littlejohn"},
	{ID: uuid.MustParse("fa2fc7df-37ed-4582-8c46-01de352b375f"), Username: "jennifer"},
}

// LocalDatasource provides mocked data for users, their portfolio, and price history.
// More details on the approach are described in README file.
type LocalDatasource struct {
}

func NewLocalDatasource() LocalDatasource {
	return LocalDatasource{}
}

func (l LocalDatasource) GetUserByUsername(username string) (*ljlib.User, error) {
	for _, u := range mockUsers {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, ljlib.NewNotFoundError("cannot find user for username [%s]", username)
}

func (l LocalDatasource) GetHistoricalPrices(ticker string, dateFrom time.Time, dateTo time.Time) ([]ljlib.HistoricalPrice, error) {
	basePrice, ok := mockRoughTickerPrices[ticker]
	if !ok {
		return nil, ljlib.NewIllegalArgumentError("invalid ticker")
	}
	if dateFrom.After(dateTo) {
		return nil, ljlib.NewIllegalArgumentError("date from cannot be after date to")
	}
	var historicalPrices []ljlib.HistoricalPrice

	dtJan012003 := time.Date(2023, 01, 01, 00, 00, 00, 0, time.UTC)
	for dt := dateTo; !dt.Before(dateFrom); dt = dt.AddDate(0, 0, -1) {
		priceDiff := float64(int(dt.Sub(dtJan012003).Hours())/24) * mockDailyPriceIncrement
		newPrice := decimal.NewFromFloat(basePrice).Add(decimal.NewFromFloat(priceDiff))
		historicalPrices = append(historicalPrices, ljlib.HistoricalPrice{Date: dt, Price: newPrice})
	}
	return historicalPrices, nil
}

func (l LocalDatasource) GetUserPortfolio(userID uuid.UUID) ([]ljlib.TickerPrice, error) {
	return l.getUserTickers(userID)
}

func (l LocalDatasource) UserHasTicker(userID uuid.UUID, ticker string) (bool, error) {
	tickers, err := l.getUserTickers(userID)
	if err != nil {
		return false, err
	}
	for _, t := range tickers {
		if t.Ticker == ticker {
			return true, nil
		}
	}
	return false, nil
}

func (l LocalDatasource) getUserTickers(userID uuid.UUID) ([]ljlib.TickerPrice, error) {
	var user *ljlib.User
	for _, u := range mockUsers {
		if u.ID == userID {
			user = &u
			break
		}
	}
	if user == nil {
		return nil, ljlib.NewNotFoundError("user not found: %s", userID)
	}

	//generate user tickers
	var tickerNames []string
	for tickerName := range mockRoughTickerPrices {
		tickerNames = append(tickerNames, tickerName)
	}
	sort.Slice(tickerNames, func(i, j int) bool {
		return tickerNames[i] < tickerNames[j]
	})
	var userTickers []ljlib.TickerPrice
	var alreadyUsedTickers = make(map[string]bool)
	today := time.Now()
	for _, c := range user.Username {
		ticker := tickerNames[int(c)%len(tickerNames)]
		if _, ok := alreadyUsedTickers[ticker]; ok {
			continue
		}
		todayPrice, err := l.GetHistoricalPrices(ticker, today, today)
		if err != nil {
			return nil, fmt.Errorf("cannot get today's price for ticket [%s]: %w", ticker, err)
		}
		userTickers = append(userTickers, ljlib.TickerPrice{
			Ticker: ticker,
			Price:  todayPrice[0].Price,
		})
		alreadyUsedTickers[ticker] = true
	}

	return userTickers, nil
}
