"use client"

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { locales, defaultLocale } from '../../next-intl.config';
import { languageEvents } from '@/utils/languageEvents';

type LanguageContextType = {
  locale: string;
  setLocale: (locale: string) => void;
};

const LanguageContext = createContext<LanguageContextType>({
  locale: defaultLocale,
  setLocale: () => {},
});

export const useLanguage = () => useContext(LanguageContext);

type LanguageProviderProps = {
  children: ReactNode;
  initialLocale?: string;
};

export const LanguageProvider: React.FC<LanguageProviderProps> = ({ 
  children, 
  initialLocale = defaultLocale 
}) => {
  const [locale, setLocaleState] = useState(initialLocale);

  // Load saved locale from localStorage on component mount
  useEffect(() => {
    try {
      const savedLocale = localStorage.getItem('locale');
      if (savedLocale && locales.includes(savedLocale)) {
        setLocaleState(savedLocale);
        // Emit the event to notify other components
        languageEvents.emit(savedLocale);
      }
    } catch (error) {
      console.error('Error accessing localStorage:', error);
    }
  }, []);

  const setLocale = (newLocale: string) => {
    if (locales.includes(newLocale)) {
      try {
        // Save to localStorage
        localStorage.setItem('locale', newLocale);
        
        // Update state
        setLocaleState(newLocale);
        
        // Emit the event to notify other components
        languageEvents.emit(newLocale);
        
      } catch (error) {
        console.error('Error setting locale:', error);
      }
    }
  };

  return (
    <LanguageContext.Provider value={{ locale, setLocale }}>
      {children}
    </LanguageContext.Provider>
  );
};

export default LanguageProvider; 