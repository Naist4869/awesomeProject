package workwx

import (
	"net/http"
	"os"
	"testing"

	"github.com/Naist4869/awesomeProject/config"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

var testWorkApp *WorkwxApp

func TestMain(m *testing.M) {
	config.Init()
	appConfig, err := config.LoadAppConfig()
	if err != nil {
		panic(err)
	}
	workwx := New(appConfig.UseCase.WorkWx.CorpID, WithHTTPClient(&http.Client{}))
	testWorkApp = workwx.WithApp(appConfig.UseCase.WorkWx.CorpSecret, appConfig.UseCase.WorkWx.AgentID)
	run := m.Run()
	os.Exit(run)
}
func TestWorkwxApp_execUserGet(t *testing.T) {
	get, err := testWorkApp.execUserGet(reqUserGet{UserID: "LanJingCheng"})
	spew.Dump(get)
	require.NoError(t, err)
}
