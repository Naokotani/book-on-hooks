import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import "./normalize.css";
import "./styles.css";
import Layout from "./components/layout/Layout";
import About from "./pages/about/AboutPage";
import Books from "./pages/books/BooksPage";
import BookLocation from "./pages/books/BookLocationPage";
import HowItWorks from "./pages/how-it-works/HowItWorksPage";
import Location from "./pages/location/LocationPage";
import LocationMachineView from "./pages/location/LocationMachinePage";
import AdminLayout from "./pages/admin/AdminLayout";
import AdminHome from "./pages/admin/AdminHomePage";
import AdminBooks from "./pages/admin/books/AdminBooksPage";
import AdminCreateBook from "./pages/admin/books/AdminCreateBookPage";
import AdminUpdateBook from "./pages/admin/books/AdminUpdateBookPage";
import AdminCreateMachine from "./pages/admin/machines/AdminCreateMachinePage";
import AdminMachineLoad from "./pages/admin/machines/AdminMachineLoadPage";

document.body.classList.add("dark");

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/about" replace />} />
        <Route element={<Layout />}>
          <Route path="/about" element={<About />} />
          <Route path="/books" element={<Books />} />
          <Route path="/books/:id/summary" element={<BookLocation />} />
          <Route path="/how-it-works" element={<HowItWorks />} />
          <Route path="/location" element={<Location />} />
          <Route path="/location/:id" element={<LocationMachineView />} />
        </Route>

        <Route path="/admin" element={<AdminLayout />}>
          <Route index element={<AdminHome />} />
          <Route path="books" element={<AdminBooks />} />
          <Route path="book/create" element={<AdminCreateBook />} />
          <Route path="book/:id/update" element={<AdminUpdateBook />} />
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
