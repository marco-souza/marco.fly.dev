package config

type Github struct {
	ClientId     string
	ClientSecret string
	CallbackUrl  string
	RefreshUrl   string
	LogoutUrl    string
	SignInUrl    string
	DashboardUrl string
	Scope        string
}

func GithubLoad() Github {
	conf := Github{
		ClientId:     env("GITHUB_CLIENT_ID", ""),
		ClientSecret: env("GITHUB_CLIENT_SECRET", ""),
		DashboardUrl: "/#/",
		CallbackUrl:  "/api/auth/github/callback",
		RefreshUrl:   "/api/auth/github/refresh",
		LogoutUrl:    "/api/auth/github/logout",
		SignInUrl:    "/api/auth/github",
		Scope:        "read:user",
	}

	// TODO: mock for testing
	// if conf.ClientId == "" || conf.ClientSecret == "" {
	// 	panic("github credentials not found")
	// }

	return conf
}
