import { useState } from "react";
import { useAuth } from "../contexts/AuthContext";
import api from "../api/axios";
import { useTranslation } from "react-i18next";

function ProfilePage() {
  const { user, setUser } = useAuth();

  const [name, setName] = useState(user?.name || "");
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const { t } = useTranslation();


  const handleSave = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    if (newPassword && newPassword !== confirmPassword) {
      setError(t("profile.passwordMismatch"));
      return;
    }

    setSaving(true);
    try {
      const res = await api.post("/user/edit", {
        name,
        currentPassword,
        newPassword,
        confirmPassword,
      });

      // Update context user
      setUser({ ...user, name: res.data.user.name });
      setSuccess(res.data.message)
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");
    } catch (err) {
      if (err.response) {
        setError(err.response.data.error);
      }
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="bg-white shadow-md rounded-lg p-6 min-w-[500px] m-auto">
      <h2 className="text-2xl font-bold text-gray-800 mb-6 text-center">
        {t("profile.title")}
      </h2>

      {error && (
        <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">
          {error}
        </div>
      )}

      {success && (
        <div className="bg-green-100 text-green-600 px-4 py-2 rounded mb-4">
          {success}
        </div>
      )}

      <form onSubmit={handleSave} className="space-y-4">
        {/* Name (editable) */}
        <div>
          <label className="block text-gray-600 text-sm mb-1">{t("profile.name")}</label>
          <input
            type="text"
            className="w-full px-4 py-2  rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </div>

        {/* Email (read-only) */}
        <div>
          <label className="block text-gray-600 text-sm mb-1">{t("profile.email")}</label>
          <input
            type="email"
            value={user?.email || ""}
            readOnly
            className="w-full px-4 py-2  rounded-lg bg-gray-100 text-gray-500 cursor-not-allowed"
          />
        </div>

        {/* Current Password (only if changing password) */}
        <div>
          <label className="block text-gray-600 text-sm mb-1">
            {t("profile.currentPassword")}
          </label>
          <input
            type="password"
            className="w-full px-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400"
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
          />
        </div>

        {/* New Password */}
        <div>
          <label className="block text-gray-600 text-sm mb-1">
            {t("profile.newPassword")}
          </label>
          <input
            type="password"
            className="w-full px-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
          />
        </div>

        {/* Confirm Password */}
        <div>
          <label className="block text-gray-600 text-sm mb-1">
            {t("profile.confirmPassword")}
          </label>
          <input
            type="password"
            className="w-full px-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
          />
        </div>


        {/* Save button */}
        <button
          type="submit"
          disabled={saving}
          className="bg-blue-500 text-white p-2 rounded-lg hover:bg-blue-600 transition disabled:opacity-50"
        >
          {saving ? t("profile.saving"): t("profile.saveChanges")}
        </button>
      </form>
    </div>
  );
}

export default ProfilePage;
