package duckdom

type Theme struct {
	PrimaryBg          string
	SecondaryBg        string
	SecondaryTextColor string

	ActiveBg        string
	ActiveTextColor string

	ActiveDarkBg        string
	ActiveDarkTextColor string

	TextColor string
}

// Sigma ligma super duper male theme
var PRIMARY_THEME = Theme{
	PrimaryBg:           MakeRGBBackground(6, 6, 6),
	SecondaryBg:         MakeRGBBackground(0, 244, 0),
	SecondaryTextColor:  MakeRGBTextColor(11, 11, 11),
	ActiveBg:            MakeRGBBackground(255, 0, 0),
	ActiveTextColor:     MakeRGBTextColor(8, 56, 100),
	ActiveDarkBg:        MakeRGBBackground(2, 14, 25),
	ActiveDarkTextColor: MakeRGBTextColor(2, 14, 25),
	TextColor:           MakeRGBTextColor(255, 255, 255),
}
