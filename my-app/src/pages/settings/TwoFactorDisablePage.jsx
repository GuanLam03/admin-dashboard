import { useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../../api/axios";
import { useAuth } from "../../contexts/AuthContext";

export default function TwoFactorDisablePage() {
  const [otp, setOtp] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const { user, setUser } = useAuth();

  const handleDisable = async (e) => {
    e.preventDefault();
    setError("");
    try {
      await api.post("/twofa/disable", { code: otp });
      setUser({ ...user, two_factor_enabled: false });
      navigate("/settings");
    } catch (err) {
      setError(err.response?.data?.error || "Failed to disable 2FA");
    }
  };

  return (
    <div className="max-w-md mx-auto mt-10 bg-white shadow rounded-lg p-6">
      <h2 className="text-xl font-bold mb-4">Disable Two-Factor Authentication</h2>
      <p className="text-gray-600 mb-4">
        Enter the 6-digit code from your authenticator app to confirm.
      </p>

      {error && (
        <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">{error}</div>
      )}

      <form onSubmit={handleDisable} className="space-y-4">
        <input
          type="text"
          value={otp}
          onChange={(e) => setOtp(e.target.value)}
          placeholder="123456"
          className="w-full px-3 py-2 border rounded"
          required
        />
        <div className="mt-4 flex flex-col gap-2">
            <button
                type="submit"
                className="w-full bg-red-500 text-white py-2 rounded hover:bg-red-600"
                >
                Disable 2FA
                </button>
                
                <button
                type="button"
                onClick={() => navigate("/settings")}
                className="w-full py-2 border rounded"
                >
                Cancel
            </button>
        </div>
        
      </form>
    </div>
  );
}
