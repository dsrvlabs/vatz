package types

type AliveStatus string

// AliveStatus is type that describes aliveness flags.
const (
	AliveStatusUp   AliveStatus = "UP"
	AliveStatusDown AliveStatus = "DOWN"
)
