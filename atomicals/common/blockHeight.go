package common

const (
	ATOMICALS_ACTIVATION_HEIGHT                      = 808080
	ATOMICALS_ACTIVATION_HEIGHT_DMINT                = 819181
	ATOMICALS_ACTIVATION_HEIGHT_COMMITZ              = 822800
	ATOMICALS_ACTIVATION_HEIGHT_DENSITY              = 828128
	ATOMICALS_ACTIVATION_HEIGHT_DFT_BITWORK_ROLLOVER = 828628
)

// is_dmint_activated sort_fifo
func IsDmintActivated(height int64) bool {
	return height >= ATOMICALS_ACTIVATION_HEIGHT_DMINT
}

// A realm, ticker, or container reveal is valid as long as it is within MINT_REALM_CONTAINER_TICKER_COMMIT_REVEAL_DELAY_BLOCKS of the reveal and commit
func is_within_acceptable_blocks_for_name_reveal(commitHeight, revealLocationHeight int64) bool {
	return commitHeight >= revealLocationHeight-MINT_REALM_CONTAINER_TICKER_COMMIT_REVEAL_DELAY_BLOCKS
}

// A payment for a subrealm is acceptable as long as it is within MINT_SUBNAME_COMMIT_PAYMENT_DELAY_BLOCKS of the commitHeight
func isWithinAcceptableBlocksForSubItemPayment(commitHeight, currentHeight int64) bool {
	return currentHeight <= commitHeight+MINT_SUBNAME_COMMIT_PAYMENT_DELAY_BLOCKS
}

func is_density_activated(height int64) bool {
	if height >= ATOMICALS_ACTIVATION_HEIGHT_DENSITY {
		return true
	}
	return false
}
