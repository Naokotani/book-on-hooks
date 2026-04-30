import { useEffect, useState } from "react";
import { NavLink, useLocation } from "react-router-dom";

export default function Nav() {
  const [isOpen, setIsOpen] = useState(false);
  const location = useLocation();

  useEffect(() => {
    setIsOpen(false);
  }, [location.pathname]);

  function closeMenu() {
    setIsOpen(false);
  }

  function toggleMenu() {
    setIsOpen((prev) => !prev);
  }

  return (
    <nav className="main-nav-shell">
      <div className="main-nav-top">
        <NavLink to="/about" className="nav-logo-link" aria-label="Books On Hooks" onClick={closeMenu}>
          <img src="/logo-books-on-hooks-header.svg" alt="Books On Hooks" className="nav-logo" />
        </NavLink>

        <button
          type="button"
          className="nav-toggle"
          aria-expanded={isOpen}
          aria-controls="main-nav-links"
          onClick={toggleMenu}
        >
          <span className="sr-only">Toggle navigation</span>
          <span className="nav-toggle-bar" />
          <span className="nav-toggle-bar" />
          <span className="nav-toggle-bar" />
        </button>
      </div>

      <div id="main-nav-links" className={`main-nav-links ${isOpen ? "is-open" : ""}`}>
        <NavLink to="/location" className="main-nav-link" onClick={closeMenu}>Locations</NavLink>
        <NavLink to="/books" className="main-nav-link" onClick={closeMenu}>Books</NavLink>
        <NavLink to="/how-it-works" className="main-nav-link" onClick={closeMenu}>How It Works</NavLink>
      </div>
    </nav>
  );
}
