package reformatted

// TypeRef is the reformatted version fo introspection.TypeRef.
// The main benefit of this one, is that the data is flat, there are no "infinite"
// layers deep that you'll have to go through to determine the exact typing.
type TypeRef struct {
	Kind            string // One of the root kinds e.g. kind.Scalar, kind.Enum, kind.Union
	Name            string // Name of the type, e.g. scalar.Float, scalar.ID
	NonNull         bool   // Whether or not the type is nullable
	List            bool   // Whether or not the type is in a list
	ListNonNull     bool   // Whether or not the list is nullable
	TwoDList        bool   // Whether or not the type is in a double list
	TwoDListNonNull bool   // Whether or not the double list is nullable
}

// Type is generic type that can describe a scalar, input,
// enum, interface, object, or union type.
//
// | Type      | Has                           |
// --------------------------------------------|
// | Scalar    | Name                          |
// | Input     | Name, InputFields             |
// | Enum      | Name, EnumValues              |
// | Interface | Name, Fields, PossibleTypes   |
// | Object    | Name, Fields                  |
// | Union     | Name, PossibleTypes           |
//
type Type struct {
	Name          string
	Fields        map[string]TypeRef // The key is the name of the field, the value is the value
	InputFields   map[string]TypeRef // The key is the name of the field, the value is the value
	EnumValues    []string
	PossibleTypes []TypeRef
}

// Operation is a struct that contains the data needed for a GraphQL query or mutation.
type Operation struct {
	Name      string
	Arguments map[string]TypeRef // What the operation input requires, the key is the name of the field, the value is the value
	Type      TypeRef            // What the operation returns
}

// Model is the reformatted version of the introspection.Model.
type Model struct {
	Mutations []Operation
	Queries   []Operation
	Types     map[string]Type
}
