package models

import (
	"database/sql"
)

type PersonFriend struct {
	PersonId int
	FriendId int
}

type PersonFriendsManager interface {
	Add(PersonFriend *PersonFriend) error
	CreatePersonFriends() error
}

type personFriendsManager struct {
	connection *sql.DB
}

func (p *personFriendsManager) CreatePersonFriends() error {
	_, err := p.connection.Query(`CREATE TABLE IF NOT EXISTS person_friends (
		person_id int,
		friend_id int)`)
	if err != nil {
		return err
	}
	return nil
}

func (p *personFriendsManager) Add(personFriend *PersonFriend) error {
	_, err := p.connection.Query("INSERT INTO person_friends (person_id, friend_id) SELECT $1, $2 WHERE NOT EXISTS (SELECT 1 FROM person_friends WHERE person_id=$1 AND friend_id=$2);", personFriend.PersonId, personFriend.FriendId)
	if err != nil {
		return err
	}
	return nil
}
func NewPersonFriendsManager(connection *sql.DB) PersonFriendsManager {
	return &personFriendsManager{connection}
}
