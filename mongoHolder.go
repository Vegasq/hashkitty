package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type MongoHolder struct {
	Client     *mongo.Client
	Context    *context.Context
	HDB        *mongo.Collection
	PROCESSEDB *mongo.Collection
}

func (mh *MongoHolder) Close() {
	mh.Client.Disconnect(*mh.Context)
	ctx := *mh.Context
	ctx.Done()
}

func (mh *MongoHolder) Insert(recs []Record) {
	fmt.Println("Inserting")
	start := time.Now()

	inserts := []interface{}{}
	for i := 0; i < len(recs); i++ {
		inserts = append(inserts, recs[i])
	}

	_, err := mh.HDB.InsertMany(*mh.Context, inserts)
	if err != nil {
		log.Fatal(err)
	}
	end := time.Now()

	took := end.Sub(start)

	fmt.Printf("insert took %f sec\n", took.Seconds())

	//fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}

// From AWS documentation
func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := ioutil.ReadFile(caFile)

	if err != nil {
		return tlsConfig, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		return tlsConfig, errors.New("Failed parsing pem file")
	}

	return tlsConfig, nil
}

func NewMongoHolder() *MongoHolder {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	mh := MongoHolder{}
	data, err := os.ReadFile("connect.uri")
	if err != nil {
		panic("Failed to read connection uri")
	}
	ctx := context.Background()

	var client *mongo.Client
	if strings.Contains(string(data), "localhost") == false {
		tlsConfig, err := getCustomTLSConfig("rds-combined-ca-bundle.pem")
		if err != nil {
			log.Fatalf("Failed getting TLS configuration: %v", err)
		}
		client, err := mongo.NewClient(options.Client().ApplyURI(string(data)).SetTLSConfig(tlsConfig))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}

		err = client.Connect(ctx)
		if err != nil {
			log.Fatalf("Failed to connect to cluster: %v", err)
		}
	} else {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(string(data)))
		if err != nil {
			fmt.Println(err)
		}
	}

	mh.Client = client
	mh.Context = &ctx
	mh.HDB = mh.Client.Database("hashesdb").Collection("hashes")
	mh.PROCESSEDB = mh.Client.Database("hashesdb").Collection("processed")

	return &mh
}
