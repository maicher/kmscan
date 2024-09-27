package config

type Profile struct {
	Name           string  `toml:"name"`
	Resolution     int     `toml:"resolution"`
	MinHeight      int     `toml:"min-height"`
	MinWidth       int     `toml:"min-width"`
	MinAspectRatio float64 `toml:"min-aspect-ratio"`
	MaxAspectRatio float64 `toml:"max-aspect-ratio"`
	Brightness     float64 `toml:"brightness"`
	Window         int     `toml:"window"`
	Threshold      int     `toml:"threshold"`
}

func (p *Profile) GetName() string            { return p.Name }
func (p *Profile) GetResolution() int         { return p.Resolution }
func (p *Profile) GetMinHeight() int          { return p.MinHeight }
func (p *Profile) GetMinWidth() int           { return p.MinWidth }
func (p *Profile) GetMinAspectRatio() float64 { return p.MinAspectRatio }
func (p *Profile) GetMaxAspectRatio() float64 { return p.MaxAspectRatio }
func (p *Profile) GetBrightness() float64     { return p.Brightness }
func (p *Profile) GetWindow() int             { return p.Window }
func (p *Profile) GetThreshold() int          { return p.Threshold }
