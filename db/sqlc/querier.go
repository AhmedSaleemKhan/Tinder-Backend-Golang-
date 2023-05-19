// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"context"
)

type Querier interface {
	AddAccountPictures(ctx context.Context, arg AddAccountPicturesParams) ([]string, error)
	AddFavourite(ctx context.Context, arg AddFavouriteParams) (Favourite, error)
	CheckAccountExists(ctx context.Context, arg CheckAccountExistsParams) (Account, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	DeleteAccount(ctx context.Context, id string) (string, error)
	DiscoverAccountsWithFilter(ctx context.Context, arg DiscoverAccountsWithFilterParams) ([]DiscoverAccountsWithFilterRow, error)
	DiscoverdAccountsMetadata(ctx context.Context, arg DiscoverdAccountsMetadataParams) (int64, error)
	FavoritesMetadata(ctx context.Context, id string) (int64, error)
	GetAccount(ctx context.Context, id string) (Account, error)
	GetAccountPictures(ctx context.Context, id string) ([]string, error)
	GetAllFavourites(ctx context.Context, arg GetAllFavouritesParams) ([]GetAllFavouritesRow, error)
	GetFavorite(ctx context.Context, arg GetFavoriteParams) (bool, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
}

var _ Querier = (*Queries)(nil)
