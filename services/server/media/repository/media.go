package repository

import (
	"services/server/media/model"
	"services/server/core/util"
	"log"
	"github.com/Masterminds/squirrel"
	"services/server/core/repository"
	"errors"
)

type MediaData struct {
	Media  model.Media
	Tracks []model.Track
	Tags   []model.Hashtag
}

type MediaRepository struct {
	repo *repository.Repository
}

func NewMediaRepository(repo *repository.Repository) (*MediaRepository) {
	return &MediaRepository{
		repo: repo,
	}
}

func (repo *MediaRepository) List(selectPage *util.SelectPage, where map[string]interface{}) ([]MediaData, error) {
	var data []MediaData

	query := squirrel.Select("*").From("media")
	if where != nil {
		query = query.
			Where(where)
	}
	if selectPage != nil {
		offset := selectPage.Page * selectPage.Count
		query = query.
			Offset(uint64(offset)).
			Limit(uint64(selectPage.Count))
	}

	rows, err := query.
		RunWith(repo.repo.DB()).
		Query();
	if err != nil {
		log.Print("failed to get media info for user", err)
		return nil, errors.New("failed to find media information")
	}
	parsed := repo.repo.Parse(&model.Media{}, rows)
	medias := make([]model.Media, len(parsed))
	for i, m := range parsed {
		medias[i] = m.(model.Media)
	}

	for _, media := range medias {

		// fetch the tracks for each media
		rows, err = squirrel.Select("*").
			From("track").
			Where(squirrel.Eq{"media_id": media.ID}).
			RunWith(repo.repo.DB()).
			Query()
		if err != nil {
			log.Print("failed to get media info for user", err)
			return nil, errors.New("failed to find media information")
		}
		parsed := repo.repo.Parse(&model.Track{}, rows)
		tracks := make([]model.Track, len(parsed))
		for i, m := range parsed {
			tracks[i] = m.(model.Track)
		}

		// fetch the hashtags
		rows, err = squirrel.Select("h.id, h.created_at, h.updated_at, h.deleted_at, h.tag").
			From("media_hashtag mh").
			Join("hashtag h ON mh.hashtag_id = h.id").
			Where(squirrel.Eq{"media_id": media.ID}).
			RunWith(repo.repo.DB()).
			Query()
		if err != nil {
			log.Print("failed to get media info for user", err)
			return nil, errors.New("failed to find media information")
		}
		parsed = repo.repo.Parse(&model.Hashtag{}, rows)
		hashtags := make([]model.Hashtag, len(parsed))
		for i, m := range parsed {
			hashtags[i] = m.(model.Hashtag)
		}

		data = append(data, MediaData{
			Media: media,
			Tracks: tracks,
			Tags: hashtags,
		})
	}

	return data, nil
}

func (repo *MediaRepository) GetByChannelID(channelID string, selectPage *util.SelectPage) ([]MediaData, error) {
	return repo.List(selectPage, squirrel.Eq{
		"channel_id": channelID,
	})
}

func (repo *MediaRepository) GetByVideoID(videoId string) (*MediaData, error) {
	medias, err := repo.List(nil, squirrel.Eq{
		"id": videoId,
	})
	if err != nil {
		return nil, err
	}

	if len(medias) == 0 {
		return nil, errors.New("could not find video for video id")
	}

	return &medias[0], nil
}

func (repo *MediaRepository) Save(media *MediaData) (error) {
	tx, err := repository.GetRepository().DB().Begin()
	if err != nil {
		log.Print(err)
		return err
	}

	// save the media data
	err = repository.GetRepository().Save(&media.Media)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}

	// delete all of the old associated media hashtags
	_, err = squirrel.Delete("media_hashtag").Where(squirrel.Eq{
		"media_id": media.Media.ID,
	}).RunWith(repository.GetRepository().DB()).Exec()
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}

	// this will save the tags into the database
	for _, tag := range media.Tags {
		err = repository.GetRepository().Save(&tag)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return err
		}

		// create the new relationships for the tags
		err = repository.GetRepository().Save(&model.MediaHashtag{
			MediaID: media.Media.ID,
			HashtagID: tag.ID,
		})
		if err != nil {
			tx.Rollback()
			log.Println(err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
