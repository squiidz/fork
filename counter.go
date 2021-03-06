package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Counter is a collection of documents (shards)
// to realize counter with high frequency.
type Counter struct {
	dRef      *firestore.DocumentRef
	numShards int
}

// Shard is a single counter, which is used in a group
// of other shards within Counter.
type Shard struct {
	Count int
}

// initCounter creates a given number of shards as
// subcollection of specified document.
func (c *Counter) initCounter(ctx context.Context, docRef *firestore.DocumentRef) error {
	_, err := docRef.Collection("shards").Doc("0").Get(ctx)
	if err == nil {
		c.dRef = docRef
		return nil
	}
	// Initialize each shard with count=0
	for num := 0; num < c.numShards; num++ {
		shard := Shard{0}

		if _, err := docRef.Collection("shards").Doc(strconv.Itoa(num)).Set(ctx, shard); err != nil {
			return fmt.Errorf("Set: %v", err)
		}
	}
	return nil
}

// incrementCounter increments a randomly picked shard.
func (c *Counter) incrementCounter(ctx context.Context) (*firestore.WriteResult, error) {
	docID := strconv.Itoa(rand.Intn(c.numShards))

	shardRef := c.dRef.Collection("shards").Doc(docID)
	return shardRef.Update(ctx, []firestore.Update{
		{Path: "Count", Value: firestore.Increment(1)},
	})
}

// getCount returns a total count across all shards.
func (c *Counter) getCount(ctx context.Context) (int64, error) {
	var total int64
	shards := c.dRef.Collection("shards").Documents(ctx)
	for {
		doc, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("Next: %v", err)
		}

		vTotal := doc.Data()["Count"]
		shardCount, ok := vTotal.(int64)
		if !ok {
			return 0, fmt.Errorf("firestore: invalid dataType %T, want int64", vTotal)
		}
		total += shardCount
	}
	return total, nil
}
