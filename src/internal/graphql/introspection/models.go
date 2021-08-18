package introspection

type Named struct {
	Name string `json:"name"`
}

// TypeRef stores a type reference.
// A type reference has at least one kind.Referring or kind.Base, and can
// have multiple kind.Descriptive before that.
//
// For example, here is a TypeRef that has a non-nullable list of non-nullable scalar.Integer's:
// "typeRef": {
//   "name": null,
//   "kind": "non_null", // kind.NonNull
//   "ofType": {
//     "name": null,
//     "kind": "list", // kind.List
//     "ofType": {
//       "name": null,
//       "kind": "non_null", // kind.NonNull
//       "ofType":  {
//         "name": "Integer", // scalar.Integer
//         "kind": "scalar",  // kind.Scalar
//         "ofType": null
//       }
//     }
//   }
// }
//
// Note that the name is always null until you get to a kind.Referring or kind.Base.
type TypeRef struct {
	Named
	Kind   string   `json:"kind"`
	OfType *TypeRef `json:"ofType"`
}

type NamedTypeRef struct {
	Named
	Type TypeRef `json:"type"`
}

type TypeField struct {
	Named
	Arguments []NamedTypeRef `json:"args"`
	Type      TypeRef        `json:"type"`
}

// Type defines a GraphQL type, dependent on what type it is,
// certain fields are filled, while some are left empty.
//
// In the description above reformatted.Type, you can see a
// table that describes which fields are filled for each type.
type Type struct {
	Named
	Fields        []TypeField    `json:"fields"`
	InputFields   []NamedTypeRef `json:"inputFields"`
	EnumValues    []Named        `json:"enumValues"`
	PossibleTypes []TypeRef      `json:"possibleTypes"`
}

// Model encapsulates all of the data that is returned from a GraphQL introspection query.
type Model struct {
	Data struct {
		Schema struct {
			QueryType    Named  `json:"queryType"`    // Contains the name of the type that contains all queries
			MutationType Named  `json:"mutationType"` // Contains the name of the type that contains all mutations
			Types        []Type `json:"types"`
		} `json:"__schema"`
	} `json:"data"`
}
