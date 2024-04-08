package common

const (
	Satoshi = float64(100000000) //1BTC
)

const (
	ATOMICALS_ENVELOPE_MARKER_BYTES                        = "0461746f6d"
	MINT_REALM_CONTAINER_TICKER_COMMIT_REVEAL_DELAY_BLOCKS = 3
	MINT_GENERAL_COMMIT_REVEAL_DELAY_BLOCKS                = 100
	MINT_SUBNAME_COMMIT_PAYMENT_DELAY_BLOCKS               = 15 // # ~2.5 hours.
	MINT_SUBNAME_RULES_BECOME_EFFECTIVE_IN_BLOCKS          = 1
	MAX_SUBNAME_RULE_SIZE_LEN                              = 100000
	MAX_SUBNAME_RULE_ENTRIES                               = 100
	DFT_MINT_AMOUNT_MIN                                    = 546
	DFT_MINT_AMOUNT_MAX                                    = 100000000
	DFT_MINT_MAX_MIN_COUNT                                 = 1
	DFT_MINT_MAX_MAX_COUNT_LEGACY                          = 500000
	DFT_MINT_MAX_MAX_COUNT_DENSITY                         = 21000000
	DFT_MINT_HEIGHT_MIN                                    = 0
	DFT_MINT_HEIGHT_MAX                                    = 10000000
	VOUT_EXPECT_OUTPUT_INDEX                               = 0
	DMINT_PATH                                             = "dmint"
	SUBNAME_MIN_PAYMENT_DUST_LIMIT                         = 0 // # It can be possible to do free

)
