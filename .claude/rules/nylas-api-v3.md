# Nylas API v3 Documentation Rule

## ⛔ CRITICAL: v3 API ONLY

**This SDK supports ONLY Nylas API v3. Do NOT use v1 or v2 API documentation.**

---

## Official v3 Documentation

When looking up Nylas API documentation, ONLY use these sources:

| Resource | URL |
|----------|-----|
| **v3 API Docs** | https://developer.nylas.com/docs/api/v3/ |
| **v3 API Reference** | https://developer.nylas.com/docs/api/v3/ecc/ |
| **v3 Authentication** | https://developer.nylas.com/docs/v3/auth/ |
| **v3 Quickstart** | https://developer.nylas.com/docs/v3/quickstart/ |
| **v3 Email API** | https://developer.nylas.com/docs/v3/email/ |
| **v3 Calendar API** | https://developer.nylas.com/docs/v3/calendar/ |
| **v3 Contacts API** | https://developer.nylas.com/docs/v3/contacts/ |

---

## Base URLs (v3 Only)

```
US Region: https://api.us.nylas.com
EU Region: https://api.eu.nylas.com
```

All endpoints are prefixed with `/v3/`:
- Messages: `/v3/grants/{grant_id}/messages`
- Calendars: `/v3/grants/{grant_id}/calendars`
- Events: `/v3/grants/{grant_id}/events`
- Contacts: `/v3/grants/{grant_id}/contacts`
- Webhooks: `/v3/webhooks`
- Grants: `/v3/grants`

---

## ❌ DO NOT Reference

- **Any `/v1/` or `/v2/` endpoints** - These are deprecated
- **Legacy Nylas documentation** - Pre-v3 docs are outdated
- **Deprecated authentication methods** - Only use v3 OAuth/API key auth
- **Old SDK examples** - Only reference v3-compatible code

---

## v3 Authentication

### API Key Authentication
```
Authorization: Bearer <api_key>
```

### Grant-Based Access
- All user data accessed via grants
- Grant ID required for user-specific endpoints
- Grants obtained through OAuth flow

---

## WebSearch Queries

When searching for Nylas documentation, use these queries:

```
"Nylas v3 API" [endpoint name]
"Nylas API v3" [feature]
site:developer.nylas.com/docs/api/v3 [topic]
site:developer.nylas.com/docs/v3 [topic]
```

**DO NOT search for:**
- "Nylas API" (without v3 - may return old docs)
- "Nylas v2" or "Nylas v1"

---

## Verification Checklist

Before implementing any Nylas API call:

- [ ] Verified endpoint exists in v3 API Reference
- [ ] Checked v3-specific request/response format
- [ ] Used v3 base URL (`api.us.nylas.com` or `api.eu.nylas.com`)
- [ ] Used `/v3/` prefix in all paths
- [ ] Used v3 authentication (Bearer token)
- [ ] Verified grant-based access pattern

---

## Common v3 Patterns

### List Resources
```
GET /v3/grants/{grant_id}/messages
GET /v3/grants/{grant_id}/calendars
GET /v3/grants/{grant_id}/events?calendar_id={id}
```

### Get Single Resource
```
GET /v3/grants/{grant_id}/messages/{message_id}
GET /v3/grants/{grant_id}/events/{event_id}
```

### Create Resource
```
POST /v3/grants/{grant_id}/events
POST /v3/grants/{grant_id}/messages/send
```

### Update Resource
```
PUT /v3/grants/{grant_id}/events/{event_id}
```

### Delete Resource
```
DELETE /v3/grants/{grant_id}/messages/{message_id}
```

---

## Response Format (v3)

All v3 responses follow this pattern:

```json
{
  "request_id": "uuid",
  "data": { ... }  // or [...] for lists
}
```

Paginated responses include:
```json
{
  "request_id": "uuid",
  "data": [...],
  "next_cursor": "cursor_string"
}
```

---

**Remember: Always verify against official v3 docs before implementing.**
