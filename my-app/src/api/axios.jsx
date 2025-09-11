import axios from "axios";

const api = axios.create({
  baseURL: "http://127.0.0.1:3000", // your Go backend
  withCredentials: true, // allow cookies (jwt_token) to be sent and received
  headers: {
    "Content-Type": "application/json",
  },
});

// Optional: response interceptor for handling 401 logout automatically
// api.interceptors.response.use(
//   (response) => response,
//   (error) => {
//     if (error.response && error.response.status === 401) {
//       // Example: redirect to login page or clear user
//       console.warn("Unauthorized, redirecting to login...");
//       window.location.href = "/login";
//     }
//     return Promise.reject(error);
//   }
// );

export default api;
