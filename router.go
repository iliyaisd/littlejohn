package littlejohn

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/iliyaisd/littlejohn/internal/api"
	"github.com/iliyaisd/littlejohn/ljlib"
)

//Router puts together all the API endpoints and wires them to handlers (controllers) and middlewares.
type Router struct {
	controllers Controllers
	authorizer  Authorizer
}

func NewRouter(controllers Controllers, authorizer Authorizer) Router {
	return Router{
		controllers: controllers,
		authorizer:  authorizer,
	}
}

func (r Router) PrepareHandler() http.Handler {
	router := mux.NewRouter()
	router.Use(r.jsonMiddleware)

	restrictedRoutes := router.PathPrefix("").Subrouter()
	restrictedRoutes.Use(r.authorizeRequest)

	r.controllers.HandleRestrictedRoutes(restrictedRoutes)

	routerCORS := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "Content-Type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
	)(router)

	return routerCORS
}

type Controllers struct {
	portfolioController api.PortfolioController
}

func (c Controllers) HandleRestrictedRoutes(router *mux.Router) {
	router.HandleFunc("/tickers", c.portfolioController.GetTickers).Methods("GET")
	router.HandleFunc("/tickers/{ticker}/history", c.portfolioController.GetTickerHistory).Methods("GET")
}

type Authorizer interface {
	Authorize(r *http.Request) (*http.Request, error)
}

func (r Router) authorizeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var reqWithUser = req
		var err error
		reqWithUser, err = r.authorizer.Authorize(reqWithUser)
		if err != nil {
			log.Printf("Unauthorized access by URI [%s]: %s", req.RequestURI, err)
			ljlib.ResponseHTTPForbidden(w, "forbidden")
			return
		}
		next.ServeHTTP(w, reqWithUser)
	})
}

func (r Router) jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
