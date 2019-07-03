package middlewares

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/kataras/iris"
	_ "github.com/ziutek/mymysql/godrv"
	_ "github.com/ziutek/mymysql/native"
)

type Queue struct {
	Domain   string
	Weight   int
	Priority int
}

var Que map[int][]string

var requestPonderations map[string]int

type Repository interface {
	Read() map[string]int
}

func GetQue() []string {
	if Que == nil {
		return nil
	}
	var keys []int
	for k, _ := range Que {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	var q []string
	for _, k := range keys {
		v := Que[k]
		q = append(q, v...)
	}
	return q
}

func (q *Queue) Read() map[string]int {
	var final []*Queue
	if os.Getenv("DB_ENABLED") == "true" {
		final = readFromDb()
	} else {
		final = readFromFile()
	}
	requestPonderations = make(map[string]int)
	for _, element := range final {
		requestPonderations[element.Domain] = element.Priority * element.Weight
	}
	return requestPonderations
}

func readFromDb() []*Queue {
	const query = "SELECT * FROM proxy;"
	db, err := sql.Open("mymysql", os.Getenv("DB_CONNECTION_STRING"))
	fmt.Println("err", err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(query)
	fmt.Println("rows", rows)
	if err != nil {
		log.Fatal(err)
	}
	var final []*Queue
	for rows.Next() {
		tmp := &Queue{}
		err := rows.Scan(&tmp.Domain, &tmp.Weight, &tmp.Priority)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(tmp.Domain, tmp.Weight, tmp.Priority)
		final = append(final, tmp)
	}
	return final
}

func readFromFile() []*Queue {
	path, _ := filepath.Abs("")
	jsonFile, err := os.Open(path + "/api/middlewares/domain.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var final []*Queue
	json.Unmarshal(byteValue, &final)
	return final
}

func ProxyMiddleware(c iris.Context) {
	var repo Repository
	repo = &Queue{}
	repo.Read()
	if Que == nil {
		Que = make(map[int][]string)
	}
	domain := c.GetHeader("domain")
	val, ok := requestPonderations[domain]
	if len(domain) == 0 || !ok {
		c.JSON(iris.Map{"status": 400, "result": "domain error"})
		return
	}

	Que[val] = append(Que[val], domain)
	c.Next()
}
