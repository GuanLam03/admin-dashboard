
import { createContext, useContext, useEffect, useRef, useState } from "react";
import { Centrifuge } from "centrifuge";
import api from "../api/axios";
import { useAuth } from "./AuthContext";

const CentrifugeContext = createContext({});
const WEBSOCKET_URL = import.meta.env.VITE_WEBSOCKET_URL;

export function CentrifugeProvider({ children }) {
    const { user, loading } = useAuth();
    const centrifugeRef = useRef(null);

    const [connected, setConnected] = useState(false);

    // connection tracking
    const [userConnections, setUserConnections] = useState({});

    // derived online users
    const onlineUsers = Object.keys(userConnections);

    async function getNewToken() {
        const res = await api.get("/centrifugo/token");
        return res.data.token;
    }

    useEffect(() => {
        if (loading || !user) return;

        let centrifuge;
        let sub;

        (async () => {
            const token = await getNewToken();

            centrifuge = new Centrifuge(`${WEBSOCKET_URL}/connection/websocket`, {
                token,
                refresh: async () => ({ token: await getNewToken() })
            });

            centrifuge.on("connected", () => setConnected(true));
            centrifuge.on("disconnected", () => setConnected(false));

            centrifuge.connect();
            centrifugeRef.current = centrifuge;

            // ---------- SUBSCRIBE ----------
            sub = centrifuge.newSubscription("public:online_users");

            // JOIN
            sub.on("join", (ctx) => {
                const userId = String(ctx.info.user);

                setUserConnections(prev => ({
                    ...prev,
                    [userId]: (prev[userId] || 0) + 1
                }));
            });

            // LEAVE
            sub.on("leave", (ctx) => {
                const userId = String(ctx.info.user);

                setUserConnections(prev => {
                    const count = prev[userId] || 0;

                    if (count <= 1) {
                        const updated = { ...prev };
                        delete updated[userId];
                        return updated;
                    }

                    return {
                        ...prev,
                        [userId]: count - 1
                    };
                });
            });

            sub.subscribe();

            // ---------- INITIAL PRESENCE ----------
            const presence = await centrifuge.presence("public:online_users");
            const initial = {};

            Object.values(presence.clients || {}).forEach(c => {
                const uid = String(c.user);
                initial[uid] = (initial[uid] || 0) + 1;
            });


            setUserConnections(initial);
            // console.log("Initial presence loaded:", initial);   
        })();

        return () => {
            if (sub) sub.unsubscribe();
            
        };

    }, [loading, user]);

    return (
        <CentrifugeContext.Provider
            value={{
                centrifuge: centrifugeRef,
                connected,
                onlineUsers,
            }}
        >
            {children}
        </CentrifugeContext.Provider>
    );
}

export function useCentrifuge() {
    return useContext(CentrifugeContext);
}
