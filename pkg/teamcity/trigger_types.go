package teamcity

type triggerType = string

const (
	//Vcs trigger type
	Vcs triggerType = "vcsTrigger"
)

// TriggerTypes represents possible types for build triggers
var TriggerTypes = struct {
	Vcs        triggerType
	Dependency triggerType
	Schedule   triggerType
}{
	Vcs: Vcs,
}
