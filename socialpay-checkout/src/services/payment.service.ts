import axios, { AxiosError } from 'axios';
import { mockPaymentService } from './mock/payment.mock';

export interface PaymentDetails {
  card?: {
    name: string;
    pan: string;
    expiry: string;
  };
  phone?: string;
}

interface InitiatePaymentRequest {
  to: string;
  medium: string;
  amount: number;
  details: PaymentDetails;
  redirects: {
    success: string;
    cancel: string;
    declined: string;
  };
}

export interface PaymentResponse {
  success: boolean;
  data: {
    status: {
      current: string;
      next: string;
    };
    transaction?: {
      id: string;
      pricing: {
        amount: number;
        fees: Array<{ transaction: number }>;
        total_amount: number;
      };
      status: {
        value: string;
        msg: string;
      };
    };
    checkout?: {
      type: string;
      data: string;
    };
  };
}

export class PaymentError extends Error {
  constructor(
    message: string,
    public statusCode: number,
    public type: 'NETWORK_ERROR' | 'BAD_REQUEST' | 'USER_ERROR' | 'UNKNOWN'
  ) {
    super(message);
    this.name = 'PaymentError';
  }
}

export const initiatePayment = async (payload: InitiatePaymentRequest): Promise<PaymentResponse> => {
  // Use mock service in development
  if (process.env.NEXT_MOCK_PAYMENT_SERVICE === 'true') {
    return mockPaymentService.initiatePayment(payload.medium);
  }

  try {
    const baseURL = process.env.NEXT_PUBLIC_API_URL;
    const response = await axios.post<PaymentResponse>(`${baseURL}/checkout/init`, payload);
    console.log("response raw", response)
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      const axiosError = error as AxiosError;
      if (!axiosError.response) {
        throw new PaymentError(
          'Network error occurred. Please check your connection.',
          0,
          'NETWORK_ERROR'
        );
      }
      
      switch (axiosError.response.status) {
        case 400:
          throw new PaymentError(
            'Invalid payment details provided.',
            400,
            'BAD_REQUEST'
          );
        case 500:
          throw new PaymentError(
            'Server error occurred. Please try again later.',
            500,
            'USER_ERROR'
          );
        default:
          throw new PaymentError(
            'An unexpected error occurred.',
            axiosError.response.status,
            'UNKNOWN'
          );
      }
    }
    throw new PaymentError(
      'An unexpected error occurred.',
      0,
      'UNKNOWN'
    );
  }
};

export const processPayment = async (txnId: string): Promise<PaymentResponse> => {
  // Use mock service in development
  if (process.env.NEXT_MOCK_PAYMENT_SERVICE === 'true') {
    // Extract medium from txnId (in real implementation this would come from the backend)
    const medium = txnId.includes('CARD') ? 'CYBERSOURCE' : 
                  txnId.includes('WALLET') ? 'TELEBIRR' : 'AWINETAA';
    return mockPaymentService.processPayment(medium);
  }

  try {
    const baseURL = process.env.NEXT_PUBLIC_API_URL;
    const response = await axios.post<PaymentResponse>(`${baseURL}/checkout/process?id=${txnId}`);
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      const axiosError = error as AxiosError;
      if (!axiosError.response) {
        throw new PaymentError(
          'Network error occurred. Please check your connection.',
          0,
          'NETWORK_ERROR'
        );
      }
      
      switch (axiosError.response.status) {
        case 400:
          throw new PaymentError(
            'Invalid transaction ID.',
            400,
            'BAD_REQUEST'
          );
        case 500:
          throw new PaymentError(
            'Server error occurred. Please try again later.',
            500,
            'USER_ERROR'
          );
        default:
          throw new PaymentError(
            'An unexpected error occurred.',
            axiosError.response.status,
            'UNKNOWN'
          );
      }
    }
    throw new PaymentError(
      'An unexpected error occurred.',
      0,
      'UNKNOWN'
    );
  }
}; 