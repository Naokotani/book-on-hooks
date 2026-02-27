import { Outlet } from "react-router-dom";
import AdminNav from "./AdminNav";
import Footer from "../../components/layout/Footer";

export default function AdminLayout() {
  return (
    <div className="layout">
      <header>
        <div className="page-container">
          <AdminNav />
        </div>
      </header>
      <main className="main-layout">
        <div className="page-container">
          <Outlet />
        </div>
      </main>
      <Footer />
    </div>
  );
}
