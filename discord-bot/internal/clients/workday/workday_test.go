package workday

import (
	"testing"
)

func TestGetWorkdayJobPostings(t *testing.T) {
	path := "companies.json"

	companies, err := loadCompanies(path)
	if err != nil {
		t.Log("Unable to load companies")
		t.Fail()
	}

	t.Logf("Loaded %d companies from config", len(companies))

	for _, company := range companies {
		t.Logf("Processing company: %s", company.Name)

		jobs, err := getWorkdayJobPostings(
			company.WorkdayRequestURL,
			company.JobFamily,
			company.JobFamilyGroup,
			company.LocationCountry,
			company.Locations,
		)
		if err != nil {
			t.Errorf("Unable to get job postings for %s", company.Name)
			t.Fail()
			continue
		}

		t.Logf("Found %d jobs for %s", len(jobs), company.Name)

		for _, job := range jobs {
			t.Logf("%s: %s", company.Name, job.Title)
		}
	}
}
