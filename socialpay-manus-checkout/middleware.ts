import createMiddleware from 'next-intl/middleware';
import { locales, defaultLocale } from './next-intl.config';

export default createMiddleware({
  // A list of all locales that are supported
  locales,
  
  // The default locale to use when visiting a non-localized route
  defaultLocale,
  
  // This will make sure the locale is not included in the URL
  localePrefix: 'never',
  
  // This will make the middleware not redirect based on locale
  localeDetection: false
});

export const config = {
  // Skip all paths that should not be internationalized
  matcher: ['/((?!api|_next|.*\\..*).*)']
}; 