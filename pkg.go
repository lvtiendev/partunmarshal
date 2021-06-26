package partunmarhal

import "errors"

const (
	tagU    = "u"
	tagJSON = "json"
)

const (
	updatable = "true"
)

// Errors exposed by the package
var (
	ErrPointerExpected  = errors.New("pointer expected")
	ErrFieldCannotBeSet = errors.New("field cannot be set")
	ErrNoTagJSON        = errors.New("field does not have json tag")
)
