import type { Metadata } from "next";
import "./globals.css";
import { inter } from './fonts';
import { LanguageProvider } from '../components/LanguageProvider';
import { defaultLocale } from '../../next-intl.config';
import { getMessages } from "@/i18n/getMessages";
import DynamicIntlProvider from "@/components/DynamicIntlProvider";


export const metadata: Metadata = {
  title: "SocialPay - Connect. Accept. Grow.",
  description: "Simple, secure payment processing for Ethiopian businesses. Connect. Accept. Grow.",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  // For server rendering, use the default locale
  const locale = defaultLocale;
  
  // Get messages for the current locale
  const messages = await getMessages(locale);

  return (
    <html lang={locale} className={inter.variable}>
      <body className={`${inter.className} transition-opacity duration-300`}>
        <DynamicIntlProvider initialMessages={messages}>
          <LanguageProvider initialLocale={locale}>
            {children}
          </LanguageProvider>
        </DynamicIntlProvider>
      </body>
    </html>
  );
}
