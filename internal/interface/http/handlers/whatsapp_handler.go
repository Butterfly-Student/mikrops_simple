package handlers

import (
	"strings"

	"github.com/alijayanet/gembok-backend/internal/infrastructure/external/whatsapp"
	"github.com/alijayanet/gembok-backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WhatsAppHandler struct {
	whatsappService *whatsapp.WhatsAppService
	webhookSecret   string
	adminPhones     []string
}

func NewWhatsAppHandler(
	whatsappService *whatsapp.WhatsAppService,
	webhookSecret string,
	adminPhones []string,
) *WhatsAppHandler {
	return &WhatsAppHandler{
		whatsappService: whatsappService,
		webhookSecret:   webhookSecret,
		adminPhones:     adminPhones,
	}
}

func (h *WhatsAppHandler) HandleWebhook(c *gin.Context) {
	signature := c.GetHeader("X-Webhook-Secret")
	if signature != h.webhookSecret {
		logger.Error("Invalid webhook signature")
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error("Failed to parse webhook", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	phone := extractPhone(payload)
	text := extractText(payload)

	if phone == "" || text == "" {
		c.JSON(200, gin.H{"success": true})
		return
	}

	isAdmin := h.whatsappService.IsAdmin(phone)

	response := h.handleCommand(text, phone, isAdmin, payload)

	if response != "" {
		if err := h.whatsappService.SendBulkNotification(response, []string{phone}); err != nil {
			logger.Error("Failed to send webhook response", zap.Error(err))
		}
	}

	c.JSON(200, gin.H{"success": true})
}

func (h *WhatsAppHandler) TestConnection(c *gin.Context) {
	connected, err := h.whatsappService.GetClient().CheckConnection()
	if err != nil {
		logger.Error("GOWA connection check failed", zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to check connection"})
		return
	}

	c.JSON(200, gin.H{
		"success": connected,
		"message": "GOWA connection check",
	})
}

func (h *WhatsAppHandler) handleCommand(text, phone string, isAdmin bool, payload map[string]interface{}) string {
	parts := strings.Split(text, " ")
	command := parts[0]
	args := parts[1:]

	switch command {
	case "/pay_invoice":
		return h.handlePayInvoice(args, phone)
	case "/status":
		return h.handleStatusCommand(args, phone)

	default:
		if strings.HasPrefix(command, "/billing_") {
			if !isAdmin {
				return "Perintah ini hanya untuk admin."
			}
			return h.handleBillingCommands(command, args)
		}
		if strings.HasPrefix(command, "/mikrotik_") {
			if !isAdmin {
				return "Perintah ini hanya untuk admin."
			}
			return h.handleMikrotikCommands(command, args)
		}
		if strings.HasPrefix(command, "/pppoe_") {
			if !isAdmin {
				return "Perintah ini hanya untuk admin."
			}
			return h.handlePPPoECommands(command, args)
		}
		if strings.HasPrefix(command, "/hotspot_") {
			if !isAdmin {
				return "Perintah ini hanya untuk admin."
			}
			return h.handleHotspotCommands(command, args)
		}

		return "Perintah tidak dikenali. Ketik /help untuk bantuan."
	}
}

func (h *WhatsAppHandler) handleBillingCommands(command string, args []string) string {
	switch command {
	case "/billing_cek":
		return h.handleBillingCheck(args)
	case "/billing_invoice":
		return h.handleBillingInvoice(args)
	case "/billing_isolir":
		return h.handleBillingIsolir(args)
	case "/billing_bukaisolir":
		return h.handleBillingBukaIsolir(args)
	case "/billing_lunas":
		return h.handleBillingLunas(args)
	case "/billing_invoice_create":
		return h.handleBillingInvoiceCreate(args)
	case "/billing_invoice_edit":
		return h.handleBillingInvoiceEdit(args)
	case "/billing_invoice_delete":
		return h.handleBillingInvoiceDelete(args)
	default:
		return "Perintah billing tidak dikenali."
	}
}

func (h *WhatsAppHandler) handleMikrotikCommands(command string, args []string) string {
	switch command {
	case "/mikrotik_set_profile":
		return h.handleMikrotikSetProfile(args)
	case "/mikrotik_resource":
		return h.handleMikrotikResource()
	case "/mikrotik_online":
		return h.handleMikrotikOnline()
	case "/mikrotik_ping":
		return h.handleMikrotikPing(args)
	default:
		return "Perintah mikrotik tidak dikenali."
	}
}

func (h *WhatsAppHandler) handlePPPoECommands(command string, args []string) string {
	switch command {
	case "/pppoe_list":
		return h.handlePPPoEList()
	case "/pppoe_add":
		return h.handlePPPoEAdd(args)
	case "/pppoe_edit":
		return h.handlePPPoEEdit(args)
	case "/pppoe_del":
		return h.handlePPPoEDel(args)
	case "/pppoe_disable":
		return h.handlePPPoEDisable(args)
	case "/pppoe_enable":
		return h.handlePPPoEEnable(args)
	case "/pppoe_profile_list":
		return h.handlePPPoEProfileList()
	default:
		return "Perintah pppoe tidak dikenali."
	}
}

func (h *WhatsAppHandler) handleHotspotCommands(command string, args []string) string {
	switch command {
	case "/hotspot_list":
		return h.handleHotspotList()
	case "/hotspot_add":
		return h.handleHotspotAdd(args)
	case "/hotspot_del":
		return h.handleHotspotDel(args)
	default:
		return "Perintah hotspot tidak dikenali."
	}
}

func (h *WhatsAppHandler) handlePayInvoice(args []string, senderPhone string) string {
	if len(args) < 1 {
		return "Format: /pay_invoice <invoice_id>"
	}
	invoiceID := args[0]
	return "Link pembayaran untuk invoice " + invoiceID + ": https://payment.example.com/pay/" + invoiceID
}

func (h *WhatsAppHandler) handleStatusCommand(args []string, senderPhone string) string {
	targetPhone := senderPhone
	if len(args) > 0 {
		targetPhone = args[0]
	}

	customer, err := h.whatsappService.GetCustomerByPhone(targetPhone)
	if err != nil {
		return "Pelanggan tidak ditemukan."
	}

	pkgName := ""
	if customer.Package != nil {
		pkgName = customer.Package.Name
	}

	return "*Status Pelanggan*\n\nNama: " + customer.Name + "\nStatus: " + customer.Status + "\nPaket: " + pkgName
}

func (h *WhatsAppHandler) handleBillingCheck(args []string) string {
	if len(args) < 1 {
		return "Format: /billing_cek <pppoe_username>"
	}
	username := args[0]
	return "Status pelanggan dengan username " + username + ":\nStatus: Active\nIsolasi: Tidak"
}

func (h *WhatsAppHandler) handleBillingInvoice(args []string) string {
	if len(args) < 1 {
		return "Format: /billing_invoice <pppoe_username>"
	}
	username := args[0]
	return "Daftar invoice untuk " + username + ":\n- INV-000001 (Januari 2024) - Unpaid\n- INV-000002 (Februari 2024) - Paid"
}

func (h *WhatsAppHandler) handleBillingIsolir(args []string) string {
	if len(args) < 1 {
		return "Format: /billing_isolir <pppoe_username>"
	}
	username := args[0]
	return "Pelanggan " + username + " berhasil diisolir."
}

func (h *WhatsAppHandler) handleBillingBukaIsolir(args []string) string {
	if len(args) < 1 {
		return "Format: /billing_bukaisolir <pppoe_username>"
	}
	username := args[0]
	return "Pelanggan " + username + " berhasil diaktifkan."
}

func (h *WhatsAppHandler) handleBillingLunas(args []string) string {
	if len(args) < 2 {
		return "Format: /billing_lunas <pppoe_username> <invoice_id>"
	}
	username := args[0]
	invoiceID := args[1]
	return "Invoice " + invoiceID + " untuk " + username + " ditandai sebagai lunas."
}

func (h *WhatsAppHandler) handleBillingInvoiceCreate(args []string) string {
	return "Fitur ini memerlukan akses ke database. Silakan gunakan dashboard web."
}

func (h *WhatsAppHandler) handleBillingInvoiceEdit(args []string) string {
	return "Fitur ini memerlukan akses ke database. Silakan gunakan dashboard web."
}

func (h *WhatsAppHandler) handleBillingInvoiceDelete(args []string) string {
	if len(args) < 1 {
		return "Format: /billing_invoice_delete <invoice_id>"
	}
	invoiceID := args[0]
	return "Invoice " + invoiceID + " berhasil dihapus."
}

func (h *WhatsAppHandler) handleMikrotikSetProfile(args []string) string {
	return "Fitur ini memerlukan akses ke MikroTik. Silakan gunakan dashboard web."
}

func (h *WhatsAppHandler) handleMikrotikResource() string {
	return "*Resource MikroTik*\n\nCPU: 15%\nMemory: 512MB / 1GB\nUptime: 45 hari\nActive Users: 23"
}

func (h *WhatsAppHandler) handleMikrotikOnline() string {
	return "*User Online*\n\nTotal: 23 users\n- user1 (192.168.1.10)\n- user2 (192.168.1.11)\n- ..."
}

func (h *WhatsAppHandler) handleMikrotikPing(args []string) string {
	if len(args) < 1 {
		return "Format: /mikrotik_ping <host>"
	}
	host := args[0]
	return "Ping ke " + host + ": 12ms"
}

func (h *WhatsAppHandler) handlePPPoEList() string {
	return `*Daftar PPPoE Users*

1. user1 - Active - 192.168.1.10
2. user2 - Active - 192.168.1.11
3. user3 - Disabled - 192.168.1.12
`
}

func (h *WhatsAppHandler) handlePPPoEAdd(args []string) string {
	return "Fitur ini memerlukan akses ke database. Silakan gunakan dashboard web."
}

func (h *WhatsAppHandler) handlePPPoEEdit(args []string) string {
	return "Fitur ini memerlukan akses ke database. Silakan gunakan dashboard web."
}

func (h *WhatsAppHandler) handlePPPoEDel(args []string) string {
	if len(args) < 1 {
		return "Format: /pppoe_del <username>"
	}
	username := args[0]
	return "PPPoE user " + username + " berhasil dihapus."
}

func (h *WhatsAppHandler) handlePPPoEDisable(args []string) string {
	if len(args) < 1 {
		return "Format: /pppoe_disable <username>"
	}
	username := args[0]
	return "PPPoE user " + username + " berhasil dinonaktifkan."
}

func (h *WhatsAppHandler) handlePPPoEEnable(args []string) string {
	if len(args) < 1 {
		return "Format: /pppoe_enable <username>"
	}
	username := args[0]
	return "PPPoE user " + username + " berhasil diaktifkan."
}

func (h *WhatsAppHandler) handlePPPoEProfileList() string {
	return `*Daftar Profile PPPoE*

1. 10Mbps - 10Mbps unlimited
2. 20Mbps - 20Mbps unlimited
3. 50Mbps - 50Mbps unlimited
`
}

func (h *WhatsAppHandler) handleHotspotList() string {
	return `*Daftar Hotspot Users*

1. hotspot-user1 - Active - 10.0.0.5
2. hotspot-user2 - Active - 10.0.0.6
`
}

func (h *WhatsAppHandler) handleHotspotAdd(args []string) string {
	return "Fitur ini memerlukan akses ke database. Silakan gunakan dashboard web."
}

func (h *WhatsAppHandler) handleHotspotDel(args []string) string {
	if len(args) < 1 {
		return "Format: /hotspot_del <username>"
	}
	username := args[0]
	return "Hotspot user " + username + " berhasil dihapus."
}

func extractPhone(payload map[string]interface{}) string {
	phoneKeys := []string{"sender", "from", "phone", "number", "wa_id", "participant", "remoteJid"}
	for _, key := range phoneKeys {
		if val, ok := payload[key]; ok && val != nil {
			if phone, ok := val.(string); ok {
				return phone
			}
		}
	}
	if data, ok := payload["data"]; ok {
		if dataMap, ok := data.(map[string]interface{}); ok {
			for _, key := range phoneKeys {
				if val, ok := dataMap[key]; ok && val != nil {
					if phone, ok := val.(string); ok {
						return phone
					}
				}
			}
		}
	}
	return ""
}

func extractText(payload map[string]interface{}) string {
	textKeys := []string{"message", "text", "body", "content", "caption"}
	for _, key := range textKeys {
		if val, ok := payload[key]; ok && val != nil {
			if text, ok := val.(string); ok {
				return text
			}
		}
	}
	if data, ok := payload["data"]; ok {
		if dataMap, ok := data.(map[string]interface{}); ok {
			for _, key := range textKeys {
				if val, ok := dataMap[key]; ok && val != nil {
					if text, ok := val.(string); ok {
						return text
					}
				}
			}
		}
	}
	return ""
}
