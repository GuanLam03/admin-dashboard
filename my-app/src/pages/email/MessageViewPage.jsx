import { useEffect, useState } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import api from "../../api/axios";
import ReplyIcon from '@mui/icons-material/Reply';
import ReplyBox from "./ReplyBox";

export default function MessageViewPage() {
  const { folder,id } = useParams(); 
  const [searchParams] = useSearchParams();
  const email = searchParams.get("email"); // query param

  const [thread, setThread] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

   const [showReply, setShowReply] = useState(false);

  useEffect(() => {
    async function fetchThread() {
      try {
        const res = await api.get(`/gmail/${folder}/messages/${id}?email=${email}`);
        setThread(res.data);
      } catch (err) {
        setError("Failed to load thread");
      } finally {
        setLoading(false);
      }
    }

    if (id && email) {
      fetchThread();
    }
  }, [id, email]);

  if (loading) return <p className="p-4">Loading...</p>;
  if (error) return <p className="p-4 text-red-500">{error}</p>;
  if (!thread || !thread.messages || thread.messages.length === 0) return <p className="p-4">No messages found</p>;

  return (
    <div className="space-y-6">
      {thread.messages.map((message, index) => (
        <div key={message.id} className="p-6 bg-white rounded-lg shadow-sm">
          <h2 className="text-lg font-bold mb-2">{message.subject}</h2>
          <div className="text-sm text-gray-500 mb-4">
            <p>
              <span className="font-semibold">From:</span> {message.from}
            </p>
            <p>
              <span className="font-semibold">Date:</span> {message.date}
            </p>
          </div>

          <div dangerouslySetInnerHTML={{ __html: message.body }} />

          {/* Show Reply button only on the last message */}
          {index === thread.messages.length - 1 && (
            <>
              {!showReply && (
                <div
                  className="flex justify-between items-center mt-4 border border-gray-400 rounded-full px-4 py-2 w-[120px] text-gray-600 font-semibold hover:bg-gray-100 cursor-pointer"
                  onClick={() => setShowReply(true)}
                >
                  <ReplyIcon />
                  <span className="ml-2">Reply</span>
                </div>
              )}

              {showReply && (
                <ReplyBox
                  mailbox={folder}
                  messageId={message.id}
                  email={email}
                  onSent={() => {
                    setShowReply(false);
                    fetchThread(); // refresh thread after reply
                  }}
                  onCancel={() => setShowReply(false)}
                />
              )}
            </>
          )}
        </div>
      ))}
    </div>
  );
}
