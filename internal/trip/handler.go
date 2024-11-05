package trip

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type TripHandler struct {
	Service *TripService
}

func (h *TripHandler) GetTrips(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	trips, err := h.Service.GetTrips(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(trips)
}

func (h *TripHandler) GetTrip(w http.ResponseWriter, r *http.Request) {
	tripID := mux.Vars(r)["tripID"]

	tripData, err := h.Service.GetTrip(tripID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tripData)
}

func (h *TripHandler) CreateTrip(w http.ResponseWriter, r *http.Request) {
	var createTripData Trip
	if err := json.NewDecoder(r.Body).Decode(&createTripData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	tripData, err := h.Service.CreateTrip(createTripData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tripData)
}

func (h *TripHandler) EditTrip(w http.ResponseWriter, r *http.Request) {
	var editTripData Trip
	if err := json.NewDecoder(r.Body).Decode(&editTripData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	tripData, err := h.Service.EditTrip(editTripData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tripData)
}

func (h *TripHandler) DeleteTrip(w http.ResponseWriter, r *http.Request) {
	tripID := mux.Vars(r)["tripID"]

	response, err := h.Service.DeleteTrip(tripID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *TripHandler) GetItineraries(w http.ResponseWriter, r *http.Request) {
	tripID := mux.Vars(r)["tripID"]

	itineraries, err := h.Service.GetItineraries(tripID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(itineraries)
}

func (h *TripHandler) GetItinerary(w http.ResponseWriter, r *http.Request) {
	itineraryID := mux.Vars(r)["itineraryID"]

	itinerary, err := h.Service.GetItinerary(itineraryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(itinerary)
}

func (h *TripHandler) CreateItinerary(w http.ResponseWriter, r *http.Request) {
	var createItineraryData Itinerary
	if err := json.NewDecoder(r.Body).Decode(&createItineraryData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	itineraryData, err := h.Service.CreateItinerary(createItineraryData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(itineraryData)
}

func (h *TripHandler) EditItinerary(w http.ResponseWriter, r *http.Request) {
	var editItineraryData Itinerary
	if err := json.NewDecoder(r.Body).Decode(&editItineraryData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	itineraryData, err := h.Service.EditItinerary(editItineraryData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(itineraryData)
}

func (h *TripHandler) DeleteItinerary(w http.ResponseWriter, r *http.Request) {
	itineraryID := mux.Vars(r)["itineraryID"]

	response, err := h.Service.DeleteItinerary(itineraryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *TripHandler) GetItineraryItems(w http.ResponseWriter, r *http.Request) {
	itineraryID := mux.Vars(r)["itineraryID"]

	itineraryItems, err := h.Service.GetItineraryItems(itineraryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(itineraryItems)
}

func (h *TripHandler) GetItineraryItem(w http.ResponseWriter, r *http.Request) {
	itineraryItemID := mux.Vars(r)["itineraryItemID"]

	itineraryItem, err := h.Service.GetItineraryItem(itineraryItemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(itineraryItem)
}

func (h *TripHandler) CreateItineraryItem(w http.ResponseWriter, r *http.Request) {
	var createItineraryItemData ItineraryItem
	if err := json.NewDecoder(r.Body).Decode(&createItineraryItemData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	itineraryItemData, err := h.Service.CreateItineraryItem(createItineraryItemData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(itineraryItemData)
}

func (h *TripHandler) EditItineraryItem(w http.ResponseWriter, r *http.Request) {
	var editItineraryItemData ItineraryItem
	if err := json.NewDecoder(r.Body).Decode(&editItineraryItemData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	itineraryItemData, err := h.Service.EditItineraryItem(editItineraryItemData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(itineraryItemData)
}

func (h *TripHandler) DeleteItineraryItem(w http.ResponseWriter, r *http.Request) {
	itineraryItemID := mux.Vars(r)["itineraryItemID"]

	response, err := h.Service.DeleteItineraryItem(itineraryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h * TripHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	planID := mux.Vars(r)["planID"]

	plan, err := h.Service.GetPlan(planID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(plan)
}

func (h * TripHandler) GetPlanItems(w http.ResponseWriter, r *http.Request) {
	planID := mux.Vars(r)["planID"]

	planItems, err := h.Service.GetPlanItems(planID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(planItems)
}

func (h * TripHandler) GetPlanItem(w http.ResponseWriter, r *http.Request) {
	planItemID := mux.Vars(r)["planItemID"]

	planItem, err := h.Service.GetPlanItem(planItemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(planItem)
}

func (h *TripHandler) CreatePlanItem() (w http.ResponseWriter, r *http.Request) {
	var createPlanItemData PlanItem
	if err := json.NewDecoder(r.Body).Decode(&createPlanItemData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	planItemData, err := h.Service.CreatePlanItem(createPlanItemData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(planItemData)
}

func (h *TripHandler) EditPlanItem() (w http.ResponseWriter, r *http.Request) {
	var editPlanItemData PlanItem
	if err := json.NewDecoder(r.Body).Decode(&editPlanItemData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	planItemData, err := h.Service.EditPlanItem(editPlanItemData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(planItemData)
}

func (h *TripHandler) DeletePlanItem() (w http.ResponseWriter, r *http.Request) {
	planItemID := mux.Vars(r)["planItemID"]
	
	response, err := h.Service.DeletePlanItem(planItemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}


