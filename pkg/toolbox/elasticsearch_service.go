package toolbox

import (
	"context"
	"sort"

	elastic "github.com/olivere/elastic/v7"
)

func getESService(esConfig ESConfig) (res *ESService, err error) {
	key := "elasticsearch-" + esConfig.Url
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s *ESService
		s, err = CreateESService(esConfig)
		if err != nil {
			return
		}
		_, err = s.GetClient()
		if err != nil {
			return
		}
		res = s
		return
	})
	if err != nil {
		return
	}
	res = service.(*ESService)
	return
}

func CreateESService(esConfig ESConfig) (*ESService, error) {
	service := &ESService{
		url: esConfig.Url,
	}
	err := service.init()
	return service, err
}

//ESService 注册处理器在线信息等
type ESService struct {
	url         string
	lastUseTime int64
}

func (this_ *ESService) init() error {
	var err error
	return err
}
func (this_ *ESService) GetClient() (client *elastic.Client, err error) {
	defer func() {
		this_.lastUseTime = GetNowTime()
	}()
	client, err = elastic.NewClient(
		elastic.SetURL(this_.url),
		//docker
		elastic.SetSniff(false),
	)
	return
}

func (this_ *ESService) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *ESService) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *ESService) Stop() {
}

func (this_ *ESService) DeleteIndex(indexName string) (err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	_, err = client.DeleteIndex(indexName).Do(context.Background())
	if err != nil {
		return
	}
	return
}

func (this_ *ESService) CreateIndex(indexName string, bodyJSON map[string]interface{}) (err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	_, err = client.CreateIndex(indexName).BodyJson(bodyJSON).Do(context.Background())
	if err != nil {
		return
	}
	return
}

func (this_ *ESService) IndexNames() (res []string, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	res, err = client.IndexNames()
	if err != nil {
		return
	}

	sort.Strings(res)
	return
}

func (this_ *ESService) GetMapping(indexName string) (res interface{}, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	mappingMap, err := client.GetMapping().Index(indexName).Do(context.Background())
	if err != nil {
		return
	}
	for key, value := range mappingMap {
		if key == indexName {
			res = value
		}
	}
	return
}

func (this_ *ESService) PutMapping(indexName string, bodyJSON map[string]interface{}) (err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	_, err = client.PutMapping().Index(indexName).BodyJson(bodyJSON).Do(context.Background())
	if err != nil {
		return
	}
	return
}

func (this_ *ESService) SetFieldType(indexName string, fieldName string, fieldType string) (err error) {
	bodyJSON := map[string]interface{}{}
	bodyJSON["properties"] = map[string]interface{}{
		fieldName: map[string]interface{}{
			"type": fieldType,
		},
	}
	err = this_.PutMapping(indexName, bodyJSON)
	if err != nil {
		return
	}
	return
}

func (this_ *ESService) Search(indexName string, pageIndex int, pageSize int) (res *elastic.SearchResult, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()

	doer := client.Search(indexName)
	query := elastic.NewBoolQuery()
	res, err = doer.Query(query).Size(pageSize).From((pageIndex - 1) * pageSize).Do(context.Background())
	if err != nil {
		return
	}

	return
}

func (this_ *ESService) Insert(indexName string, id string, doc interface{}) (res *elastic.IndexResponse, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	doer := client.Index()
	res, err = doer.Index(indexName).Id(id).BodyJson(doc).Refresh("wait_for").Do(context.Background())
	if err != nil {
		return
	}

	return
}

func (this_ *ESService) Update(indexName string, id string, doc interface{}) (res *elastic.UpdateResponse, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()

	doer := client.Update()
	res, err = doer.Index(indexName).Id(id).Doc(doc).Refresh("wait_for").Do(context.Background())
	if err != nil {
		return
	}

	return
}

func (this_ *ESService) Delete(indexName string, id string) (res *elastic.DeleteResponse, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()

	doer := client.Delete()
	res, err = doer.Index(indexName).Id(id).Refresh("wait_for").Do(context.Background())
	if err != nil {
		return
	}

	return
}

func (this_ *ESService) Reindex(sourceIndexName string, toIndexName string) (res *elastic.BulkIndexByScrollResponse, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()

	doer := client.Reindex()
	res, err = doer.Source(elastic.NewReindexSource().Index(sourceIndexName)).DestinationIndex(toIndexName).Refresh("true").Do(context.Background())
	if err != nil {
		return
	}

	return
}

func (this_ *ESService) Scroll(indexName string, scrollId string, pageSize int) (res *elastic.SearchResult, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()

	doer := client.Scroll(indexName)
	query := elastic.NewBoolQuery()
	res, err = doer.Query(query).Size(pageSize).ScrollId(scrollId).Do(context.Background())
	if err != nil {
		return
	}

	return
}
