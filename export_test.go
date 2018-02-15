package env

// ResetForTesting
func ResetForTesting() {
	CmdVar = NewVarSet("test")
}
