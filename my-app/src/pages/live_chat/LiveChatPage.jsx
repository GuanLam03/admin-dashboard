import { useEffect, useRef, useState } from "react";
import DataTable from "datatables.net-dt";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";

import { useCentrifuge } from "../../contexts/CentrifugeContext";
import dataTableLocales from "../../utils/i18n/datatableLocales";
import api from "../../api/axios";
import { useAuth } from "../../contexts/AuthContext";

function LiveChatPage() {
  const { user } = useAuth();
  const { t, i18n } = useTranslation();
  const navigate = useNavigate();
  const tableRef = useRef(null);

  const [error, setError] = useState("");
  const [selectedUser, setSelectedUser] = useState(null); // chat target
  const [showChatBox, setShowChatBox] = useState(false);

  const { centrifuge } = useCentrifuge();
  const [messagesMap, setMessagesMap] = useState({});

  const [input, setInput] = useState("");
  const chatBoxRef = useRef(null);





  // global online user list
  const { onlineUsers } = useCentrifuge();

  /** ------------------------------------------
   * 1️⃣ Initialize DataTable
   * -------------------------------------------*/
  useEffect(() => {
    const dt = new DataTable("#usersTable", {
      data: [],
      destroy: true,
      language: dataTableLocales[i18n.language],
      columns: [
        { data: "role", title: t("userManagement.fields.role") },
        { data: "name", title: t("userManagement.fields.name") },
        { data: "email", title: t("userManagement.fields.email") },
        {
          data: "online",
          title: t("userManagement.status"),
          render: (online) =>
            online
              ? '<span class="inline-block w-3 h-3 bg-green-500 rounded-full"></span>'
              : '<span class="inline-block w-3 h-3 bg-red-500 rounded-full"></span>',
        },
        {
          data: "id",
          title: t("common.labels.action"),
          render: (id, type, row) => {
            if (String(id) === String(user.id)) {
              return `<span class="text-gray-400">-</span>`;
            }

            const isOnline = row.online;
            if (!isOnline) return `<span class="text-gray-400">Offline</span>`;

            return `
              <button 
                class="chat-btn bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700"
                data-id="${id}" 
                data-name="${row.name}"
              >
                ${t("common.buttons.chat")}
              </button>
            `;
          },
        },
      ],
      responsive: true,

      createdRow: (row, rowData) => {
        const btn = row.querySelector(".chat-btn");
        if (!btn) return;

        btn.addEventListener("click", () => {
          setSelectedUser({
            id: btn.dataset.id,
            name: btn.dataset.name,
          });
          setShowChatBox(true);
        });
      },
    });

    tableRef.current = dt;
    return () => dt.destroy();
  }, [i18n.language, t]);

  /** ------------------------------------------
   * 2️⃣ Fetch users
   * -------------------------------------------*/
  const fetchUsers = async () => {
    try {
      const res = await api.get("/users");

      const data = res.data.users.map((u) => ({
        id: u.id,
        role: u.role || "",
        name: u.name,
        email: u.email,
        online: onlineUsers.includes(String(u.id)), // presence from Centrifuge
      }));

      const dt = tableRef.current;
      dt.clear();
      dt.rows.add(data);
      dt.draw();
    } catch (err) {
      setError(`Error fetching users: ${err.message}`);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []); 

  useEffect(() => {
    const dt = tableRef.current;

    dt.rows().every(function () {
      const rowData = this.data();
      const isOnline = onlineUsers.includes(String(rowData.id));

      if (rowData.online !== isOnline) {
        rowData.online = isOnline;
        this.data(rowData);
      }
    });

    dt.draw(false);
  }, [onlineUsers]);

  /** ------------------------------------------
   * 3️⃣ Auto-update online status
   * -------------------------------------------*/
  useEffect(() => {
    if (!showChatBox || !selectedUser || !centrifuge) return;

    const chan = `private:chat.${[user.id, selectedUser.id].sort().join(".")}`;

    let sub = centrifuge.current.getSubscription(chan);

    if (!sub) {
      sub = centrifuge.current.newSubscription(chan);

      sub.on("publication", (ctx) => {
        const msg = ctx.data;

        setMessagesMap((prev) => ({
          ...prev,
          [chan]: [...(prev[chan] || []), msg],
        }));

        setTimeout(() => {
          if (chatBoxRef.current) {
            chatBoxRef.current.scrollTop = chatBoxRef.current.scrollHeight;
          }
        }, 50);
      });

      sub.subscribe();
    }

    // Optionally unsubscribe when chat closes
    return () => {
      // sub.unsubscribe();
    };
  }, [showChatBox, selectedUser, centrifuge]);



  const sendMessage = () => {
    if (!input.trim() || !selectedUser || !centrifuge) return;

    const chan = `private:chat.${[user.id, selectedUser.id].sort().join(".")}`;

    const messagePayload = {
      from: user.id,
      to: selectedUser.id,
      text: input,
      time: new Date().toISOString(),
    };

    // Publish to Centrifugo
    centrifuge.current.publish(chan, messagePayload);

    setInput("");

    setTimeout(() => {
      if (chatBoxRef.current) {
        chatBoxRef.current.scrollTop = chatBoxRef.current.scrollHeight;
      }
    }, 50);
  };


  return (
    <>
      {error && (
        <div className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4">
          {error}
        </div>
      )}

      <div className="w-full">
        <h2 className="text-2xl font-bold mb-4">
          {t("userManagement.userManagement")}
        </h2>

        <div className="bg-white shadow-md rounded-lg p-4">
          <table id="usersTable" className="display w-full"></table>
        </div>
      </div>

      {/* Floating Chat Box */}
      {showChatBox && selectedUser && (
        <div className="fixed bottom-4 right-4 w-80 bg-white shadow-xl border rounded-lg flex flex-col">
          {/* Header */}
          <div className="flex justify-between items-center bg-blue-600 text-white px-4 py-2 rounded-t-lg">
            <span>Chat with {selectedUser.name}</span>
            <button
              className="text-white hover:text-gray-200"
              onClick={() => setShowChatBox(false)}
            >
              ✕
            </button>
          </div>

          {/* Chat Messages */}
          <div
            ref={chatBoxRef}
            className="h-64 overflow-y-auto p-3 bg-gray-50 flex flex-col gap-2"
          >
            {(messagesMap[`private:chat.${[user.id, selectedUser.id].sort().join(".")}`] || []).length === 0 && (
              <p className="text-gray-400 text-sm">No messages yet...</p>
            )}

            {(messagesMap[`private:chat.${[user.id, selectedUser.id].sort().join(".")}`] || []).map((msg, idx) => (
              <div
                key={idx}
                className={`max-w-[75%] px-3 py-2 rounded-lg text-sm shadow
        ${msg.from == user.id ? "bg-blue-600 text-white self-end" : "bg-white border self-start"}
      `}
              >
                {msg.text}
                <div className="text-[10px] opacity-70 mt-1">
                  {new Date(msg.time).toLocaleTimeString()}
                </div>
              </div>
            ))}
          </div>



          {/* Input */}
          <div className="p-3 border-t flex gap-2">
            <input
              type="text"
              className="flex-1 border rounded px-2 py-1"
              placeholder="Type your message..."
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && sendMessage()}
            />
            <button
              onClick={sendMessage}
              className="bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700"
            >
              Send
            </button>
          </div>

        </div>
      )}
    </>
  );
}

export default LiveChatPage;
