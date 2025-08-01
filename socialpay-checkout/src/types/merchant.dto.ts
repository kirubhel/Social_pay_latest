export interface MerchantResponse {
  success: boolean;
  data: Merchant;
}

export interface Merchant {
  id: string;
  user_id: string;
  legal_name: string;
  trading_name: string;
  business_registration_number: string;
  tax_identification_number: string;
  business_type: string;
  industry_category: string;
  is_betting_company: boolean;
  lottery_certificate_number: string;
  website_url: string;
  established_date: string;
  created_at: string;
  updated_at: string;
  status: string;
}

export interface PaymentInitiationRequest {
  medium: string;
  phone_number: string;
  amount: number;
  merchant_id: string;
}

export interface PaymentInitiationResponse {
  socialpay_transaction_id: string;
  message: string;
  payment_url: string;
  reference_id: string;
  status: string;
  success: boolean;
}

export interface TransactionStatusResponse {
  amount: number;
  created_at: string;
  currency: string;
  description: string;
  details: string;
  fee_amount: number;
  id: string;
  medium: string;
  reference: string;
  status: string;
  total_amount: number;
  type: string;
  updated_at: string;
}

export interface ReceiptResponse {
  id: string;
  merchant_id: string;
  phone_number: string;
  user_id: string;
  type: string;
  medium: string;
  reference: string;
  comment: string;
  verified: boolean;
  details: unknown;
  created_at: string;
  updated_at: string;
  reference_number: string;
  test: boolean;
  status: string;
  description: string;
  token: string;
  amount: number;
  webhook_received: boolean;
  fee_amount: number;
  admin_net: number;
  vat_amount: number;
  merchant_net: number;
  total_amount: number;
  currency: string;
  callback_url: string;
  success_url: string;
  failed_url: string;
  merchant: Merchant;
}

export interface CheckoutDetailsResponse {
  id: string;
  amount: number;
  currency: string;
  description: string;
  reference: string;
  supported_mediums: string[];
  phone_number: string;
  success_url: string;
  failed_url: string;
  status: string;
  created_at: string;
  expires_at: string;
  merchant: Merchant;
  accept_tip: boolean;
}

export interface CheckoutPaymentRequest {
  hosted_checkout_id: string;
  medium: string;
  phone_number: string;
  tip_amount?: number;
  tipee_phone?: string;
  tip_medium?: string;
}

export interface CheckoutPaymentResponse {
  success: boolean;
  status: string;
  message: string;
  payment_url: string;
  reference_id: string;
  socialpay_transaction_id: string;
} 