package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"restapi/database"
	"restapi/model"
	"strconv"
)

type Service struct {
	Cars database.CarsDB
}

func NewService(Cars database.CarsDB) *Service {
	return &Service{Cars}
}

func response(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Println(err)
		}
	}
}

func responseError(w http.ResponseWriter, code int, err error) {
	response(w, code, map[string]string{"error :": err.Error()})
}

type CreateRequest struct {
	Make           string `json:"make"`
	Model          string `json:"model"`
	Mileage        int    `json:"mileage"`
	NumberOfOwners int    `json:"number_of_owners"`
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	req := new(CreateRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}
	r.Body.Close()

	car := model.Car{
		ID:             -1,
		Make:           req.Make,
		Model:          req.Model,
		Mileage:        req.Mileage,
		NumberOfOwners: req.NumberOfOwners,
	}

	if err := s.Cars.Create(car); err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusCreated, nil)
}

type GetResponse struct {
	ID             int    `json:"id"`
	Make           string `json:"make"`
	Model          string `json:"model"`
	Mileage        int    `json:"mileage"`
	NumberOfOwners int    `json:"number_of_owners"`
}

func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	car, err := s.Cars.Get(id)
	switch {
	case err == nil:
		response(w, http.StatusOK, car)
	case errors.Is(err, sql.ErrNoRows):
		responseError(w, http.StatusNotFound, err)
	default:
		responseError(w, http.StatusInternalServerError, err)
	}
}

type GetAllResponse struct {
	Result []GetResponse `json:"results"`
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request) {
	cars, err := s.Cars.List()
	if err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	result := make([]GetResponse, len(cars))
	for i, car := range cars {
		result[i] = GetResponse{
			ID:             car.ID,
			Make:           car.Make,
			Model:          car.Model,
			Mileage:        car.Mileage,
			NumberOfOwners: car.NumberOfOwners,
		}
	}
	response(w, http.StatusOK, GetAllResponse{
		Result: result,
	})
}

type UpdateRequest struct {
	Make           string `json:"make"`
	Model          string `json:"model"`
	Mileage        int    `json:"mileage"`
	NumberOfOwners int    `json:"number_of_owners"`
}

func (s *Service) Update(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	req := new(UpdateRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}
	r.Body.Close()

	car := model.Car{
		ID:             id,
		Make:           req.Make,
		Model:          req.Model,
		Mileage:        req.Mileage,
		NumberOfOwners: req.NumberOfOwners,
	}

	if err := s.Cars.Update(car); err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusOK, nil)
}

type UpdateSomethingRequest struct {
	Make           *string `json:"make,omitempty"`
	Model          *string `json:"model,omitempty"`
	Mileage        *int    `json:"mileage,omitempty"`
	NumberOfOwners *int    `json:"number_of_owners,omitempty"`
}

func (s *Service) UpdateSomething(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	req := new(UpdateSomethingRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}
	r.Body.Close()

	carPrev, err := s.Cars.Get(id)
	if err != nil {
		responseError(w, http.StatusNotFound, err)
		return
	}

	switch {
	case req.Make != nil:
		carPrev.Make = *req.Make
	case req.Model != nil:
		carPrev.Model = *req.Model
	case req.Mileage != nil:
		carPrev.Mileage = *req.Mileage
	case req.NumberOfOwners != nil:
		carPrev.NumberOfOwners = *req.NumberOfOwners
	}

	carAfter := model.Car{
		ID:             id,
		Make:           carPrev.Make,
		Model:          carPrev.Model,
		Mileage:        carPrev.Mileage,
		NumberOfOwners: carPrev.NumberOfOwners,
	}

	if err := s.Cars.Update(carAfter); err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusNoContent, nil)
}

type DeleteRequest struct {
	ID int `json:"id"`
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	if err := s.Cars.Delete(id); err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}
	response(w, http.StatusNoContent, nil)
}
