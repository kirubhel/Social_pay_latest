import axios, { AxiosInstance, AxiosError } from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_V2_URL || 'http://196.190.251.194:8082/api/v2';

export class V2ClientError extends Error {
  constructor(
    message: string,
    public status: number,
    public type: string
  ) {
    super(message);
    this.name = 'V2ClientError';
  }
}

class V2Client {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
      timeout: 30000, // 30 seconds timeout
    });

    // Add response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response) {
          // Server responded with error status
          const errorData = error.response.data as { message?: string };
          const message = errorData?.message || error.message;
          throw new V2ClientError(
            message,
            error.response.status,
            'HTTP_ERROR'
          );
        } else if (error.request) {
          // Request was made but no response received
          throw new V2ClientError(
            'Network error occurred',
            0,
            'NETWORK_ERROR'
          );
        } else {
          // Something else happened
          throw new V2ClientError(
            error.message,
            0,
            'REQUEST_ERROR'
          );
        }
      }
    );
  }

  async get<T>(endpoint: string): Promise<T> {
    try {
      const response = await this.client.get<T>(endpoint);
      return response.data;
    } catch (error) {
      if (error instanceof V2ClientError) {
        throw error;
      }
      throw new V2ClientError('Unexpected error occurred', 500, 'UNKNOWN_ERROR');
    }
  }

  async post<T>(endpoint: string, data?: unknown): Promise<T> {
    try {
      const response = await this.client.post<T>(endpoint, data);
      return response.data;
    } catch (error) {
      if (error instanceof V2ClientError) {
        throw error;
      }
      throw new V2ClientError('Unexpected error occurred', 500, 'UNKNOWN_ERROR');
    }
  }

  async put<T>(endpoint: string, data?: unknown): Promise<T> {
    try {
      const response = await this.client.put<T>(endpoint, data);
      return response.data;
    } catch (error) {
      if (error instanceof V2ClientError) {
        throw error;
      }
      throw new V2ClientError('Unexpected error occurred', 500, 'UNKNOWN_ERROR');
    }
  }

  async delete<T>(endpoint: string): Promise<T> {
    try {
      const response = await this.client.delete<T>(endpoint);
      return response.data;
    } catch (error) {
      if (error instanceof V2ClientError) {
        throw error;
      }
      throw new V2ClientError('Unexpected error occurred', 500, 'UNKNOWN_ERROR');
    }
  }
}

export const v2Client = new V2Client(); 