package main

import (
	"github.com/gin-gonic/gin"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/vinhut/graphservice/services"

	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var SERVICE_NAME = "graph-service"

type UserAuthData struct {
	Uid     string
	Email   string
	Role    string
	Created string
}

func checkUser(authservice services.AuthService, token string) (*UserAuthData, error) {

	data := &UserAuthData{}
	user_data, auth_error := authservice.Check(SERVICE_NAME, token)
	if auth_error != nil {
		return data, auth_error
	}

	if err := json.Unmarshal([]byte(user_data), data); err != nil {
		fmt.Println(err)
		return data, err
	}

	return data, nil

}

func setupRouter(authservice services.AuthService) *gin.Engine {

	router := gin.Default()

	driver := bolt.NewDriver()
	conn, _ := driver.OpenNeo(os.Getenv("NEO4J_SERVICE_URL"))

	router.GET(SERVICE_NAME+"/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.GET(SERVICE_NAME+"/follow", func(c *gin.Context) {

		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		user_data, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		follow_uid, _ := c.GetQuery("uid")

		query_is_following := "match (n:Person { uid: {id} })-[:FOLLOW]->(p:Person {uid: {followuid}}) return 'ok'"
		data := map[string]interface{}{
			"followuid": follow_uid,
			"id":        user_data.Uid,
		}
		result, _, _, result_err := conn.QueryNeoAll(query_is_following, data)
		if result_err != nil {
			panic(result_err)
		}

		fmt.Println("result : ", result)
		if len(result) > 0 {
			c.String(200, "ok")
		} else {
			c.String(200, "nok")
		}

	})

	router.POST(SERVICE_NAME+"/follow", func(c *gin.Context) {

		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		user_data, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		follow_uid, _ := c.GetQuery("uid")

		query_follow := `
                MERGE (p:Person { uid: {followuid} })
		MERGE (n:Person { uid: {id}, name: {name} })
                MERGE (n)-[:FOLLOW]->(p)
		`
		data := map[string]interface{}{
			"followuid": follow_uid,
			"id":        user_data.Uid,
			"name":      user_data.Email,
		}
		_, result_err := conn.ExecNeo(query_follow, data)
		if result_err != nil {
			panic(result_err)
		}

		c.String(200, "OK")

	})

	router.POST(SERVICE_NAME+"/unfollow", func(c *gin.Context) {

		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		user_data, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		follow_uid, _ := c.GetQuery("uid")

		query_unfollow := `
                MATCH (n:Person { uid: {id} })-[r:FOLLOW]->(p:Person {uid: {followuid}})
		DELETE r
		`
		data := map[string]interface{}{
			"followuid": follow_uid,
			"id":        user_data.Uid,
		}
		_, result_err := conn.ExecNeo(query_unfollow, data)
		if result_err != nil {
			panic(result_err)
		}

		c.String(200, "OK")

	})

	router.GET(SERVICE_NAME+"/following", func(c *gin.Context) {

		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		user_data, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		query_following := `
                MATCH (n:Person { uid: {id} })-->(p)
		RETURN p.uid
		`
		data := map[string]interface{}{
			"id": user_data.Uid,
		}
		result, _, _, result_err := conn.QueryNeoAll(query_following, data)
		if result_err != nil {
			panic(result_err)
		}

		if len(result) == 0 {
			c.String(200, "[]")
		}

		fmt.Println("result = ", result)
		uid_list := make([]string, len(result))
		for idx, row := range result {
			fmt.Println("row = ", row)
			uid_list[idx] = row[0].(string)
		}

		result_json, _ := json.Marshal(uid_list)
		c.String(200, string(result_json))

	})

	router.GET(SERVICE_NAME+"/followers", func(c *gin.Context) {

		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		user_data, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		query_followers := `
                MATCH (n:Person { uid: {id} })<--(p)
		RETURN p.uid
		`
		data := map[string]interface{}{
			"id": user_data.Uid,
		}
		result, _, _, result_err := conn.QueryNeoAll(query_followers, data)
		if result_err != nil {
			panic(result_err)
		}

		if len(result) == 0 {
			c.String(200, "[]")
		} else {

			fmt.Println("result = ", len(result))
			uid_list := make([]string, len(result))
			for idx, row := range result {
				fmt.Println("row = ", row)
				uid_list[idx] = row[0].(string)
			}

			result_json, _ := json.Marshal(uid_list)
			c.String(200, string(result_json))
		}
	})

	router.GET(SERVICE_NAME+"/followers_count", func(c *gin.Context) {

		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		user_data, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		query_followers_count := `
                MATCH (n:Person { uid: {id} })<--(p)
		RETURN count(p.uid) as count
		`
		data := map[string]interface{}{
			"id": user_data.Uid,
		}
		result, _, _, result_err := conn.QueryNeoAll(query_followers_count, data)
		if result_err != nil {
			panic(result_err)
		}

		fmt.Println("result = ", result)
		c.String(200, strconv.FormatInt((result[0][0].(int64)), 10))
	})

	router.GET(SERVICE_NAME+"/following_count", func(c *gin.Context) {

		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		user_data, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		query_following_count := `
                MATCH (n:Person { uid: {id} })-->(p)
		RETURN count(p.uid) as count
		`
		data := map[string]interface{}{
			"id": user_data.Uid,
		}
		result, _, _, result_err := conn.QueryNeoAll(query_following_count, data)
		if result_err != nil {
			panic(result_err)
		}

		fmt.Println("result = ", result)
		c.String(200, strconv.FormatInt((result[0][0].(int64)), 10))
	})

	return router
}

func main() {

	authservice := services.NewUserAuthService()
	router := setupRouter(authservice)
	router.Run(":8080")

}
