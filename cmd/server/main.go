package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/genieacs"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/gowa"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/mikrotik"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/tripay"
	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/whatsapp"
	impl "github.com/alijayanet/gembok-backend/internal/infrastructure/repositories"
	http "github.com/alijayanet/gembok-backend/internal/interface/http"
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
		logger.Fatal("Failed to connect to database")
	}

	db := database.GetDB()

	// ── Repositories ──────────────────────────────────────────────
	adminRepo := impl.NewAdminRepository(db)
	routerRepo := impl.NewRouterRepository(db)
	customerRepo := impl.NewCustomerRepository(db)
	invoiceRepo := impl.NewInvoiceRepository(db)
	packageRepo := impl.NewPackageRepository(db)
	onuRepo := impl.NewONULocationRepository(db)
	ticketRepo := impl.NewTroubleTicketRepository(db)
	settingRepo := impl.NewSettingRepository(db)

	// ── External clients ──────────────────────────────────────────
	genieacsClient := genieacs.NewGenieACSClient(
		cfg.GenieACS.URL,
		cfg.GenieACS.Username,
		cfg.GenieACS.Password,
	)

	mikrotikClient := mikrotik.NewMikroTikClient(routerRepo)
	if err := mikrotikClient.ConnectAll(); err != nil {
		logger.Warn("Failed to connect to all routers on startup", zap.Error(err))
	}

	mikrotikService := mikrotik.NewMikroTikService(mikrotikClient, customerRepo, packageRepo, routerRepo)

	gowaClient := gowa.NewGOWAClient(
		cfg.WhatsApp.APIURL,
		cfg.WhatsApp.APIKey,
		cfg.WhatsApp.DeviceID,
		cfg.WhatsApp.UseMock,
	)

	whatsappService := whatsapp.NewWhatsAppService(
		gowaClient,
		customerRepo,
		invoiceRepo,
		cfg.WhatsApp.AdminPhones,
	)

	tripayClient := tripay.NewTripayClient(
		cfg.Tripay.APIKey,
		cfg.Tripay.PrivateKey,
		cfg.Tripay.MerchantCode,
		cfg.Tripay.Mode,
	)

	// ── Use cases ────────────────────────────────────────────────
	authUsecase := usecase.NewAuthUsecase(adminRepo, cfg.JWT.Secret, cfg.JWT.Expiration)
	dashboardUsecase := usecase.NewDashboardUsecase(customerRepo, invoiceRepo, packageRepo)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo, mikrotikService, whatsappService)
	invoiceUsecase := usecase.NewInvoiceUsecase(invoiceRepo, settingRepo, whatsappService)
	routerUsecase := usecase.NewRouterUsecase(routerRepo, mikrotikClient)
	mikrotikUsecase := usecase.NewMikroTikUsecase(mikrotikService)
	genieacsUsecase := usecase.NewGenieACSUsecase(genieacsClient)
	paymentUsecase := usecase.NewPaymentUsecase(invoiceRepo, customerRepo, tripayClient, mikrotikService, cfg.App.URL)
	onuUsecase := usecase.NewONUUsecase(onuRepo, genieacsClient)
	ticketUsecase := usecase.NewTroubleTicketUsecase(ticketRepo, customerRepo)
	portalUsecase := usecase.NewPortalUsecase(customerRepo, invoiceRepo, ticketRepo, cfg.JWT.Secret)

	// ── Handlers ─────────────────────────────────────────────────
	authHandler := handlers.NewAuthHandler(authUsecase)
	dashboardHandler := handlers.NewDashboardHandler(dashboardUsecase)
	customerHandler := handlers.NewCustomerHandler(customerUsecase)
	invoiceHandler := handlers.NewInvoiceHandler(invoiceUsecase)
	routerHandler := handlers.NewRouterHandler(routerUsecase)
	mikrotikHandler := handlers.NewMikroTikHandler(mikrotikUsecase)
	genieacsHandler := handlers.NewGenieACSHandler(genieacsUsecase)
	paymentHandler := handlers.NewPaymentHandler(paymentUsecase)
	onuHandler := handlers.NewONUHandler(onuUsecase)
	ticketHandler := handlers.NewTroubleTicketHandler(ticketUsecase)
	portalHandler := handlers.NewPortalHandler(portalUsecase)
	whatsappHandler := handlers.NewWhatsAppHandler(whatsappService, cfg.WhatsApp.WebhookSecret, cfg.WhatsApp.AdminPhones)

	// ── Router ───────────────────────────────────────────────────
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

	database.Close()
	logger.Info("Server stopped")
}
