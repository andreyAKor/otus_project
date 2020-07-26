package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cucumber/godog"
)

type serviceTest struct {
	responseStatusCode int
	responseTimeout    time.Duration
}

func (test *serviceTest) iSendRequestTo(httpMethod, addr string) (err error) {
	var r *http.Response

	start := time.Now()

	switch httpMethod {
	case http.MethodGet:
		r, err = http.Get(addr)
	default:
		err = fmt.Errorf("unknown method: %s", httpMethod)
	}

	if err != nil {
		return
	}

	test.responseTimeout = time.Now().Sub(start)
	test.responseStatusCode = r.StatusCode

	return
}

func (test *serviceTest) theResponseCodeShouldBe(code int) error {
	if test.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseStatusCode, code)
	}
	return nil
}

func (test *serviceTest) iSendRequestToAndSendSecondRequestTo(httpMethod1, addr1, httpMethod2, addr2 string) error {
	if err := test.iSendRequestTo(httpMethod1, addr1); err != nil {
		return fmt.Errorf("first request has error: %s", err.Error())
	}
	if err := test.iSendRequestTo(httpMethod1, addr1); err != nil {
		return fmt.Errorf("second request has error: %s", err.Error())
	}

	return nil
}

func (test *serviceTest) theResponseCodeShouldBeAndResonseTimeoutLess(code int, timeout string) error {
	if err := test.theResponseCodeShouldBe(code); err != nil {
		return err
	}

	timeoutDur, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("response timeout parsing fail: %s", err.Error())
	}

	if test.responseTimeout > timeoutDur {
		return fmt.Errorf("the response more than: %s > %s", test.responseTimeout, timeoutDur)
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	test := new(serviceTest)

	s.Step(`^I send "([^"]*)" request to "([^"]*)"$`, test.iSendRequestTo)
	s.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)

	s.Step(`^I send "([^"]*)" request to "([^"]*)" and send "([^"]*)" second request to "([^"]*)"$`, test.iSendRequestToAndSendSecondRequestTo)
	s.Step(`^The response code should be (\d+) and resonse timeout less "([^"]*)"$`, test.theResponseCodeShouldBeAndResonseTimeoutLess)
}
