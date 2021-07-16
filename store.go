package main

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
)

var (
	projectID              = "our-axon-319717"
	linkColName            = "link"
	counterColName         = "counter"
	countDocName           = "count"
	linkAlreadyExistsError = errors.New("Link Already exists")
)

type Store struct {
	*Counter
	db *firestore.Client
}

func NewStore() (*Store, error) {
	c, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	s := &Store{db: c}
	if err := s.getCounter(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) getCounter() error {
	countDocRef := s.db.Collection(counterColName).Doc(countDocName)
	c := &Counter{numShards: 5}
	if err := c.initCounter(context.Background(), countDocRef); err != nil {
		return err
	}
	s.Counter = c
	return nil
}

func (s *Store) AddURL(ctx context.Context, l *Link) error {
	if !s.linkExist(ctx, l.Short) {
		_, err := s.db.Collection(linkColName).Doc(l.Short).Set(ctx, l)
		return err
	}
	return linkAlreadyExistsError
}

func (s *Store) GetURL(ctx context.Context, id string) (*Link, error) {
	snap, err := s.db.Collection(linkColName).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	l := mapToLink(snap.Data())
	go s.increaseClick(ctx, id)
	return l, nil
}

func (s *Store) GetURLInfo(ctx context.Context, id string) (*Link, error) {
	snap, err := s.db.Collection(linkColName).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	l := mapToLink(snap.Data())
	return l, nil
}

func (s *Store) UpdateURL(ctx context.Context, id, url string) error {
	_, err := s.db.Collection(linkColName).Doc(id).Update(ctx, []firestore.Update{
		{
			Path:  "URL",
			Value: url,
		},
		{
			Path:  "LastUpdated",
			Value: time.Now().Unix(),
		},
	})
	return err
}

func (s *Store) increaseClick(ctx context.Context, id string) error {
	_, err := s.db.Collection(linkColName).Doc(id).Update(ctx, []firestore.Update{
		{
			Path:  "Click",
			Value: firestore.Increment(1),
		},
	})
	return err
}

func (s *Store) updateLastViewed(ctx context.Context, id string) error {
	_, err := s.db.Collection(linkColName).Doc(id).Update(ctx, []firestore.Update{
		{
			Path:  "LastViewed",
			Value: time.Now().Unix(),
		},
	})
	return err
}

func (s *Store) linkExist(ctx context.Context, id string) bool {
	d, err := s.db.Collection(linkColName).Doc(id).Get(ctx)
	if err != nil {
		return false
	}
	return d.Exists()
}

func mapToLink(m map[string]interface{}) *Link {
	l := &Link{}
	l.URL = m["URL"].(string)
	l.Short = m["Short"].(string)
	l.Click = m["Click"].(int64)
	l.LastUpdated = m["LastUpdated"].(int64)
	l.LastViewed = m["LastViewed"].(int64)
	return l
}
