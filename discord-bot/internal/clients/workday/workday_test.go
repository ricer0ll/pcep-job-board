package workday

import (
	"testing"
)

func TestGetWorkdayJobPostings(t *testing.T) {
	companies, err := loadCompanies()
	if err != nil {
		t.Errorf("Unable to load companies")
	}

	for _, company := range companies {
		_, err := getWorkdayJobPostings(
			company.WorkdayRequestURL,
			company.JobFamilyGroupIDs,
			company.Locations,
		)
		if err != nil {
			t.Errorf("Unable to get job postings for %s", company.Name)
			continue
		}
	}
}
