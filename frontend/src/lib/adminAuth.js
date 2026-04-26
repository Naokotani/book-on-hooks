export async function fetchAdminSession() {
  return fetch("/api/admin/me", {
    credentials: "include",
  });
}

export async function adminFetch(input, init = {}) {
  const response = await fetch(input, {
    ...init,
    credentials: "include",
  });

  if (
    response.status === 401 &&
    typeof window !== "undefined" &&
    window.location.pathname !== "/admin/login"
  ) {
    window.location.assign("/admin/login");
  }

  return response;
}
