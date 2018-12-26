package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"X-Token", "Content-Type"},
		//ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://127.0.0.1:8000"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.Use(static.Serve("/", static.LocalFile("../dist", true)))

	tokenLocalPool := make(map[string]int)

	r.POST("/user/login", func(c *gin.Context) {
		type requestBody struct {
			USERNAME string `json:"username"`
			PASSWORD string `json:"password"`
		}
		var rb requestBody
		var status int

		c.BindJSON(&rb)

		row := dbQuery(fmt.Sprintf("SELECT username FROM user_info WHERE username = '%s' AND password = md5('%s')", rb.USERNAME, rb.PASSWORD))
		defer row.Close()
		if row.Next() {
			// checking user status
			isLock := dbQuery(fmt.Sprintf("SELECT status FROM user_info WHERE username = '%s'", rb.USERNAME))
			defer isLock.Close()
			if isLock.Next() {
				err := isLock.Scan(&status)
				if err != nil {
					fmt.Println("check user status error...")
					c.JSON(200, gin.H{
						"code":          50011,
						"returnMessage": "check user status error...",
					})
					return
				}
				if status == 0 {
					c.JSON(200, gin.H{
						"code":          50012,
						"returnMessage": "user is locked...",
					})
					return
				}
			}

			// checking is already login
			isLogin := dbQuery(fmt.Sprintf("SELECT token FROM token_pool WHERE username = '%s'", rb.USERNAME))
			defer isLogin.Close()
			if isLogin.Next() {
				if !dbModify("DELETE FROM token_pool WHERE username = ?", rb.USERNAME) {
					fmt.Println("someone has logined and clear user token error...")
					c.JSON(200, gin.H{
						"code":          50005,
						"returnMessage": "someone has logined and clear user token error...",
					})
					return
				}
				var expireToken string
				err := isLogin.Scan(&expireToken)
				if err != nil {
					fmt.Println("someone has logined and delete local token error...")
					c.JSON(200, gin.H{
						"code":          50007,
						"returnMessage": "someone has logined and delete local token error...",
					})
					return
				}
				if _, ex := tokenLocalPool[expireToken]; ex {
					delete(tokenLocalPool, expireToken)
				}
			}

			type token struct {
				TOKEN string `json:"token"`
			}
			var t token
			t.TOKEN = generateToken()
			if dbModify("INSERT INTO token_pool(token, username) VALUES(?, ?)", t.TOKEN, rb.USERNAME) {
				// save token in local
				tokenLocalPool[t.TOKEN[:len(t.TOKEN)-1]] = 0
				c.JSON(200, gin.H{
					"code": 20000,
					//"returnMessage": "authorized success",
					"data": t,
				})
				return
			} else {
				c.JSON(200, gin.H{
					"code":          50002,
					"returnMessage": "save token failure...",
				})
				return
			}
		} else {
			c.JSON(200, gin.H{
				"code":          50001,
				"returnMessage": "authorized failure...",
			})
			return
		}
	})

	r.GET("/user/info", func(c *gin.Context) {
		type userInfo struct {
			ROLES  []string `json:"roles"`
			NAME   string   `json:"name"`
			AVATAR string   `json:"avatar"`
		}
		var (
			u        userInfo
			username string
			role     string
		)
		token := c.Query("token")
		sql := fmt.Sprintf("SELECT a.username, a.role FROM user_info a, token_pool b WHERE b.token = '%s' and a.username = b.username", token[:len(token)-1])
		rows := dbQuery(sql)
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&username, &role)
			if err != nil {
				fmt.Println("query user info error...")
				c.JSON(200, gin.H{
					"code": 50003,
					"data": nil,
				})
				return
			}
		}

		u.AVATAR = "https://cdn1.iconfinder.com/data/icons/ninja-things-1/1772/ninja-simple-512.png"
		u.NAME = username
		u.ROLES = []string{role}
		if len(username) > 0 && len(role) > 0 {
			c.JSON(200, gin.H{
				"code": 20000,
				"data": u,
			})
			return
		} else {
			c.JSON(200, gin.H{
				"code": 50008,
				"data": "ilegal token...",
			})
			return
		}
	})

	r.POST("/user/logout", func(c *gin.Context) {
		type token struct {
			TOKEN string `json:"token"`
		}
		var t token
		c.BindJSON(&t)
		if dbModify("DELETE FROM token_pool WHERE token = ?", t.TOKEN[:len(t.TOKEN)-1]) {
			// remove token from local token pool
			if _, e := tokenLocalPool[t.TOKEN[:len(t.TOKEN)-1]]; e {
				delete(tokenLocalPool, t.TOKEN[:len(t.TOKEN)-1])
			}
			//fmt.Println(tokenLocalPool)
			c.JSON(200, gin.H{
				"code": 20000,
				"data": "success",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"code": 50004,
				"data": "delete token error...",
			})
			return
		}
	})

	r.GET("/test", func(c *gin.Context) {
		token := c.GetHeader("X-Token")
		fmt.Println(token)
		if _, e := tokenLocalPool[token]; e {
			c.JSON(200, gin.H{
				"code": 20000,
			})
			return
		} else {
			c.JSON(200, gin.H{
				"code": 50008,
				"data": "ilegal token...",
			})
			return
		}
	})

	r.GET("/user", func(c *gin.Context) {
		token := c.GetHeader("X-Token")
		if _, e := tokenLocalPool[token]; !e {
			c.JSON(200, gin.H{
				"code": 50008,
				"data": "ilegal token...",
			})
			return
		}

		type items struct {
			USERNAME         string `json:"username"`
			ROLE             string `json:"role"`
			STATUS           string `json:"status"`
			LAST_MODIFY_TIME string `json:"last_modify_time"`
			EMAIL            string `json:"email"`
		}
		type data struct {
			ITEMS []items `json:"items"`
		}

		var (
			i                items
			s                []items
			d                data
			username         string
			role             string
			status           string
			last_modify_time string
			email            string
		)

		rows := dbQuery("SELECT username, role, status, last_modify_time, email FROM user_info")
		for rows.Next() {
			err := rows.Scan(&username, &role, &status, &last_modify_time, &email)
			if err != nil {
				fmt.Println("query user info error...")
				c.JSON(200, gin.H{
					"code": 50003,
					"data": nil,
				})
				return
			}
			i.USERNAME = username
			i.ROLE = role
			i.STATUS = status
			i.LAST_MODIFY_TIME = last_modify_time
			i.EMAIL = email
			s = append(s, i)
		}

		d.ITEMS = s

		c.JSON(200, gin.H{
			"code": 20000,
			"data": d,
		})
	})

	r.POST("/user/new", func(c *gin.Context) {
		type items struct {
			USERNAME         string `json:"username"`
			PASSWORD         string `json:"password"`
			ROLE             string `json:"role"`
			STATUS           string `json:"status"`
			LAST_MODIFY_TIME string `json:"last_modify_time"`
			EMAIL            string `json:"email"`
		}
		token := c.GetHeader("X-Token")

		if _, e := tokenLocalPool[token]; !e {
			c.JSON(200, gin.H{
				"code": 50008,
				"data": "ilegal token...",
			})
			return
		}

		var i items
		c.BindJSON(&i)
		i.LAST_MODIFY_TIME = time.Now().Format("2006-01-02 15:04:05")
		if dbModify("INSERT INTO user_info VALUES(null, ?, md5(?), ?, ?, ?, ?)",
			i.USERNAME,
			i.PASSWORD,
			i.ROLE,
			i.STATUS,
			i.LAST_MODIFY_TIME,
			i.EMAIL) {
			c.JSON(200, gin.H{
				"code": 20000,
				"data": "add user success",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"code": 50006,
				"data": "add user error...",
			})
			return
		}
	})

	r.POST("/user/edit", func(c *gin.Context) {
		type items struct {
			USERNAME         string `json:"username"`
			PASSWORD         string `json:"password"`
			ROLE             string `json:"role"`
			STATUS           string `json:"status"`
			LAST_MODIFY_TIME string `json:"last_modify_time"`
			EMAIL            string `json:"email"`
			ORIGIN_USERNAME  string `json:"origin_username"`
		}
		token := c.GetHeader("X-Token")

		if _, e := tokenLocalPool[token]; !e {
			c.JSON(200, gin.H{
				"code": 50008,
				"data": "ilegal token...",
			})
			return
		}

		var i items
		c.BindJSON(&i)
		i.LAST_MODIFY_TIME = time.Now().Format("2006-01-02 15:04:05")
		if dbModify("UPDATE user_info SET username = ?, password = md5(?), role = ?, status = ?, last_modify_time = ?, email = ? where username = ?",
			i.USERNAME,
			i.PASSWORD,
			i.ROLE,
			i.STATUS,
			i.LAST_MODIFY_TIME,
			i.EMAIL,
			i.ORIGIN_USERNAME) {
			c.JSON(200, gin.H{
				"code": 20000,
				"data": "modify user info success",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"code": 50009,
				"data": "edit user info error...",
			})
			return
		}
	})

	r.POST("/user/delete", func(c *gin.Context) {
		token := c.GetHeader("X-Token")
		if _, e := tokenLocalPool[token]; !e {
			c.JSON(200, gin.H{
				"code": 50008,
				"data": "ilegal token...",
			})
			return
		}
		type u struct {
			USERNAME string `json:"username"`
		}
		var user u
		c.BindJSON(&user)
		if dbModify("DELETE FROM user_info WHERE username = ?", user.USERNAME) {
			c.JSON(200, gin.H{
				"code": 20000,
				"data": "delete user success",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"code": 50010,
				"data": "delete user error...",
			})
			return
		}
	})

	r.Run(":8000")
}
