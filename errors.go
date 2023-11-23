package ibweb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IBError struct {
	Err string `json:"error"`
}

func NewIBError(resp *http.Response) IBError {
	defer resp.Body.Close()
	var err IBError

	v, ioErr := io.ReadAll(resp.Body)
	if ioErr != nil {
		return IBError{
			Err: "interactive brokers did not describe error",
		}
	}

	if err := json.Unmarshal(v, &err); err != nil {
		return IBError{
			Err: fmt.Sprintf("failed to parse IB error '%s' into go error", string(v)),
		}
	}

	return err
}

func (i IBError) Error() string {
	return i.Err
}

type StatusCodeError struct {
	StatusCode int
	Err        error
}

func (s StatusCodeError) Error() string {
	return fmt.Sprintf("invalid status code '%d': %s", s.StatusCode, s.Err)
}
