package database

import (
	"context"
	"time"

	"flood-iot-kku-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SensorRepository struct {
	collection *mongo.Collection
}

func NewSensorRepository(db *MongoDB, collectionName string) *SensorRepository {
	return &SensorRepository{
		collection: db.GetCollection(collectionName),
	}
}

// Create inserts a new sensor data record
func (r *SensorRepository) Create(ctx context.Context, data *models.SensorData) error {
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	
	// Determine status based on water level
	if data.WaterLevel > 5.0 {
		data.Status = "critical"
	} else if data.WaterLevel > 3.0 {
		data.Status = "warning"
	} else {
		data.Status = "normal"
	}

	_, err := r.collection.InsertOne(ctx, data)
	return err
}

// GetByID retrieves a sensor data record by ID
func (r *SensorRepository) GetByID(ctx context.Context, id string) (*models.SensorData, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var data models.SensorData
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetAll retrieves all sensor data with pagination
func (r *SensorRepository) GetAll(ctx context.Context, page, limit int64) ([]models.SensorData, int64, error) {
	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Set up options for pagination and sorting
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Get total count
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	// Find documents
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []models.SensorData
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// GetBySensorID retrieves sensor data by sensor ID
func (r *SensorRepository) GetBySensorID(ctx context.Context, sensorID string, page, limit int64) ([]models.SensorData, int64, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{"sensor_id": sensorID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []models.SensorData
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// GetByStatus retrieves sensor data by status
func (r *SensorRepository) GetByStatus(ctx context.Context, status string, page, limit int64) ([]models.SensorData, int64, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{"status": status}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []models.SensorData
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// Update updates a sensor data record
func (r *SensorRepository) Update(ctx context.Context, id string, data *models.SensorData) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	data.UpdatedAt = time.Now()
	
	// Update status based on water level
	if data.WaterLevel > 5.0 {
		data.Status = "critical"
	} else if data.WaterLevel > 3.0 {
		data.Status = "warning"
	} else {
		data.Status = "normal"
	}

	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, data)
	return err
}

// Delete removes a sensor data record
func (r *SensorRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// GetLatestBySensorID gets the latest reading for a specific sensor
func (r *SensorRepository) GetLatestBySensorID(ctx context.Context, sensorID string) (*models.SensorData, error) {
	opts := options.FindOne().
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	var data models.SensorData
	err := r.collection.FindOne(ctx, bson.M{"sensor_id": sensorID}, opts).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
