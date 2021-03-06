package main

import (
	"github.com/gin-gonic/gin"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	transport "github.com/uber/jaeger-client-go/transport/zipkin"
	"github.com/uber/jaeger-client-go/zipkin"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/vinhut/graphservice/models"
	"github.com/vinhut/graphservice/services"

	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var SERVICE_NAME = "graph-service"

type UserAuthData struct {
	Uid      string
	Username string
	Email    string
	Role     string
	Created  string
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

func setupRouter(authservice services.AuthService, relationdb models.RelationDatabase, kafkaservice services.KafkaService) *gin.Engine {

	var JAEGER_COLLECTOR_ENDPOINT = os.Getenv("JAEGER_COLLECTOR_ENDPOINT")
	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()

	trsport, _ := transport.NewHTTPTransport(
		JAEGER_COLLECTOR_ENDPOINT,
		transport.HTTPLogger(jaeger.StdLogger),
	)
	cfg := jaegercfg.Configuration{
		ServiceName: "graph-service",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: JAEGER_COLLECTOR_ENDPOINT,
		},
	}
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory
	cfg.InitGlobalTracer(
		"graph-service",
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaegercfg.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		jaegercfg.ZipkinSharedRPCSpan(true),
		jaegercfg.Reporter(jaeger.NewRemoteReporter(trsport)),
	)
	tracer := opentracing.GlobalTracer()

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.GET(SERVICE_NAME+"/follow", func(c *gin.Context) {

		span := tracer.StartSpan("get follow status")

		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		user_data, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		follow_uid, _ := c.GetQuery("uid")
		result, find_err := relationdb.Find("uid", user_data.Uid, follow_uid)
		if find_err != nil {
			panic(find_err)
		}

		fmt.Println("result : ", result)
		if len(result) > 0 {
			c.String(200, "ok")
			span.Finish()
		} else {
			c.String(200, "nok")
			span.Finish()
		}

	})

	router.POST(SERVICE_NAME+"/follow", func(c *gin.Context) {

		span := tracer.StartSpan("follow user")

		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		user_data, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		follow_uid, _ := c.GetQuery("uid")

		result_err := relationdb.Connect(user_data.Uid, user_data.Username, follow_uid)

		if result_err != nil {
			panic(result_err)
		}

		kafka_err := kafkaservice.Send("follow", user_data.Uid+","+follow_uid+",1")
		if kafka_err != nil {
			panic(kafka_err)
		}
		c.String(200, "OK")
		span.Finish()

	})

	router.DELETE(SERVICE_NAME+"/follow", func(c *gin.Context) {

		span := tracer.StartSpan("unfollow user")

		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		user_data, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		follow_uid, _ := c.GetQuery("uid")

		result_err := relationdb.Disconnect(user_data.Uid, follow_uid)

		if result_err != nil {
			panic(result_err)
		}

		c.String(200, "OK")
		span.Finish()

	})

	router.GET(SERVICE_NAME+"/following", func(c *gin.Context) {

		span := tracer.StartSpan("get following")

		uid, _ := c.GetQuery("uid")
		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		_, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		result, result_err := relationdb.Following(uid)
		if result_err != nil {
			panic(result_err)
		}

		if len(result) == 0 {
			c.String(200, "[]")
			span.Finish()
		} else {

			fmt.Println("result = ", result)
			uid_list := make([]string, len(result))
			for idx, row := range result {
				fmt.Println("row = ", row)
				uid_list[idx] = row[0].(string)
			}

			result_json, _ := json.Marshal(uid_list)
			c.String(200, string(result_json))
			span.Finish()
		}

	})

	router.GET(SERVICE_NAME+"/followers", func(c *gin.Context) {

		span := tracer.StartSpan("get followers")

		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		_, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		uid, _ := c.GetQuery("uid")
		result, result_err := relationdb.Followers(uid)
		if result_err != nil {
			panic(result_err)
		}

		if len(result) == 0 {
			c.String(200, "[]")
			span.Finish()
		} else {

			fmt.Println("result = ", len(result))
			uid_list := make([]string, len(result))
			for idx, row := range result {
				fmt.Println("row = ", row)
				uid_list[idx] = row[0].(string)
			}

			result_json, _ := json.Marshal(uid_list)
			c.String(200, string(result_json))
			span.Finish()
		}
	})

	router.GET(SERVICE_NAME+"/followers_count", func(c *gin.Context) {

		span := tracer.StartSpan("get followers count")

		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		_, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		uid, _ := c.GetQuery("uid")
		result, result_err := relationdb.Followers(uid)
		if result_err != nil {
			panic(result_err)
		}

		c.String(200, strconv.Itoa(len(result)))
		span.Finish()
	})

	router.GET(SERVICE_NAME+"/following_count", func(c *gin.Context) {

		span := tracer.StartSpan("get following count")

		value, err := c.Cookie("token")
		if err != nil {
			panic("failed get token")
		}
		_, check_err := checkUser(authservice, value)
		if check_err != nil {
			panic("error check user")
		}

		uid, _ := c.GetQuery("uid")
		result, result_err := relationdb.Following(uid)
		if result_err != nil {
			panic(result_err)
		}

		c.String(200, strconv.Itoa(len(result)))
		span.Finish()
	})

	// internal
	router.POST("/internal/follow", func(c *gin.Context) {

		span := tracer.StartSpan("internal follow user")

		my_uid, _ := c.GetQuery("my_uid")
		my_username, _ := c.GetQuery("my_username")
		follow_uid, _ := c.GetQuery("follow_uid")

		result_err := relationdb.Connect(my_uid, my_username, follow_uid)

		if result_err != nil {
			panic(result_err)
		}

		c.String(200, "OK")
		span.Finish()

	})

	return router
}

func main() {

	authservice := services.NewUserAuthService()
	kafka_service := services.NewKafkaWriterService()
	driver := bolt.NewDriver()
	relationdb := models.NewRelationDatabase(driver)
	router := setupRouter(authservice, relationdb, kafka_service)
	router.Run(":8080")

}
