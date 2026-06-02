import { useState, useCallback, useRef, useEffect } from 'react';
import api from '../../services/api';
import '../../styles/contacts.css';

export default function UserSearch({ onStartConversation }) {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [showResults, setShowResults] = useState(false);
  const timeoutRef = useRef(null);
  const containerRef = useRef(null);

  const searchUsers = useCallback(async (searchQuery) => {
    if (searchQuery.length < 2) {
      setResults([]);
      setShowResults(false);
      return;
    }

    setLoading(true);
    try {
      const users = await api.searchUsers(searchQuery);
      setResults(users);
      setShowResults(true);
    } catch (err) {
      console.error('Search error:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  const handleInputChange = (e) => {
    const value = e.target.value;
    setQuery(value);

    if (timeoutRef.current) clearTimeout(timeoutRef.current);
    timeoutRef.current = setTimeout(() => {
      searchUsers(value);
    }, 300);
  };

  const handleSelectUser = async (user) => {
    try {
      const conversation = await api.createConversation(user.id);
      onStartConversation({
        ...conversation,
        other_user_id: user.id,
        other_user_name: user.display_name,
        other_user_avatar: user.avatar_url,
        other_user_online: user.is_online,
      });
      setQuery('');
      setShowResults(false);
      setResults([]);
    } catch (err) {
      console.error('Failed to start conversation:', err);
    }
  };

  const getInitials = (name) => {
    if (!name) return '?';
    return name.split(' ').map(w => w[0]).join('').slice(0, 2);
  };

  // Close results on outside click
  useEffect(() => {
    const handleClick = (e) => {
      if (containerRef.current && !containerRef.current.contains(e.target)) {
        setShowResults(false);
      }
    };
    document.addEventListener('mousedown', handleClick);
    return () => document.removeEventListener('mousedown', handleClick);
  }, []);

  return (
    <div className="sidebar-search" ref={containerRef}>
      <div className="search-container">
        <span className="search-icon">🔍</span>
        <input
          id="user-search"
          type="text"
          className="search-input"
          placeholder="Search users to chat..."
          value={query}
          onChange={handleInputChange}
          onFocus={() => results.length > 0 && setShowResults(true)}
        />
      </div>

      {showResults && (
        <div className="search-results">
          <div className="search-results-header">
            {loading ? 'Searching...' : `${results.length} user${results.length !== 1 ? 's' : ''} found`}
          </div>
          {results.map(user => (
            <div
              key={user.id}
              className="search-result-item"
              onClick={() => handleSelectUser(user)}
            >
              <div className="contact-avatar">
                {getInitials(user.display_name)}
              </div>
              <div>
                <div className="search-result-name">{user.display_name}</div>
                <div className="search-result-username">@{user.username}</div>
              </div>
            </div>
          ))}
          {!loading && results.length === 0 && query.length >= 2 && (
            <div className="search-no-results">No users found</div>
          )}
        </div>
      )}
    </div>
  );
}
