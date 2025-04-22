package saga

type WorkflowStep struct {
	Name             string
	Activity         interface{}
	Args             []interface{}
	FallbackActivity interface{}
	RetryThreshold   int
	NonRetriableErrs []string
}
