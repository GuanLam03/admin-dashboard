import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../api/axios";

function PermissionPage() {
  const { id } = useParams(); // 👈 role id from URL
  const [roleName, setRoleName] = useState("");
  const [permissions, setPermissions] = useState([]);
  const [allPermissions, setAllPermissions] = useState([]); // 👈 from backend

  useEffect(() => {
    // Fetch all available permissions (from backend hardcoded list)
    const fetchAllPermissions = async () => {
      try {
        const res = await api.get("/permissions"); // backend returns config.Permissions
        // console.log(res);
        setAllPermissions(res.data.permissions || []);
      } catch (err) {
        console.error("Error fetching permissions", err);
      }
    };

    // Fetch role info
    const fetchRole = async () => {
      try {
        const res = await api.get(`/roles/${id}`);
        setRoleName(res.data.role.Name);
        console.log(res);
        setPermissions(res.data.permissions || []);
      } catch (err) {
        console.error("Error fetching role", err);
      }
    };

    fetchAllPermissions();
    fetchRole();
  }, [id]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    console.log("per:",permissions)
    try {
      await api.post(`/roles/${id}`, {
        name: roleName,
        permissions,
      });
      alert("Role updated!");
    } catch (err) {
      console.error("Error updating role", err);
    }
  };

  const togglePermission = (permKey) => {
    if (permissions.includes(permKey)) {
      setPermissions(permissions.filter((p) => p !== permKey));
    } else {
      setPermissions([...permissions, permKey]);
    }
  };

  // Group permissions by section
    const groupedPermissions = allPermissions.reduce((acc, perm) => {
    const section = perm.key.split('.')[0]; // take first part as section
    if (!acc[section]) acc[section] = [];
    acc[section].push(perm);
    return acc;
    }, {});


  return (
    <div>
      <h2 className="text-xl font-bold mb-4">Edit Role</h2>
      <form
        onSubmit={handleSubmit}
        className="flex flex-col items-start gap-4 mb-4 bg-white p-4 rounded shadow-sm"
      >
        <div className="flex flex-col">
          <label>Role Name :</label>
          <input
            type="text"
            value={roleName}
            onChange={(e) => setRoleName(e.target.value)}
            className="border rounded p-2"
            placeholder="Enter role name"
          />
        </div>

        <div>
          <h4>Permissions</h4>
          {Object.entries(groupedPermissions).map(([section, perms]) => (
            <div key={section} className="mb-4">
            <h6 className="font-bold text-lg">{section.charAt(0).toUpperCase() + section.slice(1)}</h6>
            {perms.map((perm) => (
                <div key={perm.key}>
                <input
                    type="checkbox"
                    id={perm.key}
                    checked={permissions.includes(perm.key)}
                    onChange={() => togglePermission(perm.key)}
                />
                <label htmlFor={perm.key} className="ml-2">
                    {perm.label}
                </label>
                </div>
            ))}
            </div>
        ))}
        </div>

        <button
          type="submit"
          className="bg-blue-600 text-white px-4 py-2 rounded"
        >
          Update
        </button>
      </form>
    </div>
  );
}

export default PermissionPage;
