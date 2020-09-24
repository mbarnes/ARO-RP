// Code generated by github.com/jim-minter/go-cosmosdb, DO NOT EDIT.

package cosmosdb

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	pkg "github.com/Azure/ARO-RP/pkg/api"
)

type monitorDocumentClient struct {
	*databaseClient
	path string
}

// MonitorDocumentClient is a monitorDocument client
type MonitorDocumentClient interface {
	Create(context.Context, string, *pkg.MonitorDocument, *Options) (*pkg.MonitorDocument, error)
	List(*Options) MonitorDocumentIterator
	ListAll(context.Context, *Options) (*pkg.MonitorDocuments, error)
	Get(context.Context, string, string, *Options) (*pkg.MonitorDocument, error)
	Replace(context.Context, string, *pkg.MonitorDocument, *Options) (*pkg.MonitorDocument, error)
	Delete(context.Context, string, *pkg.MonitorDocument, *Options) error
	Query(string, *Query, *Options) MonitorDocumentRawIterator
	QueryAll(context.Context, string, *Query, *Options) (*pkg.MonitorDocuments, error)
	ChangeFeed(*Options) MonitorDocumentIterator
}

type monitorDocumentChangeFeedIterator struct {
	*monitorDocumentClient
	continuation string
	options      *Options
}

type monitorDocumentListIterator struct {
	*monitorDocumentClient
	continuation string
	done         bool
	options      *Options
}

type monitorDocumentQueryIterator struct {
	*monitorDocumentClient
	partitionkey string
	query        *Query
	continuation string
	done         bool
	options      *Options
}

// MonitorDocumentIterator is a monitorDocument iterator
type MonitorDocumentIterator interface {
	Next(context.Context, int) (*pkg.MonitorDocuments, error)
	Continuation() string
}

// MonitorDocumentRawIterator is a monitorDocument raw iterator
type MonitorDocumentRawIterator interface {
	MonitorDocumentIterator
	NextRaw(context.Context, int, interface{}) error
}

// NewMonitorDocumentClient returns a new monitorDocument client
func NewMonitorDocumentClient(collc CollectionClient, collid string) MonitorDocumentClient {
	return &monitorDocumentClient{
		databaseClient: collc.(*collectionClient).databaseClient,
		path:           collc.(*collectionClient).path + "/colls/" + collid,
	}
}

func (c *monitorDocumentClient) all(ctx context.Context, i MonitorDocumentIterator) (*pkg.MonitorDocuments, error) {
	allmonitorDocuments := &pkg.MonitorDocuments{}

	for {
		monitorDocuments, err := i.Next(ctx, -1)
		if err != nil {
			return nil, err
		}
		if monitorDocuments == nil {
			break
		}

		allmonitorDocuments.Count += monitorDocuments.Count
		allmonitorDocuments.ResourceID = monitorDocuments.ResourceID
		allmonitorDocuments.MonitorDocuments = append(allmonitorDocuments.MonitorDocuments, monitorDocuments.MonitorDocuments...)
	}

	return allmonitorDocuments, nil
}

func (c *monitorDocumentClient) Create(ctx context.Context, partitionkey string, newmonitorDocument *pkg.MonitorDocument, options *Options) (monitorDocument *pkg.MonitorDocument, err error) {
	headers := http.Header{}
	headers.Set("X-Ms-Documentdb-Partitionkey", `["`+partitionkey+`"]`)

	if options == nil {
		options = &Options{}
	}
	options.NoETag = true

	err = c.setOptions(options, newmonitorDocument, headers)
	if err != nil {
		return
	}

	err = c.do(ctx, http.MethodPost, c.path+"/docs", "docs", c.path, http.StatusCreated, &newmonitorDocument, &monitorDocument, headers)
	return
}

func (c *monitorDocumentClient) List(options *Options) MonitorDocumentIterator {
	continuation := ""
	if options != nil {
		continuation = options.Continuation
	}

	return &monitorDocumentListIterator{monitorDocumentClient: c, options: options, continuation: continuation}
}

func (c *monitorDocumentClient) ListAll(ctx context.Context, options *Options) (*pkg.MonitorDocuments, error) {
	return c.all(ctx, c.List(options))
}

func (c *monitorDocumentClient) Get(ctx context.Context, partitionkey, monitorDocumentid string, options *Options) (monitorDocument *pkg.MonitorDocument, err error) {
	headers := http.Header{}
	headers.Set("X-Ms-Documentdb-Partitionkey", `["`+partitionkey+`"]`)

	err = c.setOptions(options, nil, headers)
	if err != nil {
		return
	}

	err = c.do(ctx, http.MethodGet, c.path+"/docs/"+monitorDocumentid, "docs", c.path+"/docs/"+monitorDocumentid, http.StatusOK, nil, &monitorDocument, headers)
	return
}

func (c *monitorDocumentClient) Replace(ctx context.Context, partitionkey string, newmonitorDocument *pkg.MonitorDocument, options *Options) (monitorDocument *pkg.MonitorDocument, err error) {
	headers := http.Header{}
	headers.Set("X-Ms-Documentdb-Partitionkey", `["`+partitionkey+`"]`)

	err = c.setOptions(options, newmonitorDocument, headers)
	if err != nil {
		return
	}

	err = c.do(ctx, http.MethodPut, c.path+"/docs/"+newmonitorDocument.ID, "docs", c.path+"/docs/"+newmonitorDocument.ID, http.StatusOK, &newmonitorDocument, &monitorDocument, headers)
	return
}

func (c *monitorDocumentClient) Delete(ctx context.Context, partitionkey string, monitorDocument *pkg.MonitorDocument, options *Options) (err error) {
	headers := http.Header{}
	headers.Set("X-Ms-Documentdb-Partitionkey", `["`+partitionkey+`"]`)

	err = c.setOptions(options, monitorDocument, headers)
	if err != nil {
		return
	}

	err = c.do(ctx, http.MethodDelete, c.path+"/docs/"+monitorDocument.ID, "docs", c.path+"/docs/"+monitorDocument.ID, http.StatusNoContent, nil, nil, headers)
	return
}

func (c *monitorDocumentClient) Query(partitionkey string, query *Query, options *Options) MonitorDocumentRawIterator {
	continuation := ""
	if options != nil {
		continuation = options.Continuation
	}

	return &monitorDocumentQueryIterator{monitorDocumentClient: c, partitionkey: partitionkey, query: query, options: options, continuation: continuation}
}

func (c *monitorDocumentClient) QueryAll(ctx context.Context, partitionkey string, query *Query, options *Options) (*pkg.MonitorDocuments, error) {
	return c.all(ctx, c.Query(partitionkey, query, options))
}

func (c *monitorDocumentClient) ChangeFeed(options *Options) MonitorDocumentIterator {
	continuation := ""
	if options != nil {
		continuation = options.Continuation
	}

	return &monitorDocumentChangeFeedIterator{monitorDocumentClient: c, options: options, continuation: continuation}
}

func (c *monitorDocumentClient) setOptions(options *Options, monitorDocument *pkg.MonitorDocument, headers http.Header) error {
	if options == nil {
		return nil
	}

	if monitorDocument != nil && !options.NoETag {
		if monitorDocument.ETag == "" {
			return ErrETagRequired
		}
		headers.Set("If-Match", monitorDocument.ETag)
	}
	if len(options.PreTriggers) > 0 {
		headers.Set("X-Ms-Documentdb-Pre-Trigger-Include", strings.Join(options.PreTriggers, ","))
	}
	if len(options.PostTriggers) > 0 {
		headers.Set("X-Ms-Documentdb-Post-Trigger-Include", strings.Join(options.PostTriggers, ","))
	}
	if len(options.PartitionKeyRangeID) > 0 {
		headers.Set("X-Ms-Documentdb-PartitionKeyRangeID", options.PartitionKeyRangeID)
	}

	return nil
}

func (i *monitorDocumentChangeFeedIterator) Next(ctx context.Context, maxItemCount int) (monitorDocuments *pkg.MonitorDocuments, err error) {
	headers := http.Header{}
	headers.Set("A-IM", "Incremental feed")

	headers.Set("X-Ms-Max-Item-Count", strconv.Itoa(maxItemCount))
	if i.continuation != "" {
		headers.Set("If-None-Match", i.continuation)
	}

	err = i.setOptions(i.options, nil, headers)
	if err != nil {
		return
	}

	err = i.do(ctx, http.MethodGet, i.path+"/docs", "docs", i.path, http.StatusOK, nil, &monitorDocuments, headers)
	if IsErrorStatusCode(err, http.StatusNotModified) {
		err = nil
	}
	if err != nil {
		return
	}

	i.continuation = headers.Get("Etag")

	return
}

func (i *monitorDocumentChangeFeedIterator) Continuation() string {
	return i.continuation
}

func (i *monitorDocumentListIterator) Next(ctx context.Context, maxItemCount int) (monitorDocuments *pkg.MonitorDocuments, err error) {
	if i.done {
		return
	}

	headers := http.Header{}
	headers.Set("X-Ms-Max-Item-Count", strconv.Itoa(maxItemCount))
	if i.continuation != "" {
		headers.Set("X-Ms-Continuation", i.continuation)
	}

	err = i.setOptions(i.options, nil, headers)
	if err != nil {
		return
	}

	err = i.do(ctx, http.MethodGet, i.path+"/docs", "docs", i.path, http.StatusOK, nil, &monitorDocuments, headers)
	if err != nil {
		return
	}

	i.continuation = headers.Get("X-Ms-Continuation")
	i.done = i.continuation == ""

	return
}

func (i *monitorDocumentListIterator) Continuation() string {
	return i.continuation
}

func (i *monitorDocumentQueryIterator) Next(ctx context.Context, maxItemCount int) (monitorDocuments *pkg.MonitorDocuments, err error) {
	err = i.NextRaw(ctx, maxItemCount, &monitorDocuments)
	return
}

func (i *monitorDocumentQueryIterator) NextRaw(ctx context.Context, maxItemCount int, raw interface{}) (err error) {
	if i.done {
		return
	}

	headers := http.Header{}
	headers.Set("X-Ms-Max-Item-Count", strconv.Itoa(maxItemCount))
	headers.Set("X-Ms-Documentdb-Isquery", "True")
	headers.Set("Content-Type", "application/query+json")
	if i.partitionkey != "" {
		headers.Set("X-Ms-Documentdb-Partitionkey", `["`+i.partitionkey+`"]`)
	} else {
		headers.Set("X-Ms-Documentdb-Query-Enablecrosspartition", "True")
	}
	if i.continuation != "" {
		headers.Set("X-Ms-Continuation", i.continuation)
	}

	err = i.setOptions(i.options, nil, headers)
	if err != nil {
		return
	}

	err = i.do(ctx, http.MethodPost, i.path+"/docs", "docs", i.path, http.StatusOK, &i.query, &raw, headers)
	if err != nil {
		return
	}

	i.continuation = headers.Get("X-Ms-Continuation")
	i.done = i.continuation == ""

	return
}

func (i *monitorDocumentQueryIterator) Continuation() string {
	return i.continuation
}
