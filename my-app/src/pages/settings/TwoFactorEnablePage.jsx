import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import api from "../../api/axios";
import { useAuth } from "../../contexts/AuthContext";
import { useTranslation } from "react-i18next";

export default function TwoFactorSettingsPage() {
  const {t} = useTranslation();
  const [qrCode, setQrCode] = useState("");
  const [otp, setOtp] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const { user, setUser } = useAuth();

  useEffect(() => {
    // fetch QR code from backend
    api
      .get("/twofa/qrcode", { responseType: "blob" }) // ← important!
      .then((res) => {
        const imageUrl = URL.createObjectURL(res.data); // create URL from blob
        setQrCode(imageUrl); // now you can use it in <img src={qrCode} />
      })
      .catch(async (err) => {
        if (
          err.response &&
          err.response.data instanceof Blob &&
          err.response.data.type === "application/json"
        ) {
          const text = await err.response.data.text();
          const json = JSON.parse(text);
          setError(json.error || "Failed to load QR code");
        } else {
          setError("Failed to load QR code");
        }
      });


  }, []);

  const handleVerify = async (e) => {
    e.preventDefault();
    try {
      await api.post("/twofa/enable", { code: otp });
      setUser({ ...user, two_factor_enabled: true });
      navigate("/settings"); //  return to settings page
    } catch (err) {
      setError(err.response?.data?.error ? t(err.response.data.error) : "Verification failed");
    }
  };

  return (
    <div className="max-w-md mx-auto bg-white shadow rounded-lg p-6">
      <h2 className="text-xl font-bold mb-4">Two-Factor Authentication Setup</h2>
      <p className="text-gray-600 mb-4">
        Scan the QR code with Google Authenticator and enter the 6-digit code to
        enable 2FA.
      </p>

      {qrCode && <img src={qrCode} alt="2FA QR Code" className="mx-auto mb-4" />}

      {error && <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">{error}</div>}

      <form onSubmit={handleVerify} className="flex flex-col gap-4">
        <input
          type="text"
          placeholder="Enter 6-digit code"
          value={otp}
          onChange={(e) => setOtp(e.target.value)}
          className="w-full px-4 py-2 border rounded-lg"
          required
        />
        <button type="submit" className="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-600">
          Verify & Enable
        </button>
      </form>
    </div>
  );
}
