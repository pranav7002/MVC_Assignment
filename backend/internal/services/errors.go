package services

import "errors"

var (
	// server
	ErrServer = errors.New("Something bad happened on the server :/")

	// auth
	ErrUserExists   = errors.New("username is already taken")
	ErrInvalidCreds = errors.New("invalid username or password")

	// village
	ErrBuildingNotFound      = errors.New("building not found or does not belong to user")
	ErrInsufficientResources = errors.New("insufficient resources for this upgrade")
	ErrBuildingNotUnlocked   = errors.New("building not available at this town hall level")
	ErrBuildingLimitReached  = errors.New("more buildings not allowed at this town hall level")
	ErrHighestLevelReached   = errors.New("upgrade not availible at this town hall level")
	ErrCollisionDetected     = errors.New("building placement overlaps with an existing building")
	ErrOutOfBounds           = errors.New("building placement is out of village bounds")
	ErrInvalidBuildingType   = errors.New("invalid building type")

	// troops
	ErrInsufficientElixir       = errors.New("insufficient elixir to train requested troops")
	ErrInsufficientHousingSpace = errors.New("insufficient housing space to train requested troops")
	ErrInvalidTroop             = errors.New("invalid troop type")
)
