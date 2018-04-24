package teamcity

// NewTrue is a helper function to return a *bool to true
func NewTrue() *bool {
	b := true
	return &b
}

// NewFalse is a helper function to return a *bool to true
func NewFalse() *bool {
	b := false
	return &b
}
