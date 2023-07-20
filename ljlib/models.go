package ljlib

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type HistoricalPrice struct {
	Date  time.Time
	Price decimal.Decimal
}

func (h HistoricalPrice) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Date  string `json:"date"`
		Price string `json:"price"`
	}{
		Date:  h.Date.Format(time.DateOnly),
		Price: h.Price.StringFixed(2),
	})
}

type TickerPrice struct {
	Ticker string
	Price  decimal.Decimal
}

func (t TickerPrice) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Ticker string `json:"ticker"`
		Price  string `json:"price"`
	}{
		Ticker: t.Ticker,
		Price:  t.Price.StringFixed(2),
	})
}

type User struct {
	ID       uuid.UUID
	Username string
}
