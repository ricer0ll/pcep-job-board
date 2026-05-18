package workday

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Company struct {
	Name              string
	WorkdayBaseURL    string
	WorkdayRequestURL string
	JobFamilyGroupIDs []string
}

type Config struct {
	Companies []Company `json:"companies"`
}

func loadCompanies() ([]Company, error) {
	fileData, err := os.ReadFile(filepath.Join("internal", "clients", "workday", "companies.json"))
	if err != nil {
		return []Company{}, nil
	}

	var config Config
	err = json.Unmarshal(fileData, &config)
	if err != nil {
		return []Company{}, nil
	}

	return config.Companies, nil
}
