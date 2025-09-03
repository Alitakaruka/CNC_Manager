import { useState, useEffect, createContext, useContext } from 'react'
import ru from '../locales/ru'
import en from '../locales/en'

// Создаем контекст для локализации
const LocalizationContext = createContext()

// Доступные языки
export const languages = {
  ru: {
    name: 'Русский',
    flag: '🇷🇺',
    translations: ru
  },
  en: {
    name: 'English',
    flag: '🇺🇸',
    translations: en
  }
}

// Хук для использования локализации
export const useLocalization = () => {
  const context = useContext(LocalizationContext)
  if (!context) {
    throw new Error('useLocalization must be used within a LocalizationProvider')
  }
  return context
}

// Провайдер локализации
export const LocalizationProvider = ({ children }) => {
  const [currentLanguage, setCurrentLanguage] = useState('ru')
  const [translations, setTranslations] = useState(languages.ru.translations)

  // Функция для смены языка
  const changeLanguage = (lang) => {
    if (languages[lang]) {
      setCurrentLanguage(lang)
      setTranslations(languages[lang].translations)
      localStorage.setItem('language', lang)
    }
  }

  // Функция для получения перевода
  const t = (key, params = {}) => {
    const keys = key.split('.')
    let value = translations

    for (const k of keys) {
      if (value && typeof value === 'object' && k in value) {
        value = value[k]
      } else {
        console.warn(`Translation key not found: ${key}`)
        return key
      }
    }

    if (typeof value !== 'string') {
      console.warn(`Translation value is not a string: ${key}`)
      return key
    }

    // Заменяем параметры в строке
    return value.replace(/\{(\w+)\}/g, (match, param) => {
      return params[param] !== undefined ? params[param] : match
    })
  }

  // Загружаем сохраненный язык при инициализации
  useEffect(() => {
    const savedLanguage = localStorage.getItem('language')
    if (savedLanguage && languages[savedLanguage]) {
      setCurrentLanguage(savedLanguage)
      setTranslations(languages[savedLanguage].translations)
    }
  }, [])

  const value = {
    currentLanguage,
    changeLanguage,
    t,
    languages,
    translations
  }

  return (
    <LocalizationContext.Provider value={value}>
      {children}
    </LocalizationContext.Provider>
  )
}

export default useLocalization
