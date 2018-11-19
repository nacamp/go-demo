//http://marcio.io/2015/07/singleton-pattern-in-go/
package syncs

import (
	"sync"
)

type singleton struct {
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}

//Check-Lock-Check
var mu sync.Mutex

func GetInstance2() *singleton {
	if instance == nil { // <-- Not yet perfect. since it's not fully atomic
		mu.Lock()
		defer mu.Unlock()

		if instance == nil {
			instance = &singleton{}
		}
	}
	return instance
}
