import axios from 'axios';
import { getConfig } from '@/config';

const request = axios.create({
    baseURL: getConfig('VITE_API_BASE_URL', '/api'),
    timeout: 60000,
});

request.interceptors.request.use(
    (config) => {
        // Add auth token here if needed
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

request.interceptors.response.use(
    (response) => {
        return response.data;
    },
    (error) => {
        return Promise.reject(error);
    }
);

export default request;
