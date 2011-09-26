package sqs_test

import (
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

    c.Assert(resp.QueueUrl, gocheck.Equals, "http://sqs.us-east-1.amazonaws.com/123456789012/testQueue")
    c.Assert(resp.ResponseMetadata.RequestId, gocheck.Equals, "7a62c49f-347e-4fc4-9331-6e8e7a96aa73")
    c.Assert(err, gocheck.IsNil)
}
