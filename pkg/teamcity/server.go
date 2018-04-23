package teamcity

// Server holds information about the TeamCity server
type Server struct {

	// agent pools
	AgentPools *Href `json:"agentPools,omitempty"`

	// agents
	Agents *Href `json:"agents,omitempty"`

	// build date
	BuildDate string `json:"buildDate,omitempty" xml:"buildDate"`

	// build number
	BuildNumber string `json:"buildNumber,omitempty" xml:"buildNumber"`

	// build queue
	BuildQueue *Href `json:"buildQueue,omitempty"`

	// builds
	Builds *Href `json:"builds,omitempty"`

	// current time
	CurrentTime string `json:"currentTime,omitempty" xml:"currentTime"`

	// internal Id
	InternalID string `json:"internalId,omitempty" xml:"internalId"`

	// investigations
	Investigations *Href `json:"investigations,omitempty"`

	// mutes
	Mutes *Href `json:"mutes,omitempty"`

	// projects
	Projects *Href `json:"projects,omitempty"`

	// role
	Role string `json:"role,omitempty" xml:"role"`

	// start time
	StartTime string `json:"startTime,omitempty" xml:"startTime"`

	// user groups
	UserGroups *Href `json:"userGroups,omitempty"`

	// users
	Users *Href `json:"users,omitempty"`

	// vcs roots
	VcsRoots *Href `json:"vcsRoots,omitempty"`

	// version
	Version string `json:"version,omitempty" xml:"version"`

	// version major
	VersionMajor int32 `json:"versionMajor,omitempty" xml:"versionMajor"`

	// version minor
	VersionMinor int32 `json:"versionMinor,omitempty" xml:"versionMinor"`

	// web Url
	WebURL string `json:"webUrl,omitempty" xml:"webUrl"`
}

// GetServer returns information about Server
func (c *Client) GetServer() (server *Server, err error) {

	var serverData Server
	err = c.doJSONRequest("GET", "server", nil, &serverData)
	if err != nil {
		return nil, err
	}

	return &serverData, nil
}
