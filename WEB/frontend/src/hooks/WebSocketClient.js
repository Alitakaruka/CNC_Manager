class WebSocketClient {
  constructor(urlProvider) {
    this.urlProvider = urlProvider; // () => string
    this.ws = null;
    this.connected = false;
    this.listeners = new Map(); // type -> Set(fn)
    this.reqId = 0;
    this.pending = new Map(); // reqId -> {resolve,reject,timeout}
    this.backoffMs = 800; // faster first retry
    this.maxBackoffMs = 15000;
    this.jitterMs = 400; // add small random delay to avoid thundering herd
    this.heartbeatTimer = null; // watchdog disabled (protocol ping/pong not visible in JS)
    this.lastActivity = Date.now();
    this.heartbeatTimeoutMs = 0; // disabled
    this.connectTimeoutMs = 7000; // if no onopen within N sec – force close to trigger reconnect
    this._onlineHookSet = false;
    this.lastUrl = '';
  }

  get readyState() { return this.ws?.readyState }
  get isConnected() { return this.connected }
  get isConnecting() { return this.ws?.readyState === WebSocket.CONNECTING }
  get url() { return this.lastUrl }
  get transport() { return this.connected ? 'WebSocket' : 'HTTP polling' }

  on(type, handler) {
    if (!this.listeners.has(type)) this.listeners.set(type, new Set());
    this.listeners.get(type).add(handler);
    return () => this.listeners.get(type)?.delete(handler);
  }

  emit(type, data) {
    const set = this.listeners.get(type);
    if (set) for (const fn of set) try { fn(data); } catch (_) {}
  }

  connect() {
    if (this.ws && (this.ws.readyState === WebSocket.OPEN || this.ws.readyState === WebSocket.CONNECTING)) return;

    // Setup network listeners once
    if (!this._onlineHookSet && typeof window !== 'undefined') {
      this._onlineHookSet = true;
      window.addEventListener('online', () => this.connect());
      window.addEventListener('offline', () => { try { this.ws?.close(); } catch (_) {} });
    }

    const url = this.urlProvider();
    this.lastUrl = url;
    let connectTimer = null;
    try {
      this.ws = new WebSocket(url);
    } catch (e) {
      this.scheduleReconnect();
      return;
    }

    // If socket stuck in connecting – force close to trigger backoff reconnection
    connectTimer = setTimeout(() => {
      if (this.ws && this.ws.readyState === WebSocket.CONNECTING) {
        try { this.ws.close(); } catch (_) {}
      }
    }, this.connectTimeoutMs);

    this.ws.onopen = () => {
      clearTimeout(connectTimer);
      this.connected = true;
      this.backoffMs = 800; // reset backoff after success
      this.lastActivity = Date.now();
      // Watchdog intentionally disabled: server's native ping/pong is not observable in JS,
      // so we don't close the socket proactively.
      this.emit('open');
    };

    this.ws.onmessage = (evt) => {
      this.lastActivity = Date.now();
      let msg;
      try { msg = JSON.parse(evt.data); } catch { return; }
      if (msg.type === 'ack' && msg.reqId) {
        const entry = this.pending.get(msg.reqId);
        if (entry) {
          clearTimeout(entry.timeout);
          this.pending.delete(msg.reqId);
          entry.resolve(msg.data ?? null);
        }
        return;
      }
      if (msg.type === 'error' && msg.reqId) {
        const entry = this.pending.get(msg.reqId);
        if (entry) {
          clearTimeout(entry.timeout);
          this.pending.delete(msg.reqId);
          entry.reject(new Error(msg.message || 'WS error'));
        }
        return;
      }
      this.emit('message', msg);
      if (msg.type) this.emit(msg.type, msg.data);
    };

    this.ws.onclose = () => {
      clearTimeout(connectTimer);
      this.connected = false;
      clearInterval(this.heartbeatTimer);
      this.emit('close');
      for (const [id, entry] of this.pending) {
        clearTimeout(entry.timeout);
        entry.reject(new Error('WS closed'));
      }
      this.pending.clear();
      this.scheduleReconnect();
    };

    this.ws.onerror = () => {
      try { this.ws?.close(); } catch (_) {}
    };
  }

  scheduleReconnect() {
    const jitter = Math.floor(Math.random() * this.jitterMs);
    const delay = this.backoffMs + jitter;
    setTimeout(() => this.connect(), delay);
    this.backoffMs = Math.min(this.backoffMs * 2, this.maxBackoffMs);
  }

  sendRaw(obj) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(obj));
      return true;
    }
    return false;
  }

  request(type, data, { timeoutMs = 8000 } = {}) {
    const id = `${Date.now()}-${++this.reqId}`;
    const payload = { type, data, reqId: id };
    return new Promise((resolve, reject) => {
      const ok = this.sendRaw(payload);
      const to = setTimeout(() => {
        this.pending.delete(id);
        reject(new Error('WS request timeout'));
      }, timeoutMs);
      if (!ok) {
        clearTimeout(to);
        return reject(new Error('WS not connected'));
      }
      this.pending.set(id, { resolve, reject, timeout: to });
    });
  }
}

export const wsClient = new WebSocketClient(() => {
  const fromEnv = typeof import.meta !== 'undefined' && import.meta.env && import.meta.env.VITE_WS_URL;
  if (fromEnv) return fromEnv;
  const proto = location.protocol === 'https:' ? 'wss' : 'ws';
  const host = location.hostname;
  const port = 8080; // default server port
  return `${proto}://${host}:${port}/ws`;
});

export default wsClient;
