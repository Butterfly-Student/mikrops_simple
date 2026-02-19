package http

import (
	"github.com/alijayanet/gembok-backend/internal/interface/http/handlers"
	"github.com/alijayanet/gembok-backend/internal/interface/http/middleware"
	"github.com/alijayanet/gembok-backend/pkg/config"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	cfg *config.Config,
	dashboardHandler *handlers.DashboardHandler,
	customerHandler *handlers.CustomerHandler,
	invoiceHandler *handlers.InvoiceHandler,
	routerHandler *handlers.RouterHandler,
	mikrotikHandler *handlers.MikroTikHandler,
	genieacsHandler *handlers.GenieACSHandler,
	authHandler *handlers.AuthHandler,
	paymentHandler *handlers.PaymentHandler,
	onuHandler *handlers.ONUHandler,
	ticketHandler *handlers.TroubleTicketHandler,
	portalHandler *handlers.PortalHandler,
	whatsappHandler *handlers.WhatsAppHandler,
) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.LoggingMiddleware())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// ----- Public routes (no auth) -----
	public := router.Group("/api")
	{
		// Auth
		public.POST("/auth/login", authHandler.Login)

		// Portal login (customers)
		public.POST("/portal/login", portalHandler.Login)

		// Tripay callback (webhook from Tripay, verified by signature)
		public.POST("/payment/callback", paymentHandler.TripayCallback)

		// WhatsApp webhook (no auth required, verified by signature)
		public.POST("/whatsapp/webhook", whatsappHandler.HandleWebhook)
		public.GET("/whatsapp/test", whatsappHandler.TestConnection)
	}

	// ----- Admin protected routes -----
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		// Auth
		api.POST("/auth/logout", authHandler.Logout)
		api.GET("/auth/me", authHandler.Me)

		// Dashboard
		api.GET("/dashboard", dashboardHandler.GetDashboardStats)

		// Customers
		api.GET("/customers", customerHandler.GetCustomers)
		api.GET("/customers/:id", customerHandler.GetCustomerByID)
		api.POST("/customers", customerHandler.CreateCustomer)
		api.PUT("/customers/:id", customerHandler.UpdateCustomer)
		api.DELETE("/customers/:id", customerHandler.DeleteCustomer)
		api.POST("/customers/:id/isolate", customerHandler.IsolateCustomer)
		api.POST("/customers/:id/activate", customerHandler.ActivateCustomer)
		api.POST("/customers/bulk-isolate", customerHandler.BulkIsolate)
		api.POST("/customers/bulk-activate", customerHandler.BulkActivate)
		api.POST("/customers/:id/sync", customerHandler.SyncCustomer)

		// Invoices
		api.GET("/invoices", invoiceHandler.GetInvoices)
		api.GET("/invoices/:id", invoiceHandler.GetInvoiceByID)
		api.POST("/invoices", invoiceHandler.CreateInvoice)
		api.PUT("/invoices/:id", invoiceHandler.UpdateInvoice)
		api.DELETE("/invoices/:id", invoiceHandler.DeleteInvoice)

		// Routers
		api.GET("/routers", routerHandler.GetRouters)
		api.GET("/routers/active", routerHandler.GetActive)
		api.GET("/routers/:id", routerHandler.GetRouter)
		api.POST("/routers", routerHandler.CreateRouter)
		api.PUT("/routers/:id", routerHandler.UpdateRouter)
		api.DELETE("/routers/:id", routerHandler.DeleteRouter)
		api.POST("/routers/:id/test", routerHandler.TestConnection)
		api.PUT("/routers/:id/activate", routerHandler.SetActive)
		api.GET("/routers/:id/status", routerHandler.GetStatus)
		api.GET("/routers/status/all", routerHandler.GetAllStatus)

		// MikroTik PPPoE
		api.GET("/mikrotik/ppp/users", mikrotikHandler.GetPPPUsers)
		api.GET("/mikrotik/ppp/active", mikrotikHandler.GetActiveSessions)
		api.GET("/mikrotik/ppp/profiles", mikrotikHandler.GetPPPProfiles)
		api.POST("/mikrotik/ppp/users", mikrotikHandler.AddPPPUser)
		api.PUT("/mikrotik/ppp/users/:username", mikrotikHandler.UpdatePPPUser)
		api.DELETE("/mikrotik/ppp/users/:username", mikrotikHandler.RemovePPPUser)
		api.POST("/mikrotik/ppp/users/:username/disconnect", mikrotikHandler.DisconnectUser)

		// MikroTik Hotspot & Traffic
		api.GET("/mikrotik/hotspot/logs", mikrotikHandler.GetHotspotLog)
		api.GET("/mikrotik/traffic", mikrotikHandler.GetTraffic)

		// GenieACS
		api.GET("/genieacs/devices", genieacsHandler.GetDevices)
		api.GET("/genieacs/devices/:serial", genieacsHandler.GetDevice)
		api.POST("/genieacs/devices/reboot", genieacsHandler.RebootDevice)
		api.POST("/genieacs/devices/parameter", genieacsHandler.SetParameter)
		api.GET("/genieacs/devices/find", genieacsHandler.FindDeviceByPPPoE)

		// Payment
		api.GET("/payment/gateways", paymentHandler.GetGateways)
		api.POST("/payment/create", paymentHandler.CreateTransaction)

		// ONU Locations & WiFi
		api.GET("/onu-locations", onuHandler.GetLocations)
		api.POST("/onu-locations", onuHandler.UpsertLocation)
		api.POST("/onu-wifi", onuHandler.SetWiFi)

		// Trouble Tickets (admin)
		api.GET("/tickets", ticketHandler.GetAll)
		api.GET("/tickets/:id", ticketHandler.GetByID)
		api.POST("/tickets", ticketHandler.Create)
		api.PUT("/tickets/:id", ticketHandler.Update)
	}

	// ----- Customer portal protected routes -----
	portal := router.Group("/api/portal")
	portal.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		portal.GET("/profile", portalHandler.GetProfile)
		portal.PUT("/password", portalHandler.ChangePassword)
		portal.GET("/invoices", portalHandler.GetInvoices)
		portal.GET("/tickets", portalHandler.GetTickets)
		portal.POST("/tickets", portalHandler.CreateTicket)
	}

	return router
}
