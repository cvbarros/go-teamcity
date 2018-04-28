package teamcity

import "github.com/lann/builder"

type stepType = string

const (
	//Powershell step type
	Powershell stepType = "jetbrains_powershell"
	//Dotnet CLI step type
	DotnetCli stepType = "dotnet.cli"
	//Commandline (shell/cmd) step type
	CommandLine stepType = "simpleRunner"
)

// StepTypes represents most common step types for build steps
var StepTypes = struct {
	Powershell  stepType
	DotnetCli   stepType
	CommandLine stepType
}{
	Powershell:  Powershell,
	DotnetCli:   DotnetCli,
	CommandLine: CommandLine,
}

type powershellStepBuilder builder.Builder

// ScriptFile sets properties required to run the powershell step as a script file
func (b powershellStepBuilder) ScriptFile(scriptFile string) powershellStepBuilder {
	props := []*Property{
		&Property{
			Name:  "jetbrains_powershell_execution",
			Value: "PS1",
		},
		&Property{
			Name:  "jetbrains_powershell_noprofile",
			Value: "true",
		},
		&Property{
			Name:  "jetbrains_powershell_script_mode",
			Value: "FILE",
		},
		&Property{
			Name:  "teamcity.step.mode",
			Value: "default",
		},
	}

	value := NewProperties(props...)
	out := builder.Delete(b, "Properties").(powershellStepBuilder)
	out = builder.Set(out, "Type", StepTypes.Powershell).(powershellStepBuilder)
	return builder.Set(out, "Properties", value).(powershellStepBuilder)
}

// Args sets properties required for script arguments
func (b powershellStepBuilder) Args(args string) powershellStepBuilder {
	argProp := &Property{
		Name:  "jetbrains_powershell_scriptArguments",
		Value: args,
	}

	ret, exists := builder.Get(b, "Properties")
	if !exists {
		props := NewProperties(argProp)
		return builder.Set(b, "Properties", props).(powershellStepBuilder)
	}
	props := ret.(*Properties)
	props.Add(argProp)
	return builder.Set(b, "Properties", props).(powershellStepBuilder)
}

func (b powershellStepBuilder) Build(name string) *Step {
	out := builder.GetStruct(b).(Step)
	out.Type = StepTypes.Powershell
	out.Name = name
	return &out
}

// PowershellStepBuilder is a convenience class for creating powershell build steps
var PowershellStepBuilder = builder.Register(powershellStepBuilder{}, Step{}).(powershellStepBuilder)
