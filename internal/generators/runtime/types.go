package runtime

type ExecutionPlan struct {
	Command string
	Args    []string
	Dir     string
	Env     []string
}
