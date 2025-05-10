package main

import (
	"log"
	"os"

	ct "context"

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

	// Models
	chatGPTModel *model.ChatGPTModel

	// 취소 가능 여부
	activeConextModelIds map[string]ct.CancelFunc
}

func NewServer() *Server {
	app := fiber.New(fiber.Config{
		AppName: "MCP Server",
	})

	server := &Server{
		contextManager: context.NewManager(),
		// simpelModel:    model.NewSimpleModel("xxxxxxxxxxxxxxxxxxxxxxxxx"),
		chatGPTModel:         model.NewChatGPTModel(os.Getenv("CHATGPT_KEY")),
		app:                  app,
		activeConextModelIds: map[string]ct.CancelFunc{},
	}

	app.Use(logger.New())
	app.Static("/", "./web/static")
	app.Post("/message", server.handleMessage)
	app.Post("/stop", server.handleStop)

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

	// 취소 context 생성
	processCtx, cancel := ct.WithCancel(ct.Background())
	s.activeConextModelIds[message.ContextID] = cancel

	responseChan, errChan := make(chan protocol.ModelResponse), make(chan error)

	go func() {
		response, err := s.chatGPTModel.ProcessMessage(message, ctx)
		if err != nil {
			errChan <- err
		}
		responseChan <- response
	}()

	select {
	// 성공
	case resp := <-responseChan:
		delete(s.activeConextModelIds, message.ContextID)
		ctx.Messages = append(ctx.Messages, message)      // 질문 추가
		ctx.Messages = append(ctx.Messages, resp.Message) // 답변 추가

		if err := s.contextManager.UpdateContext(ctx); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update context",
			})
		}

		// utils.PrettyPrint(s.contextManager.GetMessages())
		return c.JSON(resp)

		// 실패
	case err := <-errChan:
		delete(s.activeConextModelIds, message.ContextID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
		// 끝
	case <-processCtx.Done():
		delete(s.activeConextModelIds, message.ContextID)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "stopped",
		})
	}
}

func (s *Server) handleStop(c *fiber.Ctx) error {
	message := dg.JsonParse[protocol.Message](c.Body())
	contextId := message.ContextID

	if contextId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if cancel, exists := s.activeConextModelIds[contextId]; exists {
		cancel()
		delete(s.activeConextModelIds, contextId)
		return c.JSON(fiber.Map{
			"status": "stopped",
		})
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "No active request found for this context",
	})
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
