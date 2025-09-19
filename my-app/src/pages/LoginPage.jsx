import { useState } from "react";
import { useAuth } from "../contexts/AuthContext";
import { useNavigate } from "react-router-dom";
import api from "../api/axios";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [twoFARequired, setTwoFARequired] = useState(false);
  const [userId, setUserId] = useState(null);
  const [code, setCode] = useState("");
  const { setUser } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");

    try {
      const res = await api.post("/login", { email, password });

      if (res.data?.twofa_required) {
        //  Ask for 2FA code
        setTwoFARequired(true);
        setUserId(res.data.user_id);
      } else {
        // Normal login → fetch profile
        const profile = await api.get("/profile");
        setUser(profile.data.user);
        navigate("/dashboard");
      }
    } catch (err) {
      setError(err.response?.data?.error || "Login failed");
    }
  };

  const handle2FASubmit = async (e) => {
    e.preventDefault();
    setError("");

    try {
      await api.post("/login/twofa", { user_id: userId, code });
      // If 2FA success, fetch profile
      const profile = await api.get("/profile");
      setUser(profile.data.user);
      navigate("/dashboard");
    } catch (err) {
      setError(err.response?.data?.error || "Invalid 2FA code");
    }
  };

  return (
    <div className="bg-blue-100 w-full flex items-center justify-center">
      <div className="w-full max-w-md bg-white rounded-2xl shadow-lg p-8">
        <h2 className="text-2xl font-bold text-center text-gray-800">
          Welcome Back
        </h2>
        <p className="text-gray-500 text-center mb-6">
          Sign in to continue to your dashboard
        </p>

        {error && (
          <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">
            {error}
          </div>
        )}

        {!twoFARequired ? (
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-gray-600 mb-1">Email</label>
              <input
                type="email"
                className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400"
                placeholder="you@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>

            <div>
              <label className="block text-gray-600 mb-1">Password</label>
              <input
                type="password"
                className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>

            <button
              type="submit"
              className="w-full bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600 transition"
            >
              Sign In
            </button>
          </form>
        ) : (
          <form onSubmit={handle2FASubmit} className="space-y-4">
            <div>
              <label className="block text-gray-600 mb-1">2FA Code</label>
              <input
                type="text"
                className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-green-400"
                placeholder="Enter 6-digit code"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                required
              />
            </div>

            <button
              type="submit"
              className="w-full bg-green-500 text-white py-2 rounded-lg hover:bg-green-600 transition"
            >
              Verify
            </button>
          </form>
        )}

        {!twoFARequired && (
          <p className="text-center text-gray-500 text-sm mt-6">
            Don’t have an account?{" "}
            <a href="/register" className="text-blue-500 hover:underline">
              Register
            </a>
          </p>
        )}
      </div>
    </div>
  );
}
