declare module '@/messages/*.json' {
  import type { AbstractIntlMessages } from 'next-intl';
  
  const messages: {
    checkout: {
      // All your checkout keys
      [key: string]: string;
    };
    language: {
      english: string;
      amharic: string;
    };
    loading: string;
    close: string;
    reference: string;
    secured_by_socialpay: string;
    error: {
      not_found: string;
      checkout_status: string;
      user_error: string;
      server_error: string;
    };
  } & AbstractIntlMessages;
  
  export default messages;
}