'use client';

import { useTranslations } from 'next-intl';
import type { TranslationValues, RichTranslationValues } from 'next-intl';
import React from 'react';

// Define the shape of your translations
interface CheckoutTranslations {
  choosePayment: string;
  socialPay: string;
  wallet: string;
  bank: string;
  card: string;
  cardNumber: string;
  expiryDate: string;
  cvv: string;
  merchant: string;
  how_would_you_like_to_pay_today: string;
  // Add all other checkout keys here
  your_payment_is_encrypted_and_protected: string;
}

interface ErrorTranslations {
  not_found: string;
  checkout_status: string;
  user_error: string;
  server_error: string;
}

type CheckoutKeys = keyof CheckoutTranslations;
type ErrorKeys = keyof ErrorTranslations;

export function useCheckoutTranslations() {
  // Cast useTranslations to accept our specific keys
  const t = useTranslations() as {
    (key: string, values?: TranslationValues): string;
    rich(key: string, values?: RichTranslationValues): React.ReactNode;
  };

  return {
    // Checkout namespace translations
    checkout: (key: CheckoutKeys, values?: TranslationValues): string => {
      const fullKey = `checkout.${key}`;
      try {
        return t(fullKey, values);
      } catch {
        return String(key)
          .split('_')
          .map((word: string) => word.charAt(0).toUpperCase() + word.slice(1))
          .join(' ');
      }
    },
    
    // Error namespace translations
    error: (key: ErrorKeys, values?: TranslationValues): string => {
      const fullKey = `error.${key}`;
      try {
        return t(fullKey, values);
      } catch {
        switch (key) {
          case 'not_found': return 'Not found';
          case 'user_error': return 'An error occurred';
          case 'server_error': return 'Server error';
          default: return 'An error occurred';
        }
      }
    },
    
    // Rich text translations
    rich: (key: CheckoutKeys, values?: RichTranslationValues): React.ReactNode => {
      const fullKey = `checkout.${key}`;
      try {
        return t.rich(fullKey, values);
      } catch {
        const fallbackText = String(key)
          .split('_')
          .map((word: string) => word.charAt(0).toUpperCase() + word.slice(1))
          .join(' ');
        return React.createElement('span', null, fallbackText);
      }
    }
  };
}