package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Disneyjr/OrderAPI-CleanArch/configs"
	"github.com/Disneyjr/OrderAPI-CleanArch/internal/event/handler"
	"github.com/Disneyjr/OrderAPI-CleanArch/internal/infra/graph"
	"github.com/Disneyjr/OrderAPI-CleanArch/internal/infra/grpc/pb"
	"github.com/Disneyjr/OrderAPI-CleanArch/internal/infra/grpc/service"
	"github.com/Disneyjr/OrderAPI-CleanArch/internal/infra/web/webserver"
	"github.com/Disneyjr/OrderAPI-CleanArch/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	waitForDB(db)
	runMigrations(db)

	rabbitMQChannel := getRabbitMQChannel(configs.RabbitMQUser, configs.RabbitMQPass, configs.RabbitMQHost, configs.RabbitMQPort)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddMethodHandler(http.MethodPost, "/order", webOrderHandler.Create)
	webserver.AddMethodHandler(http.MethodGet, "/order", webOrderHandler.Get)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func waitForDB(db *sql.DB) {
	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			return
		}
		fmt.Println("Waiting for database connection...")
		time.Sleep(2 * time.Second)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
}

func runMigrations(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS orders (
		id varchar(255) NOT NULL,
		price float NOT NULL,
		tax float NOT NULL,
		final_price float NOT NULL,
		PRIMARY KEY (id)
	);`
	if _, err := db.Exec(query); err != nil {
		panic(err)
	}
}

func getRabbitMQChannel(user, pass, host, port string) *amqp.Channel {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
	var conn *amqp.Connection
	var err error
	for i := 0; i < 30; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		fmt.Println("Waiting for RabbitMQ connection...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	_, _ = ch.QueueDeclare(
		"orders", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	_ = ch.QueueBind(
		"orders",       // queue name
		"",             // routing key
		"amq.direct",   // exchange
		false,
		nil,
	)
	return ch
}
