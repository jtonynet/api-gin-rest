package config

type FeatureFlags struct {
	PprofCPUEnabled           bool `mapstructure:"PPROF_CPU_FEATURE_FLAG_ENABLED"`
	PostAlunoAsMessageEnabled bool `mapstructure:"POST_ALUNO_AS_MESSAGE_FEATURE_FLAG_ENABLED"`
}
