package api

type Team struct {
	ID      string `json:"id"`
	NameEN  string `json:"name_en"`
	NameFA  string `json:"name_fa"`
	Flag    string `json:"flag"`
	FifaCode string `json:"fifa_code"`
	Iso2    string `json:"iso2"`
	Groups  string `json:"groups"`
}

type TeamsResponse struct {
	Teams []Team `json:"teams"`
}

type TeamResponse struct {
	Team Team `json:"team"`
}

type GroupTeam struct {
	TeamID string `json:"team_id"`
	MP     string `json:"mp"`
	W      string `json:"w"`
	L      string `json:"l"`
	D      string `json:"d"`
	Pts    string `json:"pts"`
	GF     string `json:"gf"`
	GA     string `json:"ga"`
	GD     string `json:"gd"`
}

type Group struct {
	Name  string      `json:"name"`
	Teams []GroupTeam `json:"teams"`
}

type GroupsResponse struct {
	Groups []Group `json:"groups"`
}

type GroupResponse struct {
	Group Group  `json:"group"`
	Teams []Team `json:"teams"`
}

type Game struct {
	ID            string `json:"id"`
	HomeTeamID    string `json:"home_team_id"`
	AwayTeamID    string `json:"away_team_id"`
	HomeScore     string `json:"home_score"`
	AwayScore     string `json:"away_score"`
	HomeScorers   string `json:"home_scorers"`
	AwayScorers   string `json:"away_scorers"`
	Group         string `json:"group"`
	Matchday      string `json:"matchday"`
	LocalDate     string `json:"local_date"`
	PersianDate   string `json:"persian_date"`
	StadiumID     string `json:"stadium_id"`
	Finished      string `json:"finished"`
	TimeElapsed   string `json:"time_elapsed"`
	Type          string `json:"type"`
	HomeTeamNameEN string `json:"home_team_name_en,omitempty"`
	HomeTeamNameFA string `json:"home_team_name_fa,omitempty"`
	AwayTeamNameEN string `json:"away_team_name_en,omitempty"`
	AwayTeamNameFA string `json:"away_team_name_fa,omitempty"`
	HomeTeamLabel  string `json:"home_team_label,omitempty"`
	AwayTeamLabel  string `json:"away_team_label,omitempty"`
}

type GamesResponse struct {
	Games []Game `json:"games"`
}

type GameResponse struct {
	Game Game `json:"game"`
}

type Stadium struct {
	ID        string `json:"id"`
	NameEN    string `json:"name_en"`
	NameFA    string `json:"name_fa"`
	FifaName  string `json:"fifa_name"`
	CityEN    string `json:"city_en"`
	CityFA    string `json:"city_fa"`
	CountryEN string `json:"country_en"`
	CountryFA string `json:"country_fa"`
	Capacity  int    `json:"capacity"`
	Region    string `json:"region"`
}

type StadiumsResponse struct {
	Stadiums []Stadium `json:"stadiums"`
}

type StadiumResponse struct {
	Stadium Stadium `json:"stadium"`
}

type User struct {
	ID    string `json:"_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type HealthResponse struct {
	Status      string         `json:"status"`
	Timestamp   string         `json:"timestamp"`
	Uptime      int64          `json:"uptime"`
	Version     string         `json:"version"`
	Environment string         `json:"environment"`
	Database    HealthDatabase `json:"database"`
	Memory      HealthMemory   `json:"memory"`
}

type HealthDatabase struct {
	Status string `json:"status"`
	Name   string `json:"name"`
}

type HealthMemory struct {
	Used  string `json:"used"`
	Total string `json:"total"`
}
