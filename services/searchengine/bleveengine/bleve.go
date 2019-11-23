package bleveengine

import (
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/mattermost/mattermost-server/jobs"
	"github.com/mattermost/mattermost-server/model"
)

type BleveEngine struct {
	idx       bleve.Index
	cfg       *model.Config
	jobServer *jobs.JobServer
}

type BlevePost struct {
	Id        string   `json:"id"`
	TeamId    string   `json:"team_id"`
	ChannelId string   `json:"channel_id"`
	UserId    string   `json:"user_id"`
	CreateAt  int64    `json:"create_at"`
	Message   string   `json:"message"`
	Type      string   `json:"type"`
	Hashtags  []string `json:"hashtags"`
}

func NewBleveEngine(cfg *model.Config, license *model.License, jobServer *jobs.JobServer) (*BleveEngine, error) {
	mapping := bleve.NewIndexMapping()
	if cfg.BleveSettings.Filename != nil {
	}
	index, err := bleve.New(*cfg.BleveSettings.Filename, mapping)
	if err != nil {
		return nil, err
	}
	return &BleveEngine{
		idx:       index,
		cfg:       cfg,
		jobServer: jobServer,
	}, nil
}

func (b *BleveEngine) Start() *model.AppError {
	return nil
}

func (b *BleveEngine) Stop() *model.AppError {
	return nil
}

func (b *BleveEngine) IsActive() bool {
	return *b.cfg.BleveSettings.EnableIndexing
}

func (b *BleveEngine) GetVersion() int {
	return 0
}

func (b *BleveEngine) GetName() string {
	return "bleve"
}

func (b *BleveEngine) IndexPost(post *model.Post, teamId string) *model.AppError {
	searchPost := BlevePost{
		Id:        post.Id,
		TeamId:    teamId,
		ChannelId: post.ChannelId,
		UserId:    post.UserId,
		CreateAt:  post.CreateAt,
		Message:   post.Message,
		Type:      post.Type,
		Hashtags:  strings.Split(post.Hashtags, " "),
	}
	b.idx.Index(post.Id, searchPost)
	return nil
}

func (b *BleveEngine) SearchPosts(channels *model.ChannelList, searchParams []*model.SearchParams, page, perPage int) ([]string, model.PostSearchMatches, *model.AppError) {
	// b.idx.Search(post.id)
	return nil, nil, nil
}

func (b *BleveEngine) DeletePost(post *model.Post) *model.AppError {
	// b.idx.Delete(post.id)
	return nil
}

func (b *BleveEngine) IndexChannel(channel *model.Channel) *model.AppError {
	return nil
}

func (b *BleveEngine) SearchChannels(teamId, term string) ([]string, *model.AppError) {
	return nil, nil
}

func (b *BleveEngine) DeleteChannel(channel *model.Channel) *model.AppError {
	return nil
}

func (b *BleveEngine) IndexUser(user *model.User, teamsIds, channelsIds []string) *model.AppError {
	return nil
}

func (b *BleveEngine) SearchUsersInChannel(teamId, channelId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) ([]string, []string, *model.AppError) {
	return nil, nil, nil
}

func (b *BleveEngine) SearchUsersInTeam(teamId string, restrictedToChannels []string, term string, options *model.UserSearchOptions) ([]string, *model.AppError) {
	return nil, nil
}

func (b *BleveEngine) DeleteUser(user *model.User) *model.AppError {
	return nil
}

func (b *BleveEngine) TestConfig(cfg *model.Config) *model.AppError {
	return nil
}

func (b *BleveEngine) PurgeIndexes() *model.AppError {
	return nil
}

func (b *BleveEngine) DataRetentionDeleteIndexes(cutoff time.Time) *model.AppError {
	return nil
}