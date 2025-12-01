import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider } from "react-router-dom";
import router from "./router";
import { AuthProvider } from "./contexts/AuthContext";
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.bundle.min.js";
import "./utils/i18n/i18n";
import { CentrifugeProvider } from "./contexts/CentrifugeContext";
import { ToastContainer } from 'react-toastify';
// Wrapper to provide OnlineUsersProvider after AuthProvider
function AppProviders() {
 
  return (
    <CentrifugeProvider>
      <RouterProvider router={router} />
      <ToastContainer />
    </CentrifugeProvider>
  );
}

ReactDOM.createRoot(document.getElementById("root")).render(
  <AuthProvider>
    <AppProviders />
  </AuthProvider>
);
