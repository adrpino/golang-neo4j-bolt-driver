package golangNeo4jBoltDriver

import (
	"github.com/adrpino/golang-neo4j-bolt-driver/errors"
	"math/rand"
	"net/url"
)

// RoutingTable gathers addresses of different members in the causal cluster
type routingTable struct {
	ttl     int64
	routers []*url.URL
	readers []*url.URL
	writers []*url.URL
}

// Creates a new routing table from a valid connection
func NewRoutingTable(c *boltConn) (*routingTable, error) {
	data, _, _, err := c.QueryNeoAll("CALL dbms.cluster.routing.getRoutingTable({})", nil)
	if err != nil {
		return nil, err
	}
	r := &routingTable{}
	r.ttl = data[0][0].(int64)
	// XXX i know this is ugly, dont need you to tell me
	// data[0][1] underlying data is an interface{}. Inside this interface there is a []interface{} var.
	members := data[0][1].([]interface{})
	var routers, readers, writers []*url.URL
	for _, mem := range members {
		m, ok := mem.(map[string]interface{})
		if !ok {
			continue
		}
		role, ok := m["role"]
		if !ok {
			continue
		}
		addr, ok := m["addresses"]
		if !ok {
			continue
		}
		// addr is interface
		addresses := addr.([]interface{})
		for _, add := range addresses {
			a := add.(string)
			parsedAddr, err := url.Parse(a)
			if err != nil {
				return nil, err
			}
			switch role {
			case "WRITE":
				writers = append(writers, parsedAddr)
			case "READ":
				readers = append(readers, parsedAddr)
			default:
				routers = append(routers, parsedAddr)
			}
		}
	}
	if len(readers) == 0 {
		return nil, errors.New("Error creating routing table, no readers received from server.")
	}
	if len(routers) == 0 {
		return nil, errors.New("Error creating routing table, no routers received from server.")
	}
	r.readers = readers
	r.writers = writers
	r.routers = routers
	return r, nil
}

func (r *routingTable) dropAddress(dropAddr *url.URL, role string) error {
	switch role {
	case "WRITE":
		for i, addr := range r.writers {
			if addr == dropAddr {
				r.writers = append(r.writers[:i], r.writers[i+1:]...)
				return nil
			}
		}
	case "READER":
		//TODO
	default:
	}
	return errors.New("Cannot delete '%v' from routing table since it's not there", dropAddr)

}

func (r *routingTable) Reader() (*url.URL, int) {
	// Select a random index
	ind := rand.Intn(len(r.readers))
	return r.readers[ind], ind
}

func (r *routingTable) Writer() (*url.URL, int) {
	// Select a random index
	ind := rand.Intn(len(r.writers))
	return r.writers[ind], ind
}

// Gets a random reader
func (r *routingTable) PopReader() *url.URL {
	res, ind := r.Reader()
	r.readers[ind] = r.readers[len(r.readers)-1]
	r.readers[len(r.readers)-1] = nil
	r.readers = r.readers[:len(r.readers)-1]
	return res
}

// writer
func (r *routingTable) PopWriter() *url.URL {
	res, ind := r.Writer()
	r.writers[ind] = r.writers[len(r.writers)-1]
	r.writers[len(r.writers)-1] = nil
	r.writers = r.writers[:len(r.readers)-1]
	return res
}
