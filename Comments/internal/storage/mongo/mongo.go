// Пакет для работы с базой данных MongoDB для сервиса комментариев.

package mongo

import (
	"Comments/internal/model"
	"Comments/internal/storage"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Константы для обращения к БД и коллекции в mongoDB
const (
	dbName         = "goComments"
	collectionName = "comments"
)

type Storage struct {
	db *mongo.Client
}

// New - конструктор подключения к БД
func New(conn string) (*Storage, error) {
	const operation = "goComments.storage.mongodb.new"

	mongoOpts := options.Client().ApplyURI(conn)
	db, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, fmt.Errorf("%s %s", operation, err)
	}

	if err := db.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("%s %s", operation, err)
	}

	collection := db.Database(dbName).Collection(collectionName)

	//Создание индекса для поля NewsID для ускорения выдачи комментариев.
	indexId := mongo.IndexModel{
		Keys: bson.D{
			{
				Key: "NewsID", Value: -1,
			},
		},
	}

	_, err = collection.Indexes().CreateOne(context.Background(), indexId)

	if err != nil {
		return nil, fmt.Errorf("%s %s", operation, err)
	}

	return &Storage{db: db}, nil
}

// AddComment - метод записывает комментарий в БД.
func (s *Storage) AddComment(ctx context.Context, comment model.Comment) (string, error) {
	const operation = "goComments.storage.mongodb.addComment"

	if t, err := primitive.ObjectIDFromHex(comment.NewsID); err != nil {
		fmt.Println(t, err)
		return "", fmt.Errorf("%s: %w", operation, storage.ErrIncorrectPostID)
	}

	id := primitive.NewObjectID()
	bsn := bson.D{
		{Key: "_id", Value: id},
		{Key: "ParentID", Value: comment.ParentID},
		{Key: "NewsID", Value: comment.NewsID},
		{Key: "PubTime", Value: primitive.NewDateTimeFromTime(time.Now())},
		{Key: "Content", Value: comment.Content},
		{Key: "Childs", Value: bson.A{}},
	}

	collection := s.db.Database(dbName).Collection(collectionName)

	//Проверка на существование родительского комментария
	if comment.ParentID != "" {
		id, err := primitive.ObjectIDFromHex(comment.ParentID)
		if err != nil {
			return "", fmt.Errorf("%s: %w", operation, storage.ErrIncorrectParentID)
		}
		filter := bson.D{
			{Key: "_id", Value: id},
		}
		res := collection.FindOne(ctx, filter)
		if res.Err() != nil {
			if res.Err() == mongo.ErrNoDocuments {
				return "", fmt.Errorf("%s: %w", operation, storage.ErrParentNotFound)
			}
			return "", fmt.Errorf("%s: %w", operation, res.Err())
		}

	}

	_, err := collection.InsertOne(ctx, bsn)
	if err != nil {
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	return id.Hex(), nil
}

// Comments - метод возвращает все деревья комментариев по переданному ID поста,
// отсортированные по дате создания.
func (s *Storage) Comments(ctx context.Context, news string) ([]model.Comment, error) {
	const operation = "goComments.storage.mongodb.AddComment"

	if news == "" {
		return nil, storage.ErrIncorrectPostID
	}

	if _, err := primitive.ObjectIDFromHex(news); err != nil {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrIncorrectPostID)
	}

	var comments []model.Comment
	collection := s.db.Database(dbName).Collection(collectionName)
	opts := options.Find().SetSort(bson.D{{Key: "PubTime", Value: -1}})
	filter := bson.D{{Key: "NewsID", Value: news}}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	err = cursor.All(ctx, &comments)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	if len(comments) == 0 {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrNoComments)
	}
	return comments, nil
}
