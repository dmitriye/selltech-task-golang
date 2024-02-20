package main

import (
	"log"
	"net/http"
	"os"
	"sdn/internal/app"
	"sdn/internal/handlers/search"
	"sdn/internal/handlers/state"
	"sdn/internal/handlers/update"
	"sdn/internal/helper"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const workersNum int = 5

func main() {
	dsn := os.Getenv("DATABASE_URL")

	db, err := helper.NewDatabaseConnection(dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	ofacApi := app.NewOfacAPI(http.DefaultClient)
	repo := app.NewEntryRepository(db)

	// init global state
	appState := app.NewAppState()
	if has, _ := repo.HasRecords(); has {
		appState.SetState(app.S_OK)
	}

	uploader := app.NewUploader(appState, repo, ofacApi, workersNum)

	http.HandleFunc("/update", update.NewUpdate(appState, uploader))
	http.HandleFunc("/state", state.NewState(appState))
	http.HandleFunc("/get_names", search.NewSearch(repo))

	log.Println("Server is running...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
