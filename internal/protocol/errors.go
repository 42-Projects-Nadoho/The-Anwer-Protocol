package protocol

// From RFC 42TAP section 8.2.
const (
	ErrNameInUse          = 201
	ErrNoExit             = 301
	ErrNotInGroup         = 401
	ErrAlreadyInGroup     = 402
	ErrItemNotFound       = 404
	ErrItemNotInInventory = 404
	ErrNPCNotFound        = 404
	ErrNPCNotHostile      = 405
	ErrNoQuestAvailable   = 406
	ErrConnectionFailed   = 900
	ErrSendFailed         = 901
)

// Not in RFC 42TAP; documented as deviations in the README.
const (
	ErrNotAuthenticated = 902
	ErrUnknownCommand   = 903
)
