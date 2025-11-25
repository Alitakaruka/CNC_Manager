# CNC Manager Pro - WebSocket API

This project now uses a WebSocket-first transport with HTTP fallback.

## Endpoint
- WS URL: `ws://<host>/ws` (or `wss://` under HTTPS)
- Authentication (optional): include token in the query `?token=...` or send `{type:'auth', data:{token}}` right after connect (server-defined).

## Message Envelope
All application-level messages are JSON objects:
```json
{
  "type": "status | systemInfo | log | command | connect | ack | error",
  "data": {"...": "..."},
  "reqId": "optional request id for command/ack"
}
```

### Server → Client events
- `status`:
```json
{ "type": "status", "data": {"online":0, "printing":0, "offline":0, "total":0} }
```
- `systemInfo`:
```json
{ "type": "systemInfo", "data": {"uptime":"0d 0h 0m", "activeConnections":0} }
```
- `log`:
```json
{ "type": "log", "data": {"id":123, "timestamp": 1710000000000, "level":"info", "message":"...", "type":"system"} }
````
- `ack` (response to a request with `reqId`):
```json
{ "type": "ack", "reqId": "1699999999999-1", "data": {"ok":true} }
```
- `error` (failed request):
```json
{ "type": "error", "reqId": "1699999999999-1", "message": "reason" }
```

### Client → Server requests
- `connect` (establish machine connection):
```json
{ "type":"connect", "reqId":"...", "data": {"TypeOfConnection":"COM|IP|USB", "ConnectionData":"..."} }
```
- `command` (send G‑code to active machine):
```json
{ "type":"command", "reqId":"...", "data": {"gcode":"G28", "uniqueKey":"printer-001"} }
```

Server must reply with `ack` or `error` including the same `reqId`.

## Heartbeat
- Используется нативный ping/pong WebSocket‑протокола (НЕ JSON‑сообщения).
- Сервер: отправляйте WS‑ping примерно раз в 30 секунд.
- Клиент: если нет активности (ни событий, ни ping/pong фреймов) более 60 секунд — закрывает сокет и переподключается.

## Frontend usage
- WebSocket client lives at `src/hooks/WebSocketClient.js` with:
  - Auto‑reconnect (exponential backoff), native ping/pong watchdog, `request(type,data)` resolved on `ack`.
  - `on(type, handler)` to subscribe to events (`status`, `systemInfo`, `log`).
- `SystemHook` listens to WS events and falls back to HTTP snapshot if WS is offline.
- `SendGCode` и `ConnectionHook` работают **только** через WebSocket (HTTP отключён для этих действий).

## HTTP Fallback Endpoints (existing)
- `GET /api/system/snapshot` → `{ connectionStatus, systemInfo }`
- `GET /api/system/status`
- `GET /api/system/info`
- `GET /api/system/logs`
- `POST /api/sendGCode?GCode=...&uniqueKey=...`

## Notes for server implementers
1. Upgrade HTTP to WebSocket at `/ws` and emit events (`status`, `systemInfo`, `log`).
2. On incoming `command`/`connect`/`auth` messages, process and reply with `{type:'ack', reqId, data}` or `{type:'error', reqId, message}`.
3. Enable native WS ping/pong (e.g., ping every ~30s) — client relies on it.
4. Consider rate limiting and message size limits.

## Error handling
- Client retries with exponential backoff and times out requests after 8s.
- If WS disconnected, UI pulls `/api/system/snapshot` at a reduced cadence.

## Credits
- WebSocket migration and API docs were implemented with the help of an AI coding assistant (GPT‑5) in Cursor.
