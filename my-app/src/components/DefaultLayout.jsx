import { Outlet, Navigate, NavLink } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";
import { useEffect, useState } from "react";
import api from "../api/axios";

export default function DefaultLayout() {
  const { setUser, user, loading, logout, setLoading } = useAuth();

  const [emailsOpen, setEmailsOpen] = useState(() => {
    // Load from localStorage on first render
    return localStorage.getItem("emailsOpen") === "true";
  });

  // Sync state to localStorage when it changes
  useEffect(() => {
    localStorage.setItem("emailsOpen", emailsOpen.toString());
  }, [emailsOpen]);


  useEffect(() => {
    api.get("/profile")
      .then(res => setUser(res.data.user))
      .catch(() => setUser(null))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div>Loading...</div>;
  if (!user) return <Navigate to="/login" />;

  // 🔹 small helper
  const hasPermission = (perm) => user?.permissions?.includes(perm);

  return (
    <div className="flex min-h-screen bg-gray-100">
      {/* Sidebar */}
      <aside className="w-64 bg-gray-50 shadow-sm flex-shrink-0">
        <div className="p-6 text-2xl font-bold text-blue-600">
          Admin Dashboard
        </div>
        <nav className="mt-6">
          <ol>
            <li>
              <NavLink to="/dashboard" end className={({ isActive }) =>
                `block px-4 py-2 hover:bg-blue-100 ${isActive ? "bg-blue-500 text-white" : "text-gray-700"}`}>
                Home
              </NavLink>
            </li>

            <li>
              <NavLink to="/profile" className={({ isActive }) =>
                `block px-4 py-2 hover:bg-blue-100 ${isActive ? "bg-blue-500 text-white" : "text-gray-700"}`}>
                Profile
              </NavLink>
            </li>

            <li>
              <NavLink to="/user-management" className={({ isActive }) =>
                `block px-4 py-2 hover:bg-blue-100 ${isActive ? "bg-blue-500 text-white" : "text-gray-700"}`}>
                User Management
              </NavLink>
            </li>

            <li>
              <NavLink to="/roles" className={({ isActive }) =>
                `block px-4 py-2 hover:bg-blue-100 ${isActive ? "bg-blue-500 text-white" : "text-gray-700"}`}>
                Roles
              </NavLink>
            </li>

            <li>
              <NavLink to="/documents" className={({ isActive }) =>
                `block px-4 py-2 hover:bg-blue-100 ${isActive ? "bg-blue-500 text-white" : "text-gray-700"}`}>
                Documents
              </NavLink>
            </li>

            <li>
              <NavLink to="/google-documents" className={({ isActive }) =>
                `block px-4 py-2 hover:bg-blue-100 ${isActive ? "bg-blue-500 text-white" : "text-gray-700"}`}>
                Google Documents
              </NavLink>
            </li>

            <li>
              <NavLink to="/schedules" className={({ isActive }) =>
                `block px-4 py-2 hover:bg-blue-100 ${isActive ? "bg-blue-500 text-white" : "text-gray-700"}`}>
                Schedules
              </NavLink>
            </li>

            {/* 🔹 Emails Section */}
            {(hasPermission("gmail.technical.read") || hasPermission("gmail.support.read")) && (
              <li>
                <div>
                  <button
                    className="w-full flex justify-between items-center text-blue-600 hover:bg-blue-100 rounded px-4 py-2"
                    onClick={() => setEmailsOpen(!emailsOpen)}
                  >
                    <span>Emails</span>
                    <svg
                      className={`w-4 h-4 transition-transform duration-200 ${emailsOpen ? "rotate-90" : ""}`}
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                      viewBox="0 0 24 24"
                    >
                      <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
                    </svg>
                  </button>
                  {emailsOpen && (
                    <ol className="ml-4 mt-2">
                      {hasPermission("gmail.technical.read") && (
                        <li>
                          <NavLink to="/emails/technical" className={({ isActive }) =>
                            `block px-4 py-2 hover:bg-blue-100 ${isActive ? "bg-blue-500 text-white" : "text-gray-700"}`}>
                            Technical Emails
                          </NavLink>
                        </li>
                      )}
                      {hasPermission("gmail.support.read") && (
                        <li>
                          <NavLink to="/emails/support" className={({ isActive }) =>
                            `block px-4 py-2 hover:bg-blue-100 ${isActive ? "bg-blue-500 text-white" : "text-gray-700"}`}>
                            Support Emails
                          </NavLink>
                        </li>
                      )}
                      <li>
                          <NavLink to="/emails/settings" className={({ isActive }) =>
                            `block px-4 py-2 hover:bg-blue-100 ${isActive ? "bg-blue-500 text-white" : "text-gray-700"}`}>
                            Emails Settings
                          </NavLink>
                        </li>
                    </ol>
                  )}
                </div>
              </li>
            )}

            <li>
              <NavLink to="/settings" className={({ isActive }) =>
                `block px-4 py-2 hover:bg-blue-100 ${isActive ? "bg-blue-500 text-white" : "text-gray-700"}`}>
                Settings
              </NavLink>
            </li>
          </ol>
        </nav>
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col">
        {/* Header */}
        <header className="bg-gray-50 shadow-sm px-6 py-4 flex justify-end items-center">
          <div className="dropdown">
            <button className="dropdown-toggle" type="button" data-bs-toggle="dropdown" aria-expanded="false">
              {user.name}
            </button>
            <ul className="dropdown-menu dropdown-menu-end">
              <li>
                <button className="dropdown-item" onClick={logout}>Logout</button>
              </li>
            </ul>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
