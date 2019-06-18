package gcp

import (
	"cloud.google.com/go/datastore"
	"context"
	"google.golang.org/api/option"
)

type GCPDatastoreDatabase struct {
	*datastore.Client
	Namespace string
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
