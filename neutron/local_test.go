package neutron_test

import (
	"log"

	gc "gopkg.in/check.v1"

	"gopkg.in/niedbalski/goose.v3/client"
	"gopkg.in/niedbalski/goose.v3/identity"
	"gopkg.in/niedbalski/goose.v3/neutron"
	"gopkg.in/niedbalski/goose.v3/testservices"
	"gopkg.in/niedbalski/goose.v3/testservices/hook"
	"gopkg.in/niedbalski/goose.v3/testservices/openstackservice"
)

func registerLocalTests() {
	gc.Suite(&localLiveSuite{})
}

// localLiveSuite runs tests from LiveTests using a fake
// neutron server that runs within the test process itself.
type localLiveSuite struct {
	LiveTests
	openstack       *openstackservice.Openstack
	noMoreIPs       bool // If true, addFloatingIP will return ErrNoMoreFloatingIPs
	ipLimitExceeded bool // If true, addFloatingIP will return ErrIPLimitExceeded
}

func (s *localLiveSuite) SetUpSuite(c *gc.C) {
	c.Logf("Using identity and neutron service test doubles")

	// Set up an Openstack service.
	s.cred = &identity.Credentials{
		User:       "fred",
		Secrets:    "secret",
		Region:     "some region",
		TenantName: "tenant",
	}
	var logMsg []string
	s.openstack, logMsg = openstackservice.New(s.cred, identity.AuthUserPass, false)
	for _, msg := range logMsg {
		c.Logf(msg)
	}
	s.openstack.UseNeutronNetworking()

	s.openstack.SetupHTTP(nil)
	s.LiveTests.SetUpSuite(c)
}

func (s *localLiveSuite) TearDownSuite(c *gc.C) {
	s.LiveTests.TearDownSuite(c)
	s.openstack.Stop()
}

func (s *localLiveSuite) setupClient(c *gc.C, logger *log.Logger) *neutron.Client {
	client := client.NewClient(s.cred, identity.AuthUserPass, logger)
	//client.SetVersionDiscoveryEnabled(false)
	return neutron.New(client)
}

func (s *localLiveSuite) addFloatingIPHook(sc hook.ServiceControl) hook.ControlProcessor {
	return func(sc hook.ServiceControl, args ...interface{}) error {
		if s.noMoreIPs {
			return testservices.NoMoreFloatingIPs
		} else if s.ipLimitExceeded {
			return testservices.IPLimitExceeded
		}
		return nil
	}
}
