package bbk

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BBKTestSuite struct {
	suite.Suite
	server *httptest.Server
}

func (suite *BBKTestSuite) SetupSuite() {
	fileSystem := http.Dir("./testdata")
	fileServer := http.FileServer(fileSystem)

	suite.server = httptest.NewServer(fileServer)
}

func (suite *BBKTestSuite) TearDownSuite() {
	suite.server.Close()
}

func TestBBKSuite(t *testing.T) {
	suite.Run(t, new(BBKTestSuite))
}
