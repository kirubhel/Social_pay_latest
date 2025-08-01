import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

// Default API configuration
const API_CONFIG: AxiosRequestConfig = {
  baseURL: process.env.NEXT_PUBLIC_API_URL || '/api',
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000,
};

// Create a singleton axios instance
class ApiClient {
  private static instance: ApiClient;
  private client: AxiosInstance;

  private constructor() {
    this.client = axios.create(API_CONFIG);
    
    // Add request interceptor for authentication, etc.
    this.client.interceptors.request.use(
      (config) => {
        // You can add auth token here
        return config;
      },
      (error) => Promise.reject(error)
    );
    
    // Add response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        // Global error handling
        return Promise.reject(error);
      }
    );
  }

  public static getInstance(): ApiClient {
    if (!ApiClient.instance) {
      ApiClient.instance = new ApiClient();
    }
    return ApiClient.instance;
  }

  public getClient(): AxiosInstance {
    return this.client;
  }
}

// Export a getter for the axios instance
export const getApiClient = (): AxiosInstance => 
  ApiClient.getInstance().getClient(); 