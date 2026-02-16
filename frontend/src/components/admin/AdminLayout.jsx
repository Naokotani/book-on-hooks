import { Outlet } from "react-router-dom";
import AdminNav from "./AdminNav";

export default function AdminLayout() {
  return (
    <>
      <header>
        <AdminNav />
      </header>
      <main>
        <Outlet />
      </main>
    </>
  );
}
