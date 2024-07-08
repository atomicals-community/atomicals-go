package utils

// is_dmint_activated sort_fifo
func IsDmintActivated(height int64) bool {
	return height >= ATOMICALS_ACTIVATION_HEIGHT_DMINT
}

func IsCustomColoring(height int64) bool {
	return height >= ATOMICALS_ACTIVATION_HEIGHT_CUSTOM_COLORING
}
