package teamcity

import (
	"encoding/json"
	"errors"
	"fmt"
)

//StepCommandLine represents a a build step of type "CommandLine"
type StepCommandLine struct {
	ID           string
	Name         string
	stepType     string
	stepJSON     *stepJSON
	isExecutable bool
	//CustomScript contains code for platform specific script, like .cmd on windows or shell script on Unix-like environments.
	CustomScript string
	//CommandExecutable is the executable program to be called from this step.
	CommandExecutable string
	//CommandParameters are additional parameters to be passed on to the CommandExecutable.
	CommandParameters string
	//ExecuteMode is the execute mode for the step. See StepExecuteMode for details.
	ExecuteMode StepExecuteMode
	//Conditions
	Conditions string
}

//NewStepCommandLineScript creates a command line build step that runs an inline platform-specific script.
func NewStepCommandLineScript(name string, script string, mode string, conditions string) (*StepCommandLine, error) {
	if script == "" {
		return nil, errors.New("script is required")
	}
	if mode == "" {
		mode = StepExecuteModeDefault
	}

	return &StepCommandLine{
		Name:         name,
		isExecutable: false,
		stepType:     StepTypeCommandLine,
		CustomScript: script,
		ExecuteMode:  mode,
		Conditions:   conditions,
	}, nil
}

//NewStepCommandLineExecutable creates a command line that invokes an external executable.
func NewStepCommandLineExecutable(name string, executable string, args string, mode string, conditions string) (*StepCommandLine, error) {
	if executable == "" {
		return nil, errors.New("executable is required")
	}
	if mode == "" {
		mode = StepExecuteModeDefault
	}

	return &StepCommandLine{
		Name:              name,
		stepType:          StepTypeCommandLine,
		isExecutable:      true,
		CommandExecutable: executable,
		CommandParameters: args,
		ExecuteMode:       mode,
		Conditions:        conditions,
	}, nil
}

//GetID is a wrapper implementation for ID field, to comply with Step interface
func (s *StepCommandLine) GetID() string {
	return s.ID
}

//GetName is a wrapper implementation for Name field, to comply with Step interface
func (s *StepCommandLine) GetName() string {
	return s.Name
}

//Type returns the step type, in this case "StepTypeCommandLine".
func (s *StepCommandLine) Type() BuildStepType {
	return StepTypeCommandLine
}

func (s *StepCommandLine) properties() *Properties {
	props := NewPropertiesEmpty()
	props.AddOrReplaceValue("teamcity.step.mode", string(s.ExecuteMode))
	props.AddOrReplaceValue("teamcity.step.conditions", string(s.Conditions))

	if s.isExecutable {
		props.AddOrReplaceValue("command.executable", s.CommandExecutable)

		if s.CommandParameters != "" {
			props.AddOrReplaceValue("command.parameters", s.CommandParameters)
		}
	} else {
		props.AddOrReplaceValue("script.content", s.CustomScript)
		props.AddOrReplaceValue("use.custom.script", "true")
	}

	return props
}

func (s *StepCommandLine) serializable() *stepJSON {
	return &stepJSON{
		ID:         s.ID,
		Name:       s.Name,
		Type:       s.stepType,
		Properties: s.properties(),
	}
}

//MarshalJSON implements JSON serialization for StepCommandLine
func (s *StepCommandLine) MarshalJSON() ([]byte, error) {
	out := s.serializable()
	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for StepCommandLine
func (s *StepCommandLine) UnmarshalJSON(data []byte) error {
	var aux stepJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Type != string(StepTypeCommandLine) {
		return fmt.Errorf("invalid type %s trying to deserialize into StepCommandLine entity", aux.Type)
	}
	s.Name = aux.Name
	s.ID = aux.ID
	s.stepType = StepTypeCommandLine

	props := aux.Properties
	if _, ok := props.GetOk("use.custom.script"); ok {
		s.isExecutable = false
		if v, ok := props.GetOk("script.content"); ok {
			s.CustomScript = v
		}
	}

	if v, ok := props.GetOk("command.executable"); ok {
		s.CommandExecutable = v
		if v, ok := props.GetOk("command.parameters"); ok {
			s.CommandParameters = v
		}
	}

	if v, ok := props.GetOk("teamcity.step.mode"); ok {
		s.ExecuteMode = StepExecuteMode(v)
	}
	if v, ok := props.GetOk("teamcity.step.conditions"); ok {
		s.Conditions = v
	}
	return nil
}
