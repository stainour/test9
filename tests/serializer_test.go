package tests

import (
	"fmt"
	"github.com/stainour/test9/list"
	"os"
	"testing"
)

var headNode *list.Node

const nodeCount = 1000000
const fileName = "test.txt"

func init() {
	headNode = &list.Node{
		Data: "0_data",
	}

	currentListNode := headNode

	for i := 1; i < nodeCount; i++ {
		next := &list.Node{
			Data: fmt.Sprintf("%d_data", i),
			Prev: currentListNode,
		}

		if i%3 == 0 {
			next.Rand = nil
		} else {
			next.Rand = currentListNode
		}

		currentListNode.Next = next
		currentListNode = next
	}
}

func TestSerializer_Serialize_Deserialize(t *testing.T) {
	serialize(t)
	node := deserialize(t)

	currentNode := node
	deserializedCurrentNode := node

	for currentNode.Next != nil {

		if deserializedCurrentNode.Data != currentNode.Data {
			validate(deserializedCurrentNode, currentNode, t)
		}

		if deserializedCurrentNode.Prev != nil {
			validate(deserializedCurrentNode.Prev, currentNode.Prev, t)
		}

		if deserializedCurrentNode.Next != nil {
			validate(deserializedCurrentNode.Next, currentNode.Next, t)
		}

		currentNode = currentNode.Next
		deserializedCurrentNode = deserializedCurrentNode.Next
	}
}

func validate(deserializedCurrentNode *list.Node, currentNode *list.Node, t *testing.T) {
	if deserializedCurrentNode.Data != currentNode.Data {
		t.Errorf("%s != %s", deserializedCurrentNode.Data, currentNode.Data)
	}
}

func deserialize(t *testing.T) *list.Node {
	reader, err := os.Open(fileName)
	if err != nil {
		t.Error(err)
	}

	node, err := list.Deserialize(reader)
	if err != nil {
		t.Error(err)
	}

	if err = reader.Close(); err != nil {
		t.Error(err)
	}
	return node
}

func serialize(t *testing.T) {
	writer, err := os.OpenFile(fileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Error(err)
	}
	err = list.Serialize(headNode, writer)

	if err != nil {
		t.Error(err)
	}

	if err = writer.Close(); err != nil {
		t.Error(err)
	}
}
