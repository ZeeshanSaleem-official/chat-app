const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

class ApiService {
  constructor() {
    this.token = localStorage.getItem('token');
  }

  setToken(token) {
    this.token = token;
    if (token) {
      localStorage.setItem('token', token);
    } else {
      localStorage.removeItem('token');
    }
  }

  getToken() {
    return this.token || localStorage.getItem('token');
  }

  async request(endpoint, options = {}) {
    const headers = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    const token = this.getToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers,
    });

    if (response.status === 401) {
      this.setToken(null);
      window.location.href = '/login';
      throw new Error('Unauthorized');
    }

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || 'Something went wrong');
    }

    return data;
  }

  // Auth
  async register(userData) {
    const data = await this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
    this.setToken(data.token);
    return data;
  }

  async login(credentials) {
    const data = await this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
    this.setToken(data.token);
    return data;
  }

  logout() {
    this.setToken(null);
    localStorage.removeItem('user');
  }

  // Users
  async getProfile() {
    return this.request('/users/me');
  }

  async searchUsers(query) {
    return this.request(`/users/search?q=${encodeURIComponent(query)}`);
  }

  // Conversations
  async getConversations() {
    return this.request('/conversations');
  }

  async createConversation(userId) {
    return this.request('/conversations', {
      method: 'POST',
      body: JSON.stringify({ user_id: userId }),
    });
  }

  // Messages
  async getMessages(conversationId, limit = 50, offset = 0) {
    return this.request(`/conversations/${conversationId}/messages?limit=${limit}&offset=${offset}`);
  }
}

const api = new ApiService();
export default api;
