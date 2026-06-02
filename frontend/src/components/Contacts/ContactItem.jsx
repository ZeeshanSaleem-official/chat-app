import '../../styles/contacts.css';

export default function ContactItem({ conversation, isActive, onClick }) {
  const getInitials = (name) => {
    if (!name) return '?';
    return name.split(' ').map(w => w[0]).join('').slice(0, 2);
  };

  const formatTime = (dateStr) => {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    const now = new Date();
    const diff = now - date;
    const dayMs = 86400000;

    if (diff < dayMs) {
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }
    if (diff < dayMs * 7) {
      return date.toLocaleDateString([], { weekday: 'short' });
    }
    return date.toLocaleDateString([], { month: 'short', day: 'numeric' });
  };

  return (
    <div
      className={`contact-item ${isActive ? 'active' : ''}`}
      onClick={onClick}
      id={`contact-${conversation.id}`}
    >
      <div className="contact-avatar-wrapper">
        <div className="contact-avatar">
          {getInitials(conversation.other_user_name)}
        </div>
        {conversation.other_user_online && (
          <div className="online-dot pulse" />
        )}
      </div>

      <div className="contact-info">
        <div className="contact-name">{conversation.other_user_name}</div>
        <div className="contact-last-message">
          {conversation.last_message || 'No messages yet'}
        </div>
      </div>

      <div className="contact-meta">
        <span className="contact-time">
          {formatTime(conversation.last_message_time)}
        </span>
        {conversation.unread_count > 0 && (
          <span className="unread-badge">
            {conversation.unread_count > 9 ? '9+' : conversation.unread_count}
          </span>
        )}
      </div>
    </div>
  );
}
