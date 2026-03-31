# e-conomic MCP Connector

An MCP (Model Context Protocol) server that connects Claude to [e-conomic](https://www.e-conomic.com/) via the REST API. Manage customers, invoices, products, suppliers, orders, journals, and accounting data through natural conversation.

## Prerequisites

- An **e-conomic developer account** with an app registered in the [Developer Network](https://secure.e-conomic.com/developer)
- An **Agreement Grant Token** from an e-conomic customer who has authorized your app

### Getting Your Tokens

1. Register as a developer at [secure.e-conomic.com/developer](https://secure.e-conomic.com/developer)
2. Create a new app to receive your **App Secret Token**
3. Have the e-conomic customer grant access to their agreement — this produces the **Agreement Grant Token**
4. Both tokens are sent as headers on every API call (`X-AppSecretToken` and `X-AgreementGrantToken`)

> **Demo mode:** Set both tokens to `demo` to explore the API with sample data (read-only).

## Setup

### Install as Claude Desktop Extension

Download the latest `.mcpb` from [Releases](../../releases) and install it in Claude Desktop. It will prompt for both tokens through the settings UI.

### Install from Source

```bash
git clone <repo-url>
cd e-conomic-connector
cp .env.example .env
# Edit .env with your tokens
make build
```

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `ECONOMIC_APP_SECRET_TOKEN` | Yes | App Secret Token from the Developer Network |
| `ECONOMIC_AGREEMENT_GRANT_TOKEN` | Yes | Agreement Grant Token granting access to a specific e-conomic agreement |

## Available Tools

### Auth

| Tool | Description |
|------|-------------|
| `economic-auth-status` | Check authentication status and get company info |

### Customers

| Tool | Description |
|------|-------------|
| `economic-list-customers` | List customers with filtering and pagination |
| `economic-get-customer` | Get a specific customer by number |
| `economic-create-customer` | Create a new customer |
| `economic-update-customer` | Update an existing customer (full PUT) |
| `economic-delete-customer` | Delete a customer |

### Invoices

| Tool | Description |
|------|-------------|
| `economic-list-invoices-drafts` | List draft invoices |
| `economic-get-invoice-draft` | Get a specific draft invoice |
| `economic-create-invoice-draft` | Create a new draft invoice |
| `economic-book-invoice` | Book (finalize) a draft invoice |
| `economic-list-invoices-booked` | List booked invoices |
| `economic-get-invoice-booked` | Get a specific booked invoice |

### Products

| Tool | Description |
|------|-------------|
| `economic-list-products` | List products |
| `economic-get-product` | Get a specific product |
| `economic-create-product` | Create a new product |
| `economic-update-product` | Update a product (full PUT) |
| `economic-delete-product` | Delete a product |

### Suppliers

| Tool | Description |
|------|-------------|
| `economic-list-suppliers` | List suppliers |
| `economic-get-supplier` | Get a specific supplier |
| `economic-create-supplier` | Create a new supplier |
| `economic-update-supplier` | Update a supplier (full PUT) |
| `economic-delete-supplier` | Delete a supplier |

### Orders

| Tool | Description |
|------|-------------|
| `economic-list-orders-drafts` | List draft orders |
| `economic-get-order-draft` | Get a specific draft order |
| `economic-create-order-draft` | Create a new draft order |

### Journals & Entries

| Tool | Description |
|------|-------------|
| `economic-list-journals` | List journals (daybooks) |
| `economic-get-journal` | Get a specific journal |
| `economic-create-journal-entry` | Create entries in a journal |
| `economic-list-entries` | List booked account entries |

### Accounting Years & Periods

| Tool | Description |
|------|-------------|
| `economic-list-accounting-years` | List accounting years |
| `economic-get-accounting-year` | Get a specific accounting year |
| `economic-list-accounting-year-periods` | List periods (months) in an accounting year |
| `economic-list-period-entries` | List entries for a specific period (e.g. March 2026) |
| `economic-list-accounting-year-entries` | List all entries for an entire accounting year |
| `economic-list-accounting-year-totals` | Account totals/balances for an accounting year |

### Reference Data

| Tool | Description |
|------|-------------|
| `economic-list-accounts` | Chart of accounts |
| `economic-get-account` | Get a specific account |
| `economic-list-payment-terms` | Payment terms |
| `economic-list-vat-zones` | VAT zones |
| `economic-list-customer-groups` | Customer groups |
| `economic-list-supplier-groups` | Supplier groups |
| `economic-list-product-groups` | Product groups |
| `economic-list-units` | Units of measure |
| `economic-list-departments` | Departments |
| `economic-list-currencies` | Currencies |
| `economic-list-layouts` | Document layouts |
| `economic-self` | Current agreement/company info |

### Filtering & Pagination

All list endpoints support:

- **`filter`** — e.g. `name$like:*john*`, `customerNumber$gte:100`. Operators: `$eq`, `$ne`, `$gt`, `$gte`, `$lt`, `$lte`, `$like`, `$in`, `$nin`. Combine with `$and`, `$or`.
- **`sort`** — e.g. `name`, `-name` (descending), `-name,age` (multiple).
- **`pagesize`** — Items per page (default 20, max 1000).
- **`skippages`** — Pages to skip (0-based).

## Usage Examples

Once connected, you can ask Claude things like:

- *"List all customers whose name contains 'ApS'"*
- *"Create a draft invoice for customer 123 with two line items"*
- *"Book draft invoice 456"*
- *"Show me all products in product group 2"*
- *"What's our chart of accounts?"*
- *"Create a journal entry debiting account 1000 and crediting account 2000 for 5000 DKK"*
- *"Get all entries for March 2026 (accounting year 2026, period 3)"*
- *"Show me the account totals for 2025"*

## Running Locally

**Stdio mode** (default, used by Claude Desktop):
```bash
./go/server/e-conomic-connector
```

**Verify connection:**
```bash
./go/server/e-conomic-connector --verify
```

**Verbose logging:**
```bash
./go/server/e-conomic-connector --verbose
```

Logs are written to `logs/economic-mcp.log`.

## Development

```bash
make build              # Cross-compile to linux/darwin/windows
make package            # Build + zip into .mcpb
make test               # Run tests
make version            # Print current version
make bump-version V=x.y.z  # Set version in manifest.json
make clean              # Remove build artifacts
```

### Releasing

Push a version tag to trigger the GitHub Actions release workflow:

```bash
git tag v1.0.0
git push origin v1.0.0
```

This cross-compiles, packages the `.mcpb`, and creates a GitHub Release with the binary attached.

## API Reference

This connector uses the [e-conomic REST API](https://restdocs.e-conomic.com/) at `https://restapi.e-conomic.com`. See the [full API documentation](https://restdocs.e-conomic.com/) for details on request/response schemas.

## Security

- Tokens are read from environment variables, never hardcoded
- **Never commit `.env`** — it is in `.gitignore`
- The App Secret Token identifies your application; the Agreement Grant Token grants access to a specific company's data
- Both tokens are required on every request

## License

MIT
