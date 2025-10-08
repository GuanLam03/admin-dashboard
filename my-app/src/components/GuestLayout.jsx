import { Outlet, Navigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";
import { useEffect } from "react";
import api from "../api/axios";

export default function GuestLayout() {
  const { user} = useAuth();
 

  if (user) {
    return <Navigate to="/dashboard" />;
  }

  return (
    <div className="w-[100vw] h-[100vh] flex justify-center align-center">
      <Outlet />
    </div>
  );
}
