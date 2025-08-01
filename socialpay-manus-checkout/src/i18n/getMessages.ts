import { locales } from './config';

export async function getMessages(locale: string) {
  const safeLocale = locales.includes(locale as any) ? locale : 'en';
  
  try {
    // Using absolute path with @/ alias
    const messages = (await import(`@/messages/${safeLocale}.json`)).default;
    return {
      messages,
      timeZone: 'Africa/Addis_Ababa',
      now: new Date()
    };
  } catch (error) {
    console.error(`Failed to load messages for locale: ${safeLocale}`, error);
    // Fallback to English
    const fallbackMessages = (await import('@/messages/en.json')).default;
    return {
      messages: fallbackMessages,
      timeZone: 'Africa/Addis_Ababa',
      now: new Date()
    };
  }
}