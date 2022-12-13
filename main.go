package main

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/gin-gonic/gin"
)

var mySql *sql.DB

const TOTAL_GRID = 91

func main() {
	mySql = ConnectMysql()
	ginServer := gin.Default()
	ginServer.POST("/api/toss", tossDice)
	/*
		request:
		{
			"game_id":"123",
			"player_id":"abc1"
		}
		response:{
			"retcode":0,
			"msg":"Success",
			"data":{
				"game_id":"123",
				"player_id":"abc1",
				"position":50
			}
		}
	*/
}

func ConnectMysql() *sql.DB {
	mysqlUrl := "用户名:密码@(地址:端口)/数据库名称"
	db, _ := sql.Open("mysql", mysqlUrl)
	defer db.Close()
	err := db.Ping()
	if err != nil {
		fmt.Println("mysql连接失败，错误日志为：", err.Error())
		return nil
	}
	return db
}

func tossDice(context *gin.Context) {
	rNum := rand.Intn(6) + 1
	gameID := context.PostForm("game_id")
	playerID := context.PostForm("player_id")
	curPosition, err := getPosition(gameID, playerID)
	if err != nil {
		context.JSON(200, gin.H{
			"retcode": -1,
			"msg":     err.Error(),
		})
	}
	finalPosition := curPosition + rNum
	if curPosition+rNum > TOTAL_GRID {
		finalPosition = TOTAL_GRID - (curPosition + rNum - TOTAL_GRID)
	}
	err = setPosition(gameID, playerID, finalPosition)
	if err != nil {
		context.JSON(200, gin.H{
			"retcode": -1,
			"msg":     err.Error(),
		})
	}
	context.JSON(200, gin.H{
		"retcode":  0,
		"msg":      "sucess",
		"position": finalPosition,
	})
}

func setPosition(gameID, playerID string, position int) error {
	_, err := mySql.Exec("update game_tab where game_id = ? and playerID = ? set position = ?", gameID, playerID, position)
	return err
}

func getPosition(gameID, playerID string) (int, error) {
	tab := Game{}
	err := mySql.QueryRow("select from game_tab where game_id=?,player_id=?", gameID, playerID).Scan(&tab)
	if err != nil {
		return 0, err
	}
	return tab.Position, nil
}

type Game struct {
	GameID   string
	PlayerID string
	Position int
}

type GameRecord struct {
	GameID     string
	PlayID     string
	SequenceID int
	Step       int
}
