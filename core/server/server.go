package server

import (
	"cal-salary/core/cache"
	"cal-salary/core/config"
	"cal-salary/core/database"
	inmemcache "cal-salary/core/inmem_cache"
	"cal-salary/core/logger"
	"cal-salary/core/messaging"
	"cal-salary/core/middleware"
	"cal-salary/core/seed"
	coreStorage "cal-salary/core/storage"
	"cal-salary/core/utils"
	"cal-salary/core/google"
	"cal-salary/modules/activity_log"
	"cal-salary/modules/auth"
	"cal-salary/modules/meeting"
	"cal-salary/modules/payroll"
	"cal-salary/modules/pdf"
	"cal-salary/modules/timekeeping"
	"cal-salary/core/notification"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo       *echo.Echo
	addr       string
	cache      *cache.Cache
	db         database.Database
	natsClient *messaging.NatsClient
}

func initEnvironment() (config.Environment, error) {
	env := flag.String("env", "dev", "Environment (dev/prod)")
	flag.Parse()

	switch *env {
	case "dev":
		return config.DevEnvironment, nil
	case "prod":
		return config.ProdEnvironment, nil
	default:
		return "", fmt.Errorf("invalid environment. Use 'dev' or 'prod'")
	}
}

func initServer() (*Server, error) {
	environment, err := initEnvironment()
	if err != nil {
		return nil, err
	}

	if errInitConfig := config.Init(environment); errInitConfig != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", errInitConfig)
	}

	// Get config safely
	cfg, isInitialized := config.GetSafe()
	if !isInitialized {
		return nil, fmt.Errorf("config was not properly initialized")
	}

	// Initialize logger first, before validation
	if errInitLogger := logger.Init(logger.LogConfig{
		Level:         logger.LogLevelDebug,
		EnableFile:    true,
		JSONFormat:    true, // Thêm dòng này để log dạng JSON
		DailyRotation: true, // Có thể thêm daily rotation nếu muốn
	}); errInitLogger != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", errInitLogger)
	}

	// Validate configuration after logger is initialized
	if err = cfg.Validate(); err != nil {
		logger.Error("Configuration validation failed", "error", err)
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Initialize database
	db, err := database.InitDB(database.DatabaseConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize Redis cache
	redisCache := cache.NewCache(
		cfg.Redis.Address,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)

	// Initialize NATS JetStream client (optional for local dev: log warning if missing)
	var natsClient *messaging.NatsClient
	natsClient, err = messaging.NewNatsClient(cfg.Nats.Url)
	if err != nil {
		logger.Warn("NATS connection failed, running in local development mode without NATS JetStream", "error", err)
		natsClient = nil
	}

	// Initialize in-memory cache
	inmemCache := inmemcache.NewInMemoryCache()
	if inmemCache == nil {
		return nil, fmt.Errorf("failed to initialize in-memory cache")
	}

	// Initialize Email Config
	emailConfig := utils.EmailConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
		FromName: cfg.SMTP.FromName,
	}
	utils.InitEmailConfig(emailConfig)

	// Initialize MinIO client
	if err := coreStorage.InitMinIOClient(); err != nil {
		logger.Warn("Failed to initialize MinIO client", "error", err)
		// Không dừng server nếu MinIO init thất bại, chỉ log warning
	}

	// Seed initial data
	seedCtx := context.Background()
	if err := seed.SeedWorkStandards(seedCtx, db); err != nil {
		logger.Warn("Failed to seed work standards", "error", err)
		// Không dừng server nếu seeding thất bại, chỉ log warning
	}
	// Seed type competencies trước vì competency dictionaries phụ thuộc vào nó
	if err := seed.SeedTypeCompetencies(seedCtx, db); err != nil {
		logger.Warn("Failed to seed type competencies", "error", err)
		// Không dừng server nếu seeding thất bại, chỉ log warning
	}
	if err := seed.SeedCompetencyDictionaries(seedCtx, db); err != nil {
		logger.Warn("Failed to seed competency dictionaries", "error", err)
		// Không dừng server nếu seeding thất bại, chỉ log warning
	}
	if err := seed.SeedDepartments(seedCtx, db); err != nil {
		logger.Warn("Failed to seed departments", "error", err)
		// Không dừng server nếu seeding thất bại, chỉ log warning
	}
	if err := seed.SeedUserProfiles(seedCtx, db); err != nil {
		logger.Warn("Failed to seed user profiles", "error", err)
		// Không dừng server nếu seeding thất bại, chỉ log warning
	}
	if err := seed.SeedAttendanceData(seedCtx, db); err != nil {
		logger.Warn("Failed to seed attendance data", "error", err)
		// Không dừng server nếu seeding thất bại, chỉ log warning
	}
	if err := seed.SeedJobPositions(seedCtx, db); err != nil {
		logger.Warn("Failed to seed job positions", "error", err)
		// Không dừng server nếu seeding thất bại, chỉ log warning
	}
	if err := seed.SeedPermissions(seedCtx, db); err != nil {
		logger.Warn("Failed to seed permissions", "error", err)
		// Không dừng server nếu seeding thất bại, chỉ log warning
	}
	if err := seed.SeedRolesAndAdminUserRole(seedCtx, db); err != nil {
		logger.Warn("Failed to seed roles and admin user role", "error", err)
		// Không dừng server nếu seeding thất bại, chỉ log warning
	}
	if err := seed.SeedEmployeeUsers(seedCtx, db); err != nil {
		logger.Warn("Failed to seed employee user accounts", "error", err)
		// Không dừng server nếu seeding thất bại, chỉ log warning
	}

	logger.Info("Server initializing",
		"environment", environment,
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
		"database_host", cfg.Database.Host,
		"database_port", cfg.Database.Port,
		"database_name", cfg.Database.DBName,
		"redis_address", cfg.Redis.Address,
		"redis_db", cfg.Redis.DB,
		"smtp_host", cfg.SMTP.Host,
		"smtp_port", cfg.SMTP.Port,
	)

	e := echo.New()

	_, ipnet1, _ := net.ParseCIDR("127.0.0.1/32")
	_, ipnet2, _ := net.ParseCIDR("10.0.0.0/8")
	_, ipnet3, _ := net.ParseCIDR("172.16.0.0/12")
	_, ipnet4, _ := net.ParseCIDR("192.168.0.0/16")

	// Cấu hình IPExtractor để lấy đúng client IP khi chạy sau Proxy/Load Balancer
	e.IPExtractor = echo.ExtractIPFromXFFHeader(
		echo.TrustIPRange(ipnet1),
		echo.TrustIPRange(ipnet2),
		echo.TrustIPRange(ipnet3),
		echo.TrustIPRange(ipnet4),
	)

	// Middleware
	e.Use(middleware.LoggerMiddleware())
	e.Use(middleware.CORSMiddleware())
	e.Use(echomiddleware.Recover()) // Tránh server crash khi có panic trong handler

	e.Use(echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == "POST" || c.Request().Method == "PUT" {
				c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, 32<<20) // 32MB
			}
			return next(c)
		}
	}))

	// Initialize modules
	authMod := auth.Init(db, *redisCache)
	middlewareInstance := middleware.NewMiddleware(authMod.Service)

	activityLogSvc := activity_log.Init(e, db, middlewareInstance)
	authMod.SetupRouter(e, middlewareInstance, activityLogSvc)

	payrollMod := payroll.Init(db, redisCache)
	payrollMod.SetupRouter(e, middlewareInstance, activityLogSvc)

	timekeepingMod := timekeeping.Init(db, redisCache, natsClient, cfg.Nats)
	timekeepingMod.SetupRouter(e, middlewareInstance, activityLogSvc)

	pdfMod := pdf.Init()
	pdfMod.SetupRouter(e, middlewareInstance)
	if err := pdfMod.StartWatcher(context.Background()); err != nil {
		logger.Warn("PDFModule: Failed to start folder watcher", "error", err)
	}

	// Initialize Meeting Module (Leader Schedule, Google Calendar & Video Room)
	gcalSvc, _ := google.NewCalendarService(context.Background())
	meetingMod := meeting.InitMeetingModule(&db, gcalSvc, natsClient)
	apiV1Group := e.Group("/api/v1")
	meetingMod.RegisterRoutes(apiV1Group, middlewareInstance.AuthMiddleware())

	// Subscribe to meeting.notification via NATS (if client is active)
	if natsClient != nil && natsClient.Conn != nil {
		go func() {
			_, err := natsClient.Conn.Subscribe("meeting.notification", func(msg *nats.Msg) {
				var event struct {
					Type        string      `json:"type"`
					MeetingID   string      `json:"meeting_id"`
					Title       string      `json:"title"`
					AttendeeIDs []uuid.UUID `json:"attendee_ids"`
					HostID      uuid.UUID   `json:"host_id"`
					UserID      uuid.UUID   `json:"user_id"`
					UserName    string      `json:"user_name"`
					Status      string      `json:"status"`
				}
				if err := json.Unmarshal(msg.Data, &event); err != nil {
					return
				}

				if event.Type == "meeting.created" {
					for _, attID := range event.AttendeeIDs {
						notification.GlobalHub.Publish(attID, notification.Notification{
							ID:        uuid.New().String(),
							Text:      fmt.Sprintf("Bạn được mời tham gia cuộc họp: %s", event.Title),
							Time:      "Vừa xong",
							Read:      false,
							MeetingID: event.MeetingID,
						})
					}
				} else if event.Type == "meeting.rsvp" {
					actionText := "từ chối"
					if event.Status == "ACCEPTED" {
						actionText = "đồng ý"
					}
					notification.GlobalHub.Publish(event.HostID, notification.Notification{
						ID:        uuid.New().String(),
						Text:      fmt.Sprintf("%s đã %s tham gia cuộc họp: %s", event.UserName, actionText, event.Title),
						Time:      "Vừa xong",
						Read:      false,
						MeetingID: event.MeetingID,
					})
				}
			})
			if err != nil {
				logger.Error("NATS: failed to subscribe to meeting.notification", err)
			} else {
				logger.Info("NATS: successfully subscribed to meeting.notification")
			}
		}()
	}


	return &Server{
		echo:       e,
		addr:       fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		cache:      redisCache,
		db:         db,
		natsClient: natsClient,
	}, nil
}

func Run() error {
	srv, err := initServer()
	if err != nil {
		// Use fmt.Printf instead of logger since logger might not be initialized
		fmt.Printf("Failed to initialize server: %v\n", err)
		return err
	}
	return srv.start()
}

func (s *Server) start() error {
	logger.Info("Starting HTTP server", "address", s.addr)

	go func() {
		if err := s.echo.Start(s.addr); err != nil {
			logger.Info("Shutting down server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server gracefully: %w", err)
	}

	// Close Redis connection
	if err := s.cache.Close(); err != nil {
		logger.Error("Failed to close Redis connection", "error", err)
	}

	// Drain NATS connection (flushes in-flight publishes/acks before closing)
	if s.natsClient != nil {
		s.natsClient.Close()
	}

	logger.Info("Server shutdown complete")
	return nil
}
