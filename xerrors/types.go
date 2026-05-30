package xerrors

// MetaMessage encapsulates static corporate fallback strings and semantic tagging rules
// required to build standardized, predictable error formatting structures.
type MetaMessage struct {
	message   string
	fieldRule string
	extraTags []string
}

// NewMetaMessage acts as a secure constructor that initializes and returns a populated MetaMessage instance.
// By exposing this factory, the package enables external domain extensions to define fallback metadata
// while strictly preserving the encapsulation of internal structural layout properties.
func NewMetaMessage(
	message string,
	fieldRule string,
	extraTags []string,
) MetaMessage {
	return MetaMessage{
		message:   message,
		fieldRule: fieldRule,
		extraTags: extraTags,
	}
}
