package teamcity

import (
	"strconv"
	"strings"
)

// FinishBuildTriggerOptions represents optional settings for a VCS Trigger type.
type FinishBuildTriggerOptions struct {
	AfterSuccessfulBuildOnly bool
	BranchFilter             []string
}

// NewFinishBuildTriggerOptions initialize a NewFinishBuildTriggerOptions
// branchFilter can be passed as "nil" to not filter on any specific branches
func NewFinishBuildTriggerOptions(afterSuccessfulBuildOnly bool, branchFilter []string) *FinishBuildTriggerOptions {
	return &FinishBuildTriggerOptions{
		AfterSuccessfulBuildOnly: afterSuccessfulBuildOnly,
		BranchFilter:             branchFilter,
	}
}

func (o *FinishBuildTriggerOptions) properties() *Properties {
	props := NewPropertiesEmpty()

	//Defaults to false, so ommit emitting the property if 'false'
	if o.AfterSuccessfulBuildOnly {
		props.AddOrReplaceValue("afterSuccessfulBuildOnly", strconv.FormatBool(o.AfterSuccessfulBuildOnly))
	}

	if o.BranchFilter != nil && len(o.BranchFilter) > 0 {
		props.AddOrReplaceValue("branchFilter", strings.Join(o.BranchFilter, "\n"))
	}

	return props
}
