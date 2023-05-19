package api

import (
	"time"

	"github.com/google/uuid"
	db "github.com/techschool/simplebank/db/sqlc"
)

type createAccountResponse struct {
	ID         string    `json:"id"`
	FirstName  string    `json:"first_name,omitempty"`
	Email      string    `json:"email,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	BirthDate  int64     `json:"birth_date,omitempty"`
	Gender     string    `json:"gender,omitempty"`
	ShowMe     string    `json:"show_me,omitempty"`
	University string    `json:"university,omitempty"`
	Nsfw       bool      `json:"nsfw,omitempty"`
	Ethnicity  string    `json:"ethnicity,omitempty"`
	Interests  []string  `json:"interests,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

func newCreateAccountResponse(account db.Account) createAccountResponse {
	return createAccountResponse{
		ID:         account.ID,
		FirstName:  account.FirstName,
		Email:      account.Email,
		Phone:      account.Phone,
		BirthDate:  account.BirthDate,
		Gender:     account.Gender,
		ShowMe:     account.ShowMe,
		University: account.University.String,
		Nsfw:       account.Nsfw,
		Ethnicity:  account.Ethnicity,
		Interests:  account.Interests,
		CreatedAt:  account.CreatedAt,
	}
}

type getAccountResponse struct {
	ID             string    `json:"id"`
	FirstName      string    `json:"first_name,omitempty"`
	Email          string    `json:"email,omitempty"`
	Phone          string    `json:"phone,omitempty"`
	BirthDate      int64     `json:"birth_date,omitempty"`
	Gender         string    `json:"gender,omitempty"`
	ShowMe         string    `json:"show_me,omitempty"`
	University     string    `json:"university,omitempty"`
	Nsfw           bool      `json:"nsfw,omitempty"`
	Ethnicity      string    `json:"ethnicity,omitempty"`
	Interests      []string  `json:"interests,omitempty"`
	Picture        []string  `json:"picture,omitempty"`
	VerifyYourself bool      `json:"verify_yourself,omitempty"`
	AboutMe        string    `json:"about_me,omitempty"`
	TimeZone       string    `json:"time_zone,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

func newGetAccountResponse(account db.Account) getAccountResponse {
	return getAccountResponse{
		ID:             account.ID,
		FirstName:      account.FirstName,
		Email:          account.Email,
		Phone:          account.Phone,
		BirthDate:      account.BirthDate,
		Gender:         account.Gender,
		ShowMe:         account.ShowMe,
		University:     account.University.String,
		Nsfw:           account.Nsfw,
		Ethnicity:      account.Ethnicity,
		Interests:      account.Interests,
		Picture:        account.Picture,
		VerifyYourself: account.VerifyYourself.Bool,
		AboutMe:        account.AboutMe.String,
		TimeZone:       account.TimeZone.String,
		CreatedAt:      account.CreatedAt,
	}
}

type modifyAccountResponse struct {
	ID             string   `json:"id,omitempty"`
	VerifyYourself bool     `json:"verify_yourself,omitempty"`
	AboutMe        string   `json:"about_me,omitempty"`
	Interests      []string `json:"interests,omitempty"`
	Gender         string   `json:"gender,omitempty"`
	TimeZone       string   `json:"time_zone,omitempty"`
	Ethnicity      string   `json:"ethnicity,omitempty"`
	Nsfw           bool     `json:"nsfw,omitempty"`
	Picture        []string `json:"picture,omitempty"`
}

func newModifyAccountResponse(account db.Account) modifyAccountResponse {
	return modifyAccountResponse{
		ID:             account.ID,
		VerifyYourself: account.VerifyYourself.Bool,
		AboutMe:        account.AboutMe.String,
		Interests:      account.Interests,
		Gender:         account.Gender,
		TimeZone:       account.TimeZone.String,
		Ethnicity:      account.Ethnicity,
		Nsfw:           account.Nsfw,
		Picture:        account.Picture,
	}
}

type createFavouriteResponse struct {
	FavID    uuid.UUID `json:"fav_id,omitempty"`
	UserID   string    `json:"user_id,omitempty"`
	TargetID string    `json:"target_id,omitempty"`
	FavAt    time.Time `json:"fav_at,omitempty"`
}

func newcreateFavouriteResponse(fav db.Favourite) createFavouriteResponse {
	return createFavouriteResponse{
		FavID:    fav.FavID,
		UserID:   fav.UserID,
		TargetID: fav.TargetID,
		FavAt:    fav.FavAt,
	}
}

type getFavouritesResponse struct {
	Favourites []db.GetAllFavouritesRow `json:"favourites,omitempty"`
	Metadata   *metadata                `json:"metadata"`
}

func newGetFavouritesResponse(fav []db.GetAllFavouritesRow, metadata *metadata) getFavouritesResponse {
	return getFavouritesResponse{
		Favourites: fav,
		Metadata:   metadata,
	}
}

type metadata struct {
	CurrentPage  int32 `json:"page_id"`
	PageSize     int32 `json:"page_size"`
	FirstPage    int   `json:"first_page"`
	LastPage     int   `json:"last_page"`
	TotalRecords int64 `json:"total_records"`
}

type getDiscoveredAccountsResponse struct {
	DiscoveredAccounts []db.DiscoverAccountsWithFilterRow `json:"discoverd_accounts,omitempty"`
	Metadata           *metadata                          `json:"metadata"`
}

func newGetDiscoveredAccountsResponse(accounts []db.DiscoverAccountsWithFilterRow, metadata *metadata) getDiscoveredAccountsResponse {
	return getDiscoveredAccountsResponse{
		DiscoveredAccounts: accounts,
		Metadata:           metadata,
	}
}

// type deletedAccountResponse struct {
// 	ID string `json:"id,omitempty"`
// }

// func newDeletedAccountResponse(id string) deletedAccountResponse {
// 	return deletedAccountResponse{
// 		ID: id,
// 	}
// }

type uploadPictureResponse struct {
	URL string `json:"url,omitempty"`
}

func newUploadPictureResponse(url string) uploadPictureResponse {
	return uploadPictureResponse{
		URL: url,
	}
}

type modifyPictureResponse struct {
	Message     string   `json:"message"`
	PicturesURL []string `json:"pictures,omitempty"`
}

func newModifyPictureResponse(urls []string, messsage string) modifyPictureResponse {
	return modifyPictureResponse{
		Message:     messsage,
		PicturesURL: urls,
	}
}
