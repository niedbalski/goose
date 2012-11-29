package swift_test

import (
	. "launchpad.net/gocheck"
	"launchpad.net/goose/identity"
	"launchpad.net/goose/testing/httpsuite"
	"launchpad.net/goose/testservices/identityservice"
	"launchpad.net/goose/testservices/swiftservice"
	"net/http"
)

func registerLocalTests(cred *identity.Credentials) {
	Suite(&localLiveSuite{
		LiveTests: LiveTests{
			cred: cred,
		},
	})
}

const (
	baseURL = "/object-store"
)

// localLiveSuite runs tests from LiveTests using a fake
// swift server that runs within the test process itself.
type localLiveSuite struct {
	LiveTests
	// The following attributes are for using testing doubles.
	httpsuite.HTTPSuite
	identityDouble http.Handler
	swiftDouble    http.Handler
}

func (s *localLiveSuite) SetUpSuite(c *C) {
	c.Logf("Using identity service test double")
	c.Logf("Using swift service test double")
	s.HTTPSuite.SetUpSuite(c)
	s.cred.URL = s.Server.URL
	s.identityDouble = identityservice.NewUserPass()
	token := s.identityDouble.(*identityservice.UserPass).AddUser(s.cred.User, s.cred.Secrets)
	s.swiftDouble = swiftservice.New("localhost", baseURL+"/", token)
	ep := identityservice.Endpoint{
		s.Server.URL + baseURL, //admin
		s.Server.URL + baseURL, //internal
		s.Server.URL + baseURL, //public
		s.cred.Region,
	}
	service := identityservice.Service{"swift", "object-store", []identityservice.Endpoint{ep}}
	s.identityDouble.(*identityservice.UserPass).AddService(service)
	s.LiveTests.SetUpSuite(c)
}

func (s *localLiveSuite) TearDownSuite(c *C) {
	s.LiveTests.TearDownSuite(c)
	s.HTTPSuite.TearDownSuite(c)
}

func (s *localLiveSuite) SetUpTest(c *C) {
	s.HTTPSuite.SetUpTest(c)
	s.Mux.Handle(baseURL+"/", s.swiftDouble)
	s.Mux.Handle("/", s.identityDouble)
	s.LiveTests.SetUpTest(c)
}

func (s *localLiveSuite) TearDownTest(c *C) {
	s.LiveTests.TearDownTest(c)
	s.HTTPSuite.TearDownTest(c)
}

// Additional tests to be run against the service double only go here.
