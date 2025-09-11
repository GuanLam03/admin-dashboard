
import { useEffect, useState } from "react";
import DataTable from "datatables.net-dt";
import api from "../api/axios";
import { useNavigate } from "react-router-dom";

function RolePage() {
  const [roleName, setRoleName] = useState("");
  const [roles, setRoles] = useState([]);
  const navigate = useNavigate();

  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const fetchRoles = async () => {
    try {
      const res = await api.get("/roles");
      // console.log(res);
      setRoles(res.data.message.map(r => ({
        id: r.id,
        name: r.Name, // convert to lowercase key
        created_at: r.created_at,
        updated_at: r.updated_at,
      })));

    } catch (err) {
      console.log(err);
      setError(`Error fetching roles ${err.message}`);
    }
  };

  useEffect(() => {
    fetchRoles();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!roleName) {
      alert("Please enter role name");
      return;
    }

    try {
      await api.post("/roles", { name: roleName });
      setSuccess("Role created!");
      setRoleName("");
      fetchRoles(); // reload table
    } catch (err) {
      console.error("Error creating role", err);
    }
  };

  useEffect(() => {
  if (!roles) return;

  const table = new DataTable("#rolesTable", {
    data: roles,
    destroy: true,
    columns: [
      { data: "name", title: "Name", className: "dt-left" },
      { data: "created_at", title: "Created At", className: "dt-left" },
      {
        data: "id",
        title: "Action",
        render: (id) =>
          `<button class="edit-btn btn btn-primary" data-id="${id}">Edit</button>`,
      },
    ],
  });

  const handleClick = (e) => {
    if (e.target && e.target.classList.contains("edit-btn")) {
      const roleId = e.target.getAttribute("data-id");
      navigate(`/roles/edit/${roleId}`);
    }
  };

  document.addEventListener("click", handleClick);

  return () => {
    table.destroy();
    document.removeEventListener("click", handleClick);
  };
}, [roles, navigate]);


  return (
    <div>

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

      <h2 className="text-xl font-bold mb-4">Create Role</h2>

      <form onSubmit={handleSubmit} className="flex items-center gap-4 mb-4 justify-between bg-white p-4 rounded shadow-sm">
        
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
       
        <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded">
          Create Role
        </button>
      </form>

      <div className="bg-white shadow-md rounded-lg p-4">
        <table id="rolesTable" className="display" style={{ width: "100%" }}></table>
      </div>
    </div>
  );
}

export default RolePage;
