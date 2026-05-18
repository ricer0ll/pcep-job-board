package workday

import "testing"

const StandardWorkdayRequestURL = "https://standard.wd1.myworkdayjobs.com/wday/cxs/standard/Search/jobs"
const ApexWorkdayRequestURL = "https://peak6group.wd1.myworkdayjobs.com/apexfintechsolutions"

var standardJobFamilyGroupIDs []string = []string{"4b3d59d7ab731002007ca34c13c90000"}
var apexJobFamilyGroupIDs []string = []string{"9319b7dfa5ee10212c5612fee7de0000"}

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
