import { NavLink, useNavigate } from "react-router-dom";
import { adminFetch } from "../../lib/adminAuth";

export default function AdminNav() {
  const navigate = useNavigate();

  async function handleLogout() {
    try {
      await adminFetch("/api/admin/logout", {
        method: "POST",
      });
    } finally {
      navigate("/admin/login", { replace: true });
    }
  }

  return (
    <nav className="admin-nav">
      <NavLink to="/">Main Site</NavLink>
      <NavLink to="/admin">Machines</NavLink>
      <NavLink to="/admin/metrics">Metrics</NavLink>
      <NavLink to="/admin/books">Books</NavLink>
      <NavLink to="/admin/book/create">Create Book</NavLink>
      <NavLink to="/admin/machine/create">Create Machine</NavLink>
      <button className="admin-nav-button" type="button" onClick={handleLogout}>
        Sign Out
      </button>
    </nav>
  );
}
