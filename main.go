package main

import (
	"github.com/MMeier05/gator/internal/config"
	"github.com/MMeier05/gator/internal/database"
	_ "github.com/lib/pq"
	"database/sql"
	"log"
	"os"
)
type state struct {
	config *config.Config
	db *database.Queries
}


func main() {

	conf, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", conf.DbUrl)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	s := &state{
		config : &conf,
		db : database.New(db),
	}

	commands := commands{
		functions : map[string]func(*state, command) error{},
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", middlewareLoggedIn(handlerAgg))
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", handlerFollowing)
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatalf("Not enough arguments!")
	}
	cmd := command{
		name : os.Args[1],
		args : os.Args[2:],
	}
	err = commands.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
