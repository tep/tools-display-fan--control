package main

type fanState int

const (
	fsNone fanState = iota
	fsOff
	fsOn
	fsChange
)

func (fs fanState) String() string {
	switch fs {
	case fsOff:
		return "FAN_OFF"
	case fsOn:
		return "FAN_ON"
	case fsChange:
		return "FAN_CHANGE"
	default:
		return ""
	}
}
