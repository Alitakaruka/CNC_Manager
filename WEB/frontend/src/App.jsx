import MainPage from './components/Pages/MainPage'
import { LocalizationProvider } from './hooks/useLocalization.jsx'
import { Toaster } from 'react-hot-toast'

function App() {
  return (
    <LocalizationProvider>
      <div className="App min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100 dark:from-gray-900 dark:via-gray-900 dark:to-gray-800">
        <MainPage />
        <Toaster 
          position="top-right"
          toastOptions={{
            duration: 4000,
            style: {
              background: '#363636',
              color: '#fff',
            },
            success: {
              duration: 3000,
              iconTheme: {
                primary: '#22c55e',
                secondary: '#fff',
              },
            },
            error: {
              duration: 5000,
              iconTheme: {
                primary: '#ef4444',
                secondary: '#fff',
              },
            },
          }}
        />
      </div>
    </LocalizationProvider>
  )
}

export default App
