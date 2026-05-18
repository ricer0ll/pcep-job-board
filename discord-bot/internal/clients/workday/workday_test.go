package workday

import "testing"

func TestGetWorkdayJobPostings(t *testing.T) {
	_, err := getWorkdayJobPostings(StandardWorkdayRequestURL, standardJobFamilyGroupIDs)
	if err != nil {
		t.Errorf("Unable to get jobs from the Standard's Workday: %v", err)
	}

	_, err = getWorkdayJobPostings(ApexWorkdayRequestURL, apexJobFamilyGroupIDs)
	if err != nil {
		t.Errorf("Unable to get jobs from the Apex's Workday: %v", err)
	}
}
