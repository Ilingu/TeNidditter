package cron_routes

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"teniditter-server/cmd/api/routes"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"
	"teniditter-server/cmd/services/teddit"
	teddinitter "teniditter-server/cmd/services/tedinitter"

	"github.com/labstack/echo/v4"
)

func CronListener(cr *echo.Group) {
	// Receiver
	cr.POST("/", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		// Auth
		srvTrustKey := c.Request().Header.Get("Authorization")
		if utils.IsEmptyString(srvTrustKey) || utils.Hash(srvTrustKey) != os.Getenv("SERVER_TRUST_KEY") {
			return res.HandleResp(http.StatusUnauthorized, "Cannot Trust this Source")
		}

		console.Log("Receive New Cron Update", console.Info)
		go func() {
			// Refetch all ressource:
			for _, FeedType := range []string{"hot", "new", "top", "rising", "controversial"} {
				go teddit.GetHomePosts(FeedType, "", true) // teddit home page
			}

			// All User Feeds
			go teddinitter.GenerateFeeds()
		}()

		// Revalidation
		res.Response().Header().Set("Access-Control-Allow-Origin", "https://cronapi.up.railway.app") // Cors
		res.Response().Header().Set("Access-Control-Expose-Headers", "Continue")
		res.Response().Header().Set("Continue", "true")
		return res.String(200, "")
	})
}

type PayloadShape struct {
	Frequency   string `json:"Frequency"`
	CallbackUrl string `json:"CallbackUrl"`
}

func RegisterCron() {
	url, callbackUrl := "https://cronapi.up.railway.app/addJob", "https://tedinitterapi.up.railway.app/cron/"

	bodyPayload := PayloadShape{
		Frequency:   "@every 8h", // 3 refresh/day
		CallbackUrl: callbackUrl,
	}

	body, err := json.Marshal(bodyPayload)
	if err != nil {
		log.Println("Failed to Register Cron ❌")
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Failed To Create Cron Register Request ❌")
		return
	}
	req.Header = http.Header{"Content-Type": []string{"application/json"}, "Authorization": []string{os.Getenv("CRON_API_KEY")}}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Cron Register Request Failed ❌", err)
		return
	}

	if resp.StatusCode < 400 {
		log.Println("Cron Register Request Succeed ✅")
	} else {
		log.Println("Cron Register Request Failed ❌", resp.StatusCode)
	}
}
