import { useEffect, useState } from "react";
import api from "../../api/axios";

export default function InboxPage({ folder, emailAddress }) {
  const [emails, setEmails] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [nextPageToken, setNextPageToken] = useState(null);
  const [loadingMore, setLoadingMore] = useState(false);

  useEffect(() => {
    fetchEmails(null, true); // first load
  }, [folder, emailAddress]);

  async function fetchEmails(pageToken = null, reset = false) {
    try {
      if (reset) {
        setLoading(true);
        setEmails([]);
      } else {
        setLoadingMore(true);
      }

      const url = `/gmail/${folder}/messages?email=${emailAddress}${pageToken ? `&pageToken=${pageToken}` : ""}`;
      const res = await api.get(url);

      const sorted = res.data.messages.sort((a, b) => parseGmailDate(b.date) - parseGmailDate(a.date));

      setEmails(prev =>
        reset
          ? sorted.map(e => ({ ...e, date: formatDate(e.date) }))
          : [...prev, ...sorted.map(e => ({ ...e, date: formatDate(e.date) }))]
      );
      setNextPageToken(res.data.nextPageToken || null);
    } catch (err) {
      setError("Failed to fetch emails. Please try again later.");
      console.error(err);
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
      second: "2-digit",
      hour12: false,
      timeZone: "UTC",
    }).format(date).replace(",", "");
  }

  if (loading) return <div className="p-4">Loading emails...</div>;
  if (error) return <div className="p-4 text-red-600">{error}</div>;

  return (
    <div className="p-4 bg-white rounded-lg shadow-sm">
      <h2 className="text-xl font-bold mb-4">Inbox: {emailAddress}</h2>
      <table className="w-full border-collapse">
        <tbody>
          {emails.map((email, index) => {
            const replyCount = email.replyCount ? `(${email.replyCount + 1})` : "";
            return (
              <tr
                key={email.id}
                onClick={() =>
                  (window.location.href = `/emails/${folder}/${email.id}?email=${emailAddress}`)
                }
                className="border-b cursor-pointer bg-white hover:!bg-gray-100 hover:shadow-lg hover:scale-[1.001] transition"
              >
                <td className="px-2 py-2 w-8">{index + 1}.</td>
                <td className="px-2 py-2 font-normal">
                  {email.from} {replyCount}
                </td>
                <td className="px-2 py-2">
                  <span className="font-normal">{email.subject}</span>{" "}
                  <span className="text-gray-500">{email.snippet}</span>
                </td>
                <td className="px-2 py-2 text-right text-sm text-gray-600 font-normal">
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
          {loadingMore ? "Loading..." : "Load More"}
        </button>
      )}
    </div>
  );
}
