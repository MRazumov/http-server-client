package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/MRazumov/http-server-client/client/models"
	_ "github.com/lib/pq"
)

type Person struct {
	Name string
	Age  int
}
type Package struct {
	Person        *Person
	PersonFriends []*Person
}

func main() {

	//Отсылка запроса с клиента
	client := http.Client{}
	resp, err := client.Get("http://localhost:9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	//Считываем ответ от сервера
	var data Package
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Подключаемся к БД
	url := fmt.Sprintf("%s://%s:%s@%s:%d/%s", "postgres", "test", "1", "127.0.0.1", 5432, "test")
	db, err := sql.Open("postgres", url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	//Инициализируем менеджер (создаём пустую таблицу persons)
	persons := models.NewPersonManager(db)
	err = persons.CreatePersons()
	if err != nil {
		fmt.Println(err)
		return
	}

	//Инициализируем менеджер(создаём пустую таблицу person_friends)
	personFriends := models.NewPersonFriendsManager(db)
	err = personFriends.CreatePersonFriends()
	if err != nil {
		fmt.Println(err)
		return
	}

	//Заполняем таблицу persons информацией Person
	myid, err := persons.Add(&models.Person{Name: data.Person.Name, Age: data.Person.Age})
	if err != nil {
		fmt.Println(err)
		return
	}

	//Заполняем таблицу persons информацией PersonFriends
	for _, personFriend := range data.PersonFriends {
		id, err := persons.Add(&models.Person{Name: personFriend.Name, Age: personFriend.Age})
		if err != nil {
			fmt.Println(err)
			return
		}
		//Заполняем таблицу person_friends связями по id
		err = personFriends.Add(&models.PersonFriend{PersonId: myid, FriendId: id})
		if err != nil {
			fmt.Println(err)
			return
		}
		continue
	}
}
