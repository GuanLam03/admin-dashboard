import { Outlet, Navigate, NavLink } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";
import { useEffect } from "react";
import api from "../api/axios";


export default function DefaultLayout() {
  const { setUser,user, loading, logout,setLoading } = useAuth();

  

  useEffect(() => {
    api.get("/profile")
      .then(res => setUser(res.data.user))
      .catch(() => setUser(null))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!user) {
    return <Navigate to="/login" />;
  }

  return (
    <div className="flex min-h-screen bg-gray-100">
      {/* Sidebar */}
      <aside className="w-64 bg-white shadow-sm flex-shrink-0">
        <div className="p-6 text-2xl font-bold text-blue-600">
          Admin Dashboard
        </div>
        <nav className="mt-6">
          <ol>
            <li>
              <NavLink
                to="/dashboard"
                end
                className={({ isActive }) =>
                  `block px-4 py-2 hover:bg-blue-100 ${
                    isActive ? "bg-blue-500 text-white" : "text-gray-700"
                  }`
                }
              >
                Home
              </NavLink>
            </li>
            <li>
              <NavLink
                to="/profile"
                className={({ isActive }) =>
                  `block px-4 py-2 hover:bg-blue-100 ${
                    isActive ? "bg-blue-500 text-white" : "text-gray-700"
                  }`
                }
              >
                Profile
              </NavLink>
            </li>
            <li>
              <NavLink
                to="/user-management"
                className={({ isActive }) =>
                  `block px-4 py-2 hover:bg-blue-100 ${
                    isActive ? "bg-blue-500 text-white" : "text-gray-700"
                  }`
                }
              >
                User Management
              </NavLink>
            </li>
            <li>
              <NavLink
                to="/roles"
                className={({ isActive }) =>
                  `block px-4 py-2 hover:bg-blue-100 ${
                    isActive ? "bg-blue-500 text-white" : "text-gray-700"
                  }`
                }
              >
                Roles
              </NavLink>
            </li>
            <li>
              <NavLink
                to="/documents"
                className={({ isActive }) =>
                  `block px-4 py-2 hover:bg-blue-100 ${
                    isActive ? "bg-blue-500 text-white" : "text-gray-700"
                  }`
                }
              >
                Documents
              </NavLink>
            </li>
            <li>
              <NavLink
                to="/google-documents"
                className={({ isActive }) =>
                  `block px-4 py-2 hover:bg-blue-100 ${
                    isActive ? "bg-blue-500 text-white" : "text-gray-700"
                  }`
                }
              >
                Google Documents
              </NavLink>
            </li>
            <li>
              <NavLink
                to="/schedules"
                className={({ isActive }) =>
                  `block px-4 py-2 hover:bg-blue-100 ${
                    isActive ? "bg-blue-500 text-white" : "text-gray-700"
                  }`
                }
              >
                Schedueles
              </NavLink>
            </li>
            <li>
              <NavLink
                to="/settings"
                className={({ isActive }) =>
                  `block px-4 py-2 hover:bg-blue-100 ${
                    isActive ? "bg-blue-500 text-white" : "text-gray-700"
                  }`
                }
              >
                Settings
              </NavLink>
            </li>
          </ol>
        </nav>
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col">
        {/* Header */}
        <header className="bg-white shadow-sm px-6 py-4 flex justify-end items-center">
          
          {/* Bootstrap Dropdown */}
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
                <button
                  className="dropdown-item" 
                  onClick={logout}
                >
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
