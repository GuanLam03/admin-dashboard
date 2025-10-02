import { useEffect, useState } from "react";
import { FcGoogle } from "react-icons/fc";
import { RiLogoutBoxLine } from "react-icons/ri";
import { HiPlus, HiPencil, HiTrash } from "react-icons/hi";
import api from "../../../api/axios";
import { MdExpandMore } from "react-icons/md";
import { MdExpandLess } from "react-icons/md";
import CloseIcon from '@mui/icons-material/Close';
import Editor from "../../../tools/Editor";

export default function EmailsSettingsPage() {
  const [accounts, setAccounts] = useState([]);
  const [teams, setTeams] = useState({});
  const [templates, setTemplates] = useState({});
  const [expanded, setExpanded] = useState({}); // collapsible control

  // Modal state
  const [showModal, setShowModal] = useState(false);
  const [editingTemplate, setEditingTemplate] = useState(null);

  useEffect(() => {
    fetchAccounts();
    fetchTeams();
    fetchTemplates();
  }, []);

  async function fetchAccounts() {
    const res = await api.get("/gmail/accounts");
    setAccounts(res.data);
  }

  async function fetchTeams() {
    const res = await api.get("/gmail/accounts/teams");
    setTeams(res.data);
   
  }

  async function fetchTemplates() {
    const res = await api.get("/gmail/templates");
    // group by team
    const grouped = res.data.reduce((acc, t) => {
      if (!acc[t.team]) acc[t.team] = [];
      acc[t.team].push(t);
      return acc;
    }, {});
    setTemplates(grouped);
    
  }

  async function handleLogin(team) {
    const res = await api.get(`/gmail/auth?team=${team}`);
    const link = document.createElement("a");
    link.href = res.data;
    link.target = "_blank";
    link.click();
  }

  async function handleLogout(email) {
    await api.post(`/gmail/remove-accounts/${email}`);
    fetchAccounts();
  }

  const handleAddTemplate = (team) => {
    setEditingTemplate({ id: null, team, name: "", content: "" });
    setShowModal(true);
  };

  const handleEditTemplate = (team, template) => {
    setEditingTemplate({ ...template, team });
    setShowModal(true);
  };

  const handleDeleteTemplate = async (team, id) => {
    await api.post(`/gmail/templates/remove/${id}`);
    fetchTemplates();
  };

  const handleSaveTemplate = async () => {
    if (editingTemplate.id) {
      // update
      await api.post(
        `/gmail/templates/edit/${editingTemplate.id}`,
        editingTemplate
      );
    } else {
      // create
      await api.post(`/gmail/templates`, editingTemplate);
    }
    setShowModal(false);
    setEditingTemplate(null);
    fetchTemplates();
  };

  return (
    <div className="space-y-10">
      {/* Accounts Section */}
      <div className="bg-white p-4 rounded-lg shadow-sm">
        <h2 className="text-xl font-bold mb-4">Email Accounts</h2>
        <table className="w-full border-collapse">
          <thead>
            <tr className="border-b ">
              <th className="text-left p-2">Teams</th>
              <th className="text-left p-2">Email</th>
              <th className="text-left p-2">Action</th>
            </tr>
          </thead>
          <tbody>
            {Object.entries(teams).map(([key, display]) => {
              const account = accounts.find((a) => a.team === key);

              return (
                <tr key={key} className="border-b">
                  <td className="p-2">{display}</td>
                  <td className="p-2">
                    {account ? account.email : "Not connected"}
                  </td>
                  <td className="p-2">
                    {account ? (
                      <button
                        onClick={() => handleLogout(account.email)}
                        className="flex gap-2 justify-between items-center px-4 py-2 bg-gray-100 text-black font-semibold !rounded-full hover:bg-gray-200 transition"
                      >
                        <RiLogoutBoxLine size={25} className="text-blue-500" />
                        Logout
                      </button>
                    ) : (
                      <button
                        onClick={() => handleLogin(key)}
                        className="flex gap-2 items-center justify-between px-4 py-2 bg-gray-100 text-black font-semibold !rounded-full hover:bg-gray-200 transition"
                      >
                        <FcGoogle size={25} />
                        <span>Connect</span>
                      </button>
                    )}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {/* Templates Section */}
      <div className="bg-white p-4 rounded-lg shadow-sm">
        <h2 className="text-xl font-bold mb-4">Templates</h2>
        {Object.entries(teams).map(([teamKey, display]) => (
          <div key={teamKey} className="mb-6 border rounded-lg">
            {/* Team Header */}
            <div className="flex justify-between items-center bg-gray-100 p-3 cursor-pointer"
              onClick={() => setExpanded(prev => ({...prev, [teamKey]: !prev[teamKey]}))}
            >
              <span className="font-semibold">{display}</span>
              <div className="flex items-center gap-3">
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    handleAddTemplate(teamKey);
                  }}
                  className="flex items-center gap-1 px-3 py-1 text-sm bg-blue-600 text-white !rounded-full hover:bg-blue-700"
                >
                  <HiPlus /> Add
                </button>
                <span>{expanded[teamKey] ?  <MdExpandMore />: <MdExpandLess />}</span>
              </div>
            </div>

            {/* Template List */}
            {expanded[teamKey] && (
              <div className="p-3 space-y-2">
                {(templates[teamKey] || []).length === 0 && (
                  <p className="text-gray-500 text-sm">No templates yet</p>
                )}
                {(templates[teamKey] || []).map((t) => (
                  <div
                    key={t.id}
                    className="flex justify-between items-center border p-2 rounded"
                  >
                    {/* <div>
                      <p className="font-medium">{t.name}</p>
                      <p className="text-gray-600 text-sm line-clamp-2">
                        {t.content}
                      </p>
                    </div> */}
                    <div>
                      <p className="font-medium">{t.name}</p>
                      <div
                        className="text-gray-600 text-sm line-clamp-2"
                        dangerouslySetInnerHTML={{ __html: t.content }}
                      />
                    </div>
                    <div className="flex gap-2">
                      <button
                        onClick={() => handleEditTemplate(teamKey, t)}
                        className="p-2 rounded hover:bg-gray-100"
                      >
                        <HiPencil />
                      </button>
                      <button
                        onClick={() => handleDeleteTemplate(teamKey, t.id)}
                        className="p-2 rounded hover:bg-gray-100 text-red-500"
                      >
                        <HiTrash />
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        ))}
      </div>

      {/* Template Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 size-full">
          <div className="bg-white rounded-lg shadow-lg p-4">
            <div className="flex justify-between items-center mb-2">
              <h3 className="text-lg font-semibold">
                {editingTemplate?.id ? "Edit Template" : "Add Template"}
              </h3>
              <button onClick={() => setShowModal(false)}><CloseIcon /></button>
            </div>
            
            <input
              type="text"
              className="w-full border p-2 mb-3 rounded"
              placeholder="Template name"
              value={editingTemplate?.name || ""}
              onChange={(e) =>
                setEditingTemplate({ ...editingTemplate, name: e.target.value })
              }
            />
            {/* <textarea
              className="w-full border p-2 rounded h-[50vh]"
              placeholder="Template content"
              value={editingTemplate?.content || ""}
              onChange={(e) =>
                setEditingTemplate({
                  ...editingTemplate,
                  content: e.target.value,
                })
              }
            /> */}
            <Editor
              value={editingTemplate?.content || ""}
              onChange={(content) =>
                setEditingTemplate({ ...editingTemplate, content })
              }
            />

            <div className="flex justify-end gap-2 mt-4">
              <button
                className="px-4 py-2 bg-gray-300 rounded hover:bg-gray-400"
                onClick={() => setShowModal(false)}
              >
                Cancel
              </button>
              <button
                className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                onClick={handleSaveTemplate}
              >
                Save
              </button>
            </div>
          </div>
        </div>
      )}

 
    </div>
  );
}
