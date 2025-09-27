package contextkeys

type CtxKey string

const RequestIDCtxKey CtxKey = "request_id"
const TraceIDCtxKey CtxKey = "trace_id"
const UserIDCtxKey CtxKey = "user_id"