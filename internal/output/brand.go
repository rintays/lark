package output

import "github.com/charmbracelet/lipgloss"

const larkBrandBlue = "#3370FF"

// BrandColor returns the Lark/Feishu primary brand color for UI accents.
func BrandColor() lipgloss.Color {
	return lipgloss.Color(larkBrandBlue)
}
