package utils

// is_dmint_activated sort_fifo
func IsDmintActivated(height int64) bool {
	return height >= ATOMICALS_ACTIVATION_HEIGHT_DMINT
}

func IsCustomColoring(height int64) bool {
	return height >= ATOMICALS_ACTIVATION_HEIGHT_CUSTOM_COLORING
}

func Is_within_acceptable_blocks_for_sub_item_payment(commit_height, current_height int64) bool {
	return current_height <= commit_height+MINT_SUBNAME_COMMIT_PAYMENT_DELAY_BLOCKS
}

func Is_density_activated(height int64) bool {
	return height >= ATOMICALS_ACTIVATION_HEIGHT_DENSITY
}
