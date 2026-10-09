package errors

import (
	"errors"
)

var NotFound = errors.New("requested resource not found")
var NoChecks = errors.New("no checks stored in database yet")
var NoSites = errors.New("no sites stored in database yet")
var SiteExists = errors.New("site already exists")
