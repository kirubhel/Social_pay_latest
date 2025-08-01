import { v2Client, V2ClientError } from './v2client.service';
import { 
  QRLinkResponse, 
  QRPaymentRequest, 
  QRPaymentResponse 
} from '../types/qr.dto';

export { V2ClientError };

export const getQRLink = async (qrLinkId: string): Promise<QRLinkResponse> => {
  try {
    return await v2Client.get<QRLinkResponse>(`/qr/payment/link/${qrLinkId}`);
  } catch (error) {
    if (error instanceof V2ClientError) {
      throw error;
    }
    throw new V2ClientError('Failed to fetch QR link', 500, 'QR_LINK_FETCH_ERROR');
  }
};

export const processQRPayment = async (qrLinkId: string, request: QRPaymentRequest): Promise<QRPaymentResponse> => {
  try {
    return await v2Client.post<QRPaymentResponse>(`/qr/payment/link/${qrLinkId}`, request);
  } catch (error) {
    if (error instanceof V2ClientError) {
      throw error;
    }
    throw new V2ClientError('Failed to process QR payment', 500, 'QR_PAYMENT_ERROR');
  }
}; 