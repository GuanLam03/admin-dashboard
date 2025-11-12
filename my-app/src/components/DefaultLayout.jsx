import { Outlet, Navigate, NavLink } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";
import { useEffect, useState } from "react";
import api from "../api/axios";

export default function DefaultLayout() {
  const { setUser, user, loading, logout, setLoading } = useAuth();
  const [sidebar, toggleSidebar] = useSidebarState();

  const [locale, setLocale] = useState(localStorage.getItem("locale-admin") || "en");

  useEffect(() => {
      localStorage.setItem("locale-admin", locale);
  }, [locale]);

  const changeLanguage = (code) => {
    setLocale(code);
  };



  useEffect(() => {
    api.get("/profile")
      .then(res => setUser(res.data.user))
      .catch(() => setUser(null))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div>Loading...</div>;
  if (!user) return <Navigate to="/login" />;

  // small helper
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

            {/*  Emails Section */}
            {(hasPermission("gmail.technical.read") || hasPermission("gmail.support.read")) && (
              <li>
                <div>
                  <button
                    className="w-full flex justify-between items-center text-blue-600 hover:bg-blue-100 rounded px-4 py-2"
                    onClick={() => toggleSidebar("emails")}
                  >
                    <span>Emails</span>
                    <svg
                      className={`w-4 h-4 transition-transform duration-200 ${sidebar.emails ? "rotate-90" : ""}`}
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                      viewBox="0 0 24 24"
                    >
                      <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
                    </svg>
                  </button>
                  {sidebar.emails && (
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

         
            {/*  Ads Tracking Section */}
            <li>
              <div>
                <button
                  className="w-full flex justify-between items-center text-blue-600 hover:bg-blue-100 rounded px-4 py-2"
                  onClick={() => toggleSidebar("ads")}
                >
                  <span>Ads Tracking</span>
                  <svg
                    className={`w-4 h-4 transition-transform duration-200 ${sidebar.ads ? "rotate-90" : ""}`}
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    viewBox="0 0 24 24"
                  >
                    <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
                  </svg>
                </button>
                {sidebar.ads && (
                  <ol className="ml-4 mt-2">
  
                      <li>
                        <NavLink to="/ads-tracking/campaign" className={({ isActive }) =>
                          `block px-4 py-2 hover:bg-blue-100 ${isActive ? "bg-blue-500 text-white" : "text-gray-700"}`}>
                          Ads Campaign
                        </NavLink>
                      </li>
                    

                      <li>
                        <NavLink to="/ads-tracking/log" className={({ isActive }) =>
                          `block px-4 py-2 hover:bg-blue-100 ${isActive ? "bg-blue-500 text-white" : "text-gray-700"}`}>
                          Ads Log
                        </NavLink>
                      </li>

                       <li>
                        <NavLink to="/ads-tracking/report" className={({ isActive }) =>
                          `block px-4 py-2 hover:bg-blue-100 ${isActive ? "bg-blue-500 text-white" : "text-gray-700"}`}>
                          Ads Report
                        </NavLink>
                      </li>
                    
                  
                  </ol>
                )}
              </div>
            </li>
            


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
        {/* <header className="bg-gray-50 shadow-sm px-6 py-4 flex justify-end items-center">
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
        </header> */}
        <header className="bg-gray-50 shadow-sm px-6 py-4 flex justify-end items-center gap-4">
          {/* Language Switcher */}
          <div className="dropdown">
            <button
              className="dropdown-toggle"
              type="button"
              data-bs-toggle="dropdown"
              aria-expanded="false"
            >
              {locale}
            </button>
            <ul className="dropdown-menu dropdown-menu-end">
              <li>
                <button className="dropdown-item" onClick={() => changeLanguage("en")}>
                  English
                </button>
              </li>
              <li>
                <button className="dropdown-item" onClick={() => changeLanguage("cn")}>
                  中文
                </button>
              </li>
            </ul>
          </div>

          {/* User Dropdown */}
          <div className="dropdown">
            <button
              className="dropdown-toggle"
              type="button"
              data-bs-toggle="dropdown"
              aria-expanded="false"
            >
              {user.name}
            </button>
            <ul className="dropdown-menu dropdown-menu-end">
              <li>
                <button className="dropdown-item" onClick={logout}>
                  Logout
                </button>
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



function useSidebarState() {
  const [state, setState] = useState(() => {
    const stored = localStorage.getItem("sidebarStates");
    return stored ? JSON.parse(stored) : {};
  });

  useEffect(() => {
    localStorage.setItem("sidebarStates", JSON.stringify(state));
  }, [state]);

  const toggle = (key) => {
    setState(prev => ({ ...prev, [key]: !prev[key] }));
  };

  return [state, toggle];
}
