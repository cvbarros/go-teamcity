package teamcity

import (
	"fmt"
	"net/url"
)

//Locator represents a arbitraty locator to be used when querying resources, such as id:, type:, or key:
//These are used in GET requests within the URL so must be properly escaped
type Locator string

//LocatorID creates a locator for a Project/BuildType by Id
func LocatorID(id string) Locator {
	return Locator(url.QueryEscape("id:") + id)
}

//LocatorIDInt creates a locator for a Project/BuildType by Id where the Id's an integer
func LocatorIDInt(id int) Locator {
	return Locator(url.QueryEscape("id:") + fmt.Sprint(id))
}

//LocatorName creates a locator for User/Project/BuildType by Name
func LocatorName(name string) Locator {
	return Locator(url.QueryEscape("name:") + url.PathEscape(name))
}

//LocatorUsername creates a locator for User by Username
func LocatorUsername(name string) Locator {
	return Locator(url.QueryEscape("username:") + url.PathEscape(name))
}

//LocatorKey creates a locator for Group by Key
func LocatorKey(key string) Locator {
	return Locator(url.QueryEscape("key:") + url.PathEscape(key))
}

//LocatorType creates a locator for a Project Feature by Type
func LocatorType(id string) Locator {
	return Locator(url.QueryEscape("type:") + id)
}

//LocatorStart creates a locator to set offset
func LocatorStart(start int) Locator {
	return Locator(url.QueryEscape("start:") + fmt.Sprint(start))
}

//LocatorCount creates a locator to set number of answers
func LocatorCount(count int) Locator {
	return Locator(url.QueryEscape("count:") + fmt.Sprint(count))
}

func (l Locator) String() string {
	return string(l)
}
