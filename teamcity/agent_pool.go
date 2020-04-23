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
	Id   int    `json:"id,omitempty" xml:"id"`
	Name string `json:"name,omitempty" xml:"name"`
}

// AgentPool contains information about the Agent Pool
type AgentPool struct {
	Href      string `json:"href,omitempty" xml:"href"`
	Id        int    `json:"id,omitempty" xml:"id"`
	Name      string `json:"name,omitempty" xml:"name"`
	MaxAgents *int   `json:"maxAgents,omitempty" xml:"maxAgents"`
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

// Create will create an Agent Pool - which must have a unique name
func (s *AgentPoolsService) Create(pool AgentPool) (*AgentPool, error) {
	var created AgentPool

	err := s.restHelper.post("", pool, &created, "Agent Pool")
	if err != nil {
		return nil, err
	}

	return &created, nil
}

// Delete will delete an Agent Pool based on it's ID
func (s *AgentPoolsService) Delete(id int) error {
	locator := LocatorIDInt(id).String()
	err := s.restHelper.delete(locator, "Agent Pool")
	if err != nil {
		return err
	}

	return nil
}

// Get will return an Agent Pool based on it's ID
func (s *AgentPoolsService) Get(id int) (*AgentPool, error) {
	var out AgentPool
	locator := LocatorIDInt(id).String()
	err := s.restHelper.get(locator, &out, "Agent Pool")
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// List returns all of the available Agent Pools
func (s *AgentPoolsService) List() (*ListAgentPools, error) {
	var out ListAgentPools
	err := s.restHelper.get("", &out, "Agent Pools")
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// NOTE: Update support was investigated but is intentionally omitted - as the TC Documentation is incorrect
//		 PUT /app/rest/agentPools/id:4 at the time of writing returns a 405 Method Not Allowed
//		 POST /app/rest/agentPools with href & ID set also creates a new node pool
