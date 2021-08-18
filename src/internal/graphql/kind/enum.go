package kind

// Base, the kind that every Referring kind will eventually refer to, "it represents the leaves of a type".
type Base = string

const Scalar Base = "scalar"

// Referring, refers to another Referring kind or a Base kind.
type Referring = string

const (
	InputObject Referring = "input_object"
	Enum        Referring = "enum"
	Interface   Referring = "interface"
	Object      Referring = "object"
	Union       Referring = "union"
)

// Descriptive, describes whether another kind is non-nullable, or in a list.
type Descriptive = string

const (
	NonNull Descriptive = "non_null"
	List    Descriptive = "list"
)
