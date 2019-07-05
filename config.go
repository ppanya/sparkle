package sparkle

var (
	LocalMongoDBURL   = "mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs"
	LocalDatastoreURL = ""
)

type SparkleOption struct {
	DB Database
}

type Config struct{}
