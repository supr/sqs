package sqs_test

import (
    "crypto/md5"
    "hash"
    "fmt"
    "launchpad.net/gocheck"
    "launchpad.net/goamz/aws"
    "launchpad.net/goamz/sqs"
)

var _ = gocheck.Suite(&S{})

type S struct {
    HTTPSuite
    sqs *sqs.SQS
}

func (s *S) SetUpSuite(c *gocheck.C) {
    s.HTTPSuite.SetUpSuite(c)
    auth := aws.Auth{"abc", "123"}
    s.sqs = sqs.New(auth, aws.Region{SQSEndpoint: testServer.URL})
}

func (s *S) TestCreateQueue(c *gocheck.C) {
    testServer.PrepareResponse(200, nil, TestCreateQueueXmlOK)

    resp, err := s.sqs.CreateQueue("testQueue")
    req := testServer.WaitRequest()

    c.Assert(req.Method, gocheck.Equals, "GET")
    c.Assert(req.URL.Path, gocheck.Equals, "/")
    c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

    c.Assert(resp.Url, gocheck.Equals, "http://sqs.us-east-1.amazonaws.com/123456789012/testQueue")
    c.Assert(err, gocheck.IsNil)
}

func (s *S) TestListQueues(c *gocheck.C) {
    testServer.PrepareResponse(200, nil, TestListQueuesXmlOK)

    resp, err := s.sqs.ListQueues("")
    req := testServer.WaitRequest()

    c.Assert(req.Method, gocheck.Equals, "GET")
    c.Assert(req.URL.Path, gocheck.Equals, "/")
    c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

    c.Assert(len(resp.QueueUrl), gocheck.Not(gocheck.Equals), 0)
    c.Assert(resp.QueueUrl[0], gocheck.Equals, "http://sqs.us-east-1.amazonaws.com/123456789012/testQueue")
    c.Assert(resp.ResponseMetadata.RequestId, gocheck.Equals, "725275ae-0b9b-4762-b238-436d7c65a1ac")
    c.Assert(err, gocheck.IsNil)
}

func (s *S) TestDeleteQueue(c *gocheck.C) {
    testServer.PrepareResponse(200, nil, TestDeleteQueueXmlOK)

    q := &sqs.Queue{s.sqs, testServer.URL + "/123456789012/testQueue/"}
    resp,err := q.Delete()
    req := testServer.WaitRequest()

    c.Assert(req.Method, gocheck.Equals, "GET")
    c.Assert(req.URL.Path, gocheck.Equals, "/123456789012/testQueue/")
    c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

    c.Assert(resp.ResponseMetadata.RequestId, gocheck.Equals, "6fde8d1e-52cd-4581-8cd9-c512f4c64223")
    c.Assert(err, gocheck.IsNil)
}

func (s *S) TestSendMessage(c *gocheck.C) {
    testServer.PrepareResponse(200, nil, TestSendMessageXmlOK)
    
    q := &sqs.Queue{s.sqs, testServer.URL + "/123456789012/testQueue/"}
    resp,err := q.SendMessage("This is a test message")
    req := testServer.WaitRequest()

    c.Assert(req.Method, gocheck.Equals, "GET")
    c.Assert(req.URL.Path, gocheck.Equals, "/123456789012/testQueue/")
    c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

    msg := "This is a test message"
    var h hash.Hash = md5.New()
    h.Write([]byte(msg))
    c.Assert(resp.MD5, gocheck.Equals, fmt.Sprintf("%x", h.Sum()))
    c.Assert(resp.Id, gocheck.Equals, "5fea7756-0ea4-451a-a703-a558b933e274")
    c.Assert(err, gocheck.IsNil)
}
