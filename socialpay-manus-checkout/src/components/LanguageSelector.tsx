"use client"

import React, { useEffect, useState } from 'react';
import { useLanguage } from './LanguageProvider';
import Image from 'next/image';
import { languageEvents } from '@/utils/languageEvents';

const LanguageSelector: React.FC = () => {
  const { locale, setLocale } = useLanguage();
  const [currentLocale, setCurrentLocale] = useState(locale);
  const [showLanguageDropdown, setShowLanguageDropdown] = useState(false);

  // Listen for language changes from other components
  useEffect(() => {
    const unsubscribe = languageEvents.subscribe((newLocale) => {
      setCurrentLocale(newLocale);
    });
    
    return () => {
      unsubscribe();
    };
  }, []);

  const toggleLanguageDropdown = () => {
    setShowLanguageDropdown(!showLanguageDropdown);
  };

  const changeLanguage = (lang: string) => {
    setLocale(lang.toLowerCase());
    setCurrentLocale(lang.toLowerCase());
    setShowLanguageDropdown(false);
  };

  return (
    <div className="relative">
      <button 
        onClick={toggleLanguageDropdown}
        className="flex flex-row items-center justify-center gap-1 px-2 py-1 border border-[#BBBBBB] rounded-[11px]"
      >
        <span className="text-[#707070] font-medium text-sm">
          {currentLocale === 'en' ? 'EN' : 'AM'}
        </span>
        <Image src="/lang-selector-downarrow.svg" alt="Arrow down" width={19} height={19} />
      </button>
      
      {showLanguageDropdown && (
        <div className="absolute top-full right-0 mt-1 bg-white border border-[#BBBBBB] rounded-[11px] hover:rounded-[11px] shadow-md z-10">
          <button 
            className="w-full text-left px-2 py-1 hover:bg-gray-50 text-[#707070] flex items-center justify-center gap-1"
            onClick={() => changeLanguage('en')}
          >
            <Image src="/usa-flag.svg" alt="USA flag" width={17} height={13} />
            EN
          </button>
          <button 
            className="w-full text-left px-2 py-1 hover:bg-gray-50 hover:rounded-[11px] text-[#707070] flex items-center justify-center gap-1"
            onClick={() => changeLanguage('am')}
          >
            <Image src="/ethiopia-flag.svg" alt="Ethiopian flag" width={17} height={13} />
            AM
          </button>
        </div>
      )}
    </div>
  );
};

export default LanguageSelector; 