import { useState, useRef } from 'react'
import { motion } from 'framer-motion'
import { Wifi, Usb, Monitor, CheckCircle, AlertCircle } from 'lucide-react'
import ConnectionTypes from './TypesOfConnections'
import { ConnectionBuilder } from './TypesOfConnections'
import ConnectButton from '../Buttons/ConnectButton'
import ConnectionHook from '../../hooks/ConnectionHook'
import { useLocalization } from '../../hooks/useLocalization.jsx'
import toast from 'react-hot-toast'

const connectionIcons = {
  COM: Monitor,
  IP: Wifi,
  USB: Usb
}

export default function ConnectionSelector() {
  const [selectedConnection, setSelectedConnection] = useState(ConnectionTypes.COM)
  const [isConnecting, setIsConnecting] = useState(false)
  const [connectionData, setConnectionData] = useState('')
  const dataRef = useRef({ connection: '' })

  // Используем локализацию
  const { t } = useLocalization()

  const handleConnectionChange = (event) => {
    setSelectedConnection(event.target.value)
    setConnectionData('')
    dataRef.current.connection = ''
  }

  const handleDataChange = (value) => {
    setConnectionData(value)
    dataRef.current.connection = value
  }

  const handleConnect = async () => {
    if (!connectionData.trim()) {
      toast.error(t('connections.errors.emptyData'))
      return
    }

    setIsConnecting(true)
    
    try {
      await ConnectionHook(selectedConnection, connectionData)
      toast.success(t('connections.errors.success'))
    } catch (error) {
      toast.error(`${t('connections.errors.connectionFailed')}: ${error.message}`)
    } finally {
      setIsConnecting(false)
    }
  }

  const IconComponent = connectionIcons[selectedConnection]

  return (
    <motion.div 
      className="w-full max-w-2xl mx-auto space-y-8"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
    >
      {/* Header */}
      <div className="text-center">
        <motion.div
          className="inline-flex items-center justify-center w-20 h-20 bg-primary-100 rounded-full mb-6"
          whileHover={{ scale: 1.1 }}
          transition={{ type: "spring", stiffness: 400, damping: 10 }}
        >
          <IconComponent className="h-10 w-10 text-primary-600" />
        </motion.div>
        <h3 className="text-3xl font-bold text-gray-900 mb-3">{t('connections.title')}</h3>
        <p className="text-gray-600 text-lg">{t('connections.subtitle')}</p>
      </div>

      {/* Connection Type Selector */}
      <div className="space-y-4">
        <label className="block text-lg font-medium text-gray-700">
          {t('connections.type')}
        </label>
        <div className="grid grid-cols-3 gap-4">
          {Object.entries(ConnectionTypes).map(([key, value]) => {
            const Icon = connectionIcons[key]
            const isSelected = selectedConnection === value
            
            return (
              <motion.button
                key={key}
                onClick={() => handleConnectionChange({ target: { value } })}
                className={`p-6 rounded-xl border-2 transition-all duration-200 ${
                  isSelected
                    ? 'border-primary-500 bg-primary-50 text-primary-700 shadow-lg'
                    : 'border-gray-200 bg-white text-gray-600 hover:border-primary-300 hover:bg-primary-50'
                }`}
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <div className="flex flex-col items-center space-y-3">
                  <Icon className={`h-8 w-8 ${isSelected ? 'text-primary-600' : 'text-gray-400'}`} />
                  <span className="text-base font-medium">{value}</span>
                </div>
              </motion.button>
            )
          })}
        </div>
      </div>

      {/* Connection Data Input */}
      <div className="space-y-4">
        <label className="block text-lg font-medium text-gray-700">
          {t('connections.data')}
        </label>
        <ConnectionBuilder 
          BuilderType={selectedConnection} 
          ConnectionRef={dataRef}
          value={connectionData}
          onChange={handleDataChange}
        />
      </div>

      {/* Connection Status */}
      {connectionData && (
        <motion.div
          initial={{ opacity: 0, height: 0 }}
          animate={{ opacity: 1, height: 'auto' }}
          className="flex items-center space-x-3 p-4 bg-success-50 border border-success-200 rounded-lg"
        >
          <CheckCircle className="h-6 w-6 text-success-600" />
          <span className="text-base text-success-700">
            {t('connections.ready')}: {selectedConnection} - {connectionData}
          </span>
        </motion.div>
      )}

      {/* Connect Button */}
      <div className="pt-6">
        <ConnectButton 
          onClick={handleConnect}
          disabled={!connectionData.trim() || isConnecting}
          loading={isConnecting}
        />
      </div>

      {/* Help Text */}
      <div className="text-sm text-gray-500 text-center space-y-2 bg-gray-50 p-4 rounded-lg">
        <p className="font-medium text-gray-700 mb-2">{t('connections.help.title')}</p>
        <p>{t('connections.help.com')}</p>
        <p>{t('connections.help.ip')}</p>
        <p>{t('connections.help.usb')}</p>
      </div>
    </motion.div>
  )
}