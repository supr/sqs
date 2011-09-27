package sqs_test

var TestCreateQueueXmlOK = `
<CreateQueueResponse>
  <CreateQueueResult>
    <QueueUrl>http://sqs.us-east-1.amazonaws.com/123456789012/testQueue</QueueUrl>
  </CreateQueueResult>
  <ResponseMetadata>
    <RequestId>7a62c49f-347e-4fc4-9331-6e8e7a96aa73</RequestId>
  </ResponseMetadata>
</CreateQueueResponse>
`

var TestListQueuesXmlOK = `
<ListQueuesResponse>
  <ListQueuesResult>
    <QueueUrl>http://sqs.us-east-1.amazonaws.com/123456789012/testQueue</QueueUrl>
  </ListQueuesResult>
  <ResponseMetadata>
    <RequestId>725275ae-0b9b-4762-b238-436d7c65a1ac</RequestId>
  </ResponseMetadata>
</ListQueuesResponse>
`

var TestDeleteQueueXmlOK = `
<DeleteQueueResponse>
  <ResponseMetadata>
    <RequestId>6fde8d1e-52cd-4581-8cd9-c512f4c64223</RequestId>
  </ResponseMetadata>
</DeleteQueueResponse>
`
