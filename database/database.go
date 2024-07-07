package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"restapi/model"
)

type CarsDB struct {
	carList  []model.Car
	fileName string
}

func InitDB() (*CarsDB, error) {
	dbFile, err := os.Open("database.json")
	if err != nil {
		return nil, err
	}
	dbFile.Close()
	carList := make([]model.Car, 0, 10)
	return &CarsDB{
		carList:  carList,
		fileName: "database.json",
	}, nil
}

func (db *CarsDB) saveToFile() error {
	dbFile, err := os.Open("database.json")
	if err != nil {
		fmt.Println("Error with opening file:", err)
		return err
	}
	defer dbFile.Close()

	data, err := json.Marshal(db.carList)
	if err != nil {
		fmt.Println("Error with marshalling:", err)
		return err
	}

	_, err = dbFile.Write(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	return nil
}

func (db *CarsDB) Create(car model.Car) error {
	car.ID = len(db.carList)
	db.carList = append(db.carList, car)
	if err := db.saveToFile(); err != nil {
		return err
	}
	return nil
}

func (db *CarsDB) Get(id int) (model.Car, error) {

	if id < 0 || id > len(db.carList) {
		return model.Car{}, errors.New("index out of range")
	}
	car := db.carList[id]
	return car, nil
}

func (db *CarsDB) List() ([]model.Car, error) {
	copySlice := make([]model.Car, len(db.carList))
	copy(copySlice, db.carList)
	return copySlice, nil
}

func (db *CarsDB) Update(car model.Car) error {

	if car.ID < 0 || car.ID > len(db.carList) {
		return errors.New("index out of range")
	}

	db.carList[car.ID] = car
	if err := db.saveToFile(); err != nil {
		return err
	}
	return nil
}

func (db *CarsDB) Delete(id int) error {

	if id < 0 || id > len(db.carList) {
		return errors.New("index out of range")
	}

	newdb := make([]model.Car, len(db.carList)-1, cap(db.carList))
	tmp := 0
	for i, car := range db.carList {
		if i != id {
			newdb[tmp] = car
		}
		tmp++
	}

	if err := db.saveToFile(); err != nil {
		return err
	}
	return nil
}
