// Пакет для работы с базой данных MongoDB.

package mongo

import (
	"GoNews/internal/model"
	"GoNews/internal/storage"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Константы для обращения к БД и коллекции в mongoDB
const (
	dbName         = "goNews"
	collectionName = "posts"
)

type Storage struct {
	db *mongo.Client
}

// New - конструктор подключения к БД
func New(conn string) (*Storage, error) {
	const operation = "storage.mongodb.new"

	mogngoOpts := options.Client().ApplyURI(conn)
	db, err := mongo.Connect(context.Background(), mogngoOpts)
	if err != nil {
		return nil, fmt.Errorf("%s %s", operation, err)
	}

	if err := db.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("%s %s", operation, err)
	}

	collection := db.Database(dbName).Collection(collectionName)

	//Создание уникального индекса по полю title, для избежания записи уже существующих постов.
	indexUniq := mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: -1}},
		Options: options.Index().SetUnique(true),
	}

	//Создание индекса текстового поиска по полю title.
	indexId := mongo.IndexModel{
		Keys: bson.D{{Key: "title", Value: "text"}},
	}

	_, err = collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{indexUniq, indexId})
	if err != nil {
		return nil, fmt.Errorf("%s %s", operation, err)
	}

	return &Storage{db: db}, nil
}

// AddPost - записывает посты в БД. На вход метода подаётся слайс постов.
// Метод возвращает количество успешно записанных постов и ошибку, отличную от duplicate key error.
func (s *Storage) AddPost(ctx context.Context, posts []model.Post) (int, error) {
	const operation = "GoNews.storage.mongodb.addPost"

	var inputDB []interface{}
	for _, post := range posts {
		bsn := bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "title", Value: post.Title},
			{Key: "content", Value: post.Content},
			{Key: "pubTime", Value: primitive.NewDateTimeFromTime(post.PubTime)},
			{Key: "link", Value: post.Link},
		}
		inputDB = append(inputDB, bsn)
	}

	collection := s.db.Database(dbName).Collection(collectionName)
	res, err := collection.InsertMany(ctx, inputDB)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return len(res.InsertedIDs), fmt.Errorf("%s %s", operation, err)
	}
	return len(res.InsertedIDs), nil
}

// GetPosts - возвращает посты из БД в соответствии с переданными опциями.
// Опции для возвращения постов: запрос на текстовый поиск, лимит числа постов на страницу,
// оффсет для пагинации. Если параметр опции nil, то вернет все посты, отсортированные по
// дате публикации.
func (s *Storage) GetPosts(ctx context.Context, op ...*storage.Options) ([]model.Post, error) {
	const operation = "GoNews.storage.mongodb.getPosts"

	collection := s.db.Database(dbName).Collection(collectionName)
	filter := bson.D{}
	sort := bson.D{{"pubTime", -1}}
	opts := options.Find()

	var query string
	var lim, off int64

	if op[0] != nil {
		query = op[0].SearchQuery
		lim = int64(op[0].Count)
		off = int64(op[0].Offset)
	}

	if query != "" {
		filter = bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: query}}}}
		sort = bson.D{{Key: "score", Value: bson.D{{Key: "$meta", Value: "textScore"}}}}
	}

	opts = opts.SetSort(sort)

	if lim > 0 {
		opts = opts.SetLimit(lim)
	}

	if off > 0 {
		opts = opts.SetSkip(off)
	}

	cur, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	defer cur.Close(ctx)

	var posts []model.Post
	err = cur.All(ctx, &posts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	if len(posts) == 0 {
		return nil, fmt.Errorf("%s: %w", operation, storage.ErrNotFound)
	}

	return posts, nil
}

// CountPosts - метод возвращает число постов, соответствующих условиям поиска.
func (s *Storage) CountPosts(ctx context.Context, op ...*storage.Options) (int64, error) {
	const operation = "GoNews.storage.mongodb.Count"

	filter := bson.D{}
	opts := options.Count().SetHint("_id_")

	if op[0] != nil && op[0].SearchQuery != "" {
		filter = bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: op[0].SearchQuery}}}}
		opts = nil
	}

	collection := s.db.Database(dbName).Collection(collectionName)
	res, err := collection.CountDocuments(ctx, filter, opts)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}
	return res, nil
}

// PostByID - метод возвращает пост по переданному ID.
func (s *Storage) PostByID(ctx context.Context, id string) (model.Post, error) {
	const operation = "storage.mongodb.PostById"
	var post model.Post

	if id == "" {
		return post, fmt.Errorf("%s: %w", operation, "Пустой id")
	}

	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return post, fmt.Errorf("%s: %w", operation, "{storage.ErrIncorrectId}")
	}

	collection := s.db.Database(dbName).Collection(collectionName)
	filter := bson.D{{Key: "_id", Value: obj}}
	res := collection.FindOne(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		return post, fmt.Errorf("%s: %w", operation, storage.ErrNotFound)
	}
	if res.Err() != nil {
		return post, fmt.Errorf("%s: %w", operation, res.Err())
	}

	err = res.Decode(&post)
	if err != nil {
		return post, fmt.Errorf("%s: %w", operation, err)
	}

	return post, nil
}

// Close - обертка для закрытия пула подключений.
func (s *Storage) Close() error {
	return s.db.Disconnect(context.Background())
}
