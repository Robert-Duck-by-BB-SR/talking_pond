package duckdom

type Theme struct {
	PrimaryBg       string
	SecondaryBg     string
	ActiveBg        string
	ActiveTextColor string
	StatusPanelBg   string
	TextColor       string

	ActiveBorderColor   string
	UnactiveBorderColor string
}

// Sigma ligma balls (got em) super duper male theme
var PRIMARY_THEME = Theme{
	PrimaryBg:           MakeRGBBackground(25, 25, 25),
	SecondaryBg:         MakeRGBBackground(34, 32, 32),
	StatusPanelBg:       MakeRGBBackground(48, 48, 48),
	ActiveBg:            MakeRGBBackground(251, 206, 44),
	TextColor:           MakeRGBTextColor(192, 192, 192),
	ActiveTextColor:     MakeRGBTextColor(251, 206, 44),
	ActiveBorderColor:   MakeRGBTextColor(251, 206, 44),
	UnactiveBorderColor: MakeRGBTextColor(54, 52, 52),
}
