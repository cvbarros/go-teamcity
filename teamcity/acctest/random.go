package acctest

import (
	"fmt"
	"math/rand"
	"time"
)

// RandomWithPrefix is used to generate a unique name with a prefix, for
// randomizing names in acceptance tests
func RandomWithPrefix(name string) string {
	reseed()
	return fmt.Sprintf("%s-%d", name, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}

func reseed() {
	rand.Seed(time.Now().UTC().UnixNano())
}
