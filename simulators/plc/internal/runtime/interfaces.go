package plc_runtime

type OPCUAClient interface {
	ReadVar(opcNs uint16, path string, vary string) (interface{}, error)
	WriteVar(opcNs uint16, path string, vars map[string]interface{}) error
}
