package wechat

import (
	"fmt"
	"gotrue/facility/http_util"
	"gotrue/facility/log"

	"github.com/mileusna/crontab"
)

type TokenManager struct {
	at      string
	atExpIn int64
	crontab *crontab.Crontab
}

func NewTokenManager(startSchedule bool) *TokenManager {
	m := &TokenManager{}
	m.crontab = crontab.New()
	if startSchedule {
		m.ScheduleTasks()
	}
	return m
}

// ScheduleTasks crontab syntax https://github.com/mileusna/crontab
func (m *TokenManager) ScheduleTasks() {
	log.Log.Debug("Shedule token refreshing task")
	m.RefreshAccessToken()
	m.crontab.MustAddJob("*/10 * * * *", m.RefreshAccessToken)
	// run imediately when start
	m.crontab.RunAll()
}

func (m *TokenManager) AccessToken() string {
	return m.at
}

func (m *TokenManager) RefreshAccessToken() {
	act := AccessTokenResonse{}
	url := fmt.Sprintf(accessToken_url, "client_credential", appID, secret)
	err := http_util.DoGet(&act, url)
	if err != nil {
		log.Log.Debug("Access token refresh error %+v\n", err)
	}
	if act.isOk() {
		m.at = act.AccessToken
		m.atExpIn = act.ExpiresIn
	}
	log.Log.Debug("Access token : %s\n", act.AccessToken)
}
