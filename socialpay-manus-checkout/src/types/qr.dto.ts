export interface QRLinkResponse {
  id: string;
  user_id: string;
  merchant_id: string;
  type: 'DYNAMIC' | 'STATIC';
  amount?: number; // null for DYNAMIC, required for STATIC
  supported_methods: string[];
  tag: 'RESTAURANT' | 'DONATION' | 'SHOP';
  title: string;
  description: string;
  image_url: string;
  is_tip_enabled: boolean;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  qr_code_url: string;
  payment_url: string;
}

export interface QRPaymentRequest {
  amount?: number; // Required for DYNAMIC QR, ignored for STATIC
  medium: string; // Payment method
  phone_number: string; // Payer's phone
  tip_amount?: number; // Optional tip
  tipee_phone?: string; // Required if tip_amount > 0
  tip_medium?: string; // Required if tip_amount > 0
}

export interface QRPaymentResponse {
  success: boolean;
  status: string; // PENDING, SUCCESS, FAILED
  message: string;
  payment_amount: number;
  tip_amount: number;
  socialpay_transaction_id: string;
  tip_transaction_id?: string;
  payment_url?: string; // Optional payment URL for redirection
} 