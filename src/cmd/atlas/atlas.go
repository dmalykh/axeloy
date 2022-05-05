package atlas

import (
	"ariga.io/atlas/cmd/action"
	"ariga.io/atlas/sql/schema"
	"context"
	"fmt"
	"github.com/dmalykh/axeloy/repository/db"
	"io/ioutil"
)

type atlas struct {
	driver *Driver
	schema string
}

// See sql.Open for list of driveers
func NewAtlas(ctx context.Context, driverName string, dsn string) (*atlas, error) {
	//Connect to database
	conn, err := db.Connect(ctx, driverName, dsn)
	if err != nil {
		return nil, err
	}

	// Open with driver
	driver, err := open(driverName, conn)
	if err != nil {
		return nil, err
	}
	// Get schema
	s, err := action.SchemaNameFromURL(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &atlas{
		driver: driver,
		schema: s,
	}, nil
}

// Inspect schema
func (a *atlas) Inspect(ctx context.Context) ([]byte, error) {
	realm, err := a.driver.InspectRealm(ctx, &schema.InspectRealmOption{
		Schemas: []string{a.schema},
	})
	if err != nil {
		return nil, err
	}
	ddl, err := a.driver.MarshalSpec(realm)
	if err != nil {
		return nil, err
	}
	return ddl, nil
}

// Apply syncs schema and returns true when ensure that schema is synced
func (a *atlas) Apply(ctx context.Context, file string) (bool, error) {
	realm, err := a.driver.InspectRealm(ctx, &schema.InspectRealmOption{
		Schemas: []string{a.schema},
	})
	if err != nil {
		return false, err
	}

	f, err := ioutil.ReadFile(file)
	if err != nil {
		return false, err
	}
	// Unmarshal spec
	var desired = new(schema.Realm)
	if err := a.driver.UnmarshalSpec(f, desired); err != nil {
		return false, err
	}

	if len(realm.Schemas) > 0 {
		// Validate all schemas in file were selected by user.
		sm := make(map[string]bool, len(realm.Schemas))
		for _, s := range realm.Schemas {
			sm[s.Name] = true
		}
		for _, s := range desired.Schemas {
			if !sm[s.Name] {
				return false, fmt.Errorf("schema %q from file %q was not selected %+v, all schemas defined in file must be selected\n", s.Name, file, realm.Schemas)
			}
		}
	}

	// Get changes
	changes, err := a.driver.RealmDiff(realm, desired)
	if err != nil {
		return false, err
	}

	if len(changes) == 0 {
		return true, nil
	}

	// Apply
	if err := a.driver.ApplyChanges(ctx, changes); err != nil {
		return false, err
	}
	return true, nil
}
