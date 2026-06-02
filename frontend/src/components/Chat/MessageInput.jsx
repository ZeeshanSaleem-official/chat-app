import { useState, useRef } from 'react';
import '../../styles/chat.css';

export default function MessageInput({ onSend, onTyping, disabled }) {
  const [message, setMessage] = useState('');
  const typingTimeoutRef = useRef(null);

  const handleSubmit = (e) => {
    e.preventDefault();
    const trimmed = message.trim();
    if (!trimmed) return;
    onSend(trimmed);
    setMessage('');
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  const handleChange = (e) => {
    setMessage(e.target.value);

    // Debounced typing indicator
    if (onTyping) {
      if (typingTimeoutRef.current) clearTimeout(typingTimeoutRef.current);
      onTyping();
      typingTimeoutRef.current = setTimeout(() => {}, 2000);
    }
  };

  return (
    <div className="message-input-area">
      <form className="message-input-container" onSubmit={handleSubmit}>
        <input
          id="message-input"
          type="text"
          className="message-input"
          placeholder="Type a message..."
          value={message}
          onChange={handleChange}
          onKeyDown={handleKeyDown}
          disabled={disabled}
          autoComplete="off"
        />
        <button
          type="submit"
          className="send-button"
          disabled={!message.trim() || disabled}
          id="send-button"
        >
          ➤
        </button>
      </form>
    </div>
  );
}
