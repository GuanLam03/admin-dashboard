import axios from "axios";

const PUBLIC_ROUTES = ["/login", "/register"];

const api = axios.create({
  baseURL: "http://127.0.0.1:3000",
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
