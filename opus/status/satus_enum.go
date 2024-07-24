package status

//go:generate stringer -type Status

type Status int

const (
	Undefined Status = iota
	Upcoming
	Live
	Archived
)
