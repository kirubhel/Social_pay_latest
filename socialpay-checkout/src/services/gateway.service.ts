
export interface Gateway {
  can_process: boolean;
  can_settle: boolean;
  created_at: string;
  icon: string;
  id: string;
  key: string;
  name: string;
  short_name: string;
  type: 'WALLET' | 'BANK' | 'CARD';
  updated_at: string;
}

export interface GatewayResponse {
  success: boolean;
  data: Gateway[];
}

import axios, { AxiosError } from 'axios';

export interface Gateway {
  can_process: boolean;
  can_settle: boolean;
  created_at: string;
  icon: string;
  id: string;
  key: string;
  name: string;
  short_name: string;
  type: 'WALLET' | 'BANK' | 'CARD';
  updated_at: string;
}

export interface GatewayResponse {
  success: boolean;
  data: Gateway[];
}

// Map gateway keys to payment method types
export const gatewayKeyToPaymentMethod: Record<string, 'wallet' | 'bank' | 'card' | 'socialpay'> = {
  TELEBIRR: 'wallet',
  CBE: 'wallet',
  CYBERSOURCE: 'card',
  AWINETAA: 'bank',
  BUNAETAA: 'bank',
  ETHSWITCH:'card'
};

export class GatewayError extends Error {
  constructor(
    message: string,
    public statusCode: number,
    public type: 'NETWORK_ERROR' | 'BAD_REQUEST' | 'USER_ERROR' | 'NOT_FOUND'
  ) {
    super(message);
    this.name = 'GatewayError';
  }
}

const baseURL = process.env.NEXT_PUBLIC_API_URL;

export const getGateways = async (): Promise<Gateway[]> => {
  try {
    const response = await axios.get<GatewayResponse>(`${baseURL}/checkout/gateways`);
    if (!response.data.success) {
      throw new GatewayError(
        'Failed to fetch payment gateways',
        response.status,
        'USER_ERROR'
      );
    }
    return response.data.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      const axiosError = error as AxiosError;
      if (!axiosError.response) {
        throw new GatewayError(
          'Network error. Please check your internet connection.',
          0,
          'NETWORK_ERROR'
        );
      }
      
      switch (axiosError.response.status) {
        case 400:
          throw new GatewayError(
            'Invalid request to payment gateway',
            400,
            'BAD_REQUEST'
          );
        case 404:
          throw new GatewayError(
            'Payment gateway service not found',
            404,
            'NOT_FOUND'
          );
        case 500:
        default:
          throw new GatewayError(
            'Internal server error. Please try again later.',
            500,
            'USER_ERROR'
          );
      }
    }
    throw new GatewayError(
      'An unexpected error occurred',
      500,
      'USER_ERROR'
    );
  }
};