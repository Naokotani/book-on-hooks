const QR_KEY = "book_on_hooks:is_qr";
const SESSION_ID_KEY = "book_on_hooks:session_id";

function persistQrFromSearch(search) {
  const params = new URLSearchParams(search);
  const value = params.get("is_qr");

  if (value === "true") {
    sessionStorage.setItem(QR_KEY, "true");
  } else if (value === "false") {
    sessionStorage.setItem(QR_KEY, "false");
  }
}

function getSessionQr() {
  return sessionStorage.getItem(QR_KEY) === "true";
}

function createUuid() {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") {
    return crypto.randomUUID();
  }

  const ts = Date.now().toString(16);
  const rand = Math.random().toString(16).slice(2);
  return `${ts}-${rand}`;
}

function getSessionId() {
  const existing = sessionStorage.getItem(SESSION_ID_KEY);
  if (existing) {
    return existing;
  }

  const id = createUuid();
  sessionStorage.setItem(SESSION_ID_KEY, id);
  return id;
}

export function initMetricsSession(search) {
  const sessionId = getSessionId();
  persistQrFromSearch(search);

  return {
    sessionId,
    isQr: getSessionQr(),
  };
}

export function getMetricsSession() {
  return {
    sessionId: getSessionId(),
    isQr: getSessionQr(),
  };
}
