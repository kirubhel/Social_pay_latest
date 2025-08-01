import { v2Client, V2ClientError } from './v2client.service';
import { 
  MerchantResponse, 
  PaymentInitiationRequest, 
  PaymentInitiationResponse, 
  TransactionStatusResponse,
  ReceiptResponse,
  CheckoutDetailsResponse,
  CheckoutPaymentRequest,
  CheckoutPaymentResponse
} from '../types/merchant.dto';

export { V2ClientError };

export const getMerchant = async (merchantId: string): Promise<MerchantResponse> => {
  try {
    return await v2Client.get<MerchantResponse>(`/merchants/${merchantId}`);
  } catch (error) {
    if (error instanceof V2ClientError) {
      throw error;
    }
    throw new V2ClientError('Failed to fetch merchant', 500, 'MERCHANT_FETCH_ERROR');
  }
};

export const initiatePayment = async (request: PaymentInitiationRequest): Promise<PaymentInitiationResponse> => {
  try {
    return await v2Client.post<PaymentInitiationResponse>('/qr/payment/merchant', request);
  } catch (error) {
    if (error instanceof V2ClientError) {
      throw error;
    }
    throw new V2ClientError('Failed to initiate payment', 500, 'PAYMENT_INITIATION_ERROR');
  }
};

export const getTransactionStatus = async (transactionId: string): Promise<TransactionStatusResponse> => {
  try {
    return await v2Client.get<TransactionStatusResponse>(`/payment/transaction/${transactionId}`);
  } catch (error) {
    if (error instanceof V2ClientError) {
      throw error;
    }
    throw new V2ClientError('Failed to get transaction status', 500, 'TRANSACTION_STATUS_ERROR');
  }
};

export const getReceipt = async (transactionId: string): Promise<ReceiptResponse> => {
  try {
    return await v2Client.get<ReceiptResponse>(`/payment/receipt/${transactionId}`);
  } catch (error) {
    if (error instanceof V2ClientError) {
      throw error;
    }
    throw new V2ClientError('Failed to fetch receipt', 500, 'RECEIPT_FETCH_ERROR');
  }
};

export const getCheckoutDetails = async (checkoutId: string): Promise<CheckoutDetailsResponse> => {
  try {
    return await v2Client.get<CheckoutDetailsResponse>(`/checkout/${checkoutId}`);
  } catch (error) {
    if (error instanceof V2ClientError) {
      throw error;
    }
    throw new V2ClientError('Failed to fetch checkout details', 500, 'CHECKOUT_FETCH_ERROR');
  }
};

export const makeCheckoutPayment = async (request: CheckoutPaymentRequest): Promise<CheckoutPaymentResponse> => {
  try {
    return await v2Client.post<CheckoutPaymentResponse>('/checkout/makepayment', request);
  } catch (error) {
    if (error instanceof V2ClientError) {
      throw error;
    }
    throw new V2ClientError('Failed to process checkout payment', 500, 'CHECKOUT_PAYMENT_ERROR');
  }
}; 