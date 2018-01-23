package golangNeo4jBoltDriver

import (
	"fmt"
	"net/url"
)

// RoutingTable gathers addresses of different members in the causal cluster
type RoutingTable struct {
	ttl     int64
	routers []*url.URL
	readers []*url.URL
	writers []*url.URL
}

// Creates a new routing table from a valid connection
func NewRoutingTable(c *boltConn) (*RoutingTable, error) {
	data, _, _, err := c.QueryNeoAll("CALL dbms.cluster.routing.getRoutingTable({})", nil)
	if err != nil {
		return nil, err
	}
	r := &RoutingTable{}
	r.ttl = data[0][0].(int64)
	// XXX i know this is ugly, dont need you to tell me
	// data[0][1] underlying data is an interface{}. Inside this interface there is a []interface{} var.
	members := data[0][1].([]interface{})
	var routers []*url.URL
	var readers []*url.URL
	var writers []*url.URL
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
		//		addresses := addr
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
		r.readers = readers
		r.writers = writers
		r.routers = routers
	}
	return r, nil
}

func (r *RoutingTable) dropAddress(dropAddr *url.URL, role string) error {
	switch role {
	case "WRITE":
		for i, addr := range r.writers {
			if addr == dropAddr {
				r.writers = append(r.writers[:i], r.writers[i+1:]...)
				return nil
			}
		}
	default:
	}
	return fmt.Errorf("Cannot delete '%v' from routing table since it's not there", dropAddr)

}
