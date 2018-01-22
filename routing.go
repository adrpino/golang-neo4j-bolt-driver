package golangNeo4jBoltDriver

import (
	"url"
)

// RoutingTable gathers addresses of different members in the causal cluster
type RoutingTable struct {
	routers []*url.URL
	readers []*url.URL
	writers []*url.URL
}
