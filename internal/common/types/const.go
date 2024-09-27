package types

import (
	"github.com/ginx-contribs/ginx/constant/status"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
)

// custom code is composed of three parts: Order_Status_Code, it will be shown in the response body.
// Order just represents order of package create time, it is used to avoid duplicates error code in different packages.
// Status represents the error will be occurred in which situation, it is corresponds to http status.
// Code is the true error code, whose max capacity is 999.
const customCode = 0_000_000

var (
	ErrBadParams = statuserr.Errorf("bad parameters").SetCode(400_001).SetStatus(status.BadRequest)

	ErrInternal = statuserr.Errorf("internal server error").SetCode(500_000).SetStatus(status.InternalServerError)
)
