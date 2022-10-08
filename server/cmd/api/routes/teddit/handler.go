package teddit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"teniditter-server/cmd/api/routes"
	"teniditter-server/cmd/global/console"
	"teniditter-server/cmd/global/utils"

	"github.com/labstack/echo/v4"
)

func TedditHandler(t *echo.Group) {
	t.GET("/home", func(c echo.Context) error {
		res := routes.EchoWrapper{Context: c}

		homeFeedType := utils.SafeString(c.QueryParam("type"))
		afterId := utils.SafeString(c.QueryParam("afterId"))

		if !utils.IsEmptyString(homeFeedType) {
			valid := false
			for _, ftype := range []string{"hot", "new", "top", "rising", "controversial"} {
				if ftype == homeFeedType {
					valid = true
					break
				}
			}
			if !valid {
				return res.HandleResp(http.StatusBadRequest, "'type' query parameter is invalid")
			}
		} else {
			homeFeedType = "hot"
		}

		url := fmt.Sprintf("https://teddit.net/%s?api&raw_json=1", homeFeedType)
		if !utils.IsEmptyString(afterId) {
			url += fmt.Sprintf("&t=&after=t3_%s", afterId)
		}

		if !utils.IsValidURL(url) {
			return res.HandleResp(http.StatusForbidden, "This request has returned nothing")
		}

		resp, err := http.Get(url)
		if err != nil {
			return res.HandleResp(http.StatusForbidden, "This request has returned nothing")
		}
		defer resp.Body.Close()

		jsonBlob, err := io.ReadAll(resp.Body)
		if err != nil {
			return res.HandleResp(http.StatusForbidden, "This request has returned a corrupted response")
		}

		var jsonDatas map[string]interface{}
		err = json.Unmarshal(jsonBlob, &jsonDatas)
		if err != nil {
			return res.HandleResp(http.StatusForbidden, "This request has returned a corrupted response")
		}

		return res.HandleRespBlob(http.StatusOK, jsonDatas)
	})

	console.Log("TedditHandler Registered", console.Success)
}
