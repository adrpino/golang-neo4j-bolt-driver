package main

import (
	"fmt"
	bolt "github.com/adrpino/golang-neo4j-bolt-driver"
)

func main() {
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo("bolt+router://neo4j:changeme@localhost:7687")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Start by creating a node
	result, err := conn.ExecNeo("CREATE (n:NODE {foo: {foo}, bar: {bar}})", map[string]interface{}{"foo": 1, "bar": 2.2})
	if err != nil {
		panic(err)
	}
	numResult, _ := result.RowsAffected()
	fmt.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 1

	// Lets get the node
	data, rowsMetadata, _, _ := conn.QueryNeoAll("MATCH (n:NODE) RETURN n.foo, n.bar", nil)
	fmt.Printf("COLUMNS: %#v\n", rowsMetadata["fields"].([]interface{}))    // COLUMNS: n.foo,n.bar
	fmt.Printf("FIELDS: %d %f\n", data[0][0].(int64), data[0][1].(float64)) // FIELDS: 1 2.2

	//	// oh cool, that worked. lets blast this baby and tell it to run a bunch of statements
	//	// in neo concurrently with a pipeline
	//	results, _ := conn.ExecPipeline([]string{
	//		"MATCH (n:NODE) CREATE (n)-[:REL]->(f:FOO)",
	//		"MATCH (n:NODE) CREATE (n)-[:REL]->(b:BAR)",
	//		"MATCH (n:NODE) CREATE (n)-[:REL]->(z:BAZ)",
	//		"MATCH (n:NODE) CREATE (n)-[:REL]->(f:FOO)",
	//		"MATCH (n:NODE) CREATE (n)-[:REL]->(b:BAR)",
	//		"MATCH (n:NODE) CREATE (n)-[:REL]->(z:BAZ)",
	//	}, nil, nil, nil, nil, nil, nil)
	//	for _, result := range results {
	//		numResult, _ := result.RowsAffected()
	//		fmt.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 2 (per each iteration)
	//	}
	//
	//	data, _, _, _ = conn.QueryNeoAll("MATCH (n:NODE)-[:REL]->(m) RETURN m", nil)
	//	for _, row := range data {
	//		fmt.Printf("NODE: %#v\n", row[0].(graph.Node)) // Prints all nodes
	//	}
	//
	//	result, _ = conn.ExecNeo(`MATCH (n) DETACH DELETE n`, nil)
	//	numResult, _ = result.RowsAffected()
	//	fmt.Printf("Rows Deleted: %d", numResult) // Rows Deleted: 13
}
