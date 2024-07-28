package csv

// RepoData represents a single record in the CSV file
type RepoData struct {
	RepoURL                  string
	RepoLanguage             string
	RepoLicense              string
	RepoCreatedAt            string
	RepoUpdatedAt            string
	CollectionDate           string
	WorkerCommitID           string
	RepoStarCount            int
	LegacyCreatedSince       int
	LegacyUpdatedSince       int
	LegacyContributorCount   int
	LegacyOrgCount           int
	LegacyCommitFrequency    float64
	LegacyRecentReleaseCount int
	LegacyUpdatedIssuesCount int
	LegacyClosedIssuesCount  int
	LegacyIssueCommentFreq   float64
	LegacyGithubMentionCount int
	DepsDevDependentCount    int
	DefaultScore             float64
}
