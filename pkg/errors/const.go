package errors

var (
	// error in witness package
	ErrInvalidWitnessScriptLength = New(10001, "ErrInvalidWitnessScriptLength")
	ErrInvalidWitnessScriptPkFlag = New(10002, "ErrInvalidWitnessScriptPkFlag")
	ErrInvalidPayLoad             = New(10003, "ErrInvalidPayLoad")
	ErrOptionNotFound             = New(10004, "ErrOptionNotFound")
	ErrInvalidAtomicalsData       = New(10005, "ErrInvalidAtomicalsData")

	// error in atomicals package
	ErrInvalidFtCommitHeight  = New(20000, "ErrInvalidFtCommitHeight")
	ErrInvalidFtCurrentHeight = New(20001, "ErrInvalidFtCurrentHeight")
	ErrInvalidTicker          = New(20002, "ErrInvalidTicker")
	ErrTickerHasExist         = New(20003, "ErrTickerHasExist")
	ErrInvalidBitWork         = New(20004, "ErrInvalidBitWork")
	ErrInvalidRealm           = New(20005, "ErrInvalidRealm")
	ErrInvalidContainer       = New(20006, "ErrInvalidContainer")
	ErrNotDeployFt            = New(20007, "ErrNotDeployFt")
	ErrInvalidMintAmount      = New(20008, "ErrInvalidMintAmount")
	ErrInvalidMintHeight      = New(20009, "ErrInvalidMintHeight")
	ErrInvalidMaxMints        = New(20010, "ErrInvalidMaxMints")
)
