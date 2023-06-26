package datastore

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
)

const ctxTimeOut = 10 * time.Second

var dbInstance *driverBasic

var _ DriverBasic = driverBasic{}

type DriverBasic interface {
	GetAll()
	Get() (*Animal, error)
	Put(a Animal) error
	Close()
}

func NewDriverBasic(projectId string) (DriverBasic, error) {
	if dbInstance == nil {
		dbInstance = new(driverBasic)
		c, err := datastore.NewClient(context.Background(), projectId)
		if err != nil {
			return nil, err
		}
		dbInstance.client = c
	}
	return dbInstance, nil
}

type Animal struct {
	Name     string
	Legs     int
	Sound    string
	FoodType string
}

type driverBasic struct {
	client *datastore.Client
}

func (d driverBasic) Get() (*Animal, error) {
	k := datastore.NameKey("Animal", "5634161670881280", nil)
	a := &Animal{}
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeOut)
	defer cancel()

	if err := d.client.Get(ctx, k, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (d driverBasic) GetAll() {}

func (d driverBasic) Put(_ Animal) error {
	// TODO: implement me
	return nil
}

func (d driverBasic) Close() {
	err := d.client.Close()
	if err != nil {
		fmt.Printf("ERROR: datastore client close failure: %v\n", err)
		// TODO: implement logger
		// logger.Log().Err(err).Msg("datastore client close failure")
	}
}
