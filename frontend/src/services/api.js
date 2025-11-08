import axios from 'axios';

// In Docker, nginx proxies /api to backend, so use relative URL
// For local dev, use full URL if VITE_API_URL is set
const API_BASE_URL = import.meta.env.VITE_API_URL || '/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const itemsAPI = {
  create: (data) => api.post('/items', data),
  getAll: () => api.get('/items'),
  getById: (id) => api.get(`/items/${id}`),
  delete: (id) => api.delete(`/items/${id}`),
  getRelated: (id) => api.get(`/items/${id}/related`),
  refreshSummary: (id) => api.post(`/items/${id}/refresh-summary`),
};

export const searchAPI = {
  search: (query, limit = 10) => api.get('/search', { params: { q: query, limit } }),
};

export default api;

