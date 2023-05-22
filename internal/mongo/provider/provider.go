package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/usa4ev/gotracer/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// provider performs queries to mongodb 
type provider struct {
	conn *mongo.Client
}

func New(c *mongo.Client)*provider{
	return &provider{conn: c}
}

// Count qeries mongodb to get hourly request rate to each endpoint.
// Returns endpoint URI as a map key.
func (pp *provider) Count(date time.Time) (map[string][]model.Entry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	coll := pp.conn.Database("test_db").Collection("random_coll")

	endDate := date.AddDate(0,0,1)

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{
		{{"$match", bson.M{
					"Time": bson.M{
						"$gte": date,
						"$lt": endDate,
					},
					"Message": "request started",
					"path": bson.M{"$in": bson.A{"/track","/fastest","/slowest"}},
				},
		}},
		{{"$group", bson.M{
				"_id": bson.M{
						"Date": bson.M{
								"$dateTrunc": bson.M{
									"date": "$Time", 
									"unit": "hour", 
									"binSize": 1,
								},
							},
						"path": "$path"},
				"count": bson.M{"$count": bson.D{}},
			},
		}},
		{{"$sort", bson.M{"Date": -1},}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute mongodb query: %v", err)
	}

	type MongoResult struct {
		ID struct{
			Date  time.Time `bson:"date" json:"date"`
			URI   string    `bson:"path" json:"path"`
		}`bson:"_id" json:"_id"`
		Count int       `bson:"count" json:"count"`
	}

	var results []MongoResult
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode query results: %v", err)
	}

	res := make(map[string][]model.Entry)
	for _, result := range results {
		if res[result.ID.URI] == nil {
			res[result.ID.URI] = make([]model.Entry, 0)
		}

		res[result.ID.URI] = append(res[result.ID.URI], 
			model.Entry{Time: result.ID.Date.Local(), Count: result.Count})
	}

	return res, nil
}