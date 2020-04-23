package teamcity

import (
	"net/http"

	"github.com/dghubble/sling"
)

// ListAgentPools is the response object when listing Agent Pools
type ListAgentPools struct {
	Count      int                  `json:"count,omitempty" xml:"count"`
	Href       string               `json:"href,omitempty" xml:"href"`
	AgentPools []AgentPoolReference `json:"agentPool,omitempty" xml:"agentPool"`
}

// AgentPoolReference is a reference to an Agent Pool
type AgentPoolReference struct {
	Href string `json:"href,omitempty" xml:"href"`
	Id   string `json:"id,omitempty" xml:"id"`
	Name string `json:"name,omitempty" xml:"name"`
}

// AgentPoolsService has operations for handling agent pools
type AgentPoolsService struct {
	sling      *sling.Sling
	httpClient *http.Client
	restHelper *restHelper
}

func newAgentPoolsService(base *sling.Sling, client *http.Client) *AgentPoolsService {
	sling := base.Path("agentPools/")
	return &AgentPoolsService{
		sling:      sling,
		httpClient: client,
		restHelper: newRestHelperWithSling(client, sling),
	}
}

func (s *AgentPoolsService) List() (*ListAgentPools, error) {
	var out ListAgentPools
	err := s.restHelper.get("", &out, "Agent Pools")
	if err != nil {
		return nil, err
	}

	return &out, nil
}
