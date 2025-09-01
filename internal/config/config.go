package config

// Redis server configuration
var Protocol = "tcp"
var Port = ":6379"
var MaxConnection = 20000

// HTTP Gateway configuration
var HTTPPort = ":8080"
var HTTPReadTimeout = 15  // seconds
var HTTPWriteTimeout = 15 //seconds

// Gateway to Redis connection configuration
var RedisConnTimeout = 5 // seconds
var RedisRWTimeout = 10  // seconds
