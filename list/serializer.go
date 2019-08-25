package list

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

var endian = binary.LittleEndian

const (
	nilRandomId = int32(-1)
)

func Serialize(head *Node, writer io.Writer) error {
	nodeIds := map[*Node]int32{}
	nodeIdCounter := int32(0)
	setNodeId := func(node *Node) {
		if _, ok := nodeIds[node]; !ok {
			nodeIds[node] = nodeIdCounter
			nodeIdCounter++
		}
	}

	getNodeId := func(node *Node) (int32, error) {
		if node == nil {
			return nilRandomId, nil
		}

		if index, ok := nodeIds[node]; ok {
			return index, nil
		}
		return math.MinInt32, errors.New("bad link list structure")
	}

	var nodesToSerialize [] *Node = nil
	currentNode := head

	for currentNode != nil {
		setNodeId(currentNode)
		nodesToSerialize = append(nodesToSerialize, currentNode)
		currentNode = currentNode.Next
	}

	err := writeInt32(writer, nodeIdCounter)
	if err != nil {
		return err
	}

	currentNode = head

	for currentNode != nil {
		nodeId, err := getNodeId(currentNode.Rand)
		if err != nil {
			return err
		}

		err = writeInt32(writer, nodeId)
		if err != nil {
			return err
		}

		err = writeString(writer, currentNode.Data)
		if err != nil {
			return err
		}

		currentNode = currentNode.Next
	}

	return nil
}

func writeString(writer io.Writer, string string) error {
	stringBytes := []byte(string)

	err := writeInt32(writer, int32(len(stringBytes)))
	if err != nil {
		return err
	}

	writtenBytes, err := writer.Write(stringBytes)

	if err != nil {
		return err
	}

	if len(stringBytes) != writtenBytes {
		return errors.New("error writing to writer")
	}
	return nil
}

func writeInt32(file io.Writer, int int32) error {
	return binary.Write(file, endian, int)
}

func readInt32(reader io.Reader) (int32, error) {
	var number int32
	err := binary.Read(reader, endian, &number)
	return number, err
}

func readString(reader io.Reader) (string, error) {
	stringLength, err := readInt32(reader)

	if err != nil {
		return "", err
	}

	bytes := make([]byte, stringLength)
	readBytes, err := reader.Read(bytes)

	if err != nil {
		return "", err
	}
	if stringLength != int32(readBytes) {
		return "", errors.New("error reading from reader")
	}

	return string(bytes), nil
}

func Deserialize(file io.Reader) (*Node, error) {

	nodeCount, err := readInt32(file)

	if err != nil {
		return nil, err
	}

	if nodeCount > 0 {
		nodes := make([]*Node, nodeCount)
		randomNodeIds := make([]int32, nodeCount)
		for i := int32(0); i < nodeCount; i++ {
			randomNodeId, err := readInt32(file)
			if err != nil {
				return nil, err
			}
			randomNodeIds[i] = randomNodeId

			data, err := readString(file)
			if err != nil {
				return nil, err
			}

			nodes[i] = &Node{
				Data: data,
			}

		}
		return restoreNodeRandomReferences(nodes, randomNodeIds), nil

	} else {
		return &Node{}, nil
	}

}

func restoreNodeRandomReferences(nodes []*Node, randomNodeIds []int32) *Node {
	for i := range nodes {
		node := nodes[i]
		randomNodeId := randomNodeIds[i]

		if randomNodeId != nilRandomId {
			node.Rand = nodes[randomNodeId]
		}

		prevIndex := i - 1
		nextNode := i + 1
		if prevIndex >= 0 {
			node.Prev = nodes[prevIndex]
		}

		if nextNode < len(nodes) {
			node.Next = nodes[nextNode]
		}
	}
	return nodes[0]
}
