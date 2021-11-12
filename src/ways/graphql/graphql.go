package graphql

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/message/payload"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/way/driver"
	"github.com/dmalykh/axeloy/axeloy/way/model"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"net/http"
)

type Field struct {
	Key   string `graphql:",optional"`
	Value []string
}

type Message struct {
	Publisher    []Field       `graphql:"publisher"`
	Destinations []Destination `graphql:"destinations,optional"`
	Payload      []Field
}

type Destination struct {
	Way      string
	Consumer []Field
}

type GraphQl struct {
	Addr string
}

func (g *GraphQl) ValidateProfile(ctx context.Context, p profile.Profile) (map[string]string, error) {
	panic("implement me")
}

func (g *GraphQl) SetWayParams(params driver.Params) {
	panic("implement me")
}

func (g *GraphQl) SetConfig(config driver.DriverConfig) {
	panic("implement me")
}

func (g *GraphQl) Stop() error {
	panic("implement me")
}

func (g *GraphQl) Listen(ctx context.Context, f func(ctx context.Context, message driver.Message) error) error {
	var schema = g.buildSchema(ctx, f)

	introspection.AddIntrospectionToSchema(schema)

	// Expose schema and graphiql.
	http.Handle("/graphql", graphql.Handler(schema))
	http.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiql.Handler()))
	return http.ListenAndServe(g.Addr, nil)
}

func (g *GraphQl) buildSchema(ctx context.Context, f func(ctx context.Context, message driver.Message) error) *graphql.Schema {
	var schema = schemabuilder.NewSchema()
	schema.Object("Field", Field{})
	schema.Object("Message", Message{})
	schema.Object("Destination", Destination{})

	var obj = schema.Mutation()
	obj.FieldFunc("send", func(args Message) bool {
		err := f(ctx, &model.Message{
			Payload: func(fields []Field) payload.Payload {
				var p = make(payload.Payload)
				for _, field := range fields {
					if _, exists := p[field.Key]; !exists {
						p[field.Key] = make([]string, len(field.Value))
					}
					p[field.Key] = append(p[field.Key], field.Value...)
				}
				return p
			}(args.Payload),
			Publisher: func(fields []Field) profile.Fields {
				var p = make(profile.Fields)
				for _, field := range fields {
					if _, exists := p[field.Key]; !exists {
						p[field.Key] = make([]string, len(field.Value))
					}
					p[field.Key] = append(p[field.Key], field.Value...)
				}
				return p
			}(args.Publisher),
			Destinations: func(destinations []Destination) []driver.Destination {
				var dests = make([]driver.Destination, len(destinations))
				for i, destination := range destinations {
					dests[i] = &model.Destination{
						Way: destination.Way,
						Consumer: func(fields []Field) profile.Fields {
							var p = make(profile.Fields)
							for _, field := range fields {
								if _, exists := p[field.Key]; !exists {
									p[field.Key] = make([]string, len(field.Value))
								}
								p[field.Key] = append(p[field.Key], field.Value...)
							}
							return p
						}(destination.Consumer),
					}
				}
				return dests
			}(args.Destinations),
		})
		//@TODO: response must be struct
		return err != nil
	})
	return schema.MustBuild()
}
