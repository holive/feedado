package rss

type Processor struct {
	updater   Repository
	userAgent string
}

type ProcessorConfig struct {
	UserAgent string
}
