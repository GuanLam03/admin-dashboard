import { useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../../api/axios";
import { useAuth } from "../../contexts/AuthContext";
import { useTranslation } from "react-i18next";

export default function TwoFactorDisablePage() {
  const {t} = useTranslation();
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
      setError(err.response?.data?.error ? t(err.response.data.error) : t("settings.twoFactor.failedToDisable"));
    }
  };

  return (
    <div className="max-w-md mx-auto mt-10 bg-white shadow rounded-lg p-6">
      <h2 className="text-xl font-bold mb-4">{t("settings.twoFactor.disableTitle")}</h2>
      <p className="text-gray-600 mb-4">
        {t("settings.twoFactor.disableDescription")}
      </p>

      {error && (
        <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">{error}</div>
      )}

      <form onSubmit={handleDisable} className="space-y-4">
        <input
          type="text"
          value={otp}
          onChange={(e) => setOtp(e.target.value)}
          placeholder={t("settings.twoFactor.codePlaceholder")}
          className="w-full px-3 py-2 border rounded"
          required
        />
        <div className="mt-4 flex flex-col gap-2">
            <button
                type="submit"
                className="w-full bg-red-500 text-white py-2 rounded hover:bg-red-600"
                >
                {t("settings.twoFactor.disableButton")}
                </button>
                
                <button
                type="button"
                onClick={() => navigate("/settings")}
                className="w-full py-2 border rounded"
                >
                {t("common.buttons.cancel")}
            </button>
        </div>
        
      </form>
    </div>
  );
}
