package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tabichanorg/tabichan-server/internal/db"
	"github.com/tabichanorg/tabichan-server/internal/healthcheck"
	middleware "github.com/tabichanorg/tabichan-server/internal/middleware/session"
	"github.com/tabichanorg/tabichan-server/internal/trip"
	"github.com/tabichanorg/tabichan-server/internal/user"
)

func SetupRoutes(mux *mux.Router) *mux.Router {
	mux.HandleFunc("/healthcheck", healthcheck.HealthCheck).Methods("GET")

	userHandler := initUserHandler()
	initRoute(mux, "/signup", userHandler.Signup, false, "POST")
	initRoute(mux, "/login", userHandler.Login, false, "POST")
	initRoute(mux, "/user/details", userHandler.GetUser, true, "GET")

	initTripRoutes(mux)

	return mux
}

func initTripRoutes(mux *mux.Router) {
	tripHandler := initTripHandler()
	initRoute(mux, "/trips", tripHandler.GetTrips, true, "GET")
	initRoute(mux, "/trips/{tripID}", tripHandler.GetTrip, true, "GET")
	initRoute(mux, "/trips", tripHandler.CreateTrip, true, "POST")
	// initRoute(mux, "/trips/{tripID}", tripHandler.EditTrip, true, "PUT")
	initRoute(mux, "/trips/{tripID}", tripHandler.DeleteTrip, true, "DELETE")

	initRoute(mux, "/itineraries/{planID}", tripHandler.GetItineraries, true, "GET")
	initRoute(mux, "/itineraries", tripHandler.CreateItinerary, true, "POST")
	initRoute(mux, "/itineraries/{itineraryID}", tripHandler.GetItinerary, true, "GET")
	// initRoute(mux, "/itineraries/{itineraryID}", tripHandler.EditItinerary, true, "PUT")
	initRoute(mux, "/itineraries/{itineraryID}", tripHandler.DeleteItinerary, true, "DELETE")

	// initRoute(mux, "/itineraries/{itineraryID}/items", tripHandler.GetItineraryItems, true, "GET")
	// initRoute(mux, "/itineraries/{itineraryID}/{itineraryItemID}", tripHandler.GetItineraryItem, true, "GET")
	// initRoute(mux, "/itineraries/{itineraryID}", tripHandler.CreateItineraryItem, true, "POST")
	// initRoute(mux, "/itineraries/{itineraryID}/{itineraryItemID}", tripHandler.EditItineraryItem, true, "PUT")
	// initRoute(mux, "/itineraries/{itineraryID}/{itineraryItemID}", tripHandler.DeleteItineraryItem, true, "DELETE")

	// initRoute(mux, "/plans/{planID}", tripHandler.GetPlan, true, "GET") // get itinerary (+items), get planitems
	// initRoute(mux, "/plans/{planID}/items", tripHandler.GetPlanItems, true, "GET")
	// initRoute(mux, "/plans/{planID}/{planItemID}", tripHandler.GetPlanItem, true, "GET")
	// initRoute(mux, "/plans/{planID}", tripHandler.CreatePlanItem, true, "POST")
	// initRoute(mux, "/plans/{planID}/{planItemID}", tripHandler.EditPlanItem, true, "PUT")
	// initRoute(mux, "/plans/{planID}/{planItemID}", tripHandler.DeletePlanItem, true, "DELETE")
}

func initUserHandler() *user.UserHandler {
	userRepo := &user.UserRepository{Client: db.DynamoClient}
	userService := &user.UserService{Repo: userRepo}
	return &user.UserHandler{Service: userService}
}

func initTripHandler() *trip.TripHandler {
	tripRepo := &trip.TripRepository{Client: db.DynamoClient}
	tripService := &trip.TripService{Repo: tripRepo}
	return &trip.TripHandler{Service: tripService}
}

func initMiddleware() *middleware.MiddlewareService {
	middlewareRepo := &middleware.MiddlewareRepository{Client: db.DynamoClient}
	return &middleware.MiddlewareService{Repo: middlewareRepo}
}

func initRoute(mux *mux.Router, endpoint string, handlerFunc func(http.ResponseWriter, *http.Request), isSecureRoute bool, method string) {
	middlewareService := initMiddleware()
	if isSecureRoute {
		mux.HandleFunc(endpoint, middlewareService.SessionMiddleware(handlerFunc)).Methods(method)
		return
	}
	mux.HandleFunc(endpoint, handlerFunc).Methods(method)
}
