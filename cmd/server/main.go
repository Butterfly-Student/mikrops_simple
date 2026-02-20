package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/audit"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/casbin"
	eventbus "github.com/alijayanet/gembok-backend/internal/infrastructure/external/eventbus"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/genieacs"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/gowa"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/mikrotik"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/tripay"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/whatsapp"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/repositories"
	"github.com/alijayanet/gembok-backend/internal/interface/http"
	"github.com/alijayanet/gembok-backend/internal/interface/http/handlers"
	"github.com/alijayanet/gembok-backend/internal/usecase"
	"github.com/alijayanet/gembok-backend/pkg/config"
	"github.com/alijayanet/gembok-backend/pkg/database"
	"github.com/alijayanet/gembok-backend/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := logger.InitLogger(cfg.Server.Mode); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db := database.GetDB()

	adminRepo := impl.NewAdminRepository(db)
	routerRepo := impl.NewRouterRepository(db)
	customerRepo := impl.NewCustomerRepository(db)
	invoiceRepo := impl.NewInvoiceRepository(db)
	packageRepo := impl.NewPackageRepository(db)
	onuLocationRepo := impl.NewONULocationRepository(db)
	ticketRepo := impl.NewTroubleTicketRepository(db)
	settingRepo := impl.NewSettingRepository(db)

	genieacsClient := genieacs.NewGenieACSClient(
		cfg.GenieACS.URL,
		cfg.GenieACS.Username,
		cfg.GenieACS.Password,
	)

	gowaClient := gowa.NewGOWAClient(
		cfg.WhatsApp.APIURL,
		cfg.WhatsApp.Token,
		cfg.WhatsApp.DeviceID,
		cfg.WhatsApp.UseMock,
	)

	tripayClient := tripay.NewTripayClient(
		cfg.Tripay.APIKey,
		cfg.Tripay.PrivateKey,
		cfg.Tripay.MerchantCode,
		cfg.Tripay.Mode,
	)

	whatsappService := whatsapp.NewWhatsAppService(
		gowaClient,
		customerRepo,
		invoiceRepo,
		cfg.WhatsApp.AdminPhones,
	)

	mikrotikClient := mikrotik.NewMikroTikClient(routerRepo)
	mikrotikService := mikrotik.NewMikroTikService(
		mikrotikClient,
		customerRepo,
		packageRepo,
		routerRepo,
	)

	authUsecase := usecase.NewAuthUsecase(adminRepo, cfg.JWT.Secret, cfg.JWT.Expiration)
	dashboardUsecase := usecase.NewDashboardUsecase(
		customerRepo,
		invoiceRepo,
		packageRepo,
	)
	customerUsecase := usecase.NewCustomerUsecase(
		customerRepo,
		mikrotikService,
		whatsappService,
	)
	invoiceUsecase := usecase.NewInvoiceUsecase(invoiceRepo, settingRepo, whatsappService)
	routerUsecase := usecase.NewRouterUsecase(routerRepo, mikrotikClient)
	onuUsecase := usecase.NewONUUsecase(onuLocationRepo, genieacsClient)
	ticketUsecase := usecase.NewTroubleTicketUsecase(ticketRepo, customerRepo)
	portalUsecase := usecase.NewPortalUsecase(
		customerRepo,
		invoiceRepo,
		ticketRepo,
		cfg.JWT.Secret,
	)
	tripayUsecase := usecase.NewPaymentUsecase(
		invoiceRepo,
		customerRepo,
		tripayClient,
		mikrotikService,
		cfg.App.URL,
	)
	genieacsUsecase := usecase.NewGenieACSUsecase(genieacsClient)
	mikrotikUsecase := usecase.NewMikroTikUsecase(mikrotikService)

	hotspotWrapper := mikrotik.NewHotspotClientWrapper(mikrotikClient, routerRepo)
	hotspotProfileUC := usecase.NewHotspotProfileUsecase(hotspotWrapper)
	hotspotUserUC := usecase.NewHotspotUserUsecase(hotspotWrapper)
	voucherUC := usecase.NewVoucherUsecase(hotspotWrapper)
	hotspotSaleUC := usecase.NewHotspotSaleUsecase(hotspotWrapper)
	hotspotSessionUC := usecase.NewHotspotSessionUsecase(hotspotWrapper)

	authHandler := handlers.NewAuthHandler(authUsecase)
	dashboardHandler := handlers.NewDashboardHandler(dashboardUsecase)
	customerHandler := handlers.NewCustomerHandler(customerUsecase)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceUsecase)
	routerHandler := handlers.NewRouterHandler(routerUsecase)
	mikrotikHandler := handlers.NewMikroTikHandler(mikrotikUsecase)
	genieacsHandler := handlers.NewGenieACSHandler(genieacsUsecase)
	paymentHandler := handlers.NewPaymentHandler(tripayUsecase)
	onuHandler := handlers.NewONUHandler(onuUsecase)
	ticketHandler := handlers.NewTroubleTicketHandler(ticketUsecase)
	portalHandler := handlers.NewPortalHandler(portalUsecase)
	whatsappHandler := handlers.NewWhatsAppHandler(whatsappService, cfg.WhatsApp.WebhookSecret, cfg.WhatsApp.AdminPhones)
	hotspotHandler := handlers.NewHotspotHandler(
		hotspotProfileUC,
		hotspotUserUC,
		voucherUC,
		hotspotSaleUC,
		hotspotSessionUC,
	)

	eventBus, err := eventbus.NewEventBus(&eventbus.EventConfig{
		Type: cfg.RBAC.EventSystem.Type,
		RabbitMQ: &eventbus.RabbitMQConfig{
			URL:        cfg.RBAC.EventSystem.RabbitMQ.URL,
			Exchange:   cfg.RBAC.EventSystem.RabbitMQ.Exchange,
			Queue:      cfg.RBAC.EventSystem.RabbitMQ.Queue,
			RoutingKey: cfg.RBAC.EventSystem.RabbitMQ.RoutingKey,
			Durable:    cfg.RBAC.EventSystem.RabbitMQ.Durable,
		},
	})
	if err != nil {
		log.Fatalf("Failed to initialize event bus: %v", err)
	}
	logger.Info("Event bus initialized", zap.String("type", cfg.RBAC.EventSystem.Type))

	auditService := audit.NewAuditService(db, &cfg.RBAC.Audit)
	logger.Info("Audit service initialized")

	if err := casbin.SeedDefaultSuperadmin(db, &cfg.RBAC.DefaultSuperAdmin); err != nil {
		logger.Warn("Failed to seed default superadmin", zap.Error(err))
	}

	casbinService, err := casbin.NewCasbinService(
		db,
		eventBus,
		auditService,
		&cfg.RBAC.Casbin,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Casbin: %v", err)
	}
	logger.Info("Casbin service initialized")

	casbinHandler := handlers.NewCasbinHandler(casbinService)

	router := http.SetupRouter(
		cfg,
		dashboardHandler,
		customerHandler,
		invoiceHandler,
		routerHandler,
		mikrotikHandler,
		genieacsHandler,
		authHandler,
		paymentHandler,
		onuHandler,
		ticketHandler,
		portalHandler,
		whatsappHandler,
		hotspotHandler,
		casbinHandler,
	)

	logger.Info("Starting server", zap.String("port", cfg.Server.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := router.Run(":" + cfg.Server.Port); err != nil {
			logger.Fatal("Failed to start server")
		}
	}()

	<-quit

	logger.Info("Shutting down server...")

	if err := eventBus.Close(); err != nil {
		logger.Error("Failed to close event bus", zap.Error(err))
	}

	database.Close()
	logger.Info("Server stopped")
}
