package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/yevhenshymotiuk/toad-bot-automator/scheduler"
	"github.com/zelenin/go-tdlib/client"
)

func main() {
	authorizer := client.ClientAuthorizer()
	go client.CliInteractor(authorizer)

	apiId, err := strconv.Atoi(os.Getenv("API_ID"))
	if err != nil {
		log.Fatalln(err)
	}

	apiHash := os.Getenv("API_HASH")

	authorizer.TdlibParameters <- &client.TdlibParameters{
		UseTestDc:              false,
		DatabaseDirectory:      filepath.Join(".tdlib", "database"),
		FilesDirectory:         filepath.Join(".tdlib", "files"),
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  int32(apiId),
		ApiHash:                apiHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "Server",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
		EnableStorageOptimizer: true,
		IgnoreFileNames:        false,
	}

	logVerbosity := client.WithLogVerbosity(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 0,
	})

	tdlibClient, err := client.NewClient(authorizer, logVerbosity)
	if err != nil {
		log.Fatalf("NewClient error: %s", err)
	}

	optionValue, err := tdlibClient.GetOption(&client.GetOptionRequest{
		Name: "version",
	})
	if err != nil {
		log.Fatalf("GetOption error: %s", err)
	}

	log.Printf(
		"TDLib version: %s",
		optionValue.(*client.OptionValueString).Value,
	)

	me, err := tdlibClient.GetMe()
	if err != nil {
		log.Fatalf("GetMe error: %s", err)
	}

	log.Printf("Me: %s %s [%s]", me.FirstName, me.LastName, me.Username)

	l, err := time.LoadLocation("Europe/Kiev")
	if err != nil {
		log.Fatalln(err)
	}
	// msgs, err := scheduler.ScheduleFeedings(
	// 	tdlibClient,
	// 	3,
	// 	6,
	// 	time.Date(2020, 12, 19, 17, 34, 0, 0, l),
	// )
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	msgs, err := scheduler.ScheduleWork(
		tdlibClient,
		3,
		time.Date(2020, 12, 19, 21, 3, 0, 0, l),
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(msgs)
}
