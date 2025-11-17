import { useState, useEffect } from "react";
import { FcGoogle } from "react-icons/fc";
import { RiLogoutBoxLine } from "react-icons/ri";
import { useAuth } from "../../contexts/AuthContext";
import api from "../../api/axios";
import { useTranslation } from "react-i18next";

export default function ScheduleSettingsPage() {
  const {t} = useTranslation();
  const { user } = useAuth();
  const [account, setAccount] = useState(null);
  const [loading, setLoading] = useState(true);

  // Load connected account info from backend
  useEffect(() => {
    async function fetchAccount() {
      try {
        const res = await api.get("/google/schedule/account"); // backend route
        setAccount(res.data.account || null);
      } catch (err) {
        console.error("Failed to load Google account:", err);
      } finally {
        setLoading(false);
      }
    }
    fetchAccount();
  }, []);

  const handleLogin = async () => {
    try {
      const res = await api.get("/google/auth/schedule/url");
      if (res.data.url) {
        window.location.href = res.data.url; // redirect to Google OAuth
      }
    } catch (err) {
      console.error("Failed to start Google login:", err);
    }
  };

  const handleLogout = async (id) => {
    if (!window.confirm("Disconnect Google account?")) return;
    try {
      await api.post(`/gmail/remove-accounts/${id}`);
      setAccount(null);
    } catch (err) {
      console.error("Logout failed:", err);
    }
  };

  if (loading) {
    return (
      <div className="p-4">
        <p>Loading...</p>
      </div>
    );
  }

  return (
    <div className="space-y-10">
      {/* Accounts Section */}
      <div className="bg-white p-4 rounded-lg shadow-sm">
        <h2 className="text-xl font-bold mb-4">{t("scheduleManagement.settingsPage.gmailAccountForGoogleCalendar")}</h2>
        <table className="w-full border-collapse">
          <thead>
            <tr className="border-b">
              <th className="text-left p-2">{t("gmail")}</th>
              <th className="text-left p-2">{t("common.labels.action")}</th>
            </tr>
          </thead>
          <tbody>
            <tr className="border-b">
              <td className="p-2">{account ? account.email : t("scheduleManagement.settingsPage.notConnected")}</td>
              <td className="p-2">
                {account ? (
                  <button
                    onClick={() => handleLogout(account.id)}
                    className="flex gap-2 items-center px-4 py-2 bg-gray-100 text-black font-semibold !rounded-full hover:bg-gray-200 transition"
                  >
                    <RiLogoutBoxLine size={22} className="text-blue-500" />
                    {t("logout")}
                  </button>
                ) : (
                  <button
                    onClick={handleLogin}
                    className="flex gap-2 items-center px-4 py-2 bg-gray-100 text-black font-semibold !rounded-full hover:bg-gray-200 transition"
                  >
                    <FcGoogle size={22} />
                    <span>{t("scheduleManagement.settingsPage.connect")}</span>
                  </button>
                )}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
}
