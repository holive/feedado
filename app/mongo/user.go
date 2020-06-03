package mongo

import (
	"context"
	"strconv"

	"github.com/holive/feedado/app/user"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func (ur *UserRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	resp, err := ur.collection.InsertOne(ctx, u)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a user")
	}

	var newUser user.User

	if err = ur.collection.FindOne(ctx, bson.M{"_id": resp.InsertedID}).Decode(&newUser); err != nil {
		return nil, errors.Wrap(err, "could not find the new feed")
	}

	return &newUser, nil
}

func (ur *UserRepository) Update(ctx context.Context, newUser *user.User) error {
	update, err := bson.Marshal(newUser)
	if err != nil {
		return errors.Wrap(err, "could not marshal bson")
	}

	opts := options.Replace().SetUpsert(false)
	filter := bson.M{"email": bson.M{"$eq": newUser.Email}}

	resp, err := ur.collection.ReplaceOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if resp.MatchedCount != 1 || resp.ModifiedCount != 1 {
		return errors.New("document not found or not updated")
	}

	return nil
}

func (ur *UserRepository) Delete(ctx context.Context, email string) error {
	filter := bson.M{"email": bson.M{"$eq": email}}

	_, err := ur.collection.DeleteOne(ctx, filter)

	return err
}

func (ur *UserRepository) Find(ctx context.Context, email string) (*user.User, error) {
	var u user.User

	filter := bson.M{"email": bson.M{"$eq": email}}

	if err := ur.collection.FindOne(ctx, filter).Decode(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (ur *UserRepository) FindAll(ctx context.Context, limit string, offset string) (*user.SearchResult, error) {
	intLimit, intOffset, err := ur.getLimitOffset(limit, offset)
	if err != nil {
		return &user.SearchResult{}, errors.Wrap(err, "could not get limit or offset")
	}

	findOptions := options.Find().SetLimit(intLimit).SetSkip(intOffset)

	cur, err := ur.collection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		return &user.SearchResult{}, err
	}

	total, err := ur.collection.CountDocuments(ctx, bson.D{{}})
	if err != nil {
		return nil, errors.Wrap(err, "could not count documents")
	}

	results, err := ur.resultFromCursor(ctx, cur)
	if err != nil {
		return &user.SearchResult{}, errors.Wrap(err, "could not get results from cursor")
	}

	return &user.SearchResult{
		Users: results,
		Result: struct {
			Offset int64 `json:"offset"`
			Limit  int64 `json:"limit"`
			Total  int64 `json:"total"`
		}{
			Offset: intOffset,
			Limit:  intLimit,
			Total:  total,
		},
	}, nil
}

func (ur *UserRepository) getLimitOffset(limit string, offset string) (int64, int64, error) {
	if offset == "" {
		offset = "0"
	}

	if limit == "" {
		limit = "24"
	}

	intOffset, err := strconv.Atoi(offset)
	if err != nil {
		return 0, 0, err
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		return 0, 0, err
	}

	return int64(intLimit), int64(intOffset), nil
}

func (ur *UserRepository) resultFromCursor(ctx context.Context, cur *mongo.Cursor) ([]user.User, error) {
	var results []user.User
	for cur.Next(ctx) {
		var elem user.User
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)

	return results, nil
}

func NewUserRepository(conn *Client) *UserRepository {
	return &UserRepository{
		collection: conn.db.Collection(UserCollection),
	}
}
