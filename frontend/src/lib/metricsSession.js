const QR_KEY = "book_on_hooks:is_qr";

export function persistQrFromSearch(search) {
  const params = new URLSearchParams(search);
  const value = params.get("is_qr");

  if (value === "true") {
    sessionStorage.setItem(QR_KEY, "true");
  } else if (value === "false") {
    sessionStorage.setItem(QR_KEY, "false");
  }
}

export function getSessionQr() {
  return sessionStorage.getItem(QR_KEY) === "true";
}
