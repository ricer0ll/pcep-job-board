package workday

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/ricer0ll/pcep-job-board/discord-bot/api/workday/dto"
)

const StandardWorkdayBaseURL = "https://standard.wd1.myworkdayjobs.com"
const ApexWorkdayBaseURL = "https://https://peak6group.wd1.myworkdayjobs.com"

const StandardWorkdayRequestURL = "https://standard.wd1.myworkdayjobs.com/wday/cxs/standard/Search/jobs"
const ApexWorkdayRequestURL = "https://peak6group.wd1.myworkdayjobs.com/wday/cxs/peak6group/apexfintechsolutions/jobs"

var standardJobFamilyGroupIDs []string = []string{"4b3d59d7ab731002007ca34c13c90000"}
var apexJobFamilyGroupIDs []string = []string{"9319b7dfa5ee10212c5612fee7de0000"}

var jobsCache map[string][]dto.JobPosting = make(map[string][]dto.JobPosting)

// list of channels to notify
var channelIDs []string = []string{
	"1505733135357313115", // dev server
}

func InitJobsCache() {
	standardJobs, err := getWorkdayJobPostings(StandardWorkdayRequestURL, standardJobFamilyGroupIDs)
	if err != nil {
		slog.Error("Unable to get job postings from The Standard")
	}

	apexJobs, err := getWorkdayJobPostings(ApexWorkdayRequestURL, apexJobFamilyGroupIDs)
	if err != nil {
		slog.Error("Unable to get job postings from The Apex")
	}

	jobsCache["Standard"] = append(jobsCache["Standard"], standardJobs...)
	jobsCache["Apex"] = append(jobsCache["Apex"], apexJobs...)
}

func GetNewJobPostings(client *bot.Client) {
	standardJobs, err := getWorkdayJobPostings(StandardWorkdayRequestURL, standardJobFamilyGroupIDs)
	if err != nil {
		slog.Error("Unable to get job postings from The Standard")
	}

	apexJobs, err := getWorkdayJobPostings(ApexWorkdayRequestURL, apexJobFamilyGroupIDs)
	if err != nil {
		slog.Error("Unable to get job postings from The Standard")
	}

	if len(standardJobs) != len(jobsCache["Standard"]) {
		jobPosting := standardJobs[0]
		notifyNewJob(client, &jobPosting, "The Standard", StandardWorkdayBaseURL)
		jobsCache["Standard"] = standardJobs
	}

	if len(apexJobs) != len(jobsCache["Apex"]) {
		jobPosting := standardJobs[0]
		notifyNewJob(client, &jobPosting, "Apex Fintech Solutions", ApexWorkdayBaseURL)
		jobsCache["Apex"] = standardJobs
	}
}

func notifyNewJob(client *bot.Client, jobPosting *dto.JobPosting, company string, workdayURL string) {
	embed := generateNewJobPostingEmbed(jobPosting, company, workdayURL)
	for _, channelID := range channelIDs {
		client.Rest.CreateMessage(
			snowflake.MustParse(channelID),
			discord.NewMessageCreate().WithEmbeds(embed),
		)
	}
}

func getWorkdayJobPostings(url string, jobFamilyGroup []string) ([]dto.JobPosting, error) {
	jobPostings := []dto.JobPosting{}

	request := dto.JobPostingRequest{
		AppliedFacets: dto.AppliedFacet{
			JobFamilyGroup: jobFamilyGroup,
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return jobPostings, err
	}

	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return jobPostings, err
	}
	if resp.StatusCode != http.StatusOK {
		return jobPostings, err
	}
	defer resp.Body.Close()

	responseObj, err := mapJobPostingResponse(resp)
	if err != nil {
		return jobPostings, err
	}

	jobPostings = responseObj.JobPostings
	return jobPostings, nil
}

func mapJobPostingResponse(resp *http.Response) (*dto.JobPostingResponse, error) {
	var responseObj dto.JobPostingResponse

	err := json.NewDecoder(resp.Body).Decode(&responseObj)
	if err != nil {
		return nil, err
	}

	return &responseObj, nil
}

func generateNewJobPostingEmbed(jobPosting *dto.JobPosting, company string, workdayURL string) discord.Embed {
	var title string = fmt.Sprintf("New Job Posting from %s!", company)
	var description string = fmt.Sprintf("Position: **%s**\nLocation: %s", jobPosting.Title, jobPosting.LocationsText)
	var url string = workdayURL + fmt.Sprintf("/en-US/Search%s", jobPosting.ExternalPath)

	embed := discord.NewEmbed().
		WithTitle(title).
		WithDescription(description).
		WithURL(url)

	return embed
}
