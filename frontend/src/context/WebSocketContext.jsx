import { createContext, useContext, useEffect, useState, useCallback, useRef } from 'react';
import wsService from '../services/websocket';

const WebSocketContext = createContext(null);

export function WebSocketProvider({ children }) {
  const [isConnected, setIsConnected] = useState(false);
  const [typingUsers, setTypingUsers] = useState({});
  const [onlineUsers, setOnlineUsers] = useState(new Set());
  const messageListenersRef = useRef([]);
  const conversationListenersRef = useRef([]);
  const typingTimeoutsRef = useRef({});

  useEffect(() => {
    const unsubConnection = wsService.on('connection', ({ status }) => {
      setIsConnected(status === 'connected');
    });

    const unsubMessage = wsService.on('chat_message', (msg) => {
      messageListenersRef.current.forEach(listener => listener(msg));
      conversationListenersRef.current.forEach(listener => listener());
    });

    const unsubTyping = wsService.on('typing', (msg) => {
      const key = msg.conversation_id;
      setTypingUsers(prev => ({ ...prev, [key]: true }));

      // Clear previous timeout for this conversation
      if (typingTimeoutsRef.current[key]) {
        clearTimeout(typingTimeoutsRef.current[key]);
      }

      // Clear typing after 3 seconds
      typingTimeoutsRef.current[key] = setTimeout(() => {
        setTypingUsers(prev => {
          const next = { ...prev };
          delete next[key];
          return next;
        });
      }, 3000);
    });

    const unsubOnline = wsService.on('user_online', (msg) => {
      setOnlineUsers(prev => new Set([...prev, msg.data.user_id]));
    });

    const unsubOffline = wsService.on('user_offline', (msg) => {
      setOnlineUsers(prev => {
        const next = new Set(prev);
        next.delete(msg.data.user_id);
        return next;
      });
    });

    const unsubReadReceipt = wsService.on('read_receipt', () => {
      conversationListenersRef.current.forEach(listener => listener());
    });

    return () => {
      unsubConnection();
      unsubMessage();
      unsubTyping();
      unsubOnline();
      unsubOffline();
      unsubReadReceipt();
      Object.values(typingTimeoutsRef.current).forEach(clearTimeout);
    };
  }, []);

  const onMessage = useCallback((callback) => {
    messageListenersRef.current.push(callback);
    return () => {
      messageListenersRef.current = messageListenersRef.current.filter(l => l !== callback);
    };
  }, []);

  const onConversationUpdate = useCallback((callback) => {
    conversationListenersRef.current.push(callback);
    return () => {
      conversationListenersRef.current = conversationListenersRef.current.filter(l => l !== callback);
    };
  }, []);

  const sendMessage = useCallback((conversationId, content) => {
    wsService.sendMessage(conversationId, content);
  }, []);

  const sendTyping = useCallback((conversationId) => {
    wsService.sendTyping(conversationId);
  }, []);

  const sendReadReceipt = useCallback((conversationId) => {
    wsService.sendReadReceipt(conversationId);
  }, []);

  const isUserOnline = useCallback((userId) => {
    return onlineUsers.has(userId);
  }, [onlineUsers]);

  const value = {
    isConnected,
    typingUsers,
    onlineUsers,
    onMessage,
    onConversationUpdate,
    sendMessage,
    sendTyping,
    sendReadReceipt,
    isUserOnline,
  };

  return (
    <WebSocketContext.Provider value={value}>
      {children}
    </WebSocketContext.Provider>
  );
}

export function useWebSocket() {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error('useWebSocket must be used within a WebSocketProvider');
  }
  return context;
}

export default WebSocketContext;
