// src/i18n/config.ts
import enMessages from '@/messages/en.json';

type CheckoutKeys = keyof typeof enMessages.checkout;

export const locales = ['en', 'am'] as const;
export const defaultLocale = 'en';
export const localePrefix = 'as-needed';

export type Locale = typeof locales[number];
export type TranslationKey = CheckoutKeys; // Now fully type-safe