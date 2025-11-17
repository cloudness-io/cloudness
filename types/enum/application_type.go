package enum

type ApplicationType string

const (
	ApplicationTypeStateless  ApplicationType = "Stateless"
	ApplicationTypeStateful   ApplicationType = "Stateful"
	ApplicationTypePostgresHA ApplicationType = "PostgresHA"
)
