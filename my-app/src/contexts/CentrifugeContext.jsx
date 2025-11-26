// import { createContext, useContext, useEffect, useRef } from "react";
// import { Centrifuge } from "centrifuge";
// import api from "../api/axios";


// const CentrifugeContext = createContext(null);

// export function CentrifugeProvider({ children }) {

//   const centrifugeRef = useRef(null);

//   async function getNewToken() {
//     const res = await api.get("/centrifugo/token");
//     return res.data.token;
//   }

//   useEffect(() => {
//     async function init() {
//       const token = await getNewToken();

//       const centrifuge = new Centrifuge("ws://localhost:8000/connection/websocket", {
//         token,
//         refresh: async () => {
//           const newToken = await getNewToken();
//           return { token: newToken };
//         }
//       });

//       centrifuge.on("connect", () => console.log("Connected to Centrifugo"));
//       centrifuge.on("disconnect", () => console.log("Disconnected from Centrifugo"));

//       centrifuge.connect();

//       centrifugeRef.current = centrifuge;


//     }

//     init();

//     return () => {
//       if (centrifugeRef.current) {
//         centrifugeRef.current.disconnect();
//       }
//     };
//   }, []);

//   return (
//     <CentrifugeContext.Provider value={centrifugeRef}>
//       {children}
//     </CentrifugeContext.Provider>
//   );
// }

// export function useCentrifuge() {
//   return useContext(CentrifugeContext);
// }


// import { createContext, useContext, useEffect, useRef, useState } from "react";
// import { Centrifuge } from "centrifuge";
// import api from "../api/axios";

// const CentrifugeContext = createContext({});

// export function CentrifugeProvider({ children }) {
//     const centrifugeRef = useRef(null);
//     const [onlineUsers, setOnlineUsers] = useState([]);

//     async function getNewToken() {
//         const res = await api.get("/centrifugo/token");
//         return res.data.token;
//     }

//     useEffect(() => {
//         (async () => {
//             const token = await getNewToken();

//             const centrifuge = new Centrifuge("ws://localhost:8000/connection/websocket", {
//                 token,
//                 refresh: async () => {
//                     const newToken = await getNewToken();
//                     return { token: newToken };
//                 }
//             });

//             centrifuge.on('connecting', function (ctx) {
//                 console.log(`connecting: ${ctx.code}, ${ctx.reason}`);
//             }).on('connected', function (ctx) {
//                 console.log(`connected over ${ctx.transport}`);
//             }).on('disconnected', function (ctx) {
//                 console.log(`disconnected: ${ctx.code}, ${ctx.reason}`);
//             }).connect();

//             const sub = centrifuge.newSubscription("online_users");

//               sub.on("join", (ctx) => {
//                 console.log("User joined:", ctx.user);
//                 setOnlineUsers((prev) => {
//                   if (prev.includes(ctx.user)) return prev;
//                   return [...prev, ctx.user];
//                 });
//               });

//               sub.on("leave", (ctx) => {
//                 console.log("User left:", ctx.user);
//                 setOnlineUsers((prev) => prev.filter((u) => u !== ctx.user));
//               });


//             sub.on("publication", (ctx) => {
//                 console.log("Publication on online_users:", ctx.data);
//                 const userId = ctx.data.user;

//                 setOnlineUsers(prev => {
//                     if (prev.includes(userId)) return prev; // already marked online
//                     return [...prev, userId]; // add new online user
//                 });
//             });

//             sub.subscribe();

//             const centrifuge2 = new Centrifuge("ws://localhost:8000/api/presence", {
//                 token,
//                 refresh: async () => {
//                     const newToken = await getNewToken();
//                     return { token: newToken };
//                 }
//             });

//             const members = await centrifuge2.presence("public:online_users");
//             console.log("Presence members:", members);

//             const userIds = Object.keys(members.presence || {});
//             console.log("User IDs:", userIds);

//             setOnlineUsers(userIds);



//             centrifugeRef.current = centrifuge;
//         })();

//         return () => {
//             if (centrifugeRef.current) {
//                 centrifugeRef.current.disconnect();
//             }
//         };
//     }, []);

//     return (
//         <CentrifugeContext.Provider value={{ centrifugeRef, onlineUsers }}>
//             {children}
//         </CentrifugeContext.Provider>
//     );
// }

// export function useCentrifuge() {
//     return useContext(CentrifugeContext);
// }


// import { createContext, useContext, useEffect, useRef, useState } from "react";
// import { Centrifuge } from "centrifuge";
// import api from "../api/axios";
// import { useAuth } from "./AuthContext";

// const CentrifugeContext = createContext({});

// export function CentrifugeProvider({ children }) {
//     const { user, loading } = useAuth();
//     const centrifugeRef = useRef(null);
//     const [onlineUsers, setOnlineUsers] = useState([]);




//     async function getNewToken() {
//         const res = await api.get("/centrifugo/token");
//         return res.data.token;
//     }
//     useEffect(() => {
//     if (loading) return;   // wait until auth finishes
//     if (!user) return;     // user not logged in

//     console.log("Centrifuge start for user:", user.id);

//     let sub;

//     (async () => {
//         const token = await getNewToken();

//         const centrifuge = new Centrifuge("ws://localhost:8000/connection/websocket", {
//             token,
//             refresh: async () => ({ token: await getNewToken() })
//         });

//         centrifuge.connect();

//         // subscribe AFTER connect, AFTER user ready
//         sub = centrifuge.newSubscription("public:online_users");

//         const presence = await centrifuge.presence("public:online_users");
//         console.log("Initial presence:", presence);

//         const clients = presence.clients || {};
//         const userIds = Object.values(clients).map(c => c.user);

//         // include self
//         if (!userIds.includes(String(user.id))) {
//             userIds.push(String(user.id));
//         }

//         setOnlineUsers(userIds);

//         // JOIN
//         sub.on("join", (ctx) => {
//             console.log("User joined:", ctx.user);
//             setOnlineUsers(prev =>
//                 prev.includes(ctx.user) ? prev : [...prev, ctx.user]
//             );
//         });

//         // LEAVE
//         sub.on("leave", (ctx) => {
//             console.log("User left:", ctx.user);
//             setOnlineUsers(prev =>
//                 prev.filter((u) => u !== ctx.user)
//             );
//         });

//         sub.subscribe();

//         centrifugeRef.current = centrifuge;
//     })();

//     return () => {
//         if (sub) sub.unsubscribe();
//         if (centrifugeRef.current) centrifugeRef.current.disconnect();
//     };

// }, [loading, user]);   // <-- FIX HERE


//     return (
//         <CentrifugeContext.Provider value={{ centrifugeRef, onlineUsers }}>
//             {children}
//         </CentrifugeContext.Provider>
//     );
// }

// export function useCentrifuge() {
//     return useContext(CentrifugeContext);
// }



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
    console.log("Online users:", userConnections);
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
        })();

        return () => {
            if (sub) sub.unsubscribe();
            if (centrifugeRef.current) centrifugeRef.current.disconnect();
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
