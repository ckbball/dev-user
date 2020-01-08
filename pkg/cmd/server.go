package cmd

import (
  "context"
  //"database/sql"
  "flag"
  "fmt"
  "os"
  "strconv"

  // mysql driver
  //"github.com/Shopify/sarama"
  //"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
  "github.com/go-redis/cache/v7"
  "github.com/go-redis/redis/v7"
  _ "github.com/go-sql-driver/mysql"
  "github.com/vmihailenco/msgpack/v4"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"

  "github.com/ckbball/dev-user/pkg/logger"
  userGrpc "github.com/ckbball/dev-user/pkg/protocol/grpc"
  "github.com/ckbball/dev-user/pkg/protocol/rest"
  v1 "github.com/ckbball/dev-user/pkg/service/v1"
)

// Config is configuration for Server
type Config struct {
  // gRPC server start parameters section
  // gRPC is TCP port to listen by gRPC server
  GRPCPort string

  // the port to listen for http calls
  HTTPPort string

  // DB Datastore parameters section
  // DatastoreDBHost is host of database
  DatastoreDBHost string
  // DatastoreDBUser is username to connect to database
  DatastoreDBUser string
  // DatastoreDBPassword password to connect to database
  DatastoreDBPassword string
  // DatastoreDBSchema is schema of database
  DatastoreDBSchema string
  // address for single redis node
  RedisAddress string
  // uri for mongodb
  MongoAddress string
  // name for mongodb
  MongoName string
  // collection for mongodb
  MongoCollection string

  // Log parameters section
  // LogLevel is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
  LogLevel int
  // LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
  LogTimeFormat string
  // address for process level log server
  LoggerAddress string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
  ctx := context.Background()

  // get configuration
  var cfg Config
  flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
  flag.StringVar(&cfg.HTTPPort, "http-port", "", "http port to bind")
  flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Database host")
  flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Database user")
  flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
  flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Database schema")
  flag.StringVar(&cfg.RedisAddress, "redis-address", "", "Redis address")
  flag.StringVar(&cfg.MongoAddress, "mongouri", "", "Mongo address")
  flag.StringVar(&cfg.MongoName, "mongo-name", "", "Mongo db name")
  flag.StringVar(&cfg.MongoCollection, "mongo-collection", "", "Mongo db collection")
  flag.Parse()

  if len(cfg.GRPCPort) == 0 {
    cfg.GRPCPort = os.Getenv("GRPC_PORT")
    cfg.HTTPPort = os.Getenv("HTTP_PORT")
    cfg.DatastoreDBHost = os.Getenv("DB_HOST")
    cfg.DatastoreDBUser = os.Getenv("DB_USER")
    cfg.DatastoreDBPassword = os.Getenv("DB_PASSWORD")
    cfg.DatastoreDBSchema = os.Getenv("DB_SCHEMA")
    cfg.RedisAddress = os.Getenv("REDIS_ADDRESS")
    cfg.MongoAddress = os.Getenv("MONGO_URI")
    cfg.MongoName = os.Getenv("MONGO_NAME")
    cfg.MongoCollection = os.Getenv("MONGO_COLLECTION")
    cfg.LogLevel, _ = strconv.Atoi(os.Getenv("LOG_LEVEL"))
    cfg.LogTimeFormat = os.Getenv("LOG_TIME")
    cfg.LoggerAddress = os.Getenv("PROCESS_LOGGER")
  }

  if len(cfg.GRPCPort) == 0 {
    return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
  }

  if len(cfg.HTTPPort) == 0 {
    return fmt.Errorf("invalid TCP port for http server: '%s'", cfg.HTTPPort)
  }

  // initialize logger
  if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
    return fmt.Errorf("failed to initialize logger: %v", err)
  }

  // add MySQL driver specific parameter to parse date/time
  // Drop it for another database
  // param := "parseTime=true"

  // for non localhost db %s:%s@tcp(%s)/%s?%s
  // currently set for localhost
  /*
     dsn := fmt.Sprintf("%s:%s@/%s?%s",
       cfg.DatastoreDBUser,
       cfg.DatastoreDBPassword,
       // cfg.DatastoreDBHost,
       cfg.DatastoreDBSchema,
       param)
     db, err := sql.Open("mysql", dsn)
     if err != nil {
       return fmt.Errorf("failed to open database: %v", err)
     }
     defer db.Close()
     // create repository
     repository := &v1.JournalRepository{
       Db: db,
     }
  */
  // init pool of connections to redis cluster
  // redisPool := initRedis(cfg.RedisAddress)

  // Make subscriber config here
  /*
     saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
     saramaSubscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

     // Make subscriber pointer here
     subscriber := v1.InitSubscriber(saramaSubscriberConfig)

     // Make publisher pointer here
     publisher := v1.InitPublisher()
  */

  // SET up mongo client
  // retry := false
  clientOptions := options.Client().ApplyURI(cfg.MongoAddress)
  // clientOptions := options.MergeClientOptions(&options.ClientOptions{RetryWrites: &retry, RetryReads: &retry}).ApplyURI(cfg.MongoAddress)
  client, err := mongo.Connect(context.TODO(), clientOptions)
  if err != nil {
    return err
  }
  collection := client.Database(cfg.MongoName).Collection(cfg.MongoCollection)

  // create repository
  repository := v1.NewUserRepository(collection)

  // pass in fields of handler directly to method
  // v1API := v1.NewUserServiceServer(repository, subscriber, publisher)
  v1API := v1.NewUserServiceServer(repository, cfg.LoggerAddress)

  // run http gateway
  go func() {
    _ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
  }()

  return userGrpc.RunServer(ctx, v1API, cfg.GRPCPort)
}

func initRedis(address string) *cache.Codec {
  ring := redis.NewRing(&redis.RingOptions{
    Addrs: map[string]string{
      "server1": ":" + address,
    },
  })

  codec := &cache.Codec{
    Redis: ring,

    Marshal: func(v interface{}) ([]byte, error) {
      return msgpack.Marshal(v)
    },
    Unmarshal: func(b []byte, v interface{}) error {
      return msgpack.Unmarshal(b, v)
    },
  }

  return codec
}
