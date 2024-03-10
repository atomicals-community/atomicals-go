package common

// is_dmint_activated
func IsDmintActivated(height int64) bool {
	return height >= ATOMICALS_ACTIVATION_HEIGHT_DMINT
}
