package events

// Dependency inversion interface
type Message interface {
	Type() string
}
