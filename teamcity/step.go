package teamcity

import (
	"encoding/json"
	"fmt"
)

// BuildStepType represents most common step types for build steps
type BuildStepType = string

const (
	//StepTypePowershell step type
	StepTypePowershell BuildStepType = "jetbrains_powershell"
	//StepTypeDotnetCli step type
	StepTypeDotnetCli BuildStepType = "dotnet.cli"
	//StepTypeCommandLine (shell/cmd) step type
	StepTypeCommandLine          BuildStepType = "simpleRunner"
	StepTypeOctopusPushPackage   BuildStepType = "octopus.push.package"
	StepTypeOctopusCreateRelease BuildStepType = "octopus.create.release"
)

//StepExecuteMode represents how a build configuration step will execute regarding others.
type StepExecuteMode = string

const (
	//StepExecuteModeDefault executes the step only if all previous steps finished successfully.
	StepExecuteModeDefault = "default"
	//StepExecuteModeOnlyIfBuildIsSuccessful executes the step only if the whole build is successful.
	StepExecuteModeOnlyIfBuildIsSuccessful = "execute_if_success"
	//StepExecuteModeEvenWhenFailed executes the step even if previous steps failed.
	StepExecuteModeEvenWhenFailed = "execute_if_failed"
	//StepExecuteAlways executes even if build stop command was issued.
	StepExecuteAlways = "execute_always"
)

// Step interface represents a a build configuration/template build step. To interact with concrete step types, see the Step* types.
type Step interface {
	GetID() string
	GetName() string
	Type() string

	serializable() *stepJSON
}

type stepJSON struct {
	Disabled   *bool       `json:"disabled,omitempty" xml:"disabled"`
	Href       string      `json:"href,omitempty" xml:"href"`
	ID         string      `json:"id,omitempty" xml:"id"`
	Inherited  *bool       `json:"inherited,omitempty" xml:"inherited"`
	Name       string      `json:"name,omitempty" xml:"name"`
	Properties *Properties `json:"properties,omitempty"`
	Type       string      `json:"type,omitempty" xml:"type"`
}

type stepsJSON struct {
	Count int32       `json:"count,omitempty" xml:"count"`
	Items []*stepJSON `json:"step"`
}

//OperatorStrings - All possible condition strings. Do not change the values.
var OperatorStrings = []string{
	"EXISTS",
	"NOT_EXISTS",
	"EQUALS",
	"DOES_NOT_EQUAL",
	"MORE_THAN",
	"NO_MORE_THAN",
	"LESS_THAN",
	"NO_LESS_THAN",
	"STARTS_WITH",
	"CONTAINS",
	"DOES_NOT_CONTAIN",
	"ENDS_WITH",
	"MATCHES",
	"DOES_NOT_MATCH",
	"VER_MORE_THAN",
	"VER_NO_MORE_THAN",
	"VER_LESS_THAN",
	"VER_NO_LESS_THAN",
}

//Operators - Possible conditions for build step. Do not change the values.
var Operators = struct {
	Exists            string
	NotExists         string
	Equals            string
	DoesNotEqual      string
	MoreThan          string
	NoMoreThan        string
	LessThan          string
	NoLessThan        string
	StartsWith        string
	Contains          string
	DoesNotContain    string
	EndsWith          string
	Matches           string
	DoesNotMatch      string
	VersionMoreThan   string
	VersionNoMoreThan string
	VersionLessThan   string
	VersionNoLessThan string
}{
	Exists:            OperatorStrings[0],
	NotExists:         OperatorStrings[1],
	Equals:            OperatorStrings[2],
	DoesNotEqual:      OperatorStrings[3],
	MoreThan:          OperatorStrings[4],
	NoMoreThan:        OperatorStrings[5],
	LessThan:          OperatorStrings[6],
	NoLessThan:        OperatorStrings[7],
	StartsWith:        OperatorStrings[8],
	Contains:          OperatorStrings[9],
	DoesNotContain:    OperatorStrings[10],
	EndsWith:          OperatorStrings[11],
	Matches:           OperatorStrings[12],
	DoesNotMatch:      OperatorStrings[13],
	VersionMoreThan:   OperatorStrings[14],
	VersionNoMoreThan: OperatorStrings[15],
	VersionLessThan:   OperatorStrings[16],
	VersionNoLessThan: OperatorStrings[17],
}

var stepsReadingFunc = func(dt []byte, out interface{}) error {
	var payload stepsJSON
	if err := json.Unmarshal(dt, &payload); err != nil {
		return err
	}

	var steps = make([]Step, payload.Count)
	for i := 0; i < int(payload.Count); i++ {
		sdt, err := json.Marshal(payload.Items[i])
		if err != nil {
			return err
		}
		err = stepReadingFunc(sdt, &steps[i])
		if err != nil {
			return err
		}
	}
	replaceValue(out, &steps)
	return nil
}

var stepReadingFunc = func(dt []byte, out interface{}) error {
	var payload stepJSON
	if err := json.Unmarshal(dt, &payload); err != nil {
		return err
	}

	var step Step
	var err error
	switch payload.Type {
	case string(StepTypePowershell):
		var ps StepPowershell
		err = ps.UnmarshalJSON(dt)
		step = &ps
	case string(StepTypeCommandLine):
		var cmd StepCommandLine
		err = cmd.UnmarshalJSON(dt)
		step = &cmd
	case string(StepTypeOctopusPushPackage):
		var opp StepOctopusPushPackage
		err = opp.UnmarshalJSON(dt)
		step = &opp
	case string(StepTypeOctopusCreateRelease):
		var ocr StepOctopusCreateRelease
		err = ocr.UnmarshalJSON(dt)
		step = &ocr
	default:
		return fmt.Errorf("Unsupported step type: '%s' (id:'%s')", payload.Type, payload.ID)
	}
	if err != nil {
		return err
	}

	replaceValue(out, &step)
	return nil
}
