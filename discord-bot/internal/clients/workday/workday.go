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

var jobsCache map[string][]dto.JobPosting = make(map[string][]dto.JobPosting)

// list of channels to notify
var channelIDs []string = []string{
	"1505733135357313115", // dev server
}

func InitJobsCache() {
	companies, err := loadCompanies()
	if err != nil {
		panic("Unable to load companies")
	}

	for _, company := range companies {
		jobs, err := getWorkdayJobPostings(
			company.WorkdayRequestURL,
			company.JobFamilyGroupIDs,
			company.Locations,
		)
		if err != nil {
			slog.Error(fmt.Sprintf("Unable to get job postings from %s", company.Name))
			continue
		}

		jobsCache[company.Name] = append(jobsCache[company.Name], jobs...)
	}
}

func GetNewJobPostings(client *bot.Client) {
	companies, err := loadCompanies()
	if err != nil {
		panic("Unable to load companies")
	}

	for _, company := range companies {
		// get workday job postings
		liveJobs, err := getWorkdayJobPostings(
			company.WorkdayRequestURL,
			company.JobFamilyGroupIDs,
			company.Locations,
		)
		if err != nil {
			slog.Error(fmt.Sprintf("Unable to get job postings from %s", company.Name))
			continue
		}

		// add to id cache to check later (basically a set)
		cachedIDs := make(map[string]struct{}) // Bullet Fields = ID (sorta)
		for _, job := range jobsCache[company.Name] {
			cachedIDs[job.BulletFields[0]] = struct{}{}
		}

		for _, job := range liveJobs {
			_, ok := cachedIDs[job.BulletFields[0]]
			if !ok {
				notifyNewJob(client, &job, company.Name, company.WorkdayBaseURL) // notify on discord if new job
			}
		}

		jobsCache[company.Name] = liveJobs
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

func getWorkdayJobPostings(url string, jobFamilyGroup []string, locations []string) ([]dto.JobPosting, error) {
	jobPostings := []dto.JobPosting{}

	request := dto.JobPostingRequest{
		AppliedFacets: dto.AppliedFacet{
			JobFamilyGroup: jobFamilyGroup,
			Locations:      locations,
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
	var url string = workdayURL + fmt.Sprintf("%s", jobPosting.ExternalPath)

	embed := discord.NewEmbed().
		WithTitle(title).
		WithDescription(description).
		WithURL(url)

	return embed
}
