package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type question struct {
	username   string
	tureName   string
	likeFood   string
	lovePeople string
}

type User struct {
	username string
	password string
}

func main() {
	var dns = "root:040818@tcp(127.0.0.1:3306)/user?charset=utf8mb4&parseTime=True&loc=Local"
	var tem User
	var ques question
	users := make(map[string]string)

	db, err := sql.Open("mysql", dns)
	if err != nil {
		log.Fatalf("open mysql error : %v", err)
	}

	r := gin.Default()

	r.GET("/register", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		rows, _ := db.Query("select * from user")
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&tem.username, &tem.password)
			users[tem.username] = tem.password
		}
		_, ok := users[username]
		if ok {
			c.String(200, "注册失败（用户名已存在），请返回上一页面重新注册")
		} else {
			db.Exec("insert into user (username ,password) values (?,?)", username, password)
			c.String(200, "注册成功")
		}
	})

	r.GET("/password_question", func(c *gin.Context) {
		ques.username = c.PostForm("username")
		ques.tureName = c.PostForm("true_name")
		ques.likeFood = c.PostForm("like_food")
		ques.lovePeople = c.PostForm("love_people")

		res, err := db.Exec("insert into question (username, true_name ,like_food ,love_people) values (?,?,?,?)",
			ques.username, ques.tureName, ques.likeFood, ques.lovePeople)
		if err != nil {
			log.Printf("create question error : %v", err)
			c.JSON(200, "fail")
			return
		}
		_, err = res.RowsAffected()
		if err != nil {
			log.Printf("rowsAffected error : %v", err)
			c.JSON(200, "fail")
			return
		}
		c.String(200, "问题保存成功")
	})

	r.GET("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		rows, _ := db.Query("select * from user")
		defer rows.Close()
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

	r.GET("/change_password", func(c *gin.Context) {
		username := c.PostForm("username")
		tureName := c.PostForm("true_name")
		likeFood := c.PostForm("like_food")
		lovePeople := c.PostForm("love_people")
		row := db.QueryRow("select true_name,like_food,love_people from question where username=?", username)
		row.Scan(&ques.tureName, &ques.likeFood, &ques.lovePeople)
		if ques.tureName == tureName && ques.likeFood == likeFood && ques.lovePeople == lovePeople {
			c.String(200, "密保问题回答正确，请重新设置密码")
		} else {
			c.String(200, "密保问题回答错误")
		}
	})

	r.GET("/password", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		db.Exec("update user set password=? where username=?", password, username)
		c.String(200, "修改成功，请重新登录")
	})

	r.Run()
}
