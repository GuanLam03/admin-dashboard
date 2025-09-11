import { useEffect } from "react";
import DataTable from "datatables.net-dt";
import api from "../api/axios"; // your axios instance
import { useNavigate } from "react-router-dom";

function UserManagementPage() {
  const navigate = useNavigate();

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const res = await api.get("/users"); // your backend endpoint
        const data = res.data.users.map((u) => ({
          id: u.id,
          role: u.role || "",
          name: u.name,
          email: u.email,
          created_at: u.created_at,
        }));

        // initialize DataTable
        const table = new DataTable("#usersTable", {
        data,
        destroy: true,
        columns: [
            { data: "role", title: "Role" },
            { data: "name", title: "Name" },
            { data: "email", title: "Email" },
            { data: "created_at", title: "Created At" },
            {
            data: "id",
            title: "Action",
            render: (id) =>
                `<button class="edit-btn btn btn-primary" data-id="${id}">Edit</button>`,
            },
        ],
        responsive: true,
        createdRow: (row, rowData) => {
            const btn = row.querySelector(".edit-btn");
            btn.addEventListener("click", () => {
            navigate(`/user-management/edit/${rowData.id}`);
            });
        },
        });


        return () => {
          table.destroy();
          document.removeEventListener("click", handleClick);
        };
      } catch (err) {
        console.error("Error fetching users", err);
      }
    };

    fetchUsers();
  }, [navigate]);

  return (
    <div className="w-full">
      <h2 className="text-2xl font-bold mb-4">User Management</h2>
      <div className="bg-white shadow-md rounded-lg p-4">
        <table id="usersTable" className="display" style={{ width: "100%" }}></table>
      </div>
    </div>
  );
}

export default UserManagementPage;
