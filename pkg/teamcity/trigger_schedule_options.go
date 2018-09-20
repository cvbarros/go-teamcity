package teamcity

//TriggerScheduleOptions represent options for configuring a scheduled build trigger
type TriggerScheduleOptions struct {
	TriggerIfWatchedBuildChanges        bool                       `prop:"triggerBuildIfWatchedBuildChanges"`
	BuildOnAllCompatibleAgents          bool                       `prop:"triggerBuildOnAllCompatibleAgents"`
	BuildWithPendingChangedOnly         bool                       `prop:"triggerBuildWithPendingChangesOnly"`
	PromoteWatchedBuild                 bool                       `prop:"promoteWatchedBuild"`
	RevisionRuleSourceBuildID           string                     `prop:"revisionRuleDependsOn"`
	RevisionRule                        ArtifactDependencyRevision `prop:"revisionRule"`
	EnforceCleanCheckout                bool                       `prop:"enforceCleanCheckout"`
	EnforceCleanCheckoutForDependencies bool                       `prop:"enforceCleanCheckoutForDependencies"`
	QueueOptimization                   bool                       `prop:"enableQueueOptimization"`
}
