package datastore

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
)

var instance *driver

type Driver interface {
	Find(ancestorId *datastore.Key, objectType string, filter *DataFilter, sort string) *datastore.Iterator
	FindIds(ancestorId *datastore.Key, objectType string, filter *DataFilter, sort string) ([]string, error)
	Create(key *datastore.Key, object interface{}) (string, error)
	Delete(key *datastore.Key) error
	Update(key *datastore.Key, data interface{}) error
	close()
}

type driver struct {
	client *datastore.Client
}

func newDriver(projectId string) (Driver, error) {
	if instance == nil {
		instance = new(driver)
		c, err := datastore.NewClient(context.Background(), projectId)
		if err != nil {
			return nil, err
		}
		instance.client = c
	}
	return instance, nil
}

func (d *driver) close() {
	err := d.client.Close()
	if err != nil {
		fmt.Printf("ERROR: datastore client close failure: %v\n", err)
		// logger.Log().Err(err).Msg("datastore client close failure")
	}
}

func (d *driver) FindIds(ancestor *datastore.Key, objectType string, filter *DataFilter, sort string) (keys []string, err error) {
	q := datastore.NewQuery(objectType)

	if ancestor != nil {
		q = q.Ancestor(ancestor)
	}

	if filter != nil {
		f := *filter
		q = q.FilterField(f.GetField(), f.GetCondition(), f.GetValue())
	}

	if sort != "" {
		q = q.Order(sort)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resultKeys, err := d.client.GetAll(ctx, q, nil)
	if err != nil {
		return nil, fmt.Errorf("driver.Find can't GetAll: %v", err)
	}

	for _, key := range resultKeys {
		keys = append(keys, key.Encode())
	}
	return
}

func (d *driver) Find(ancestor *datastore.Key, objectType string, filter *DataFilter, sort string) *datastore.Iterator {
	q := datastore.NewQuery(objectType)

	if ancestor != nil {
		q = q.Ancestor(ancestor)
	}

	if sort != "" {
		q = q.Order(sort)
	}

	if filter != nil {
		f := *filter
		//logger.Debug(fmt.Sprintf("Find(%v.%v %v %v)", objectType, f.GetField(), f.GetCondition(), f.GetValue()))
		q = q.FilterField(f.GetField(), f.GetCondition(), f.GetValue())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = cancel

	return d.client.Run(ctx, q)
}

func (d *driver) Create(key *datastore.Key, object interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newKey, err := d.client.Put(ctx, key, object)
	if err != nil {
		return "", err
	}

	return newKey.Encode(), nil
}

func (d *driver) Delete(key *datastore.Key) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := d.client.Delete(ctx, key)
	if err != nil {
		return err
	}

	return nil
}

func (d *driver) Update(key *datastore.Key, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := d.client.Put(ctx, key, data)
	if err != nil {
		return err
	}

	return nil
}
