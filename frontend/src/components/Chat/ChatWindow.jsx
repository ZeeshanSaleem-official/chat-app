import { useState, useEffect, useRef, useCallback } from 'react';
import { useAuth } from '../../context/AuthContext';
import { useWebSocket } from '../../context/WebSocketContext';
import api from '../../services/api';
import ChatHeader from './ChatHeader';
import MessageBubble from './MessageBubble';
import MessageInput from './MessageInput';
import '../../styles/chat.css';

export default function ChatWindow({ conversation }) {
  const { user } = useAuth();
  const { onMessage, sendMessage, sendTyping, sendReadReceipt, typingUsers, isUserOnline } = useWebSocket();
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(true);
  const messagesEndRef = useRef(null);
  const messagesAreaRef = useRef(null);

  const isOnline = isUserOnline(conversation.other_user_id) || conversation.other_user_online;
  const isTyping = typingUsers[conversation.id];

  // Fetch message history
  useEffect(() => {
    let cancelled = false;

    const fetchMessages = async () => {
      setLoading(true);
      try {
        const data = await api.getMessages(conversation.id);
        if (!cancelled) {
          setMessages(data);
        }
      } catch (err) {
        console.error('Failed to fetch messages:', err);
      } finally {
        if (!cancelled) setLoading(false);
      }
    };

    fetchMessages();
    sendReadReceipt(conversation.id);

    return () => { cancelled = true; };
  }, [conversation.id, sendReadReceipt]);

  // Listen for incoming messages
  useEffect(() => {
    const unsubscribe = onMessage((msg) => {
      if (msg.conversation_id === conversation.id) {
        const newMessage = {
          id: msg.data.id,
          conversation_id: msg.conversation_id,
          sender_id: msg.sender_id,
          content: msg.data.content,
          message_type: msg.data.message_type,
          is_read: msg.data.is_read,
          created_at: msg.data.created_at,
          sender_name: msg.data.sender_name,
          sender_avatar: msg.data.sender_avatar,
        };

        setMessages(prev => {
          // Avoid duplicates
          if (prev.some(m => m.id === newMessage.id)) return prev;
          return [...prev, newMessage];
        });

        // Mark as read if we're viewing this conversation
        if (msg.sender_id !== user.id) {
          sendReadReceipt(conversation.id);
        }
      }
    });

    return unsubscribe;
  }, [conversation.id, onMessage, sendReadReceipt, user.id]);

  // Auto-scroll to bottom
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages, isTyping]);

  const handleSend = useCallback((content) => {
    sendMessage(conversation.id, content);
  }, [conversation.id, sendMessage]);

  const handleTyping = useCallback(() => {
    sendTyping(conversation.id);
  }, [conversation.id, sendTyping]);

  // Group messages by date
  const getDateLabel = (dateStr) => {
    const date = new Date(dateStr);
    const today = new Date();
    const yesterday = new Date(today);
    yesterday.setDate(yesterday.getDate() - 1);

    if (date.toDateString() === today.toDateString()) return 'Today';
    if (date.toDateString() === yesterday.toDateString()) return 'Yesterday';
    return date.toLocaleDateString([], { month: 'long', day: 'numeric', year: 'numeric' });
  };

  const renderMessages = () => {
    let lastDate = null;

    return messages.map((msg) => {
      const dateLabel = getDateLabel(msg.created_at);
      const showDate = dateLabel !== lastDate;
      lastDate = dateLabel;

      return (
        <div key={msg.id}>
          {showDate && (
            <div className="date-separator">
              <span>{dateLabel}</span>
            </div>
          )}
          <MessageBubble
            message={msg}
            isOwn={msg.sender_id === user.id}
          />
        </div>
      );
    });
  };

  return (
    <div className="chat-container" id="chat-window">
      <ChatHeader conversation={conversation} isOnline={isOnline} />

      <div className="messages-area" ref={messagesAreaRef}>
        {loading ? (
          <div className="messages-loading">
            <div className="loading-spinner" />
          </div>
        ) : (
          <>
            {messages.length === 0 && (
              <div className="empty-state" style={{ padding: '2rem' }}>
                <div style={{ fontSize: '2.5rem' }}>👋</div>
                <p style={{ textAlign: 'center', color: 'var(--text-tertiary)', fontSize: '0.85rem' }}>
                  No messages yet. Say hello to {conversation.other_user_name}!
                </p>
              </div>
            )}
            {renderMessages()}
            {isTyping && (
              <div className="typing-indicator">
                <div className="typing-dots">
                  <span />
                  <span />
                  <span />
                </div>
                <span className="typing-indicator-text">
                  {conversation.other_user_name} is typing...
                </span>
              </div>
            )}
            <div ref={messagesEndRef} />
          </>
        )}
      </div>

      <MessageInput
        onSend={handleSend}
        onTyping={handleTyping}
        disabled={loading}
      />
    </div>
  );
}
