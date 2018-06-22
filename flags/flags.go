package flags

type EtcdFlags struct {
	EtcdAddress *string
	CertFile    *string
	KeyFile     *string
	CaFile      *string
	Version     *int
	Timeout     *int64
	Path        *string
	OverWrite   *bool
}
