package file

import (
	"bexs/domain/model"
	"bexs/utils"
	"io"
	"io/ioutil"
	"testing"
)

var buffer io.ReadWriteSeeker

func init() {
	buffer = utils.NewReadWriteSeeker("GRU,BRC,10\nBRC,SCL,5")
}

func TestGetGraph(t *testing.T) {
	repo, err := NewBufferPathRepository(buffer)
	if err != nil {
		t.Fatal(err)
	}

	graph, err := repo.GetGraph()
	if err != nil {
		t.Fatal(err)
	}

	utils.Assert(t, 3, len(graph.Vertexes))
}

func TestAddRoute(t *testing.T) {
	repo, err := NewBufferPathRepository(buffer)
	utils.Assert(t, nil, err)

	repo.AddRoute(model.Route{
		Origin: "GRU",
		Dest:   "CDG",
		Price:  75,
	})

	buffer.Seek(0, io.SeekStart)

	expected := "GRU,BRC,10\nBRC,SCL,5\nGRU,CDG,75"
	data, err := ioutil.ReadAll(buffer)
	utils.Assert(t, nil, err)

	if string(data) != expected {
		utils.Assert(t, expected, string(data))
	}
}
