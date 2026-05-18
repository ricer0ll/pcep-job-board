package dto

type AppliedFacet struct {
	JobFamilyGroup []string `json:"jobFamilyGroup"`
}

type JobPosting struct {
	Title         string   `json:"title"`
	ExternalPath  string   `json:"externalPath"`
	LocationsText string   `json:"locationsText"`
	PostedOn      string   `json:"postedOn"`
	BulletFields  []string `json:"bulletFields"`
}

type JobPostingRequest struct {
	AppliedFacets AppliedFacet `json:"appliedFacets"`
}

type JobPostingResponse struct {
	Total       uint64 `json:"total"`
	JobPostings []JobPosting
}
