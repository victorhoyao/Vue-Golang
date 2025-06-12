import axios from 'axios';

export const API_BASE_URL = 'http://localhost:8036';

// Configure axios base URL
axios.defaults.baseURL = API_BASE_URL;

export default API_BASE_URL; 