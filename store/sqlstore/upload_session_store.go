// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package sqlstore

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store"
)

type SQLUploadSessionStore struct {
	*SQLStore
}

func newSQLUploadSessionStore(sqlStore *SQLStore) store.UploadSessionStore {
	s := &SQLUploadSessionStore{
		SQLStore: sqlStore,
	}
	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.UploadSession{}, "UploadSessions").SetKeys(false, "ID")
		table.ColMap("ID").SetMaxSize(26)
		table.ColMap("Type").SetMaxSize(32)
		table.ColMap("UserID").SetMaxSize(26)
		table.ColMap("ChannelID").SetMaxSize(26)
		table.ColMap("Filename").SetMaxSize(256)
		table.ColMap("Path").SetMaxSize(512)
		table.ColMap("RemoteID").SetMaxSize(26)
		table.ColMap("ReqFileID").SetMaxSize(26)
	}
	return s
}

func (us SQLUploadSessionStore) createIndexesIfNotExists() {
	us.CreateIndexIfNotExists("idx_uploadsessions_user_id", "UploadSessions", "Type")
	us.CreateIndexIfNotExists("idx_uploadsessions_create_at", "UploadSessions", "CreateAt")
	us.CreateIndexIfNotExists("idx_uploadsessions_user_id", "UploadSessions", "UserId")
}

func (us SQLUploadSessionStore) Save(session *model.UploadSession) (*model.UploadSession, error) {
	if session == nil {
		return nil, errors.New("SqlUploadSessionStore.Save: session should not be nil")
	}
	session.PreSave()
	if err := session.IsValid(); err != nil {
		return nil, errors.Wrap(err, "SqlUploadSessionStore.Save: validation failed")
	}
	if err := us.GetMaster().Insert(session); err != nil {
		return nil, errors.Wrap(err, "SqlUploadSessionStore.Save: failed to insert")
	}
	return session, nil
}

func (us SQLUploadSessionStore) Update(session *model.UploadSession) error {
	if session == nil {
		return errors.New("SqlUploadSessionStore.Update: session should not be nil")
	}
	if err := session.IsValid(); err != nil {
		return errors.Wrap(err, "SqlUploadSessionStore.Update: validation failed")
	}
	if _, err := us.GetMaster().Update(session); err != nil {
		if err == sql.ErrNoRows {
			return store.NewErrNotFound("UploadSession", session.ID)
		}
		return errors.Wrapf(err, "SqlUploadSessionStore.Update: failed to update session with id=%s", session.ID)
	}
	return nil
}

func (us SQLUploadSessionStore) Get(id string) (*model.UploadSession, error) {
	if !model.IsValidID(id) {
		return nil, errors.New("SqlUploadSessionStore.Get: id is not valid")
	}
	query := us.getQueryBuilder().
		Select("*").
		From("UploadSessions").
		Where(sq.Eq{"Id": id})
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "SqlUploadSessionStore.Get: failed to build query")
	}
	var session model.UploadSession
	if err := us.GetReplica().SelectOne(&session, queryString, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("UploadSession", id)
		}
		return nil, errors.Wrapf(err, "SqlUploadSessionStore.Get: failed to select session with id=%s", id)
	}
	return &session, nil
}

func (us SQLUploadSessionStore) GetForUser(userID string) ([]*model.UploadSession, error) {
	query := us.getQueryBuilder().
		Select("*").
		From("UploadSessions").
		Where(sq.Eq{"UserId": userID}).
		OrderBy("CreateAt ASC")
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "SqlUploadSessionStore.GetForUser: failed to build query")
	}
	var sessions []*model.UploadSession
	if _, err := us.GetReplica().Select(&sessions, queryString, args...); err != nil {
		return nil, errors.Wrap(err, "SqlUploadSessionStore.GetForUser: failed to select")
	}
	return sessions, nil
}

func (us SQLUploadSessionStore) Delete(id string) error {
	if !model.IsValidID(id) {
		return errors.New("SqlUploadSessionStore.Delete: id is not valid")
	}

	query := us.getQueryBuilder().
		Delete("UploadSessions").
		Where(sq.Eq{"Id": id})
	queryString, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "SqlUploadSessionStore.Delete: failed to build query")
	}

	if _, err := us.GetMaster().Exec(queryString, args...); err != nil {
		return errors.Wrap(err, "SqlUploadSessionStore.Delete: failed to delete")
	}

	return nil
}
