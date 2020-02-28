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
	SortField  *string
	SortValue  *int
	Limit      int
	Page       int
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
	limit := int64(paging.Limit)
	opt := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}
	if paging.SortField != nil && paging.SortValue != nil {
		opt.Sort = bson.D{{*paging.SortField, *paging.SortValue}}
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
	paginationInfo := *paginator.PaginationData()
	paginationInfo.RecordsOnPage = len(docs)
	result := PaginatedData{
		Pagination: paginationInfo,
		Data:       docs,
	}
	return &result, nil
}

// getSkip return calculated skip value for query
func getSkip(page, limit int) (skip int64) {
	if page > 0 {
		skip = int64((page - 1) * limit)
	} else {
		skip = int64(page)
	}
	return
}
