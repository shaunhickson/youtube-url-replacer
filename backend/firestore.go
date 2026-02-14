package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/sph/youtube-url-replacer/backend/resolvers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const collectionName = "video_titles"

type FirestoreCache struct {
	client *firestore.Client
}

func NewFirestoreCache(projectID string) (resolvers.Cache, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %v", err)
	}
	return &FirestoreCache{client: client}, nil
}

func hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func (f *FirestoreCache) Get(key string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	doc, err := f.client.Collection(collectionName).Doc(hashKey(key)).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return "", false
	}
	if err != nil {
		log.Printf("Error reading from Firestore: %v", err)
		return "", false
	}

	if title, ok := doc.Data()["title"].(string); ok {
		return title, true
	}
	return "", false
}

func (f *FirestoreCache) Set(key string, title string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := f.client.Collection(collectionName).Doc(hashKey(key)).Set(ctx, map[string]interface{}{
		"title":     title,
		"updatedAt": firestore.ServerTimestamp,
		"original":  key, // Store original key for debugging
	})
	if err != nil {
		log.Printf("Error writing to Firestore: %v", err)
	}
}

func (f *FirestoreCache) GetMulti(keys []string) map[string]string {
	// Firestore allows getting multiple documents by reference, but the SDK
	// GetAll API takes DocumentRefs.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var refs []*firestore.DocumentRef
	for _, key := range keys {
		refs = append(refs, f.client.Collection(collectionName).Doc(hashKey(key)))
	}

	docs, err := f.client.GetAll(ctx, refs)
	if err != nil {
		log.Printf("Error executing GetAll on Firestore: %v", err)
		return map[string]string{}
	}

	results := make(map[string]string)
	for i, doc := range docs {
		if !doc.Exists() {
			continue
		}
		if title, ok := doc.Data()["title"].(string); ok {
			results[keys[i]] = title
		}
	}
	return results
}
