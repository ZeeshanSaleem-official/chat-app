import '../../styles/chat.css';

export default function MessageBubble({ message, isOwn }) {
  const formatTime = (dateStr) => {
    const date = new Date(dateStr);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };

  return (
    <div className={`message-wrapper ${isOwn ? 'sent' : 'received'}`}>
      <div className="message-bubble">
        {message.content || message.data?.content}
      </div>
      <span className="message-time">
        {formatTime(message.created_at || message.data?.created_at)}
      </span>
    </div>
  );
}
