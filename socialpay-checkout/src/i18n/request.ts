import { notFound } from 'next/navigation';
import { getRequestConfig } from 'next-intl/server';
import { locales } from './config';

export default getRequestConfig(async ({ locale }) => {
  // Validate locale
  if (!locales.includes(locale as any)) {
    notFound();
  }

  try {
    const messages = (await import(`@/messages/${locale}.json`)).default;
    return {
      messages,
      timeZone: 'Africa/Addis_Ababa',
      now: new Date()
    };
  } catch (error) {
    console.error(`Failed to load messages for locale ${locale}`, error);
    notFound();
  }
});