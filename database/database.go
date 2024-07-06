package database

import "os"

type CarsDB struct {
	FileName string
}

func InitDB() (*CarsDB, error) {
	dbFile, err := os.Open("database.json")
	if err != nil {
		return nil, err
	}
	defer dbFile.Close()
	return &CarsDB{"database.json"}, nil
}
