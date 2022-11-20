package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	username string
	password string
}

func main() {
	var dns = "root:040818@tcp(127.0.0.1:3306)/user?charset=utf8mb4&parseTime=True&loc=Local"
	var tem user
	users := make(map[string]string)

	db, _ := sql.Open("mysql", dns)

	r := gin.Default()

	r.GET("/register", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		rows, _ := db.Query("select * from user")
		for rows.Next() {
			rows.Scan(&tem.username, &tem.password)
			users[tem.username] = tem.password
		}
		_, ok := users[username]
		if ok {
			c.String(200, "注册失败（用户名已存在），请返回上一页面重新注册")
		} else {
			db.Exec("insert into user (username ,password) value (?,?)", username, password)
			c.String(200, "注册成功")
		}
	})

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
	r.Run()
}
