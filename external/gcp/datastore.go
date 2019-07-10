package gcp

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/octofoxio/sparkle"
	"google.golang.org/api/option"
	"reflect"
)

type GCPDatastoreCollection struct {
	*GCPDatastoreDatabase
	CollectionName string
}

func (g *GCPDatastoreCollection) FindByID(ctx context.Context, ID string, value interface{}) error {
	k := datastore.NameKey(g.CollectionName, ID, nil)
	k.Namespace = g.Namespace
	err := g.Client.Get(ctx, k, value)
	return err
}

func (g *GCPDatastoreCollection) FindOne(ctx context.Context, filter, value interface{}) error {
	// WIP
	q := datastore.NewQuery(g.CollectionName)
	v := reflect.ValueOf(filter)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		fieldValue := v.FieldByName(fieldName).Interface()
		q = q.Filter(fieldName+"=", fieldValue)
	}
	q = q.Limit(1)
	var results []*struct {
		Hi string
	}
	k, err := g.Client.GetAll(ctx, q, &results)
	fmt.Print(k)
	return err
}

func (g *GCPDatastoreCollection) Save(ctx context.Context, ID string, entity interface{}) error {
	k := datastore.NameKey(g.CollectionName, ID, nil)
	k.Namespace = g.Namespace
	k, err := g.Client.Put(ctx, k, entity)
	return err
}

func (g *GCPDatastoreCollection) DeleteByID(ctx context.Context, ID string) error {
	panic("implement me")
}

type GCPDatastoreDatabase struct {
	*datastore.Client
	Namespace string
}

func (g *GCPDatastoreDatabase) Collection(name string) sparkle.Collection {
	return &GCPDatastoreCollection{
		GCPDatastoreDatabase: g,
		CollectionName:       name,
	}
}

func (g *GCPDatastoreDatabase) FindByID(ctx context.Context, Collection, ID string, value interface{}) error {
	k := datastore.NameKey(Collection, ID, nil)
	k.Namespace = g.Namespace
	err := g.Client.Get(ctx, k, value)
	return err
}

func (g *GCPDatastoreDatabase) Save(ctx context.Context, Collection, ID string, entity interface{}) error {
	k := datastore.NameKey(Collection, ID, nil)
	k.Namespace = g.Namespace
	k, err := g.Client.Put(ctx, k, entity)
	return err
}

func (g *GCPDatastoreDatabase) DeleteByID(ctx context.Context, Collection, ID string) error {
	k := datastore.NameKey(Collection, ID, nil)
	k.Namespace = g.Namespace
	err := g.Client.Delete(ctx, k)
	return err
}

func NewGCPDatastore(projectID string, Namespace string, options ...option.ClientOption) *GCPDatastoreDatabase {
	client, err := datastore.NewClient(context.Background(), projectID, options...)
	if err != nil {
		panic(err)
	}

	return &GCPDatastoreDatabase{
		Client:    client,
		Namespace: Namespace,
	}
}
