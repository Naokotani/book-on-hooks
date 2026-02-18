import { NavLink } from "react-router-dom";

export default function Nav() {
  return (
    <nav>
      <NavLink to="/about">About</NavLink>{" "}
      <NavLink to="/books">Books</NavLink>{" "}
      <NavLink to="/how-it-works">How It Works</NavLink>{" "}
      <NavLink to="/location">Location</NavLink>
    </nav>
  );
}
