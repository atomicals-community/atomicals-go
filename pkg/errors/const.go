package errors

var (
	// error in witness package
	ErrInvalidWitnessScriptLength = New(10001, "ErrInvalidWitnessScriptLength")
	ErrInvalidWitnessScriptPkFlag = New(10002, "ErrInvalidWitnessScriptPkFlag")
	ErrInvalidPayLoad             = New(10003, "ErrInvalidPayLoad")
	ErrOptionNotFound             = New(10004, "ErrOptionNotFound")
	ErrInvalidAtomicalsData       = New(10005, "ErrInvalidAtomicalsData")

	// error in atomicals package
	ErrInvalidFtCommitHeight  = New(20006, "ErrInvalidFtCommitHeight")
	ErrInvalidFtCurrentHeight = New(20006, "ErrInvalidFtCurrentHeight")
	ErrInvalidTicker          = New(20006, "ErrInvalidTicker")
	ErrInvalidBitWork         = New(20006, "ErrInvalidBitWork")
	ErrInvalidRealm           = New(20006, "ErrInvalidRealm")
	ErrInvalidContainer       = New(20006, "ErrInvalidContainer")
	ErrNotDeployFt            = New(20006, "ErrNotDeployFt")
	ErrInvalidMintAmount      = New(20006, "ErrInvalidMintAmount")
	ErrInvalidMintHeight      = New(20006, "ErrInvalidMintHeight")
	ErrInvalidMaxMints        = New(20006, "ErrInvalidMaxMints")
)
