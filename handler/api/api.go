package api

import (
	"net/http"

	"github.com/fox-one/pkg/logger"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/handler/api/task"
)

type Server struct {
	Tasks    core.TaskStore
	Commands core.CommandStore
	Courses  core.CourseStore
}

func New(
	tasks core.TaskStore,
	commands core.CommandStore,
	courses core.CourseStore,
) *Server {
	return &Server{
		Tasks:    tasks,
		Commands: commands,
		Courses:  courses,
	}
}

func (s Server) Handle() http.Handler {
	router := gin.New()
	router.Use(
		gin.Recovery(),
		cors.AllowAll(),
		logger.Handler(),
	)

	router.POST(
		"/task",
		task.HandleCreate(s.Tasks, s.Courses),
	)

	router.POST(
		"/task/:id/active",
		task.TaskRequired(s.Tasks),
		task.HandleActive(
			s.Tasks,
			s.Courses,
			s.Commands,
		),
	)

	router.GET(
		"/task/:id",
		task.TaskRequired(s.Tasks),
		task.HandleFind(s.Courses),
	)

	return router
}
