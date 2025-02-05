package schema

type Metadata struct {
	Guid string
}

type ListResponse struct {
	Pagination Pagination
	Resources  []List
}

type List struct {
	Guid string
}
type Pagination struct {
	TotalResults int `json:"total_results"`
}

type AppResource struct {
	Metadata struct {
		Url string
	}
}

type AppsResponse struct {
	Resources []AppResource
}
type Stat struct {
	Stats struct {
		Host string
		Port int
	}
}
type StatsResponse map[string]Stat

type RouteResource struct {
	Entity struct {
		Port uint16
	}
}

type OrgResource struct {
	Resources []Org
}
type Org struct {
	Guid  string
	Links Links
}
type Links struct {
	Quota Quota
}
type Quota struct {
	Href string
}

type DomainResponse struct {
	Resources []Domain
}
type Domain struct {
	Guid string
}
type RouteObject struct {
	Resources []Route
}
type Route struct {
	Guid         string        `json:"guid"`
	Destinations []Destination `json:"destinations"`
}
type Destination struct {
	App      App    `json:"app"`
	Port     uint16 `json:"port"`
	Protocol string `json:"protocol"`
}
type App struct {
	Guid    string  `json:"guid"`
	Process Process `json:"process"`
}
type Process struct {
	Type string `json:"type"`
}
