package config

type PageParams struct {
	IsDev bool
	Title string
	Repo  string
}

var conf = Load()

var DefaultPageParams = PageParams{
	conf.Env == "development",
	"Marco.labs 🚀",
	"https://github.com/marco-souza",
}
