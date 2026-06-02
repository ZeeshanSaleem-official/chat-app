import ContactItem from './ContactItem';

export default function ContactList({ conversations, activeConversationId, onSelectConversation }) {
  if (!conversations || conversations.length === 0) {
    return (
      <div className="search-no-results">
        <p>No conversations yet</p>
        <p style={{ fontSize: '0.75rem', marginTop: '4px', color: 'var(--text-tertiary)' }}>
          Search for users above to start chatting
        </p>
      </div>
    );
  }

  return (
    <div className="contact-list">
      <div className="section-label">Messages</div>
      {conversations.map(conv => (
        <ContactItem
          key={conv.id}
          conversation={conv}
          isActive={conv.id === activeConversationId}
          onClick={() => onSelectConversation(conv)}
        />
      ))}
    </div>
  );
}
