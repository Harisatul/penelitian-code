package shared

var (
	LogEventStateDecodeRequest   = "decode_request"
	LogEventStateValidateRequest = "validate_request"
	LogEventStateFetchDB         = "fetch_db"
	LogEventStateInsertDB        = "insert_db"
	LogEventStateUpdateDB        = "update_db"
	LogEventStateCreatePayment   = "create_payment"
	LogEventStateKafkaPublish    = "kafka_publish"
)
