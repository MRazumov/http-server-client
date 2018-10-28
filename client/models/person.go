package models

import (
	"database/sql"
)

type Person struct {
	Id   int `unique`
	Name string
	Age  int
}

type PersonManager interface {
	//Get(id int) (*Person, error)
	Add(Person *Person) (int, error)
	CreatePersons() error
}

type personManager struct {
	connection *sql.DB
}

func (p *personManager) CreatePersons() error {
	_, err := p.connection.Query(`CREATE TABLE IF NOT EXISTS persons (
		id serial UNIQUE NOT NULL PRIMARY KEY, 
		name varchar (100) UNIQUE NOT NULL,
		age int)`)
	if err != nil {
		return err
	}
	return nil
}

func (p *personManager) Add(person *Person) (int, error) {
	row := p.connection.QueryRow("INSERT INTO persons (name, age) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET age=$2 RETURNING id;", person.Name, person.Age)
	err := row.Scan(&person.Id)
	if err != nil {
		return 0, err
	}
	return person.Id, nil
}
func NewPersonManager(connection *sql.DB) PersonManager {
	return &personManager{connection}
}
