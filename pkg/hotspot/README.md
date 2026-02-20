# Hotspot Management Library

Reusable MikroTik Hotspot Management Library for Go.

## Features

- Profile Management (create, update, delete, list)
- User Management (create, update, delete, list, filter)
- Voucher Generation (batch generation with random credentials)
- Active Sessions Monitoring
- Sales Recording (RouterOS system scripts)
- Expiry Management (scheduler-based)

## Architecture

All data is stored directly in MikroTik RouterOS:
- Profiles: `/ip/hotspot/user/profile`
- Users: `/ip/hotspot/user`
- Sessions: `/ip/hotspot/active`
- Sales: `/system/script` (named with special format)
- Schedulers: `/system/scheduler`

No additional database required.

## Usage

```go
import "github.com/alijayanet/gembok-backend/pkg/hotspot"

// Create client
mikrotikClient := getMikroTikClient() // Your MikroTik client
client := hotspot.NewClient(routerID, mikrotikClient)

// Create profile
profile := &hotspot.Profile{
    Name: "prepaid-daily",
    Validity: "1d",
    Price: 5.00,
}
err := client.CreateProfile(profile)

// Generate vouchers
gen := &hotspot.VoucherGenerator{
    Profile: "prepaid-daily",
    Prefix: "CAFE",
    Quantity: 50,
    LengthUsername: 8,
    LengthPassword: 8,
}
vouchers, err := client.GenerateVouchers(gen)

// Get active sessions
sessions, err := client.GetActiveSessions()

// Record sale
sale := &hotspot.Sale{
    Username: "CAFE-ABCD",
    Price: 5.00,
    Date: "jan/20/2024",
    Time: "16:05:11",
}
err := client.RecordSale(sale)
```

## Configuration

### Profile On-Login Script Format
```
:local expmode "rem";:local price "5.00";:local validity "1d";:local selling "5.50";:local lock "Yes";
```

### Sales Record Script Name Format
```
date-|-time-|-username-|-price-|-address-|-mac-|-validity
```

Example: `jan/20/2024-|-16:05:11-|-CAFE-ABCD-|-5.00-|-192.168.1.100-|-AA:BB:CC:DD:EE:FF-|-1d`

## Expiry Modes

- `rem` - Remove expired users automatically
- `ntf` - Notify expired users only
- `remc` - Remove + record to script
- `ntfc` - Notify + record to script

## User Comment Formats

- `vc-PREFIX` - Voucher (username == password)
- `up-PREFIX` - User-Password (username != password)

## License

MIT
