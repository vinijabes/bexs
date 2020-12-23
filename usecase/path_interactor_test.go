package usecase

import (
	"bexs/infrastructure/repository/file"
	"bexs/interface/interactor"
	"bexs/utils"
	"io"
	"io/ioutil"
	"testing"
)

var buffer io.ReadWriteSeeker

func init() {
	buffer = utils.NewReadWriteSeeker("GRU,BRC,10\nBRC,SCL,5")
}

func TestGetPath(t *testing.T) {
	repo, err := file.NewBufferPathRepository(buffer)
	if err != nil {
		t.Fatal(err)
	}

	usecase, err := NewPathInteractor(repo)
	if err != nil {
		t.Fatal(err)
	}

	path, err := usecase.FindPath(interactor.FindPathRequest{
		Origin: "GRU",
		Dest:   "SCL",
	})

	utils.Assert(t, nil, err)
	utils.Assert(t, 15, int(path.Dist))
	utils.Assert(t, true, path.Connections[0] == "GRU" && path.Connections[1] == "BRC" && path.Connections[2] == "SCL")
}

func TestAddRoute(t *testing.T) {
	repo, err := file.NewBufferPathRepository(buffer)
	if err != nil {
		t.Fatal(err)
	}

	usecase, err := NewPathInteractor(repo)
	if err != nil {
		t.Fatal(err)
	}

	usecase.AddRoute(interactor.AddRouteRequest{
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
