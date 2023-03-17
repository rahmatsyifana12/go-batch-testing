package main

import (
	"fmt"
	"strings"
	"log"
	"github.com/jmoiron/sqlx"
	resty "github.com/go-resty/resty/v2"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID	uint	`gorm:"column:id; primaryKey"`
	Name string `gorm:"column:name" json:"name"`
	Age	int		`gorm:"column:age" json:"age"`
}

type DisposableEmailDomain struct {
	ID 			uint 	`db:"id"`
	EmailDomain	string	`db:"email_domain"`
}

func main() {
	// dsn := "root:kjsuhwkr@tcp(127.0.0.1:3306)/batch-testing?charset=utf8mb4&parseTime=True&loc=Local"
  	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	panic(err)
	// }

	// this Pings the database trying to connect
    // use sqlx.Open() for sql.Open() semantics
    // db, err := sql.Open("mysql", "root@tcp(localhost:3306)/batch-testing?parseTime=true")
	db, err := sqlx.Connect("mysql", "root:kjsuhwkr@/batch-testing")
    if err != nil {
        log.Fatalln(err)
    }
	defer db.Close()

    // exec the schema or fail; multi-statement Exec behavior varies between
    // database drivers;  pq will exec them all, sqlite3 won't, ymmv
    // db.MustExec(schema)

	dispEmailDomainList := []DisposableEmailDomain{}
	GetEmailDomainFromRepository(&dispEmailDomainList)

	fmt.Println(dispEmailDomainList)
	fmt.Println(len(dispEmailDomainList))

	// domainList := []DisposableEmailDomain{
	// 	{ EmailDomain: "yahoo.com" },
	// 	{ EmailDomain: "gmail.com" },
	// 	{ EmailDomain: "itqan.com" },
	// }

	// postgresql query
	// _, err = db.NamedExec(`INSERT INTO disposable_email_domains (email_domain) VALUES (:email_domain) ON CONFLICT (email_domain) DO NOTHING;`, dispEmailDomainList)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// mysql query
	_, err = db.NamedExec(`INSERT IGNORE INTO disposable_email_domains (email_domain) VALUES (:email_domain);`, dispEmailDomainList)
	if err != nil {
		log.Fatalln(err)
	}

	// e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })
	
	// e.Logger.Fatal(e.Start(":1234"))
}

// Get email domain data from github repository as raw content, then convert it to list of DisposableEmailDomain
func GetEmailDomainFromRepository(dispEmailDomainList *[]DisposableEmailDomain) {
	client := resty.New()
	url := "https://raw.githubusercontent.com/disposable-email-domains/disposable-email-domains/master/disposable_email_blocklist.conf"
	resp, err := client.R().EnableTrace().Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	stringRespBody := string(resp.Body())
	domainList := strings.Split(stringRespBody, "\n")

	for _, emailDomain := range domainList {
		if emailDomain != "" {
			*dispEmailDomainList = append(*dispEmailDomainList, DisposableEmailDomain{ EmailDomain: emailDomain })
		}
	}
}