package main

import (
	"encoding/json"
	"net/http"

	"github.com/VictorOlea/go-rest-api/pkg/data"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
)

func main() {

	// create the store and sensor handler
	store := data.NewMemStore()
	sensorHandler := NewSensorhandler(store)
	home := homeHandler{}

	// create the router
	router := mux.NewRouter()

	// register the routers
	router.HandleFunc("/", home.ServeHTTP)
	router.HandleFunc("/sensors", sensorHandler.ListSensors).Methods("GET")
	router.HandleFunc("/sensors", sensorHandler.CreateSensor).Methods("POST")
	router.HandleFunc("/sensors/{id}", sensorHandler.GetSensor).Methods("GET")
	router.HandleFunc("/sensors/{id}", sensorHandler.UpdateSensor).Methods("PUT")
	router.HandleFunc("/sensors/{id}", sensorHandler.DeleteSensor).Methods("DELETE")

	// start the server
	http.ListenAndServe(":8000", router)

}

func InternalServerErrorhandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request)  {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home Page Sensor API"))
}

type sensorStore interface {
	Add(name string, sensor data.Sensor) error
	Get(name string) (data.Sensor, error)
	List() (map[string]data.Sensor, error)
	Update(name string, sensor data.Sensor) error
	Remove(name string) error
}

type SensorHandler struct {
	store sensorStore
}

func NewSensorhandler(s sensorStore) * SensorHandler {
	return &SensorHandler{
		store: s,
	}
}

func (h SensorHandler) CreateSensor(w http.ResponseWriter, r *http.Request) {

	var sensor data.Sensor

	if err := json.NewDecoder(r.Body).Decode(&sensor); err != nil {
		InternalServerErrorhandler(w, r)
		return
	}

	sensorID := slug.Make(sensor.Name)

	if  err := h.store.Add(sensorID, sensor); err != nil {
		InternalServerErrorhandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h SensorHandler) ListSensors(w http.ResponseWriter, r *http.Request) {

	sensors, err := h.store.List()
	if err != nil {
		InternalServerErrorhandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(sensors)
	if err != nil {
		InternalServerErrorhandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h SensorHandler) GetSensor(w http.ResponseWriter, r *http.Request) {
	
	id := mux.Vars(r)["id"]

	sensor, err := h.store.Get(id)
	if err != nil {
		if err == data.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}

		InternalServerErrorhandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(sensor)
	if err != nil {
		InternalServerErrorhandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h SensorHandler) UpdateSensor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var sensor data.Sensor
	if err := json.NewDecoder(r.Body).Decode(&sensor);err != nil {
		InternalServerErrorhandler(w, r)
		return
	}

	if err := h.store.Update(id, sensor); err != nil {
		if err == data.NotFoundErr {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorhandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(sensor)
	if err != nil {
		InternalServerErrorhandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h SensorHandler) DeleteSensor(w http.ResponseWriter, r *http.Request)  {
	id := mux.Vars(r)["id"]

	if err := h.store.Remove(id); err != nil {
		InternalServerErrorhandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}