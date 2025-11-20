import { useEffect, useState } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import api from "../../api/axios";
import ShortcutIcon from '@mui/icons-material/Shortcut';
import MessageActionBox from "./MessageActionBox";
import { useTranslation } from "react-i18next";




export default function MessageViewPage() {
  const {t} = useTranslation();
  const { folder, id } = useParams();
  const [searchParams] = useSearchParams();
  const email = searchParams.get("email"); // query param

  const [thread, setThread] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [showReply, setShowReply] = useState(false);
  const [showForward, setShowForward] = useState(false);

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

  const handleReplyClick = () => {
    setShowReply(true);
    setShowForward(false); // Hide forward button when replying
  };

  const handleForwardClick = () => {
    setShowForward(true);
    setShowReply(false); // Hide reply button when forwarding
  };

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

          {/* Show Reply/Forward buttons only on the last message */}
          {index === thread.messages.length - 1 && (
            <div className="flex flex-col space-y-4 mt-4">
              {/* Show Reply and Forward buttons at the top */}
              {!showReply && !showForward && (
                <div className="flex space-x-4">
                  <div
                    className="flex justify-between items-center border border-gray-400 rounded-full px-4 py-2 w-[120px] text-gray-600 font-semibold hover:bg-gray-100 cursor-pointer"
                    onClick={handleReplyClick}
                  >
                    <ShortcutIcon className="scale-x-[-1]"/>
                    <span className="ml-2">{t("emailManagement.reply")}</span>
                  </div>
                  <div
                    className="flex justify-between items-center border border-gray-400 rounded-full px-4 py-2 w-[150px] text-gray-600 font-semibold hover:bg-gray-100 cursor-pointer"
                    onClick={handleForwardClick}
                  >
                    <ShortcutIcon />
                    <span className="ml-2">{t("emailManagement.forward")}</span>
                  </div>
                </div>
              )}

              {/* Show ReplyBox for Reply */}
              {showReply && (
                <div className="flex flex-col space-y-4">
                  <div className="flex justify-between items-center mb-4">
                    <h3 className="text-xl font-semibold">{t("emailManagement.reply")}</h3>
                  </div>
                  <MessageActionBox
                    mailbox={folder}
                    messageId={message.id}
                    email={email}
                    type="reply"
                    onSent={() => {
                      setShowReply(false);
                      fetchThread(); // Refresh thread after reply
                    }}
                    onCancel={() => setShowReply(false)}
                  />
                </div>
              )}

              {/* Show ReplyBox for Forward */}
              {showForward && (
                <div className="flex flex-col space-y-4">
                  <div className="flex justify-between items-center mb-4">
                    <h3 className="text-xl font-semibold">{t("emailManagement.forward")}</h3>
                  </div>
                  <MessageActionBox
                    mailbox={folder}
                    messageId={message.id}
                    email={email}
                    type="forward"
                    originalMessage={message}
                    onSent={() => {
                      setShowForward(false);
                      fetchThread(); // Refresh thread after forward
                    }}
                    onCancel={() => setShowForward(false)}
                  />
                </div>
              )}
            </div>
          )}
        </div>
      ))}
    </div>
  );
}
