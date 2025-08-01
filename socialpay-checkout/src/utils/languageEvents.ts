"use client"

type LanguageChangeListener = (locale: string) => void;

class LanguageEventEmitter {
  private listeners: LanguageChangeListener[] = [];

  subscribe(listener: LanguageChangeListener): () => void {
    this.listeners.push(listener);
    
    // Return unsubscribe function
    return () => {
      this.listeners = this.listeners.filter(l => l !== listener);
    };
  }

  emit(locale: string): void {
    this.listeners.forEach(listener => listener(locale));
  }
}

// Singleton instance
export const languageEvents = new LanguageEventEmitter(); 