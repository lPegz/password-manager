package storage

// Then import the package
import (
	"context"
	"fmt"
	immudb "github.com/codenotary/immudb/pkg/client"
	"google.golang.org/grpc/metadata"
	"log"
)

type ImmutableDBStorageStrategy struct {
	client immudb.ImmuClient
	ctx    context.Context
	md     metadata.MD
}

func NewImmuDbStorageStrategy() (*ImmutableDBStorageStrategy, error) {
	client, err := immudb.NewImmuClient(immudb.DefaultOptions())
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	// login with default username and password and storing a token
	lr, err := client.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		log.Fatal("error connecting to immudb", err)
	}
	// set up an authenticated context that will be required in future operations
	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	return &ImmutableDBStorageStrategy{client: client, md: md, ctx: ctx}, nil
}

func (i ImmutableDBStorageStrategy) StorageSave(passwordEntry PasswordEntry, output bool) {
	passwordKey := fmt.Sprintf("%v-%v", passwordEntry.Tag, passwordEntry.Username)
	_, err := i.client.Set(i.ctx, []byte(passwordKey), []byte(passwordEntry.Password))
	if err != nil {
		log.Fatal("error adding key ", err)
	}
	if output {
		log.Printf("Password: %v\n", passwordEntry.Password)
	}
}

func (i ImmutableDBStorageStrategy) StorageGet(tag string, username string, output bool) string {
	passwordKey := fmt.Sprintf("%v-%v", tag, username)
	entry, err := i.client.Get(i.ctx, []byte(passwordKey))
	if err != nil {
		log.Fatal(err)
	}
	return string(entry.Value)
}
