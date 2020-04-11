package timeout

import "time"

//SetTimeout to call a callback after duration
func SetTimeout(cabllack func(params ...interface{}), ms int, params ...interface{}) {
	//Disattach by goroutine
	go func(params ...interface{}) {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		cabllack(params...)
	}(params)
}
