import { useEffect, useState } from "react";
import api from "../../api/axios";
import { MdInbox } from "react-icons/md";
import { FaRegStar } from "react-icons/fa";
import { FaStar } from "react-icons/fa";
import { useTranslation } from "react-i18next";


const tabs = [
  { name: "Inbox", glabel: "inbox", icon: <MdInbox size={20} /> },
  { name: "Starred", glabel: "starred", icon: <FaRegStar size={20} /> },
  // Add more here, e.g.
  // { label: "Sent", folder: "sent", icon: <FaPaperPlane /> },
  // { label: "Trash", folder: "trash", icon: <FaTrash /> },
];


export default function InboxPage({ folder, emailAddress }) {
  const {t} = useTranslation();
  const [emails, setEmails] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [nextPageToken, setNextPageToken] = useState(null);
  const [loadingMore, setLoadingMore] = useState(false);
  const [activeTab, setActiveTab] = useState(tabs[0].glabel); // default is "inbox"


  useEffect(() => {
    fetchEmails(null, true); // first load
  }, [emailAddress,activeTab]);

  async function fetchEmails(pageToken = null, reset = false) {
    try {
      if (reset) {
        setLoading(true);
        setEmails([]);
      } else {
        setLoadingMore(true);
      }

      const url = `/gmail/${folder}/messages?email=${emailAddress}&label=${activeTab}${pageToken ? `&pageToken=${pageToken}` : ""}`;
      const res = await api.get(url);

  
      const sorted = res.data.messages.sort((a, b) => parseGmailDate(b.date) - parseGmailDate(a.date));

      setEmails(prev =>
        reset
          ? sorted.map(e => ({ ...e, date: formatDate(e.date) }))
          : [...prev, ...sorted.map(e => ({ ...e, date: formatDate(e.date) }))]
      );
      setNextPageToken(res.data.nextPageToken || null);
    } catch (err) {
      setError(err.response?.data?.error ? t(err.response.data.error) : "");
      
    } finally {
      setLoading(false);
      setLoadingMore(false);
    }
  }

  function parseGmailDate(str) {
    const clean = str.replace(/\s*\(.*\)$/, "");
    return new Date(clean).getTime();
  }

  function formatDate(dateString) {
  const date = new Date(dateString);

  return new Intl.DateTimeFormat("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    // second: "2-digit",
    hour12: true,
    timeZone: "Asia/Kuala_Lumpur",  // Specify MYT (Malaysia Standard Time)
  }).format(date).replace(",", "");
}


  // toggle star/unstar
  async function toggleStar(emailId, threadId, isCurrentlyStarred) {
    try {
      await api.post(`/gmail/${folder}/messages/${threadId}/star?email=${emailAddress}`);

      // Update UI immediately
      setEmails((prev) =>
        prev.map((e) =>
          e.id === emailId ? { ...e, isStarred: !isCurrentlyStarred } : e
        )
      );
    } catch (err) {
      setError(err.response?.data?.error ? t(err.response.data.error) : "Could not update star status. Try again.");
    }
  }

  if (loading) return <div className="p-4">Loading emails...</div>;
  if (error) return <div className="p-4 text-red-600">{error}</div>;

  return (
    <div className="rounded-lg shadow-sm">
      {/* <h2 className="text-xl font-bold mb-4">Inbox: {emailAddress}</h2> */}

      <ul className="nav nav-tabs">
        {tabs.map((tab,index) => (
          <li key={`${tab.folder}-${index}`} className="nav-item">
            <button
              className={`nav-link ${activeTab === tab.glabel ? "active" : ""}`}
              onClick={() => setActiveTab(tab.glabel)}
            >
              <div className="flex items-center gap-2 text-black">
                {tab.icon}
                <span>{tab.name}</span>
              </div>
            </button>
          </li>
        ))}
      </ul>

      <table className="w-full border-collapse">
        <tbody>
          {emails.map((email, index) => {
            const replyCount = email.replyCount ? `(${email.replyCount + 1})` : "";
            const isUnread = email.isUnread ? "font-bold" : "font-normal";
            const isUnreadBg =email.isUnread ? "bg-white" : "bg-gray-100";
            return (
              <tr
                key={email.id}
                onClick={() => {
                  // Navigate using React Router for SPA navigation
                  window.location.href = `/emails/${folder}/${email.id}?email=${emailAddress}`;
                }}
                //add flex so padding only will effect
                className={`border-b cursor-pointer ${isUnreadBg} hover:!bg-gray-100 hover:shadow-lg hover:scale-[1.001] transition`}
              >
                <td className="px-2 py-2 w-8">{index + 1}.</td>

                <td
                className="text-center w-6 cursor-pointer"
                onClick={(e) => {
                  e.stopPropagation(); // prevent row click
                  toggleStar(email.id, email.threadId, email.isStarred);
                }}
              >
                {email.isStarred ? (
                  <FaStar className="text-yellow-400" />
                ) : (
                  <FaRegStar className="text-gray-600 hover:text-black" />
                )}
              </td>
                
           

                <td className={`px-2 py-2 ${isUnread}`}>
                  {email.from} {replyCount}
                </td>
                <td className="px-2 py-2">
                  <span className={isUnread}>{email.subject}</span>{" - "}
                  <span className="text-gray-500">{email.snippet}</span>
                </td>
                <td className={`px-2 py-2 text-right text-sm text-gray-600 ${isUnread}`}>
                  {email.date}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>

      {nextPageToken && (
        <button
          onClick={() => fetchEmails(nextPageToken)}
          disabled={loadingMore}
          className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-lg disabled:opacity-50"
        >
          {loadingMore ? "Loading..." : t("emailManagement.loadMore")}
        </button>
      )}
    </div>
  );
}
