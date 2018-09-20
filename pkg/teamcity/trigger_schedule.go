package teamcity

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TriggerSchedulingPolicy represents the shceduling policy for a scheduled trigger. Can be 'daily', 'weekly' or 'cron'
type TriggerSchedulingPolicy = string

const (
	//TriggerSchedulingDaily triggers every day
	TriggerSchedulingDaily TriggerSchedulingPolicy = "daily"
	//TriggerSchedulingWeekly triggers at specified day + time of the week, once per week
	TriggerSchedulingWeekly TriggerSchedulingPolicy = "weekly"
	//TriggerSchedulingCron triggers by matching a cron expression
	TriggerSchedulingCron TriggerSchedulingPolicy = "cron"
)

//TriggerSchedule represents a build trigger that fires on a time-bound schedule
type TriggerSchedule struct {
	triggerJSON *triggerJSON
	buildTypeID string

	SchedulingPolicy TriggerSchedulingPolicy
	Rules            []string
	Timezone         string
	Weekday          time.Weekday
	Hour             uint
	Minute           uint
	Second           uint
}

//ID for this entity
func (t *TriggerSchedule) ID() string {
	return t.triggerJSON.ID
}

//Type returns TriggerTypes.Schedule ("schedulingTrigger")
func (t *TriggerSchedule) Type() string {
	return TriggerTypes.Schedule
}

//SetDisabled controls whether this trigger is disabled or not
func (t *TriggerSchedule) SetDisabled(disabled bool) {
	t.triggerJSON.Disabled = NewBool(disabled)
}

//Disabled gets the disabled status for this trigger
func (t *TriggerSchedule) Disabled() bool {
	return *t.triggerJSON.Disabled
}

//BuildTypeID gets the build type identifier
func (t *TriggerSchedule) BuildTypeID() string {
	return t.buildTypeID
}

//SetBuildTypeID sets the build type identifier
func (t *TriggerSchedule) SetBuildTypeID(id string) {
	t.buildTypeID = id
}

//NewDailyTriggerSchedule returns a TriggaerSchedule that fires daily on the hour/minute specified
func NewDailyTriggerSchedule(hour uint, minute uint, timezone string, rules []string) (*TriggerSchedule, error) {
	if hour > 23 {
		return nil, fmt.Errorf("Invalid hour: %d, must be between 0-23", hour)
	}
	if minute > 59 {
		return nil, fmt.Errorf("Invalid minute: %d, must be between 0-59", minute)
	}

	return &TriggerSchedule{
		SchedulingPolicy: TriggerSchedulingDaily,
		Timezone:         timezone,
		Rules:            rules,
		Weekday:          time.Sunday,
		Hour:             hour,
		Minute:           minute,
		Second:           0,
	}, nil
}

func (t *TriggerSchedule) read(dt *triggerJSON) error {
	if dt.Disabled != nil {
		t.SetDisabled(*dt.Disabled)
	}
	t.triggerJSON = dt

	if v, ok := dt.Properties.GetOk("schedulingPolicy"); ok {
		t.SchedulingPolicy = v
	} else {
		return fmt.Errorf("invalid 'schedulingPolicy' property")
	}

	if v, ok := dt.Properties.GetOk("triggerRules"); ok {
		t.Rules = strings.Split(v, "\n")
	}

	if v, ok := dt.Properties.GetOk("timezone"); ok {
		t.Timezone = v
	}

	switch t.SchedulingPolicy {
	case TriggerSchedulingDaily:
		return t.readDaily(dt)
	}
	return nil
}

func (t *TriggerSchedule) readDaily(dt *triggerJSON) error {
	if v, ok := dt.Properties.GetOk("hour"); ok {
		p, _ := strconv.ParseUint(v, 10, 0)
		t.Hour = uint(p)
	} else {
		return fmt.Errorf("invalid 'hour' property")
	}
	if v, ok := dt.Properties.GetOk("minute"); ok {
		p, _ := strconv.ParseUint(v, 10, 0)
		t.Minute = uint(p)
	} else {
		return fmt.Errorf("invalid 'minute' property")
	}

	return nil
}

func (t *TriggerSchedule) properties() *Properties {
	props := NewProperties()

	props.AddOrReplaceValue("timezone", t.Timezone)
	props.AddOrReplaceValue("schedulingPolicy", t.SchedulingPolicy)
	props.AddOrReplaceValue("triggerRules", strings.Join(t.Rules, "\n"))
	props.AddOrReplaceValue("hour", fmt.Sprint(t.Hour))
	props.AddOrReplaceValue("minute", fmt.Sprint(t.Minute))

	return props
}

//MarshalJSON implements JSON serialization for TriggerVcs
func (t *TriggerSchedule) MarshalJSON() ([]byte, error) {
	out := &triggerJSON{
		ID:         t.ID(),
		Type:       t.Type(),
		Disabled:   NewBool(t.Disabled()),
		Properties: t.properties(),
	}

	return json.Marshal(out)
}

//UnmarshalJSON implements JSON deserialization for TriggerSchedule
func (t *TriggerSchedule) UnmarshalJSON(data []byte) error {
	var aux triggerJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Type != TriggerTypes.Schedule {
		return fmt.Errorf("invalid type %s trying to deserialize into TriggerSchedule entity", aux.Type)
	}

	if err := t.read(&aux); err != nil {
		return err
	}

	return nil
}
