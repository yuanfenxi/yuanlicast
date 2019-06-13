package yfxcast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"github.com/juju/errors"
)

// Although there are many Elasticsearch clients with Go, I still want to implement one by myself.
// Because we only need some very simple usages.
type Client struct {
	Addr string
	c    *http.Client
	AppSecret string
}

type PusherMessage struct {
	Name     string
	Data     string
}



type DBEvent struct {
	Action string
	Table string
	Schema string
	ID string
	Parent string
	Data map[string]interface{}

}



func NewClient(addr string,secret string ) *Client {
	c := new(Client)
	c.Addr = addr
	c.c = &http.Client{}
	c.AppSecret = secret
	return c
}

type ResponseItem struct {
	ID      string                 `json:"_id"`
	Index   string                 `json:"_index"`
	Type    string                 `json:"_type"`
	Version int                    `json:"_version"`
	Found   bool                   `json:"found"`
	Source  map[string]interface{} `json:"_source"`
}

type Response struct {
	Code int
	ResponseItem
}

// See http://www.elasticsearch.org/guide/en/elasticsearch/guide/current/bulk.html
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionInsert = "insert"
	ActionIndex = "index"
)

type BulkRequest struct {
	Action string
	Index  string
	Type   string
	ID     string
	Parent string
	Data   map[string]interface{}
}

func (r *BulkRequest) bulk(buf *bytes.Buffer) error {

	var dbEvent DBEvent = DBEvent{}

	if len(r.ID)>0 {
		dbEvent.ID = r.ID
	}

	if len(r.Index)>0 {
		dbEvent.Schema = r.Index
	}

	if len(r.Type)>0 {
		dbEvent.Table = r.Type
	}



	if len(r.Parent) > 0 {
		dbEvent.Parent = r.Parent
	}
	dbEvent.Action = r.Action
	dbEvent.Data = r.Data


	data, err := json.Marshal(dbEvent)
	if err != nil {
		return errors.Trace(err)
	}

	var pushMessage PusherMessage
	pushMessage.Data = string(data)
	pushMessage.Name = "dbEvent"

	jsonData ,errorJson := json.Marshal(pushMessage)
	if errorJson!=nil {
		return errors.Trace(errorJson)
	}

	buf.Write(jsonData)
	buf.WriteByte('\n')
	return nil
}




func (c *Client) Do(method string, url string, body map[string]interface{}) (*Response, error) {
	bodyData, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Trace(err)
	}

	buf := bytes.NewBuffer(bodyData)

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, errors.Trace(err)
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, errors.Trace(err)
	}

	ret := new(Response)
	ret.Code = resp.StatusCode

	if resp.ContentLength > 0 {
		d := json.NewDecoder(resp.Body)
		err = d.Decode(&ret.ResponseItem)
	}

	resp.Body.Close()

	return ret, err
}

func (c *Client) ExecutePostQuery(buf *bytes.Buffer)(error){

	url_to_post :=  c.Addr
	params:=url.Values{}
	// params.Set("hello","fdsfs")
	params.Set("data",buf.String())

	resp, err := http.PostForm(url_to_post, params)
	if err != nil {
		return  errors.Trace(err)
	}



	defer resp.Body.Close()

	if 200!=resp.StatusCode {
		if resp.ContentLength > 0 {
			responseBody := make([]byte,resp.ContentLength)
			_,errorOfRead := resp.Body.Read(responseBody)
			if (errorOfRead!=nil){
				return errors.New("can't read response from server")
			}
			return errors.BadRequestf("bad request: %s",string(responseBody))
		}
		return errors.New("empty response")
	}

	if resp.ContentLength > 0 {
	}
	resp.Body.Close()
	return nil

}

func (c *Client) DoBulk( items []*BulkRequest) ( error) {
	var buf bytes.Buffer

	for _, item := range items {
		if err := item.bulk(&buf); err != nil {
			return  errors.Trace(err)
		}
	}

	return c.ExecutePostQuery(&buf)

}

func (c *Client) CreateMapping(index string, docType string, mapping map[string]interface{}) error {
	reqUrl := fmt.Sprintf("%s/%s", c.Addr,
		url.QueryEscape(index))

	r, err := c.Do("HEAD", reqUrl, nil)
	if err != nil {
		return errors.Trace(err)
	}

	// index doesn't exist, create index first
	if r.Code != http.StatusOK {
		_, err = c.Do("POST", reqUrl, nil)

		if err != nil {
			return errors.Trace(err)
		}
	}

	reqUrl = fmt.Sprintf("%s/%s/%s/_mapping", c.Addr,
		url.QueryEscape(index),
		url.QueryEscape(docType))

	_, err = c.Do("POST", reqUrl, mapping)
	return errors.Trace(err)
}

func (c *Client) DeleteIndex(index string) error {
	reqUrl := fmt.Sprintf("%s/%s", c.Addr,
		url.QueryEscape(index))

	r, err := c.Do("DELETE", reqUrl, nil)
	if err != nil {
		return errors.Trace(err)
	}

	if r.Code == http.StatusOK || r.Code == http.StatusNotFound {
		return nil
	} else {
		return errors.Errorf("Error: %s, code: %d", http.StatusText(r.Code), r.Code)
	}
}

func (c *Client) Get(index string, docType string, id string) (*Response, error) {
	reqUrl := fmt.Sprintf("%s/%s/%s/%s", c.Addr,
		url.QueryEscape(index),
		url.QueryEscape(docType),
		url.QueryEscape(id))

	return c.Do("GET", reqUrl, nil)
}

// Can use Update to create or update the data
func (c *Client) Update(index string, docType string, id string, data map[string]interface{}) error {
	reqUrl := fmt.Sprintf("%s/%s/%s/%s", c.Addr,
		url.QueryEscape(index),
		url.QueryEscape(docType),
		url.QueryEscape(id))

	r, err := c.Do("PUT", reqUrl, data)
	if err != nil {
		return errors.Trace(err)
	}

	if r.Code == http.StatusOK || r.Code == http.StatusCreated {
		return nil
	} else {
		return errors.Errorf("Error: %s, code: %d", http.StatusText(r.Code), r.Code)
	}
}

func (c *Client) Exists(index string, docType string, id string) (bool, error) {
	reqUrl := fmt.Sprintf("%s/%s/%s/%s", c.Addr,
		url.QueryEscape(index),
		url.QueryEscape(docType),
		url.QueryEscape(id))

	r, err := c.Do("HEAD", reqUrl, nil)
	if err != nil {
		return false, err
	}

	return r.Code == http.StatusOK, nil
}

func (c *Client) Delete(index string, docType string, id string) error {
	reqUrl := fmt.Sprintf("%s/%s/%s/%s", c.Addr,
		url.QueryEscape(index),
		url.QueryEscape(docType),
		url.QueryEscape(id))

	r, err := c.Do("DELETE", reqUrl, nil)
	if err != nil {
		return errors.Trace(err)
	}

	if r.Code == http.StatusOK || r.Code == http.StatusNotFound {
		return nil
	} else {
		return errors.Errorf("Error: %s, code: %d", http.StatusText(r.Code), r.Code)
	}
}

// only support parent in 'Bulk' related apis
func (c *Client) Bulk(items []*BulkRequest) ( error) {
	//reqUrl := fmt.Sprintf("%s/_bulk", c.Addr)

	return c.DoBulk(items)
}

func (c *Client) IndexBulk(index string, items []*BulkRequest) ( error) {
	//reqUrl := fmt.Sprintf("%s/%s/_bulk", c.Addr, url.QueryEscape(index))

	return c.DoBulk( items)
}

func (c *Client) IndexTypeBulk(index string, docType string, items []*BulkRequest) ( error) {
	//reqUrl := fmt.Sprintf("%s/%s/%s/_bulk", c.Addr,
	//	url.QueryEscape(index),
	//	url.QueryEscape(docType))

	return c.DoBulk(items)
}
