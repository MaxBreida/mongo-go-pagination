package mongopagination

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PagingQuery struct for holding mongo
// connection, filter needed to apply
// filter data with page, limit, sort key
// and sort value
type PagingQuery struct {
	Collection *mongo.Collection
	Filter     interface{}
	SortField  string
	SortValue  int
	Limit      int64
	Page       int64
}

// PaginatedData struct holds data and
// pagination detail
type PaginatedData struct {
	Data       []bson.Raw     `json:"data"`
	Pagination PaginationData `json:"pagination"`
}

// Find returns two value pagination data with document queried from mongodb and
// error if any error occurs during document query
func (paging *PagingQuery) Find() (paginatedData *PaginatedData, err error) {
	skip := getSkip(paging.Page, paging.Limit)
	opt := &options.FindOptions{
		Skip:  &skip,
		Sort:  bson.D{{paging.SortField, paging.SortValue}},
		Limit: &paging.Limit,
	}
	cursor, err := paging.Collection.Find(context.Background(), paging.Filter, opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var docs []bson.Raw
	for cursor.Next(context.Background()) {
		var document *bson.Raw
		if err := cursor.Decode(&document); err == nil {
			docs = append(docs, *document)
		}
	}
	paginator := Paging(paging)
	result := PaginatedData{
		Pagination: *paginator.PaginationData(),
		Data:       docs,
	}
	return &result, nil
}

// getSkip return calculated skip value for query
func getSkip(page, limit int64) (skip int64) {
	if page > 0 {
		skip = (page - 1) * limit
	} else {
		skip = page
	}
	return
}
