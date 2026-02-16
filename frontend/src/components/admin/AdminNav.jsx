import { NavLink } from "react-router-dom";

export default function AdminNav() {
  return (
    <nav>
      <NavLink to="/admin">Admin Home</NavLink>{" "}
      <NavLink to="/admin/books">List Books</NavLink>{" "}
      <NavLink to="/admin/book/create">Create Book</NavLink>{" "}
      <NavLink to="/admin/machine/create">Create Machine</NavLink>
    </nav>
  );
}
