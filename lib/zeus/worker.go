package zeus

import (
	"os"
	"time"

	"go.uber.org/zap"
	"gopkg.in/resty.v1"
)

func TokenFetcher() {
	for {
		zap.S().Info("Fetching token")
		r, err := resty.New().R().Get(os.Getenv("TOKEN_FETCHER_URL"))
		if err != nil {
			zap.S().Error(err)
			time.Sleep(time.Second * 10)
			continue
		}

		if r.String() == "" {
			time.Sleep(time.Second * 10)
			continue
		}

		token = r.String()
		time.Sleep(time.Hour * 8)
	}
}
