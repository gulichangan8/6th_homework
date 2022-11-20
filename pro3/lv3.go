package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type newUser struct {
	username string
	password string
}

type mes struct {
	writeName string
	message   string
}

func main() {
	var dns = "root:040818@tcp(127.0.0.1:3306)/user?charset=utf8mb4&parseTime=True&loc=Local"
	var tem newUser
	var mess mes
	users := make(map[string]string)
	m := make(map[string]string)
	var res string

	db, err := sql.Open("mysql", dns)
	if err != nil {
		log.Fatalf("open mysql error : %v", err)
	}

	r := gin.Default()

	r.GET("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		rows, _ := db.Query("select * from user")
		for rows.Next() {
			rows.Scan(&tem.username, &tem.password)
			users[tem.username] = tem.password
		}
		if users[username] == password {
			c.String(200, "登陆成功")
		} else {
			c.String(200, "登陆失败")
		}
	})

	r.GET("/watch_message", func(c *gin.Context) {
		username := c.PostForm("username") //username此时是信息拥有者，即owner_name
		fmt.Println(username)
		row := db.QueryRow("select write_name,message from message where owner_name=?", username)
		e := row.Err()
		fmt.Println(e)
		err := row.Scan(&mess.writeName, &mess.message)
		fmt.Println(mess)
		fmt.Println(err)
		m[mess.writeName] = mess.message
		c.JSON(200, m)
	})

	r.GET("/write_message", func(c *gin.Context) {
		username := c.PostForm("username") //username此时是写信息的人，即write_name
		writeTo := c.PostForm("writeToName")
		message := c.PostForm("message")
		_, err := db.Exec("insert into message (owner_name, write_name, message,respond) values (?,?,?,?)",
			writeTo, username, message, "")
		fmt.Println(err)
		c.String(200, "留言成功")
	})

	r.GET("/respond_message", func(c *gin.Context) {
		username := c.PostForm("username") //username此时是信息拥有者，即owner_name
		respond := c.PostForm("respond")
		db.Exec("update message set respond=? where owner_name=?", respond, username)
		c.String(200, "回复成功")
	})

	r.GET("/respond_to_me", func(c *gin.Context) {
		username := c.PostForm("username") //username此时是写信息的人，即write_name
		row := db.QueryRow("select respond from message where write_name=?", username)
		row.Scan(&res)
		c.String(200, res)
	})

	r.Run()
}
