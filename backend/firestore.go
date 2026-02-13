package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const collectionName = "video_titles"

type FirestoreCache struct {
	client *firestore.Client
}

func NewFirestoreCache(projectID string) (*FirestoreCache, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %v", err)
	}
	return &FirestoreCache{client: client}, nil
}

func (f *FirestoreCache) Get(videoID string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	doc, err := f.client.Collection(collectionName).Doc(videoID).Get(ctx)
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

func (f *FirestoreCache) Set(videoID string, title string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := f.client.Collection(collectionName).Doc(videoID).Set(ctx, map[string]interface{}{
		"title":     title,
		"updatedAt": firestore.ServerTimestamp,
	})
	if err != nil {
		log.Printf("Error writing to Firestore: %v", err)
	}
}

func (f *FirestoreCache) GetMulti(videoIDs []string) map[string]string {
	// Firestore allows getting multiple documents by reference, but the SDK
	// GetAll API takes DocumentRefs.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var refs []*firestore.DocumentRef
	for _, id := range videoIDs {
		refs = append(refs, f.client.Collection(collectionName).Doc(id))
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
			results[videoIDs[i]] = title
		}
	}
	return results
}
