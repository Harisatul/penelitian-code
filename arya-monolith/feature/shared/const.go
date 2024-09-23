package shared

var (
	LogEventStateDecodeRequest   = "decode_request"
	LogEventStateValidateRequest = "validate_request"
	LogEventStateFetchDB         = "fetch_db"
	LogEventStateInsertDB        = "insert_db"
	LogEventStateUpdateDB        = "update_db"
	LogEventStateCreatePayment   = "create_payment"
	LogEventStateInsertJobQueue  = "insert_job_queue"
	LogEventStateCompleteJob     = "complete_job"
	LogEventStateCancelJob       = "cancel_job"
	LogEventStateFetchCache      = "fetch_cache"
	LogEventStateSetCache        = "set_cache"
)
