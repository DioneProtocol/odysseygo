package indexer

import (
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	f func(reply interface{}) error
}

func (mc *mockClient) SendRequest(_ string, _ interface{}, reply interface{}) error {
	return mc.f(reply)
}

func TestIndexClient(t *testing.T) {
	assert := assert.New(t)
	client, err := NewClient("http://localhost:9650", CChain, IndexTypeBlocks, time.Minute)
	assert.NoError(err)

	// Test GetIndex
	client.EndpointRequester = &mockClient{
		f: func(reply interface{}) error {
			*(reply.(*GetIndexResponse)) = GetIndexResponse{Index: 5}
			return nil
		},
	}
	index, err := client.GetIndex(&GetIndexArgs{ContainerID: ids.Empty, Encoding: formatting.Hex})
	assert.NoError(err)
	assert.EqualValues(5, index.Index)

	// Test GetLastAccepted
	client.EndpointRequester = &mockClient{
		f: func(reply interface{}) error {
			*(reply.(*FormattedContainer)) = FormattedContainer{ID: "hi"}
			return nil
		},
	}
	container, err := client.GetLastAccepted(&GetLastAcceptedArgs{Encoding: formatting.Hex})
	assert.NoError(err)
	assert.EqualValues("hi", container.ID)

	// Test GetContainerRange
	client.EndpointRequester = &mockClient{
		f: func(reply interface{}) error {
			*(reply.(*[]FormattedContainer)) = []FormattedContainer{{ID: "yeet"}}
			return nil
		},
	}
	containers, err := client.GetContainerRange(&GetContainerRange{StartIndex: 1, NumToFetch: 10, Encoding: formatting.Hex})
	assert.NoError(err)
	assert.Len(containers, 1)
	assert.EqualValues("yeet", containers[0].ID)

	// Test IsAccepted
	client.EndpointRequester = &mockClient{
		f: func(reply interface{}) error {
			*(reply.(*bool)) = true
			return nil
		},
	}
	isAccepted, err := client.IsAccepted(&GetIndexArgs{ContainerID: ids.Empty, Encoding: formatting.Hex})
	assert.NoError(err)
	assert.True(isAccepted)
}

func TestIndexClientInvalidIndex(t *testing.T) {
	assert := assert.New(t)
	_, err := NewClient("http://localhost:9650", PChain, IndexTypeTransactions, time.Minute)
	assert.Error(err)
	_, err = NewClient("http://localhost:9650", PChain, IndexTypeVertices, time.Minute)
	assert.Error(err)
	_, err = NewClient("http://localhost:9650", CChain, IndexTypeTransactions, time.Minute)
	assert.Error(err)
	_, err = NewClient("http://localhost:9650", CChain, IndexTypeVertices, time.Minute)
	assert.Error(err)
	_, err = NewClient("http://localhost:9650", XChain, IndexTypeBlocks, time.Minute)
	assert.Error(err)
	_, err = NewClient("http://localhost:9650", IndexedChain(111), IndexTypeBlocks, time.Minute)
	assert.Error(err)
	_, err = NewClient("http://localhost:9650", XChain, IndexType(11), time.Minute)
	assert.Error(err)
}
