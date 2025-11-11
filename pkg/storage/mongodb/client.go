package mongodb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/KOJIMEISTER/it_russian_stat/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	Collection         *mongo.Collection
	existingVacancyIDS map[string]struct{}
	existingDescHashes *sync.Map
}

func NewMongoStore(ctx context.Context, config *config.MongoDBConfig) (*MongoStore, error) {
	clientOptions := options.Client().ApplyURI(config.URL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("MongoDB connection error: %w", err)
	}

	collection := client.Database(config.Database).Collection(config.Collection)
	return &MongoStore{
		Collection: collection,
	}, nil
}

func (s *MongoStore) LoadExistingData() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.existingVacancyIDS = make(map[string]struct{})
	s.existingDescHashes = &sync.Map{}

	cursor, err := s.Collection.Find(ctx, bson.D{}, options.Find().SetProjection(bson.D{
		{"id", 1},
		{"description_hash", 1},
	}))
	if err != nil {
		return fmt.Errorf("failed to fetch existing vacancies: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc struct {
			ID              string `bson:"id"`
			DescriptionHash string `bson:"description_hash"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return fmt.Errorf("failed to decode document: %w", err)
		}
		s.existingVacancyIDS[doc.ID] = struct{}{}
		if doc.DescriptionHash != "" {
			s.existingDescHashes.Store(doc.DescriptionHash, struct{}{})
		}
	}

	return cursor.Err()
}

func (s *MongoStore) VacancyExists(id string) bool {
	_, exists := s.existingVacancyIDS[id]
	return exists
}

func (s *MongoStore) DescriptionHashExists(hash string) bool {
	_, exists := s.existingDescHashes.Load(hash)
	return exists
}

func (s *MongoStore) AddDescriptionHash(hash string) {
	s.existingDescHashes.Store(hash, struct{}{})
}

func (s *MongoStore) UpsertVacancy(data map[string]interface{}) error {
	filter := bson.M{"id": data["id"]}
	update := bson.M{"$set": data}
	_, err := s.Collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	return err
}
