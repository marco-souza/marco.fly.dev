package config

type MenuItem struct {
	Href string
	Name string
}

type PageParams struct {
	IsDev     bool
	Title     string
	Repo      string
	MenuItems []MenuItem
}

var conf = Load()

var DefaultPageParams = PageParams{
	conf.Env == "development",
	"Marco.labs 🚀",
	"https://github.com/marco-souza",
	[]MenuItem{
		{"/", "Home"},
		{"/resume", "Resume"},
		// {"https://marco.deno.dev/blog", "Blog"},
		{"/login", "Login"},
	},
}

var PrivatePageParams = PageParams{
	conf.Env == "development",
	"Marco.labs 🚀",
	"https://github.com/marco-souza/marco.fly.io",
	[]MenuItem{
		{"/", "Home"},
		{"/app/", "Dashboard"},
		{"/app/playground", "Playground"},
		{"/app/orders", "Ordero"},
		{conf.Github.LogoutUrl, "Logout"},
	},
}

func MakePageParams(authhenticated bool) PageParams {
	if authhenticated {
		return PrivatePageParams
	}
	return DefaultPageParams
}
