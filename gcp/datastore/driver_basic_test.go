package datastore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const kLabProjectIDTest = "klaboratory"

type DriverBasicTestSuite struct {
	suite.Suite
	d DriverBasic
}

func TestDriverBasicTestSuite(t *testing.T) {
	suite.Run(t, new(DriverBasicTestSuite))
}

func (s *DriverBasicTestSuite) SetupSuite() {
	dsDriver, err := NewDriverBasic(kLabProjectIDTest)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), dsDriver)
	s.d = dsDriver
	s.T().Logf("Datastore driver loaded")
}

func (s *DriverBasicTestSuite) TearDownTest() {
	//s.d.Close()
}

func (s *DriverBasicTestSuite) TestGet() {
	a, err := s.d.Get()
	s.T().Logf("ERROR: %v", err)
	s.T().Logf("Animal: %#v", a)

	//assert.Nil(s.T(), err)
	//assert.NotNil(s.T(), a)
}
func (s *DriverBasicTestSuite) TestGetAll() {}
func (s *DriverBasicTestSuite) TestPut() {

}
