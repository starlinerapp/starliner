package value

type LogSource string

const (
	LogSourceWorkload LogSource = "workload"
	LogSourceIngress  LogSource = "ingress"
)
