import { useEffect } from "react";
import DataTable from "datatables.net-dt";
import api from "../api/axios"; // your axios instance
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import dataTableLocales from "../utils/i18n/datatableLocales";

function UserManagementPage() {
  const navigate = useNavigate();
  const {t,i18n} = useTranslation();

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
          language: dataTableLocales[i18n.language],
          
          columns: [
              { data: "role", title: t("userManagement.fields.role") },
              { data: "name", title: t("userManagement.fields.name") },
              { data: "email", title: t("userManagement.fields.email") },
              { data: "created_at", title: t("userManagement.fields.createdAt") },
              {
              data: "id",
              title: t("common.labels.action"),
              render: (id) =>
                  `<button class="edit-btn btn btn-primary" data-id="${id}">${t("common.buttons.edit")}</button>`,
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
      <h2 className="text-2xl font-bold mb-4">{t("userManagement.userManagement")}</h2>
      <div className="bg-white shadow-md rounded-lg p-4">
        <table id="usersTable" className="display" style={{ width: "100%" }}></table>
      </div>
    </div>
  );
}

export default UserManagementPage;
