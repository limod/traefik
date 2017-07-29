package try

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

// RepeatRequest is like Repeat, but runs a request against the given URL and applies
// the condition on the response.
// ResponseCondition may be nil, in which case only the request against the URL must
// succeed.
func RepeatRequest(req *http.Request, timeout time.Duration, conditions ...ResponseCondition) error {
	resp, err := doRepeatRequest(req, timeout, conditions...)

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	return err
}

func doRepeatRequest(request *http.Request, timeout time.Duration, conditions ...ResponseCondition) (*http.Response, error) {
	return doRequest(Repeat, timeout, request, conditions...)
}

// Repeat repeatedly executes an operation until error condition occurs or the
// given timeout is reached, whatever comes first.
func Repeat(timeout time.Duration, operation DoCondition) error {
	if timeout <= 0 {
		panic("timeout must be larger than zero")
	}

	interval := time.Duration(math.Ceil(float64(timeout) / 15.0))
	if interval > maxInterval {
		interval = maxInterval
	}

	timeout = applyCIMultiplier(timeout)

	var err error
	if err = operation(); err != nil {
		fmt.Println("-")
		return fmt.Errorf("repeat operation failed: %s", err)
	}
	fmt.Print("*")

	stopTimer := time.NewTimer(timeout)
	defer stopTimer.Stop()
	retryTick := time.NewTicker(interval)
	defer retryTick.Stop()

	for {
		select {
		case <-stopTimer.C:
			fmt.Println("+")
			return nil
		case <-retryTick.C:
			fmt.Print("*")
			if err = operation(); err != nil {
				fmt.Println("-")
				return err
			}
		}
	}
}
