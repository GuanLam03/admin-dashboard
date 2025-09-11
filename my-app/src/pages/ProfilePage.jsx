import { useState } from "react";
import { useAuth } from "../contexts/AuthContext";
import api from "../api/axios";

function ProfilePage() {
  const { user, setUser } = useAuth();

  const [name, setName] = useState(user?.name || "");
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");


  const handleSave = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    if (newPassword && newPassword !== confirmPassword) {
      setError("New passwords do not match");
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
      setSuccess("Update successfully")
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");
    } catch (err) {
      if (err.response) {
        setError(err.response.data.error || "Update failed");
      }
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="bg-white shadow-md rounded-lg p-6 min-w-[500px] m-auto">
      <h2 className="text-2xl font-bold text-gray-800 mb-6 text-center">
        Profile
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
          <label className="block text-gray-600 text-sm mb-1">Name</label>
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
          <label className="block text-gray-600 text-sm mb-1">Email</label>
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
            Current Password
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
            New Password
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
            Confirm Password
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
          {saving ? "Saving..." : "Save Changes"}
        </button>
      </form>
    </div>
  );
}

export default ProfilePage;
