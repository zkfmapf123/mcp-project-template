package main

import (
	"log"
	"os"

	"github.com/go-mcp/internal/context"
	"github.com/go-mcp/internal/model"
	"github.com/go-mcp/pkg/protocol"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	dg "github.com/zkfmapf123/donggo"
)

type Server struct {
	contextManager *context.Manager
	app            *fiber.App

	// SimpleModels
	simpelModel  *model.SimpleModel
	chatGPTModel *model.ChatGPTModel
}

func NewServer() *Server {
	app := fiber.New(fiber.Config{
		AppName: "MCP Server",
	})

	server := &Server{
		contextManager: context.NewManager(),
		simpelModel:    model.NewSimpleModel("xxxxxxxxxxxxxxxxxxxxxxxxx"),
		chatGPTModel:   model.NewChatGPTModel(os.Getenv("CHATGPT_KEY")),
		app:            app,
	}

	app.Use(logger.New())
	app.Static("/", "./web/static")
	app.Post("/message", server.handleMessage)

	return server
}

func (s *Server) handleMessage(c *fiber.Ctx) error {
	message := dg.JsonParse[protocol.Message](c.Body())

	// create meesage Id (inject)
	messageId := context.GenerateID()
	message.ID = messageId

	ctx, err := s.contextManager.GetContext(message.ContextID)
	if err != nil {
		ctx, err = s.contextManager.CreateContext()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create context",
			})
		}
		message.ContextID = ctx.ID // ctx 대화에 넣기
	}

	// response, err := s.simpelModel.ProcessMessage(message, ctx)
	response, err := s.chatGPTModel.ProcessMessage(message, ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process message",
		})
	}

	ctx.Messages = append(ctx.Messages, message)          // 질문 추가
	ctx.Messages = append(ctx.Messages, response.Message) // 답변 추가
	if err := s.contextManager.UpdateContext(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update context",
		})
	}

	// utils.PrettyPrint(s.contextManager.GetMessages())

	return c.JSON(response)
}

func (s *Server) Start(addr string) error {
	return s.app.Listen(addr)
}

func main() {
	server := NewServer()

	log.Println("Server starting on :8000")
	if err := server.Start(":8000"); err != nil {
		log.Fatal(err)
	}
}
