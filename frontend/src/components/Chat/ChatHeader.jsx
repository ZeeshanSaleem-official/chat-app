import '../../styles/chat.css';

export default function ChatHeader({ conversation, isOnline }) {
  const getInitials = (name) => {
    if (!name) return '?';
    return name.split(' ').map(w => w[0]).join('').slice(0, 2);
  };

  return (
    <div className="chat-header" id="chat-header">
      <div className="chat-header-avatar">
        {getInitials(conversation.other_user_name)}
      </div>
      <div className="chat-header-info">
        <div className="chat-header-name">{conversation.other_user_name}</div>
        <div className={`chat-header-status ${isOnline ? 'online' : ''}`}>
          <span className={`status-dot ${isOnline ? 'online' : ''}`} />
          {isOnline ? 'Online' : 'Offline'}
        </div>
      </div>
    </div>
  );
}
