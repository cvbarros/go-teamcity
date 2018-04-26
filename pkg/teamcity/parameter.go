package teamcity

import (
	"fmt"

	"github.com/dghubble/sling"
)

// ParameterService has operations for handling parameters for projects or build configurations
type ParameterService struct {
	base *sling.Sling
}

func newParameterService(base *sling.Sling) *ParameterService {
	return &ParameterService{
		base: base,
	}
}

//Project returns a parameter service that operated parameters for the project with given id
func (s *ParameterService) Project(id string) *ParameterService {
	return &ParameterService{
		base: s.base.New().Path(fmt.Sprintf("projects/%s/", LocatorProject(id))),
	}
}

// Add requires a remote call for each parameter being added, since TeamCity only support creating them via
// POST to /project/<project_locator>/parameterName, and not batch operations.
// This function is created just for convenience and batch creation of parameters
// Parameters will be created in the order they are passed and there is no guarantee that it will be an atomic operation
// On the first failure it will stop and not create any further parameters
func (s *ParameterService) Add(id string, parameters ...*Property) error {
	for _, param := range parameters {
		var out *Property
		resp, err := s.base.New().Post("parameters").BodyJSON(param).ReceiveSuccess(&out)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("Error creating parameter: %s, statusCode: %d", param.Name, resp.StatusCode)
		}
	}
	return nil
}
