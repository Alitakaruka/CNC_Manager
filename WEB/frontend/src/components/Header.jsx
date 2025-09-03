import React, { useState, useEffect, useRef } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  Settings, 
  Menu, 
  X, 
  User, 
  Palette, 
  LogOut, 
  Home, 
  BarChart3,
  Plus,
  Wifi,
  WifiOff,
  Clock,
  Activity,
  AlertCircle,
  CheckCircle,
  Info,
  ChevronDown,
  ChevronUp,
  RefreshCw
} from 'lucide-react'
import NewCon from './Buttons/NewConnectionButton'
import LanguageSwitcher from './LanguageSwitcher'
import useSystemData from '../hooks/SystemHook'
import { useLocalization } from '../hooks/useLocalization.jsx'
import CncLogo from './CncLogo'
import wsClient from '../hooks/WebSocketClient'
import toast from 'react-hot-toast'

export default function Header() {
  const [navOpen, setNavOpen] = useState(false)
  const [settingsOpen, setSettingsOpen] = useState(false)
  const [logsOpen, setLogsOpen] = useState(false)
  const [transport, setTransport] = useState(wsClient.transport)
  const [wsUrl, setWsUrl] = useState(wsClient.url)
  const logsEndRef = useRef(null)
  
  const { systemData, isLoading, updateAllData } = useSystemData()
  const { connectionStatus, systemInfo, logs } = systemData
  const { t } = useLocalization()
  const isConnected = wsClient.isConnected

  useEffect(() => {
    const onOpen = () => { setTransport(wsClient.transport); setWsUrl(wsClient.url) }
    const onClose = () => { setTransport(wsClient.transport) }
    const offOpen = wsClient.on('open', onOpen)
    const offClose = wsClient.on('close', onClose)
    onOpen()
    return () => { offOpen(); offClose() }
  }, [])

  useEffect(() => {
    if (logsOpen && logsEndRef.current) logsEndRef.current.scrollIntoView({ behavior: 'smooth' })
  }, [logs, logsOpen])

  const navItems = [
    { icon: Home, label: t('navigation.home'), active: true },
    { icon: CncLogo, label: t('navigation.printers') },
    { icon: BarChart3, label: t('navigation.reports') }
  ]

  const settingsItems = [ { icon: Palette, label: t('common.theme') } ]

  const getLogIcon = (level) => {
    switch (level) {
      case 'error': return <AlertCircle className="h-3 w-3 text-danger-500" />
      case 'warning': return <AlertCircle className="h-3 w-3 text-warning-500" />
      case 'success': return <CheckCircle className="h-3 w-3 text-success-500" />
      default: return <Info className="h-3 w-3 text-primary-500" />
    }
  }

  const formatTime = (date) => {
    const d = date instanceof Date ? date : new Date(date)
    if (Number.isNaN(d.getTime())) return '--:--'
    return d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
  }

  return (
    <header className="sticky top-2 z-50 glass-effect border-b border-white/20 shadow-lg dark:bg-gray-900/90 dark:border-gray-700">
      <div className="w-full pl-0 pr-4 sm:pr-6 lg:pr-8">
        <div className="flex items-center justify-between h-20">
          {/* Logo + status */}
          <div className="flex items-center space-x-3">
            <div className="relative">
              <CncLogo className="h-8 w-8 text-primary-600" />
              <motion.div 
                className={`absolute -top-1 -right-1 h-3 w-3 rounded-full ${isConnected ? 'bg-success-500' : 'bg-danger-500'}`}
                animate={{ scale: [1, 1.2, 1] }}
                transition={{ duration: 2, repeat: Infinity }}
                title={`${isConnected ? 'Online' : 'Offline'} • ${transport}${wsClient.url ? ` • ${wsUrl}` : ''}`}
              />
            </div>
            <div>
              <h1 className="text-xl font-bold text-gradient">{t('header.title')}</h1>
              <div className="text-xs font-semibold text-primary-600 tracking-wide" title={`Channel: ${transport}${wsUrl ? `\nURL: ${wsUrl}` : ''}`}>AliPri</div>
              <div className="flex items-center space-x-2 text-xs text-gray-500">
                {isConnected ? (
                  <>
                    <Wifi className="h-3 w-3 text-success-500" />
                    <span>{t('header.connected')} ({connectionStatus.online}/{connectionStatus.total})</span>
                  </>
                ) : (
                  <>
                    <WifiOff className="h-3 w-3 text-danger-500" />
                    <span>{t('header.disconnected')}</span>
                  </>
                )}
              </div>
            </div>
          </div>

          {/* Navigation */}
          <div className="hidden md:flex items-center space-x-4">
            {navItems.map((item) => (
              <button key={item.label}
                className={`flex items-center space-x-2 px-5 py-3 rounded-lg font-medium transition-all duration-200 ${item.active ? 'bg-primary-100 text-primary-700 shadow-md' : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900'}`}
              >
                <item.icon className="h-4 w-4" />
                <span>{item.label}</span>
              </button>
            ))}
          </div>

          {/* Right side */}
          <div className="flex items-center space-x-4">
            <div className="hidden lg:flex items-center space-x-4 text-xs text-gray-600">
              <div className="flex items-center space-x-1">
                <Activity className="h-3 w-3 text-success-500" />
                <span>{t('header.online')}: {connectionStatus.online}</span>
              </div>
              <div className="flex items-center space-x-1">
                <div className="w-2 h-2 bg-warning-500 rounded-full animate-pulse" />
                <span>{t('header.printing')}: {connectionStatus.printing}</span>
              </div>
              <div className="flex items-center space-x-1">
                <Clock className="h-3 w-3 text-gray-400" />
                <span>{systemInfo.uptime}</span>
              </div>
            </div>

            <LanguageSwitcher />

            {/* Live Logs */}
            <div className="relative">
              <motion.button
                onClick={() => setLogsOpen(!logsOpen)}
                className="flex items-center space-x-2 px-3 py-2 rounded-lg text-gray-600 hover:bg-gray-100 hover:text-gray-900 transition-all duration-200"
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                title={`Channel: ${transport}${wsUrl ? `\nURL: ${wsUrl}` : ''}`}
              >
                <div className="relative">
                  <Activity className="h-4 w-4" />
                  <motion.div 
                    className="absolute -top-1 -right-1 h-2 w-2 bg-success-500 rounded-full"
                    animate={{ scale: [1, 1.5, 1] }}
                    transition={{ duration: 1, repeat: Infinity }}
                  />
                </div>
                <span className="hidden sm:inline">{t('header.logs')}</span>
                {logsOpen ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
              </motion.button>

              <AnimatePresence>
                {logsOpen && (
                  <motion.div
                    initial={{ opacity: 0, y: -10, scale: 0.95 }}
                    animate={{ opacity: 1, y: 0, scale: 1 }}
                    exit={{ opacity: 0, y: -10, scale: 0.95 }}
                    className="absolute right-0 mt-2 w-80 glass-effect rounded-xl shadow-xl border border-white/20 max-h-96 overflow-hidden"
                  >
                    <div className="flex items-center justify-between p-3 border-b border-white/20">
                      <h3 className="text-sm font-semibold text-gray-900">{t('header.systemLogs')}</h3>
                      <motion.button
                        onClick={updateAllData}
                        className="p-1 text-gray-400 hover:text-gray-600 rounded"
                        whileHover={{ scale: 1.1 }}
                        whileTap={{ scale: 0.9 }}
                        title={`Channel: ${transport}${wsUrl ? `\nURL: ${wsUrl}` : ''}`}
                      >
                        <RefreshCw className={`h-3 w-3 ${isLoading ? 'animate-spin' : ''}`} />
                      </motion.button>
                    </div>
                    <div className="max-h-64 overflow-y-auto p-2">
                      <AnimatePresence>
                        {logs.map((log, index) => (
                          <motion.div key={log.id}
                            initial={{ opacity: 0, x: -20 }}
                            animate={{ opacity: 1, x: 0 }}
                            exit={{ opacity: 0, x: 20 }}
                            transition={{ delay: index * 0.05 }}
                            className="flex items-start space-x-2 p-2 rounded-lg hover:bg-white/50 transition-colors duration-200"
                          >
                            {getLogIcon(log.level)}
                            <div className="flex-1 min-w-0">
                              <p className="text-xs text-gray-900 truncate">{log.message}</p>
                              <div className="flex items-center justify-between mt-1">
                                <span className="text-xs text-gray-500">{formatTime(log.timestamp)}</span>
                                <span className="text-xs text-gray-400 bg-gray-100 px-1 rounded">{log.type}</span>
                              </div>
                            </div>
                          </motion.div>
                        ))}
                      </AnimatePresence>
                      <div ref={logsEndRef} />
                    </div>
                    <div className="p-2 border-t border-white/20 bg-gray-50/50">
                      <div className="flex items-center justify-between text-xs text-gray-500">
                        <span>{t('header.totalLogs')}: {logs.length}</span>
                        <span>{t('header.updated')}: {formatTime(systemInfo.lastUpdate)}</span>
                      </div>
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>

            {/* Settings */}
            <div className="relative">
              <button onClick={() => setSettingsOpen(!settingsOpen)}
                className="flex items-center space-x-2 px-3 py-2 rounded-lg text-gray-600 hover:bg-gray-100 hover:text-gray-900 transition-all duration-200"
                title={`Channel: ${transport}${wsUrl ? `\nURL: ${wsUrl}` : ''}`}
              >
                <Settings className="h-5 w-5" />
                <span className="hidden sm:inline">{t('common.settings')}</span>
              </button>
              <AnimatePresence>
                {settingsOpen && (
                  <motion.div initial={{ opacity: 0, y: -10, scale: 0.95 }}
                    animate={{ opacity: 1, y: 0, scale: 1 }}
                    exit={{ opacity: 0, y: -10, scale: 0.95 }}
                    className="absolute right-0 mt-2 w-48 glass-effect rounded-xl shadow-xl border border-white/20"
                  >
                    {settingsItems.map((item, index) => (
                      <motion.button key={item.label}
                        className="flex items-center space-x-3 w-full px-4 py-3 text-left text-gray-700 hover:bg-primary-50 hover:text-primary-700 transition-colors duration-200 first:rounded-t-xl last:rounded-b-xl"
                        whileHover={{ x: 5 }} initial={{ opacity: 0, x: -10 }} animate={{ opacity: 1, x: 0 }} transition={{ delay: index * 0.1 }}
                      >
                        <item.icon className="h-4 w-4" />
                        <span>{item.label}</span>
                      </motion.button>
                    ))}
                  </motion.div>
                )}
              </AnimatePresence>
            </div>

            <NewCon />
            <button className="md:hidden p-2 rounded-lg text-gray-600 hover:bg-gray-100 hover:text-gray-900" onClick={() => setNavOpen(!navOpen)}>
              {navOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
            </button>
          </div>
        </div>

        {/* Mobile Navigation */}
        <AnimatePresence>
          {navOpen && (
            <motion.div initial={{ opacity: 0, height: 0 }} animate={{ opacity: 1, height: 'auto' }} exit={{ opacity: 0, height: 0 }} className="md:hidden py-4 border-t border-white/20">
              <div className="flex flex-col space-y-2">
                {navItems.map((item, index) => (
                  <motion.button key={item.label}
                    className={`flex items-center space-x-3 px-4 py-3 rounded-lg font-medium transition-all duration-200 ${item.active ? 'bg-primary-100 text-primary-700' : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900'}`}
                    initial={{ opacity: 0, x: -20 }} animate={{ opacity: 1, x: 0 }} transition={{ delay: index * 0.1 }}
                  >
                    <item.icon className="h-5 w-5" />
                    <span>{item.label}</span>
                  </motion.button>
                ))}
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </header>
  )
}
