package golangNeo4jBoltDriver

import (
	"fmt"
	"reflect"
)

// Result represents a result from a query that returns no data
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
	Metadata() map[string]interface{}
}

// TODO: Would sql/driver.RowsAffected render this useless?
type boltResult struct {
	metadata map[string]interface{}
}

func newResult(metadata map[string]interface{}) boltResult {
	return boltResult{metadata: metadata}
}

// Returns the response metadata from the bolt success message
func (r boltResult) Metadata() map[string]interface{} {
	return r.metadata
}

// LastInsertId gets the last inserted id. This will always return -1.
func (r boltResult) LastInsertId() (int64, error) {
	// TODO: Is this possible? -
	// 	I think we would need to parse the query to get the number of parameters
	return -1, nil
}

// RowsAffected returns the number of nodes+rels created/deleted.  For reasons of limitations
// on the API, we cannot tell how many nodes+rels were updated, only how many properties were
// updated.  If this changes in the future, number updated will be added to the output of this
// interface.
func (r boltResult) RowsAffected() (int64, error) {
	stats, ok := r.metadata["stats"].(map[string]interface{})
	if !ok {
		return -1, fmt.Errorf("Unrecognized type for stats metadata: %#v", r.metadata)
	}

	var rowsAffected int64
	nodesCreated, ok := stats["nodes-created"]
	if ok {
		switch nodesCreated.(type) {
		case int, int8, int16, int32, int64:
			rowsAffected += reflect.ValueOf(nodesCreated).Int()
		default:
			return -1, fmt.Errorf("Unrecognized type for nodes created: %#v Metadata: %#v", nodesCreated, r.metadata)
		}
	}

	relsCreated, ok := stats["rel-created"]
	if ok {
		switch relsCreated.(type) {
		case int, int8, int16, int32, int64:
			rowsAffected += reflect.ValueOf(relsCreated).Int()
		default:
			return -1, fmt.Errorf("Unrecognized type for nodes created: %#v Metadata: %#v", relsCreated, r.metadata)
		}
	}

	nodesDeleted, ok := stats["nodes-deleted"]
	if ok {
		switch nodesDeleted.(type) {
		case int, int8, int16, int32, int64:
			rowsAffected += reflect.ValueOf(nodesDeleted).Int()
		default:
			return -1, fmt.Errorf("Unrecognized type for nodes created: %#v Metadata: %#v", nodesDeleted, r.metadata)
		}
	}

	relsDeleted, ok := stats["rel-deleted"]
	if ok {
		switch relsDeleted.(type) {
		case int, int8, int16, int32, int64:
			rowsAffected += reflect.ValueOf(relsDeleted).Int()
		default:
			return -1, fmt.Errorf("Unrecognized type for nodes created: %#v Metadata: %#v", relsDeleted, r.metadata)
		}
	}

	return rowsAffected, nil
}
