package teamcity

func buildQueryLocator(locators ...Locator) *queryStruct {
	locatorQuery := ""
	if len(locators) > 1 {
		for _, locator := range locators[1:] {
			locatorQuery += "," + locator.String()
		}
	}
	if len(locators) >= 1 {
		locatorQuery = locators[0].String() + locatorQuery
	}
	return &queryStruct{
		key:   "locator",
		value: locatorQuery,
	}
}
