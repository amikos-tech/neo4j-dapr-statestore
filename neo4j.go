package neo4j

import (
	"context"
	"errors"
	"github.com/dapr/components-contrib/state"
	"github.com/dapr/kit/logger"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"time"
)

type neo4jMetadata struct {
	Host             string
	Username         string
	Password         string
	DatabaseName     string
	OperationTimeout time.Duration
}

type Node struct {
	Key   string      `json:"_id"`
	Value interface{} `json:"value"`
	Etag  string      `json:"_etag"`
	TTL   *time.Time  `json:"_ttl,omitempty"`
}

type Neo4j struct {
	driver   *neo4j.DriverWithContext
	metadata neo4jMetadata
	logger   logger.Logger
}

func (store *Neo4j) Init(metadata state.Metadata) (err error) {
	store.metadata, err = getNeo4jMetadata(metadata)
	if err != nil {
		return err
	}
	ctx := context.Background()
	// URI examples: "neo4j://localhost", "neo4j+s://xxx.databases.neo4j.io"
	dbUri := store.metadata.Host
	dbUser := store.metadata.Username
	dbPassword := store.metadata.Password
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	store.driver = &driver

	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
	return nil
}

func getNeo4jMetadata(metadata state.Metadata) (neo4jMetadata, error) {
	m := neo4jMetadata{
		Host:         metadata.Properties["host"],
		Username:     metadata.Properties["username"],
		Password:     metadata.Properties["password"],
		DatabaseName: metadata.Properties["databaseName"],
	}
	var err error
	if val, ok := metadata.Properties["operationTimeout"]; ok && val != "" {
		m.OperationTimeout, err = time.ParseDuration(val)
		if err != nil {
			return m, errors.New("incorrect operationTimeout field from metadata")
		}
	}
	return m, nil
}

func (store *Neo4j) GetComponentMetadata() map[string]string {
	// Not used with pluggable components...
	return map[string]string{}
}

func (store *Neo4j) Features() []state.Feature {
	// Return a list of features supported by the state store...
}

func (store *Neo4j) Delete(ctx context.Context, req *state.DeleteRequest) error {
	// Delete the requested key from the state store...
}

func (store *Neo4j) Get(ctx context.Context, req *state.GetRequest) (*state.GetResponse, error) {
	// Get the requested key value from the state store, else return an empty response...
}

func (store *Neo4j) Set(ctx context.Context, req *state.SetRequest) error {
	// Set the requested key to the specified value in the state store...
}

func (store *Neo4j) BulkGet(ctx context.Context, req []state.GetRequest) (bool, []state.BulkGetResponse, error) {
	// Get the requested key values from the state store...
}

func (store *Neo4j) BulkDelete(ctx context.Context, req []state.DeleteRequest) error {
	// Delete the requested keys from the state store...
}

func (store *Neo4j) BulkSet(ctx context.Context, req []state.SetRequest) error {
	// Set the requested keys to their specified values in the state store...
}
