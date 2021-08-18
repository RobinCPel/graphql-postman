package reformatted

import (
	"github.com/RobinCPel/graphql-postman/src/internal/graphql/introspection"
	"github.com/RobinCPel/graphql-postman/src/internal/graphql/kind"
	log "github.com/sirupsen/logrus"
	"strings"
)

// reformatTypeRef reformats an introspection.TypeRef to a reformatted.TypeRef.
func reformatTypeRef(typeRef introspection.TypeRef) TypeRef {
	var t TypeRef
	nn := false
	a := &typeRef

	for a != nil {
		switch strings.ToLower(a.Kind) {

		case kind.NonNull:
			nn = true

		case kind.List:
			if t.List {
				if t.TwoDList {
					log.Warn("Lists with more than two dimensions are not supported!")
					break
				}

				t.TwoDList = true
				if nn == true {
					t.TwoDListNonNull = true
					nn = false
				}
			} else {
				t.List = true
				if nn == true {
					t.ListNonNull = true
					nn = false
				}
			}

		default:
			t.Kind = a.Kind
			t.Name = a.Name
			if nn == true {
				t.NonNull = true
				nn = false
			}
		}

		// Go to the next type
		a = a.OfType
	}

	return t
}

// Reformat takes an introspection.Model and converts it to a reformatted.Model.
// The reformatted Model is easier to parse.
func Reformat(model *introspection.Model) *Model {
	if model == nil {
		return nil
	}

	types := make([]introspection.Type, len(model.Data.Schema.Types)-2)
	var mutation, query introspection.Type

	// Divide the types into a  mutation, a query, and the rest
	for _, t := range model.Data.Schema.Types {
		if t.Name == model.Data.Schema.MutationType.Name {
			mutation = t
		} else if t.Name == model.Data.Schema.QueryType.Name {
			query = t
		} else {
			types = append(types, t)
		}
	}

	reformatted := Model{
		Mutations: make([]Operation, len(mutation.Fields)),
		Queries:   make([]Operation, len(query.Fields)),
		Types:     make(map[string]Type),
	}

	// Reformat the mutations
	for i, m := range mutation.Fields {

		arguments := make(map[string]TypeRef)
		for _, a := range m.Arguments {
			arguments[a.Name] = reformatTypeRef(a.Type)
		}

		reformatted.Mutations[i] = Operation{
			Name:      m.Name,
			Arguments: arguments,
			Type:      reformatTypeRef(m.Type),
		}
	}

	// Reformat the queries
	for i, q := range query.Fields {

		arguments := make(map[string]TypeRef)
		for _, a := range q.Arguments {
			arguments[a.Name] = reformatTypeRef(a.Type)
		}

		reformatted.Queries[i] = Operation{
			Name:      q.Name,
			Arguments: arguments,
			Type:      reformatTypeRef(q.Type),
		}
	}

	// Reformat the types
	for _, t := range types {
		rt := Type{
			Name:          t.Name,
			Fields:        make(map[string]TypeRef),
			InputFields:   make(map[string]TypeRef),
			EnumValues:    make([]string, 0),
			PossibleTypes: make([]TypeRef, 0),
		}

		for _, f := range t.Fields {
			// Field arguments are skipped because they are only present for mutations and queries, not types.
			rt.Fields[f.Name] = reformatTypeRef(f.Type)
		}

		for _, f := range t.InputFields {
			rt.InputFields[f.Name] = reformatTypeRef(f.Type)
		}

		for _, e := range t.EnumValues {
			rt.EnumValues = append(rt.EnumValues, e.Name)
		}

		for _, p := range t.PossibleTypes {
			rt.PossibleTypes = append(rt.PossibleTypes, reformatTypeRef(p))
		}

		reformatted.Types[rt.Name] = rt
	}

	return &reformatted
}
