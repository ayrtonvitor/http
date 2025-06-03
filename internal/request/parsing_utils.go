package request

import (
	"fmt"
	"strconv"

	"github.com/ayrtonvitor/http/internal/headers"
)

type handleInnerParseReturnParams struct {
	nRead      int
	nextState  reqState
	finalizer  func() error
	err        error
	reqEnd     bool
	callerDone bool
}

func (r *Request) handleInnerParseReturn(pars handleInnerParseReturnParams) (int, error) {
	if pars.err != nil {
		return 0, pars.err
	}
	if pars.reqEnd {
		if !pars.callerDone {
			fmt.Println("BEEN HERE!!!!!!")
			return 0, ErrMalformedReq
		}
		r.state = reqStateDone
		return pars.nRead, nil
	}
	if pars.callerDone {
		if pars.finalizer != nil {
			err := pars.finalizer()
			if err != nil {
				return 0, pars.err
			}
		}
		r.state = pars.nextState
		return pars.nRead, nil
	}
	return pars.nRead, nil
}

func (r *Request) prepareToParseBody() error {
	contLen, err := handleContentLength(r.Headers.Get(headers.ContentLength))
	if err != nil {
		return err
	}
	r.Body.expectedLength = contLen

	return nil
}

func handleContentLength(sContLen string) (int, error) {
	if sContLen == "" {
		return 0, nil
	}
	contLen, err := strconv.Atoi(sContLen)
	if err != nil {
		return 0, fmt.Errorf("Invalid content length: %w", err)
	}

	return contLen, nil
}
