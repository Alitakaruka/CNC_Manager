import { useState, useEffect, useRef } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  FileText, 
  AlertCircle, 
  CheckCircle, 
  Info, 
  Clock,
  Trash2,
  Download,
  RefreshCw,
  Activity,
  Filter,
  Search
} from 'lucide-react'
import toast from 'react-hot-toast'
import { useLocalization } from '../../../hooks/useLocalization.jsx'
import wsClient from '../../../hooks/WebSocketClient'

export default function PrinterLogs() {
  const [logs, setLogs] = useState([])
  const [isLoading, setIsLoading] = useState(false)
  const [autoScroll, setAutoScroll] = useState(false) // Изменили на false по умолчанию
  const [filter, setFilter] = useState('all')
  const [searchTerm, setSearchTerm] = useState('')
  const logsEndRef = useRef(null)
  const startedRef = useRef(false)
  const { t } = useLocalization()

  // Подписка на WebSocket события для получения логов
  useEffect(() => {
    if (!startedRef.current) {
      startedRef.current = true
      wsClient.connect()
    }

    // Подписка на событие 'log'
    const offLog = wsClient.on('log', (logData) => {
      console.log('Received log:', logData) // Для отладки
      
      if (logData && typeof logData === 'object') {
        // Преобразуем timestamp из uint32 (секунды) в Date объект
        // uint32 timestamp всегда в секундах, умножаем на 1000 для миллисекунд
        const timestampDate = logData.timestamp 
          ? new Date(logData.timestamp * 1000) 
          : new Date()

        // Преобразуем данные из бэкенда в формат, ожидаемый фронтендом
        const newLog = {
          id: logData.id || Date.now() + Math.random(), // Используем id из данных или генерируем уникальный
          timestamp: timestampDate,
          level: (logData.level || 'info').toLowerCase(),
          message: logData.message || '',
          printer: logData.printer || logData.printerName || 'System',
          type: logData.type || 'system'
        }

        setLogs(prev => {
          // Ограничиваем количество логов (например, последние 500)
          const updated = [...prev, newLog]
          return updated.slice(-500)
        })
      }
    })

    return () => {
      offLog()
    }
  }, [])

  useEffect(() => {
    if (autoScroll && logsEndRef.current) {
      logsEndRef.current.scrollIntoView({ behavior: 'smooth' })
    }
  }, [logs, autoScroll])

  const loadLogs = async () => {
    setIsLoading(true)
    try {
      // В реальном приложении здесь был бы API запрос
      // const response = await fetch('/api/logs')
      // const data = await response.json()
      
      // Просто возвращаем текущие логи без добавления новых
      setTimeout(() => {
        setIsLoading(false)
        toast.success(t('logs.logsUpdated'))
      }, 500)
    } catch (error) {
      console.error('Error loading logs:', error)
      toast.error(t('logs.loadError'))
      setIsLoading(false)
    }
  }

  const clearLogs = () => {
    setLogs([])
    toast.success(t('logs.logsCleared'))
  }

  const exportLogs = () => {
    const logText = logs.map(log => 
      `${log.timestamp.toISOString()} [${log.level.toUpperCase()}] ${log.printer}: ${log.message}`
    ).join('\n')
    
    const blob = new Blob([logText], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `printer-logs-${new Date().toISOString().split('T')[0]}.txt`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    
    toast.success(t('logs.logsExported'))
  }

  const getLogIcon = (level) => {
    switch (level) {
      case 'error':
        return <AlertCircle className="h-4 w-4 text-danger-500" />
      case 'warning':
        return <AlertCircle className="h-4 w-4 text-warning-500" />
      case 'success':
        return <CheckCircle className="h-4 w-4 text-success-500" />
      default:
        return <Info className="h-4 w-4 text-primary-500" />
    }
  }

  const getLogColor = (level) => {
    switch (level) {
      case 'error':
        return 'border-l-danger-500 bg-danger-50'
      case 'warning':
        return 'border-l-warning-500 bg-warning-50'
      case 'success':
        return 'border-l-success-500 bg-success-50'
      default:
        return 'border-l-primary-500 bg-primary-50'
    }
  }

  const formatTime = (date) => {
    return date.toLocaleTimeString('ru-RU', { 
      hour: '2-digit', 
      minute: '2-digit',
      second: '2-digit'
    })
  }

  // Фильтрация логов
  const filteredLogs = logs.filter(log => {
    const matchesFilter = filter === 'all' || log.level === filter
    const matchesSearch = searchTerm === '' || 
      log.message.toLowerCase().includes(searchTerm.toLowerCase()) ||
      log.printer.toLowerCase().includes(searchTerm.toLowerCase())
    return matchesFilter && matchesSearch
  })

  const getFilterCount = (level) => {
    return logs.filter(log => level === 'all' || log.level === level).length
  }

  const containerClass = `card h-fit dark:bg-gray-800/70 dark:border-gray-700`

  return (
    <motion.div
      className={containerClass}
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
    >
      {/* Header */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center space-x-2">
          <div className="relative">
            <FileText className="h-5 w-5 text-primary-600" />
            {/* Убрали анимированную точку */}
          </div>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">{t('logs.title')}</h3>
          <span className="bg-primary-100 text-primary-800 text-xs font-medium px-2 py-1 rounded-full">
            {logs.length}
          </span>
        </div>
        
        <div className="flex items-center space-x-2">
          <motion.button
            onClick={loadLogs}
            disabled={isLoading}
            className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors duration-200"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
          </motion.button>
          
          <motion.button
            onClick={exportLogs}
            className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors duration-200"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <Download className="h-4 w-4" />
          </motion.button>
          
          <motion.button
            onClick={clearLogs}
            className="p-2 text-gray-400 hover:text-danger-600 hover:bg-danger-50 rounded-lg transition-colors duration-200"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <Trash2 className="h-4 w-4" />
          </motion.button>
        </div>
      </div>

      {/* Search and Filter */}
      <div className="space-y-3 mb-4">
        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
          <input
            type="text"
            placeholder={t('logs.search')}
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all duration-200 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-100"
          />
        </div>

        {/* Filter Buttons */}
        <div className="flex items-center space-x-2">
          <Filter className="h-4 w-4 text-gray-400" />
          <div className="flex flex-wrap gap-1">
            {[
              { key: 'all', label: t('logs.filters.all'), color: 'bg-gray-500' },
              { key: 'info', label: t('logs.filters.info'), color: 'bg-primary-500' },
              { key: 'success', label: t('logs.filters.success'), color: 'bg-success-500' },
              { key: 'warning', label: t('logs.filters.warning'), color: 'bg-warning-500' },
              { key: 'error', label: t('logs.filters.error'), color: 'bg-danger-500' }
            ].map(({ key, label, color }) => (
              <motion.button
                key={key}
                onClick={() => setFilter(key)}
                className={`px-3 py-1 rounded-full text-xs font-medium transition-all duration-200 ${
                  filter === key
                    ? `${color} text-white shadow-md`
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                {label} ({getFilterCount(key)})
              </motion.button>
            ))}
          </div>
        </div>
      </div>

      {/* Auto-scroll toggle */}
      <div className="flex items-center space-x-2 mb-4">
        <input
          type="checkbox"
          id="autoScroll"
          checked={autoScroll}
          onChange={(e) => setAutoScroll(e.target.checked)}
          className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
        />
        <label htmlFor="autoScroll" className="text-sm text-gray-600">
          {t('logs.autoScroll')}
        </label>
        {/* Убрали индикатор Live */}
      </div>

      {/* Logs Container */}
      <div className="space-y-2 max-h-96 overflow-y-auto">
        <AnimatePresence>
          {filteredLogs.map((log, index) => (
            <motion.div
              key={log.id}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: 20 }}
              transition={{ delay: index * 0.02 }}
              className={`p-3 rounded-lg border-l-4 ${getLogColor(log.level)} hover:shadow-md transition-shadow duration-200 dark:bg-gray-800/60`}
            >
              <div className="flex items-start space-x-3">
                {getLogIcon(log.level)}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between">
                    <p className="text-sm font-medium text-gray-900 break-words">
                      {log.message}
                    </p>
                    <div className="flex items-center space-x-2 text-xs text-gray-500">
                      <Clock className="h-3 w-3" />
                      <span>{formatTime(log.timestamp)}</span>
                    </div>
                  </div>
                  <div className="flex items-center justify-between mt-1">
                    <p className="text-xs text-gray-600">
                      {log.printer}
                    </p>
                    <span className="text-xs text-gray-400 bg-gray-100 px-2 py-1 rounded-full">
                      {log.type}
                    </span>
                  </div>
                </div>
              </div>
            </motion.div>
          ))}
        </AnimatePresence>
        
        {filteredLogs.length === 0 && !isLoading && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="text-center py-8 text-gray-500"
          >
            <FileText className="h-12 w-12 mx-auto mb-4 text-gray-300" />
            <p>{t('logs.noLogs')}</p>
            {searchTerm && <p className="text-xs mt-2">{t('logs.changeSearch')}</p>}
          </motion.div>
        )}
        
        <div ref={logsEndRef} />
      </div>

      {/* Loading State */}
      {isLoading && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="flex items-center justify-center py-4"
        >
          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary-600"></div>
          <span className="ml-2 text-sm text-gray-600">{t('logs.updateLogs')}</span>
        </motion.div>
      )}

      {/* Footer Stats */}
      <div className="mt-4 pt-4 border-t border-gray-200">
        <div className="flex items-center justify-between text-xs text-gray-500 flex-wrap gap-2">
          <div className="flex items-center space-x-4">
            <span>{t('logs.total')}: {logs.length}</span>
            <span>{t('logs.shown')}: {filteredLogs.length}</span>
          </div>
          <div className="flex items-center space-x-2">
            <span>{t('header.manualUpdate')}</span>
          </div>
        </div>
      </div>
    </motion.div>
  )
}