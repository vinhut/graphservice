package models

import (
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"

	"os"
)

type RelationDatabase interface {
	Find(string, string, string) ([][]interface{}, error)
	Connect(string, string, string) error
	Disconnect(string, string) error
	Following(string) ([][]interface{}, error)
	Followers(string) ([][]interface{}, error)
}

type relationDatabase struct {
	conn bolt.Conn
}

func NewRelationDatabase(driver bolt.Driver) RelationDatabase {
	conn, _ := driver.OpenNeo(os.Getenv("NEO4J_SERVICE_URL"))
	return &relationDatabase{
		conn: conn,
	}
}

func (relationdb *relationDatabase) Find(key, value, value2 string) ([][]interface{}, error) {
	query_is_following := "match (n:Person { uid: {value1} })-[:FOLLOW]->(p:Person { uid: {value2} }) return 'ok'"
	data := map[string]interface{}{
		"key1":   key,
		"key2":   key,
		"value1": value,
		"value2": value2,
	}
	result, _, _, result_err := relationdb.conn.QueryNeoAll(query_is_following, data)
	if result_err != nil {
		return nil, result_err
	}
	return result, nil
}

func (relationdb *relationDatabase) Connect(uid, username, follow_uid string) error {

	query_follow := `
            MERGE (p:Person { uid: {follow_uid} })
            MERGE (n:Person { uid: {uid}, name: {name} })
            MERGE (n)-[:FOLLOW]->(p)
        `
	data := map[string]interface{}{
		"follow_uid": follow_uid,
		"uid":        uid,
		"name":       username,
	}
	_, result_err := relationdb.conn.ExecNeo(query_follow, data)

	return result_err

}

func (relationdb *relationDatabase) Disconnect(uid, follow_uid string) error {

	query_unfollow := `
        MATCH (n:Person { uid: {uid} })-[r:FOLLOW]->(p:Person {uid: {follow_uid} })
        DELETE r
        `
	data := map[string]interface{}{
		"follow_uid": follow_uid,
		"uid":        uid,
	}
	_, result_err := relationdb.conn.ExecNeo(query_unfollow, data)

	return result_err

}

func (relationdb *relationDatabase) Following(uid string) ([][]interface{}, error) {

	query_following := `
        MATCH (n:Person { uid: {uid} })-->(p)
        RETURN p.uid
        `
	data := map[string]interface{}{
		"uid": uid,
	}
	result, _, _, result_err := relationdb.conn.QueryNeoAll(query_following, data)

	return result, result_err

}

func (relationdb *relationDatabase) Followers(uid string) ([][]interface{}, error) {

	query_followers := `
        MATCH (n:Person { uid: {uid} })<--(p)
        RETURN p.uid
        `
	data := map[string]interface{}{
		"uid": uid,
	}

	result, _, _, result_err := relationdb.conn.QueryNeoAll(query_followers, data)

	return result, result_err

}
