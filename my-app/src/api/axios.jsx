import axios from "axios";

const PUBLIC_ROUTES = ["/login", "/register"];
const API_URL = import.meta.env.VITE_API_URL;

const api = axios.create({
  baseURL: API_URL,
  // baseURL: "https://admin-dashboard.test/",
 
  withCredentials: true, // allow cookies (jwt_token) to be sent and received
  headers: {
    "Content-Type": "application/json",
  },
});

// api.interceptors.request.use((config) => {
//   const locale = localStorage.getItem("locale-admin") || "en";
//   config.headers["Accept-Language"] = locale;
//   return config;
// });

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const currentPath = window.location.pathname;
    if (
      error.response &&
      error.response.status === 401 &&
      !PUBLIC_ROUTES.includes(currentPath)
    ) {
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);


export default api;
