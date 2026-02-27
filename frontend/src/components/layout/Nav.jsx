import { NavLink } from "react-router-dom";

export default function Nav() {
  return (
    <nav className="main-nav">
      <NavLink to="/about" className="nav-logo-link" aria-label="Books On Hooks">
        <img src="/logo-books-on-hooks-header.svg" alt="Books On Hooks" className="nav-logo" />
      </NavLink>
      <NavLink to="/books" className="main-nav-link">Books</NavLink>
      <NavLink to="/how-it-works" className="main-nav-link">How It Works</NavLink>
      <NavLink to="/location" className="main-nav-link">Locations</NavLink>
    </nav>
  );
}
