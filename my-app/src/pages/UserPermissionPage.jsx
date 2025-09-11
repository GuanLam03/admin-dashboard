import React, { useEffect, useState } from "react";
import api from "../api/axios";
import { useParams } from "react-router-dom";

function UserPermissionPage() {
  const { id } = useParams();
  const [roles, setRoles] = useState([]);
  const [selectedRole, setSelectedRole] = useState(null);

  useEffect(() => {
    // fetch roles
    console.log("Current userid: ", id);
    const fetchRoles = async () => {
      try {
        const res = await api.get("/user-management/roles");
        console.log(res);
        setRoles(res.data.message);
      } catch (err) {
        console.error("Failed to fetch roles", err);
      }
    };

    fetchRoles();
  }, []);

  const handleUpdate = async () => {
    if (!selectedRole) {
      alert("Please select a role");
      return;
    }

    try {
      await api.post(`/user-management/${id}/assign-role`, {
        role_id: selectedRole,
      });
      alert("Role updated!");
    } catch (err) {
      console.error("Failed to assign role", err);
      alert("Update failed");
    }
  };

  return (
    <div>
      <h1 className="text-xl font-bold mb-4">Edit User</h1>

      <div className="space-y-2 flex flex-col">
        {roles.map((role) => (
          <label key={role.id} className="flex items-center">
            <input
              type="checkbox"
              checked={selectedRole === role.id}
              onChange={() => setSelectedRole(role.id)}
            />
            <span className="ml-2">{role.Name}</span>
          </label>
        ))}
      </div>

      <button
        onClick={handleUpdate}
        className="mt-4 bg-blue-600 text-white px-4 py-2 rounded"
      >
        Update
      </button>
    </div>
  );
}

export default UserPermissionPage;
