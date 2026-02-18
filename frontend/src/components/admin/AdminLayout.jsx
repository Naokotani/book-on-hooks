import { Outlet } from "react-router-dom";
import AdminNav from "./AdminNav";
import Footer from "../Footer";

export default function AdminLayout() {
  return (
    <>
      <header>
        <AdminNav />
      </header>
      <main className="main-layout">
        <Outlet />
      </main>
      <Footer />
    </>
  );
}
