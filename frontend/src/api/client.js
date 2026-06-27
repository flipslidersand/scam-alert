import axios from 'axios';

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const client = axios.create({
  baseURL: API_BASE,
  timeout: 30000,
});

export const analyzeImage = async (file) => {
  const formData = new FormData();
  formData.append('file', file);

  try {
    const response = await client.post('/api/analyze', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  } catch (error) {
    throw error.response?.data || { error: 'Failed to analyze image' };
  }
};

export const reportPattern = async (patternId) => {
  try {
    const response = await client.post('/api/report', { patternId });
    return response.data;
  } catch (error) {
    throw error.response?.data || { error: 'Failed to report pattern' };
  }
};

export default client;
