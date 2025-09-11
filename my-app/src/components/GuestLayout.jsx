import { Outlet, Navigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";
import { useEffect } from "react";
import api from "../api/axios";

export default function GuestLayout() {
  const { user, loading, setUser} = useAuth();

  // if (loading) {
  //   return <div>Loading...</div>;
  // }

  useEffect(() => {
    // Try to fetch profile on first load
    api.get("/profile")
      .then(res => setUser(res.data.user))
      .catch(() => setUser(null))
  }, []);

  if (user) {
    return <Navigate to="/dashboard" />;
  }

  return (
    <div className="w-[100vw] h-[100vh] flex justify-center align-center">
      <Outlet />
    </div>
  );
}
