package littlejohn

import (
	"net/http"

	"github.com/iliyaisd/littlejohn/internal/api"
	"github.com/iliyaisd/littlejohn/internal/datasource"
)

// Data source constants are used as just a demo of how actual different data source can be used in real project,
// and the right one gets instantiated depending on the env var. 
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
