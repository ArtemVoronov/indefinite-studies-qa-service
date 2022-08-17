//go:build integration
// +build integration

package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/auth"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/notes"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/ping"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/tags"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/tasks"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/users"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/app"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Utils struct {
	asserts          TestAsserts
	entityGenerators TestEntityGenerators
}

var utils Utils = Utils{asserts: TestAsserts{}}

var TestRouter *gin.Engine

func TestMain(m *testing.M) {
	Setup()
	TestRouter = SetupRouter()
	code := m.Run()
	Shutdown()
	os.Exit(code)
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	authorized := r.Group("/")
	authorized.Use(app.AuthReqired())
	{
		authorized.GET("/safe-ping", ping.SafePing)
	}

	r.GET("/ping", ping.Ping)
	r.POST("/auth/login", auth.Authenicate)
	r.POST("/auth/refresh-token", auth.RefreshToken)

	r.GET("/tasks", tasks.GetTasks)
	r.GET("/tasks/:id", tasks.GetTask)
	r.POST("/tasks", tasks.CreateTask)
	r.PUT("/tasks/:id", tasks.UpdateTask)
	r.DELETE("/tasks/:id", tasks.DeleteTask)

	r.GET("/tags", tags.GetTags)
	r.GET("/tags/:id", tags.GetTag)
	r.POST("/tags", tags.CreateTag)
	r.PUT("/tags/:id", tags.UpdateTag)
	r.DELETE("/tags/:id", tags.DeleteTag)

	r.GET("/users", users.GetUsers)
	r.GET("/users/:id", users.GetUser)
	r.POST("/users", users.CreateUser)
	r.PUT("/users/:id", users.UpdateUser)
	r.DELETE("/users/:id", users.DeleteUser)

	r.GET("/notes", notes.GetNotes)
	r.GET("/notes/:id", notes.GetNote)
	r.POST("/notes", notes.CreateNote)
	r.PUT("/notes/:id", notes.UpdateNote)
	r.DELETE("/notes/:id", notes.DeleteNote)

	return r
}

func GetRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	d1 := path.Join(path.Dir(b))
	return d1[:len(d1)-len("/test/integration")]
}

func InitTestEnv() {
	if err := godotenv.Load(GetRootPath() + "/.env.test"); err != nil {
		fmt.Println("No .env.test file found")
	}
}

func RecreateTestDB() {
	// TODO: think about carelessness removing prod database
	cmd := exec.Command("docker-compose", "--env-file", "./.env.test", "--profile", "integration-tests-only", "up", "liquibase_rollback_all_and_create_db_again")
	cmd.Dir = GetRootPath()

	_ /*stdout*/, err := cmd.Output()

	if err != nil {
		fmt.Printf("error during recreating test DB: %v\n", err.Error())
		return
	}

	// uncomment for debugging
	// fmt.Println("-------------------------------------")
	// fmt.Println(string(stdout))
	// fmt.Println("-------------------------------------")
}

type TestFunc func(t *testing.T)

func RunWithRecreateDB(f TestFunc) func(t *testing.T) {
	RecreateTestDB()
	return func(t *testing.T) {
		f(t)
	}
}

func Setup() {
	InitTestEnv()
	auth.Setup()
	db.GetInstance()
}

func Shutdown() {
	defer db.GetInstance().GetDB().Close()
}
