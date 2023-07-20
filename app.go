package littlejohn

import (
	"net/http"

	"github.com/iliyaisd/littlejohn/internal/api"
	"github.com/iliyaisd/littlejohn/internal/datasource"
)

const (
	DataSourceLocal = "local"
	DataSourceYahoo = "yahoo"
)

type Config struct {
	Port       int
	DataSource string
}

type App struct {
	MainHandler http.Handler
}

func BuildApp() (App, error) {
	dataSource := datasource.NewLocalDatasource()

	authorizer := api.NewAPIKeyAuthorizer(dataSource)

	payrollController := api.NewPortfolioController(dataSource)
	router := NewRouter(Controllers{
		portfolioController: payrollController,
	}, authorizer)

	return App{
		MainHandler: router.PrepareHandler(),
	}, nil
}
