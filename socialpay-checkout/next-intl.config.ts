// next-intl.config.ts

import { getRequestConfig } from 'next-intl/server';

// Named exports for locales and defaultLocale
export const locales = ['en', 'am'] as const;
export const defaultLocale = 'en';

export default getRequestConfig(async ({ locale }) => {
  try {
    const messages = (await import(`./src/messages/${locale}.json`)).default;

    return {
      messages,
      timeZone: 'Africa/Addis_Ababa',
    };
  } catch (error) {
    console.error(`Failed to load messages for locale: ${locale}`, error);
    return {
      messages: (await import('./src/messages/en.json')).default,
      timeZone: 'Africa/Addis_Ababa',
    };
  }
});
