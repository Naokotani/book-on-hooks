import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import "./normalize.css";
import "./styles.css";
import Layout from "./components/Layout";
import About from "./components/About";
import Books from "./components/Books";
import BookLocation from "./components/BookLocation";
import HowItWorks from "./components/HowItWorks";
import Location from "./components/Location";
import LocationMachineView from "./components/LocationMachineView";
import AdminLayout from "./components/admin/AdminLayout";
import AdminHome from "./components/admin/AdminHome";
import AdminBooks from "./components/admin/AdminBooks";
import AdminCreateBook from "./components/admin/AdminCreateBook";
import AdminCreateMachine from "./components/admin/AdminCreateMachine";
import AdminMachineLoad from "./components/admin/AdminMachineLoad";

document.body.classList.add("dark");

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/about" replace />} />
        <Route element={<Layout />}>
          <Route path="/about" element={<About />} />
          <Route path="/books" element={<Books />} />
          <Route path="/books/:id/locations" element={<BookLocation />} />
          <Route path="/how-it-works" element={<HowItWorks />} />
          <Route path="/location" element={<Location />} />
          <Route path="/location/:id" element={<LocationMachineView />} />
        </Route>

        <Route path="/admin" element={<AdminLayout />}>
          <Route index element={<AdminHome />} />
          <Route path="books" element={<AdminBooks />} />
          <Route path="book/create" element={<AdminCreateBook />} />
          <Route path="machine/create" element={<AdminCreateMachine />} />
          <Route path="machine/load/:id" element={<AdminMachineLoad />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
