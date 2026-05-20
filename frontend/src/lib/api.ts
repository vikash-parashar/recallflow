import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
apiClient.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

// Handle auth errors
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      if (typeof window !== 'undefined') {
        localStorage.removeItem('auth_token');
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export default apiClient;

// API methods
export const auth = {
  register: (data: {
    organization_name: string;
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    phone: string;
  }) => apiClient.post('/auth/register', data),
  
  login: (email: string, password: string) =>
    apiClient.post('/auth/login', { email, password }),
  
  getMe: () => apiClient.get('/users/me'),
};

export const conversations = {
  list: () => apiClient.get('/conversations'),
  get: (id: string) => apiClient.get(`/conversations/${id}`),
  getMessages: (id: string) => apiClient.get(`/conversations/${id}/messages`),
  resolve: (id: string) => apiClient.post(`/conversations/${id}/resolve`),
};

export const dashboard = {
  getStats: () => apiClient.get('/dashboard/stats'),
  getAnalytics: () => apiClient.get('/dashboard/analytics'),
};

export const calls = {
  list: () => apiClient.get('/calls'),
  get: (id: string) => apiClient.get(`/calls/${id}`),
};
