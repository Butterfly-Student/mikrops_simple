package hotspot

import (
	"fmt"
	"strings"
	"time"
)

// RecordSale creates sales record as RouterOS system script
func (c *Client) RecordSale(sale *Sale) error {
	// Validate
	if sale.Username == "" || sale.Price <= 0 {
		return NewError("record sale", fmt.Errorf("username and price required"))
	}

	// Set default values
	if sale.Date == "" {
		sale.Date = time.Now().Format("Jan/02/2006")
	}
	if sale.Time == "" {
		sale.Time = time.Now().Format("15:04:05")
	}

	// Build script name
	scriptName := BuildSaleScriptName(sale)

	_, err := c.execute("/system/script/add",
		"=name="+scriptName,
		"=owner=hotspot-sales",
		"=policy=read,write,policy,test",
	)

	if err != nil {
		return NewError("record sale", err)
	}

	return nil
}

// GetAllSales retrieves all sales records from RouterOS
func (c *Client) GetAllSales(filter *SaleFilter) ([]Sale, error) {
	reply, err := c.execute("/system/script/print")
	if err != nil {
		return nil, NewError("get all sales", err)
	}

	sales := make([]Sale, 0)

	for _, re := range reply.Re {
		scriptName := re.Map["name"]

		// Check if script name matches sales format
		if !strings.Contains(scriptName, "-|-") {
			continue
		}

		sale, err := ParseSaleScriptName(scriptName)
		if err != nil {
			continue
		}

		sale.ScriptID = re.Map[".id"]

		// Apply date filter
		if filter != nil && filter.StartDate != "" && filter.EndDate != "" {
			saleDate, _ := time.Parse("Jan/02/2006", sale.Date)
			startDate, _ := time.Parse("Jan/02/2006", filter.StartDate)
			endDate, _ := time.Parse("Jan/02/2006", filter.EndDate)

			if saleDate.Before(startDate) || saleDate.After(endDate) {
				continue
			}
		}

		// Apply prefix filter
		if filter != nil && filter.Prefix != "" {
			if !strings.HasPrefix(sale.Username, filter.Prefix) {
				continue
			}
		}

		sales = append(sales, *sale)
	}

	// Apply pagination
	if filter != nil && filter.Limit > 0 {
		start := filter.Offset
		if start > len(sales) {
			start = len(sales)
		}
		end := start + filter.Limit
		if end > len(sales) {
			end = len(sales)
		}
		sales = sales[start:end]
	}

	return sales, nil
}

// GetSalesByDateRange retrieves sales within date range
func (c *Client) GetSalesByDateRange(startDate, endDate string) ([]Sale, error) {
	return c.GetAllSales(&SaleFilter{
		StartDate: startDate,
		EndDate:   endDate,
	})
}

// GetSalesByPrefix retrieves sales filtered by username prefix
func (c *Client) GetSalesByPrefix(prefix string) ([]Sale, error) {
	return c.GetAllSales(&SaleFilter{Prefix: prefix})
}

// GetTotalRevenue calculates total revenue
func (c *Client) GetTotalRevenue(startDate, endDate string) (float64, error) {
	sales, err := c.GetSalesByDateRange(startDate, endDate)
	if err != nil {
		return 0, err
	}

	totalRevenue := 0.0
	for _, sale := range sales {
		totalRevenue += sale.Price
	}

	return totalRevenue, nil
}

// DeleteSale deletes sales record
func (c *Client) DeleteSale(scriptID string) error {
	_, err := c.execute("/system/script/remove", "=.id="+scriptID)
	if err != nil {
		return NewError("delete sale", err)
	}

	return nil
}
