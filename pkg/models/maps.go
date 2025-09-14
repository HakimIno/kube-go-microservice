package models

import "time"

// HERE API Response Models
type HereAPIResponse struct {
	Items []HerePlaceItem `json:"items"`
}

type HerePlaceItem struct {
	Title      string          `json:"title"`
	ID         string          `json:"id"`
	Language   string          `json:"language"`
	ResultType string          `json:"resultType"`
	Address    HereAddress     `json:"address"`
	Position   HerePosition    `json:"position"`
	Access     []HereAccess    `json:"access"`
	Distance   int             `json:"distance"`
	Categories []HereCategory  `json:"categories"`
	References []HereReference `json:"references"`
	Contacts   []HereContact   `json:"contacts"`
}

type HereAddress struct {
	Label       string `json:"label"`
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
	County      string `json:"county"`
	City        string `json:"city"`
	District    string `json:"district"`
	Street      string `json:"street"`
	PostalCode  string `json:"postalCode"`
}

type HerePosition struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type HereAccess struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type HereCategory struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Primary bool   `json:"primary"`
}

type HereReference struct {
	Supplier struct {
		ID string `json:"id"`
	} `json:"supplier"`
	ID string `json:"id"`
}

type HereContact struct {
	Phone []HerePhone `json:"phone"`
	WWW   []HereWWW   `json:"www"`
}

type HerePhone struct {
	Value string `json:"value"`
}

type HereWWW struct {
	Value      string            `json:"value"`
	Categories []HereCategoryRef `json:"categories"`
}

type HereCategoryRef struct {
	ID string `json:"id"`
}

// Request Models
type PlaceSearchRequest struct {
	Query    string  `json:"query" binding:"required" example:"วัดพระแก้ว กรุงเทพมหานคร"`
	Lat      float64 `json:"lat" example:"13.7563"`
	Lng      float64 `json:"lng" example:"100.5018"`
	Limit    int     `json:"limit" example:"5"`
	Language string  `json:"language" example:"th"`
}

// Response Models
type PlaceSearchResponse struct {
	Places []PlaceItem `json:"places"`
	Meta   SearchMeta  `json:"meta"`
}

type PlaceItem struct {
	Title      string     `json:"title"`
	ID         string     `json:"id"`
	Language   string     `json:"language"`
	ResultType string     `json:"resultType"`
	Address    Address    `json:"address"`
	Position   Position   `json:"position"`
	Distance   int        `json:"distance"`
	Categories []Category `json:"categories"`
	Contacts   []Contact  `json:"contacts"`
}

type Address struct {
	Label       string `json:"label"`
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
	County      string `json:"county"`
	City        string `json:"city"`
	District    string `json:"district"`
	Street      string `json:"street"`
	PostalCode  string `json:"postalCode"`
}

type Position struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Category struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Primary bool   `json:"primary"`
}

type Contact struct {
	Phone []string `json:"phone"`
	WWW   []string `json:"www"`
}

type SearchMeta struct {
	Query        string    `json:"query"`
	ResultsCount int       `json:"resultsCount"`
	SearchTime   time.Time `json:"searchTime"`
	FromCache    bool      `json:"fromCache"`
}

// Cache Model
type CachedPlaceSearch struct {
	Query     string      `json:"query"`
	Lat       float64     `json:"lat"`
	Lng       float64     `json:"lng"`
	Limit     int         `json:"limit"`
	Language  string      `json:"language"`
	Results   []PlaceItem `json:"results"`
	CachedAt  time.Time   `json:"cachedAt"`
	ExpiresAt time.Time   `json:"expiresAt"`
}
