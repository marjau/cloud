package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const testProjectID = "klaboratory"

type DriverTestSuite struct {
	suite.Suite
	d Driver
}

func TestDriverTestSuite(t *testing.T) {
	suite.Run(t, new(DriverTestSuite))
}

func (s *DriverTestSuite) SetupSuite() {
	dsDriver, err := newDriver(testProjectID)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), dsDriver)
	s.d = dsDriver
}

func (s *DriverTestSuite) TearDownTest() {
	//s.T().Logf("--> Closing client connection...")
	//s.d.close()
}

func (s *DriverTestSuite) TestFind() {
	objectType := "Animal"
	i := s.d.Find(nil, objectType, nil, "")
	var entity Animal
	s.T().Logf("Starting interations...")
	for {
		key, err := i.Next(&entity)
		if err != nil {
			break
		}
		s.T().Logf("K: %v", key)
		s.T().Logf("%#v", entity)
	}
}

//func (s *DriverTestSuite) TestUpdate() {
//	k64, _ := strconv.ParseInt("5634161670881280", 10, 64)
//	k := &datastore.Key{
//		Kind: "Animal",
//		ID:   k64,
//	}
//	s.d.Update(k)
//}
