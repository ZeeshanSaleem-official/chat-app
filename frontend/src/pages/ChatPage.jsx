import { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../context/AuthContext';
import { useWebSocket } from '../context/WebSocketContext';
import api from '../services/api';
import UserSearch from '../components/Contacts/UserSearch';
import ContactList from '../components/Contacts/ContactList';
import ChatWindow from '../components/Chat/ChatWindow';
import '../styles/layout.css';

export default function ChatPage() {
  const { user, logout } = useAuth();
  const { onConversationUpdate, isConnected } = useWebSocket();
  const [conversations, setConversations] = useState([]);
  const [activeConversation, setActiveConversation] = useState(null);
  const [loading, setLoading] = useState(true);

  // Fetch conversations
  const fetchConversations = useCallback(async () => {
    try {
      const data = await api.getConversations();
      setConversations(data);
    } catch (err) {
      console.error('Failed to fetch conversations:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchConversations();
  }, [fetchConversations]);

  // Refresh conversations on WebSocket updates
  useEffect(() => {
    const unsubscribe = onConversationUpdate(() => {
      fetchConversations();
    });
    return unsubscribe;
  }, [onConversationUpdate, fetchConversations]);

  const handleSelectConversation = useCallback((conv) => {
    setActiveConversation(conv);
  }, []);

  const handleStartConversation = useCallback((conv) => {
    setActiveConversation(conv);
    fetchConversations();
  }, [fetchConversations]);

  return (
    <div className="chat-layout" id="chat-page">
      {/* Sidebar */}
      <aside className="sidebar" id="sidebar">
        <div className="sidebar-header">
          <h1 className="sidebar-title">ChatApp</h1>
          <div className="sidebar-user">
            <span className="sidebar-user-name">{user?.display_name}</span>
            <button
              className="logout-btn"
              onClick={logout}
              title="Logout"
              id="logout-btn"
            >
              ⏻
            </button>
          </div>
        </div>

        <UserSearch onStartConversation={handleStartConversation} />

        <div className="sidebar-content">
          {loading ? (
            <div style={{ display: 'flex', justifyContent: 'center', padding: '2rem' }}>
              <div className="loading-spinner" />
            </div>
          ) : (
            <ContactList
              conversations={conversations}
              activeConversationId={activeConversation?.id}
              onSelectConversation={handleSelectConversation}
            />
          )}
        </div>

        {/* Connection status */}
        <div style={{
          padding: '8px 24px',
          fontSize: '0.7rem',
          color: isConnected ? 'var(--online-color)' : 'var(--error-color)',
          display: 'flex',
          alignItems: 'center',
          gap: '6px',
          borderTop: '1px solid var(--border-color)',
        }}>
          <span style={{
            width: '6px',
            height: '6px',
            borderRadius: '50%',
            background: isConnected ? 'var(--online-color)' : 'var(--error-color)',
          }} />
          {isConnected ? 'Connected' : 'Reconnecting...'}
        </div>
      </aside>

      {/* Main Area */}
      <main className="main-area">
        {activeConversation ? (
          <ChatWindow
            key={activeConversation.id}
            conversation={activeConversation}
          />
        ) : (
          <div className="empty-state">
            <div className="empty-state-icon">💬</div>
            <h2>Welcome to ChatApp</h2>
            <p>Select a conversation or search for users to start chatting</p>
          </div>
        )}
      </main>
    </div>
  );
}
