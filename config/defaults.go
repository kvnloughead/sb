package config

// DefaultConfig holds the fallback configuration for the sb CLI tool.
var DefaultConfig = Config{
	Repo: "~/tripleten/se-canonicals_en",
	Slugs: map[string]string{
		"master": "master",
		"p1_1":   "html/sprint-1-new_stage-1",
		"p1_2":   "html/sprint-1-new_stage-2",
		"p1_3":   "html/sprint-1-new_stage-3",
	},
}
