import { motion } from 'framer-motion'
import { Wifi, Loader2, CheckCircle } from 'lucide-react'
import { useLocalization } from '../../hooks/useLocalization.jsx'

export default function ConnectButton({ onClick, disabled = false, loading = false }) {
  const { t } = useLocalization()

  return (
    <motion.button
      onClick={onClick}
      disabled={disabled || loading}
      className={`w-full py-3 px-6 rounded-xl font-semibold text-lg transition-all duration-200 flex items-center justify-center space-x-3 ${
        disabled || loading
          ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
          : 'btn-primary shadow-glow'
      }`}
      whileHover={!disabled && !loading ? { scale: 1.02, boxShadow: "0 0 30px rgba(59, 130, 246, 0.4)" } : {}}
      whileTap={!disabled && !loading ? { scale: 0.98 } : {}}
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.2 }}
    >
      {loading ? (
        <>
          <Loader2 className="h-5 w-5 animate-spin" />
          <span>{t('connections.connecting')}</span>
        </>
      ) : (
        <>
          <Wifi className="h-5 w-5" />
          <span>{t('connections.connect')}</span>
        </>
      )}
    </motion.button>
  )
}
