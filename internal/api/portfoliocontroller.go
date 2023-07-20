package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/iliyaisd/littlejohn/ljlib"
)

const (
	pageSizeDays = 90
	maxDaysBack  = 10 * 365 //roughly 10 years for simplicity
	maxPage      = maxDaysBack / pageSizeDays
)

type PortfolioController struct {
	priceDataSource DataSource
}

type DataSource interface {
	GetHistoricalPrices(ticker string, dateFrom time.Time, dateTo time.Time) ([]ljlib.HistoricalPrice, error)
	GetUserPortfolio(userID uuid.UUID) ([]ljlib.TickerPrice, error)
	UserHasTicker(userID uuid.UUID, ticker string) (bool, error)
}

func NewPortfolioController(ds DataSource) PortfolioController {
	return PortfolioController{
		priceDataSource: ds,
	}
}

func (c PortfolioController) GetTickers(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*ljlib.User)
	if !ok {
		ljlib.ResponseHTTPForbidden(w, "Forbidden")
		return
	}

	portfolio, err := c.priceDataSource.GetUserPortfolio(user.ID)
	if err != nil {
		log.Printf("cannot fetch portfolio for user [%s]: %s", user.Username, err)
		if errors.Is(err, ljlib.NotFoundError{}) {
			ljlib.ResponseHTTPNotFound(w, "Forbidden")
			return
		}
		ljlib.ResponseHTTPError(w, "Cannot get tickers")
		return
	}

	ljlib.ResponseHTTP(w, http.StatusOK, portfolio)
}

func (c PortfolioController) GetTickerHistory(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*ljlib.User)
	if !ok {
		ljlib.ResponseHTTPForbidden(w, "Forbidden")
		return
	}

	ticker := mux.Vars(r)["ticker"]
	hasTicker, err := c.priceDataSource.UserHasTicker(user.ID, ticker)
	if err != nil {
		log.Printf("cannot determine whether user [%s] has ticker: %s", user.Username, err)
		ljlib.ResponseHTTPError(w, "Cannot get ticker history")
		return
	}
	if !hasTicker {
		ljlib.ResponseHTTPNotFound(w, "Ticker not found")
		return
	}

	page := c.extractPage(r)
	dateFrom, dateTo := c.buildDatesForPage(page)

	prices, err := c.priceDataSource.GetHistoricalPrices(ticker, dateFrom, dateTo)
	if err != nil {
		log.Printf("cannot get historical prices: %s", err)
		ljlib.ResponseHTTPError(w, "Cannot get historical prices")
		return
	}

	ljlib.ResponseHTTP(w, http.StatusOK, prices)
}

func (c PortfolioController) extractPage(r *http.Request) int {
	pageStr := r.URL.Query().Get("page")
	if len(pageStr) == 0 {
		return 0
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Printf("Warning: non-int page value: %s", pageStr)
		return 0
	}
	//in request, pages are numbered from one:
	page--
	if page < 0 {
		return 0
	}
	if page > maxPage {
		return maxPage
	}
	return page
}

func (c PortfolioController) buildDatesForPage(page int) (time.Time, time.Time) {
	today := time.Now()
	dateTo := today.Add(time.Duration(-page*pageSizeDays*24) * time.Hour)
	dateFrom := dateTo.Add(-(pageSizeDays - 1) * 24 * time.Hour)
	maxBefore := today.Add(-maxDaysBack * 24 * time.Hour)
	if dateFrom.Before(maxBefore) {
		dateFrom = maxBefore
	}
	return dateFrom, dateTo
}
