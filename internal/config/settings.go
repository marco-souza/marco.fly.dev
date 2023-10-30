package config

type MenuItem struct {
	Href string
	Name string
}

type PageParams struct {
	IsEnv     bool
	Title     string
	Repo      string
	MenuItems []MenuItem
}

var conf = Load()

var DefaultPageParams = PageParams{
	conf.Env == "development",
	"Marco.labs 🚀",
	"https://github.com/marco-souza/gx",
	[]MenuItem{
		{"/", "Home"},
		{"https://marco.deno.dev/blog", "Blog"},
		{"/playground", "Playground"},
	},
}
