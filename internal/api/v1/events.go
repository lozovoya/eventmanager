package v1

import (
	"EventManager/internal/cache"
	"EventManager/internal/model"
	"EventManager/internal/views"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Event struct {
	callEvent cache.Call
}

func NewEventsController (callEvent cache.Call) *Event {
	return &Event{callEvent: callEvent}
}

func (e *Event) InEvent (writer http.ResponseWriter, request *http.Request) {
	var event *CallEventDTO
	err := json.NewDecoder(request.Body).Decode(&event)
	if err != nil {
		log.Println(fmt.Errorf("InEvent: %w", err))
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if IsEmpty(event.EventID) || IsEmpty(event.Type) {
		log.Printf("InEvent: mandatory field is empty. %s", event.EventID)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var currentCall model.Call
	currentCall.CallID = event.Data.CallID
	currentCall.CallingNumber = event.Data.CallingNumber
	currentCall.CallingLevel = event.Data.CallingLevel
	currentCall.DialedNumber = event.Data.DialedNumber
	currentCall.Queue_ID = event.Data.Queue_ID

	switch event.Type {
	case "OnQueueInEvent": {
		err = e.callEvent.CallToCache(request.Context(),  &currentCall)
		if err != nil {
			log.Println(fmt.Errorf("InEvent: %w", err))
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	case "OnQueueOutEvent": {
		err = e.callEvent.CallFromCache(request.Context(),  currentCall.Queue_ID, currentCall.CallID)
		if err != nil {
			log.Println(fmt.Errorf("InEvent: %w", err))
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	default:
		log.Printf("InEvent: wrong type. %s", event.EventID)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusOK)
	return
}

func (e *Event) GetSnapShot(writer http.ResponseWriter, request *http.Request) {
	calls, err := e.callEvent.GetCallsSnapshot(request.Context())
	if err != nil {
		log.Println(fmt.Errorf("GetSnapShot: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	snapshot, err := views.SnapShot(request.Context(), calls)
	if err != nil {
		log.Println(fmt.Errorf("GetSnapShot: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(&snapshot)
	if err != nil {
		log.Println(fmt.Errorf("GetSnapShot: %w", err))
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}