"use client"

import React, { useState, useEffect } from 'react';
import { NextIntlClientProvider } from 'next-intl';
import { defaultLocale } from '../../next-intl.config';
import { languageEvents } from '@/utils/languageEvents';
import type { AbstractIntlMessages } from 'next-intl';

type DynamicIntlProviderProps = {
  children: React.ReactNode;
  initialMessages: AbstractIntlMessages;
};

const DynamicIntlProvider: React.FC<DynamicIntlProviderProps> = ({ 
  children, 
  initialMessages 
}) => {
  const [locale, setLocale] = useState(defaultLocale);
  const [messages, setMessages] = useState<AbstractIntlMessages>(initialMessages);
  
  // Load initial locale from localStorage
  useEffect(() => {
    try {
      const savedLocale = localStorage.getItem('locale') || defaultLocale;
      if (savedLocale !== locale) {
        setLocale(savedLocale);
        loadMessages(savedLocale);
      }
    } catch (error) {
      console.error("Error accessing localStorage:", error);
    }
  }, [locale]);
  
  // Listen for language change events
  useEffect(() => {
    const unsubscribe = languageEvents.subscribe((newLocale) => {
      setLocale(newLocale);
      loadMessages(newLocale);
    });
    
    return () => {
      unsubscribe();
    };
  }, []);
  
  const loadMessages = async (localeToLoad: string) => {
    try {
      const messagesModule = await import(`../messages/${localeToLoad}.json`);
      setMessages(messagesModule.default);
      console.log(`Loaded messages for ${localeToLoad}`);
    } catch (error) {
      console.error(`Failed to load messages for locale: ${localeToLoad}`, error);
      // Fallback to English
      const fallbackModule = await import('../messages/en.json');
      setMessages(fallbackModule.default);
    }
  };
  
  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      {children}
    </NextIntlClientProvider>
  );
};

export default DynamicIntlProvider; 