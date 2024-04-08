package errors

var (
	// error in package
	ErrNotExistBlock = New(10001, "ErrNotExistBlock")

	// error in witness package
	ErrInvalidWitnessScriptLength = New(10001, "ErrInvalidWitnessScriptLength")
	ErrInvalidWitnessScriptPkFlag = New(10002, "ErrInvalidWitnessScriptPkFlag")
	ErrInvalidPayLoad             = New(10003, "ErrInvalidPayLoad")
	ErrOptionNotFound             = New(10004, "ErrOptionNotFound")
	ErrInvalidAtomicalsData       = New(10005, "ErrInvalidAtomicalsData")

	// error in atomicals package
	ErrInvalidCommitHeight          = New(20000, "ErrInvalidCommitHeight")
	ErrInvalidRevealLocationHeight  = New(20000, "ErrInvalidRevealLocationHeight")
	ErrInvalidCommitVoutIndex       = New(20000, "ErrInvalidCommitVoutIndex")
	ErrInvalidFtCurrentHeight       = New(20001, "ErrInvalidFtCurrentHeight")
	ErrInvalidTicker                = New(20002, "ErrInvalidTicker")
	ErrTickerHasExist               = New(20003, "ErrTickerHasExist")
	ErrTickerNotExist               = New(20003, "ErrTickerNotExist")
	ErrInvalidBitWork               = New(20004, "ErrInvalidBitWork")
	ErrInvalidRealm                 = New(20005, "ErrInvalidRealm")
	ErrRealmHasExist                = New(20005, "ErrRealmHasExist")
	ErrInvalidContainer             = New(20006, "ErrInvalidContainer")
	ErrInvalidContainerDmitem       = New(20006, "ErrInvalidContainerDmitem")
	ErrContainerHasExist            = New(20005, "ErrContainerHasExist")
	ErrContainerNotExist            = New(20005, "ErrContainerNotExist")
	ErrParentRealmNotExist          = New(20005, "ErrParentRealmNotExist")
	ErrSubRealmHasExist             = New(20005, "ErrSubRealmHasExist")
	ErrNotDeployFt                  = New(20007, "ErrNotDeployFt")
	ErrInvalidMintAmount            = New(20008, "ErrInvalidMintAmount")
	ErrInvalidMintHeight            = New(20009, "ErrInvalidMintHeight")
	ErrInvalidMaxMints              = New(20010, "ErrInvalidMaxMints")
	ErrCannotBeImmutable            = New(20010, "ErrCannotBeImmutable")
	ErrInvalidVinIndex              = New(20010, "ErrInvalidVinIndex")
	ErrInvalidClaimType             = New(20010, "ErrInvalidClaimType")
	ErrInvalidDftMd                 = New(20010, "ErrInvalidDftMd")
	ErrInvalidDftBv                 = New(20010, "ErrInvalidDftBv")
	ErrInvalidDftMintBitwork        = New(20010, "ErrInvalidDftMintBitwork")
	ErrInvalidDftBci                = New(20010, "ErrInvalidDftBci")
	ErrInvalidDftBsc                = New(20010, "ErrInvalidDftBsc")
	ErrInvalidDftBri                = New(20010, "ErrInvalidDftBri")
	ErrInvalidDftBrs                = New(20010, "ErrInvalidDftBrs")
	ErrInvalidDftMaxg               = New(20010, "ErrInvalidDftMaxg")
	ErrNameTypeMintMastHaveBitworkc = New(20010, "ErrNameTypeMintMastHaveBitworkc")
	ErrInvalidPerpetualBitwork      = New(20010, "ErrInvalidPerpetualBitwork")
	ErrInvalidMintedTimes           = New(20010, "ErrInvalidMintedTimes")
	ErrInvalidBitworkcPrefix        = New(20010, "ErrInvalidBitworkcPrefix")
	ErrBitworkcNeeded               = New(20010, "ErrBitworkcNeeded")
	ErrCheckRequest                 = New(20010, "ErrCheckRequest")
	ErrDmintNotStart                = New(20010, "ErrDmintNotStart")
	ErrInvalidRevealInputIndex      = New(20010, "ErrInvalidRevealInputIndex")
	ErrInvalidMerkleVerify          = New(20010, "ErrInvalidMerkleVerify")
)
