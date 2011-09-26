//
// goamz - Go packages to interact with the Amazon Web Services.
//
//   https://wiki.ubuntu.com/goamz
//
// Copyright (c) 2011 Memeo Inc.
//
// Written by Prudhvi Krishna Surapaneni <me@prudhvi.net>
//
package sqs

import (
	"http"
	"xml"
	"url"
	"os"
	"time"
	"strconv"
	"launchpad.net/goamz/aws"
)

// The SQS type encapsulates operation with an SQS region.
type SQS struct {
	aws.Auth
	aws.Region
	private byte // Reserve the right of using private data.
}

func New(auth aws.Auth, region aws.Region) *SQS {
	return &SQS{auth, region, 0}
}

type Queue struct {
	*SQS
	Url string
}

type CreateQueueResponse struct {
	QueueUrl string `xml:"CreateQueueResult>QueueUrl"`
	ResponseMetadata
}

type ListQueuesResponse struct {
	QueueUrl []string `xml:"ListQueuesResult>QueueUrl"`
	ResponseMetadata
}

type ResponseMetadata struct {
	RequestId string
	BoxUsage  float64
}

type Error struct {
	StatusCode int
	Code       string
	Message    string
	RequestId  string
}

func (err *Error) String() string {
	return err.Message
}

type xmlErrors struct {
	RequestId string
	Errors    []Error `xml:"Errors>Error"`
}

func (s *SQS) CreateQueue(queueName string) (*Queue, os.Error) {
	return s.CreateQueueWithTimeout(queueName, 30)
}

func (s *SQS) CreateQueueWithTimeout(queueName string, timeout int) (q *Queue, err os.Error) {
	resp, err := s.newQueue(queueName, timeout)
	if err != nil {
		return nil, err
	}
	q = &Queue{s, resp.QueueUrl}
	return
}

func (s *SQS) newQueue(queueName string, timeout int) (resp *CreateQueueResponse, err os.Error) {
	resp = &CreateQueueResponse{}
	params := makeParams("CreateQueue")

	params["QueueName"] = queueName
	params["DefaultVisibilityTimeout"] = strconv.Itoa(timeout)

	err = s.query("", params, resp)
	return
}

func (s *SQS) ListQueues(QueueNamePrefix string) (resp *ListQueuesResponse, err os.Error) {
	resp = &ListQueuesResponse{}
	params := makeParams("ListQueues")

	if QueueNamePrefix != "" {
		params["QueueNamePrefix"] = QueueNamePrefix
	}

	err = s.query("", params, resp)
	return
}

func (s *SQS) query(queueUrl string, params map[string]string, resp interface{}) os.Error {
	params["Timestamp"] = time.UTC().Format(time.RFC3339)
	var url_ *url.URL
	var err os.Error
	var path string
	if queueUrl != "" {
		url_, err = url.Parse(queueUrl)
		path = "/" + queueUrl[len(s.Region.SQSEndpoint):]
	} else {
		url_, err = url.Parse(s.Region.SQSEndpoint)
		path = "/"
	}
	if err != nil {
		return err
	}

	//url_, err := url.Parse(s.Region.SQSEndpoint)
	//if err != nil {
	//	return err
	//}

	sign(s.Auth, "GET", path, params, url_.Host)
	url_.RawQuery = multimap(params).Encode()
	r, err := http.Get(url_.String())
	if err != nil {
		return err
	}
	defer r.Body.Close()

	//dump, _ := http.DumpResponse(r, true)
	//println("DUMP:\n", string(dump))
	//return nil

	if r.StatusCode != 200 {
		return buildError(r)
	}
	err = xml.Unmarshal(r.Body, resp)
	return err
}

func buildError(r *http.Response) os.Error {
	errors := xmlErrors{}
	xml.Unmarshal(r.Body, &errors)
	var err Error
	if len(errors.Errors) > 0 {
		err = errors.Errors[0]
	}
	err.RequestId = errors.RequestId
	err.StatusCode = r.StatusCode
	if err.Message == "" {
		err.Message = r.Status
	}
	return &err
}

func makeParams(action string) map[string]string {
	params := make(map[string]string)
	params["Action"] = action
	return params
}

func multimap(p map[string]string) url.Values {
	q := make(url.Values, len(p))
	for k, v := range p {
		q[k] = []string{v}
	}
	return q
}
