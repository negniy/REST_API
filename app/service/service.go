package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"restapi/database"
	"strconv"
)

type Service struct {
	Cars database.CarsDB
}

func NewService(Cars database.CarsDB) *Service {
	return &Service{Cars}
}

func response(w http.ResponseWriter, code int, data any){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil{
		err:= json.NewEncoder(w).Encode(data)
		if err !=nil{
			log.Println(err)
		}
	}
}

func responseError(w http.ResponseWriter, code int, err error){
	response(w, code, map[string]string{"error :":err.Error()})
}

type CreateRequest struct {
	Make           string `json:"make"`
	Model          string `json:"model"`
	Mileage        int    `json:"mileage"`
	NumberOfOwners int    `json:"number_of_owners"`
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request){
	req := new(CreateRequest)
	if err:= json.NewDecoder(r.Body).Decode(req); err !=nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}
	r.Body.Close()

	if err := s.Cars.Create(); err!=nil{
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusCreated, nil)
}

type GetResponse struct{
	ID             int    `json:"id"`
	Make           string `json:"make"`
	Model          string `json:"model"`
	Mileage        int    `json:"mileage"`
	NumberOfOwners int    `json:"number_of_owners"`
}

func (s *Service) Get(w http.ResponseWriter, r *http.Request){
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil{
		responseError(w, http.StatusBadRequest, err)
		return
	}

	cars, err := s.Cars.Read(id)
	switch{
	case err == nil:
		response(w, http.StatusOK, cars)
	case errors.Is(err, sql.ErrNoRows):
		responseError(w, http.StatusNotFound, err)
	default:
		responseError(w, http.StatusInternalServerError, err)
	}
}

type GetAllResponse struct{
	Result []GetResponse `json:"results"`
}

func (s *Service) GetAll(w http.ResponseWriter, r *http.Request){
	cars, err := s.Cars.List()
	if err != nil{
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	result := make([]GetResponse, len(cars))
	for i, car :=range cars{
		result[i] = GetResponse{
			ID: car.ID,
			Make: car.Make,
			Model: car.Model,
			Mileage: car.Mileage,
			NumberOfOwners: car.NumberOfOwners,
		}
	}
	response(w, http.StatusOK, GetAllResponse{
		Result: result,
	})
}

type UpdateRequest struct{
	Make           string `json:"make"`
	Model          string `json:"model"`
	Mileage        int    `json:"mileage"`
	NumberOfOwners int    `json:"number_of_owners"`
}

func (s *Service) Update(w http.ResponseWriter, r *http.Request){
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil{
		responseError(w, http.StatusBadRequest, err)
		return
	}

	req := new(UpdateRequest)
	if err := 
}