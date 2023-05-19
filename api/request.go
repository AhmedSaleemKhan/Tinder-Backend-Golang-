package api

type createAccountRequest struct {
	FirstName  string   `json:"first_name" binding:"required"`
	Email      string   `json:"email" binding:"required,email"`
	Phone      string   `json:"phone" binding:"required"`
	BirthDate  int64    `json:"birth_date" binding:"required"`
	Gender     string   `json:"gender" binding:"required"`
	ShowMe     string   `json:"show_me" binding:"required"`
	University string   `json:"university"`
	Nsfw       bool     `json:"nsfw" binding:"required"`
	Ethnicity  string   `json:"ethnicity" binding:"required"`
	Interests  []string `json:"interests"`
}

type modifyAccountRequest struct {
	VerifyYourself bool     `json:"verify_yourself"`
	AboutMe        string   `json:"about_me"`
	Interests      []string `json:"interests"`
	Gender         string   `json:"gender"`
	TimeZone       string   `json:"time_zone"`
	Ethnicity      string   `json:"ethnicity"`
	Nsfw           bool     `json:"nsfw"`
	Picture        []string `json:"picture"`
}

type getFavouritesFilter struct {
	PageID   int32 `json:"page_id"`
	PageSize int32 `json:"page_size"`
}

type discoverFilters struct {
	PageID   int32 `json:"page_id"`
	PageSize int32 `json:"page_size"`
	Filters  struct {
		ShowMe    string `json:"show_me"`
		MinAge    int64  `json:"min_age"`
		MaxAge    int64  `json:"max_age"`
		Ethnicity string `json:"ethnicity"`
	} `json:"filters"`
}

type addPictureRequest struct {
	PictureURL []string `json:"url"`
}
