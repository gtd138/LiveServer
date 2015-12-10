package main

// 服务器主程序
import (
	. "framework/server"
	. "gameserver/auth/app"
	. "gameserver/connector/app"
	. "gameserver/database/app"
	. "gameserver/game/app"
	. "gameserver/gate/app"
	. "gameserver/lobby/app"
	. "gameserver/master/app"
	. "gameserver/society/app"
	"log"
	"os"
)

func main() {
	app := CreateServerApp()
	//app := DebugCreateServerApp()
	app.Run()
}

// 调试创建服务器
func DebugCreateServerApp() (app IServer) {
	app = NewLobby("lobby")
	return
}

// 创建对应的服务器应用
func CreateServerApp() (app IServer) {
	if len(os.Args) < 2 {
		log.Println("请输入服务器类型")
		os.Exit(1)
	}
	server_type := os.Args[1]
	switch server_type {
	case "connector":
		app = NewConnector(server_type)
	case "gate":
		app = NewGate(server_type)
	case "master":
		app = NewMaster(server_type)
	case "game":
		app = NewGame(server_type)
	case "database":
		app = NewDataBase(server_type)
	case "society":
		app = NewSociety(server_type)
	case "lobby":
		app = NewLobby(server_type)
	case "auth":
		app = NewAuth(server_type)
	}

	if app == nil {
		log.Println("启动的服务器类型不存在！")
		os.Exit(1)
	}
	return
}
