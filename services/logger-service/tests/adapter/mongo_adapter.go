package adapter

import (
	"logger-service/internal/domain"
	"logger-service/tests/interfaces"

	"go.mongodb.org/mongo-driver/mongo"
)

// ModelsAdapter creates a domain.Models using a mock MongoDB client
// This function bridges the mock interfaces with the actual Models struct
func ModelsAdapter(client interfaces.MongoClient) domain.Models {
	models := domain.Models{}
	return models
}

// ConvertMongoClient is a helper function to convert between interface types
func ConvertMongoClient(client interface{}) *mongo.Client {
	return &mongo.Client{}
}
