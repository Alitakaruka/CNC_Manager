import React, { useEffect, useState, useRef } from 'react'
import { 
  Printer, 
  Box,
  Wifi, 
  Sunset,
  WifiOff, 
  Play, 
  Clock,
  Thermometer,
  Zap,
  Droplet,
  Layers,
  Settings
} from 'lucide-react'
import { useLocalization } from '../../../hooks/useLocalization.jsx'
import wsClient from '../../../hooks/WebSocketClient'

function getCncUniqueKey(cnc) {
  if (!cnc || typeof cnc !== 'object') return ''
  return cnc.uniqueKey || cnc.UniqueKey || ''
}

function mergeCncIntoList(prev, cnc) {
  const key = getCncUniqueKey(cnc)
  if (!key) return prev
  const index = prev.findIndex(c => getCncUniqueKey(c) === key)
  if (index !== -1) {
    const next = [...prev]
    next[index] = { ...next[index], ...cnc }
    return next
  }
  return [...prev, cnc]
}

/** Полный снимок таблицы (массив из нескольких станков) заменяет state; один объект или [один] — мерж по uniqueKey. */
function applyPrintersWsPayload(prev, data) {
  if (Array.isArray(data)) {
    if (data.length === 0) return prev
    if (data.length === 1) return mergeCncIntoList(prev, data[0])
    return data
  }
  if (data && typeof data === 'object' && getCncUniqueKey(data)) {
    return mergeCncIntoList(prev, data)
  }
  return prev
}

function PrintersTable({ SetNowPrinter, SetDetailsIsOpen }) {
  const [cncs, setCncs] = useState([])
  const [loading, setLoading] = useState(true)
  const { t } = useLocalization()
  const startedRef = useRef(false)

  // // // HTTP fallback для загрузки CNC станков
  const fetchCncs = async () => {
    try {
      const response = await fetch('/api/Printers', {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' }
      })
      
      if (response.ok) {
        const json = await response.json()
        if (Array.isArray(json)) {
          setCncs(json)
          setLoading(false)
        }
      }
    } catch (err) {
      console.log("Error loading CNC stations: " + err.message)
      setLoading(false)
    }
  }

  // // Запрос через WebSocket
  const requestCncsViaWS = async () => {
    try {
      if (wsClient.isConnected) {
        const data = await wsClient.request('GetMachines', {})
        if (Array.isArray(data)) {
          setCncs(data)
          setLoading(false)
          return true
        }
      }
    } catch (err) {
      console.log("WS request failed: " + err.message)
    }
    return false
  }

  // WebSocket handlers и начальная загрузка
  useEffect(() => {
    if (!startedRef.current) {
      startedRef.current = true
      wsClient.connect()
    }

    // Событие printers: полный массив станков (2+) заменяет таблицу; один объект или массив из одного — только мерж этой строки.
    const offPrinters = wsClient.on('printers', (data) => {
      setLoading(false)
      setCncs(prev => applyPrintersWsPayload(prev, data))
    })

    // Подписка на обновление отдельного CNC станка
    const offPrinterUpdate = wsClient.on('printerUpdate', (cnc) => {
      if (!getCncUniqueKey(cnc)) return
      setCncs(prev => mergeCncIntoList(prev, cnc))
    })

    // Запрос данных при подключении WebSocket
    const offOpen = wsClient.on('open', async () => {
      const success = await requestCncsViaWS()
      if (!success) {
        await fetchCncs()
      }
    })

    // Начальная загрузка данных
    const loadInitialData = async () => {
      setLoading(true)
      if (wsClient.isConnected) {
        const success = await requestCncsViaWS()
        if (!success) {
          await fetchCncs()
        }
      } else {
        await fetchCncs()
      }
    }
    loadInitialData()

    // Fallback polling только если WebSocket не подключен
    let interval
    const visibilityHandler = () => {
      clearInterval(interval)
      if (!wsClient.isConnected) {
        const ms = document.visibilityState === 'visible' ? 5000 : 10000
        interval = setInterval(() => {
          if (!wsClient.isConnected) {
            fetchCncs()
          }
        }, ms)
      }
    }
    
    document.addEventListener('visibilitychange', visibilityHandler)
    visibilityHandler()
    
    return () => {
      offPrinters()
      offPrinterUpdate()
      offOpen()
      document.removeEventListener('visibilitychange', visibilityHandler)
      clearInterval(interval)
    }
  }, [])

  const showCncDetails = (cnc) => {
    SetNowPrinter(cnc)
    SetDetailsIsOpen(cnc)
  }

  const getStatusIcon = (cnc) => {
    const isWorking = cnc.isWorking !== undefined ? cnc.isWorking : cnc.Flags?.Connected
    const executingTask = cnc.Flags?.ExecutingTask
    
    if (!isWorking) {
      return <WifiOff className="h-4 w-4 text-danger-500" />
    }
    if (executingTask) {
      return <Play className="h-4 w-4 text-success-500 animate-pulse" />
    }
    return <Wifi className="h-4 w-4 text-success-500" />
  }

  const getStatusText = (cnc) => {
    const isWorking = cnc.isWorking !== undefined ? cnc.isWorking : cnc.Flags?.Connected
    const executingTask = cnc.Flags?.ExecutingTask
    
    if (!isWorking) return t('printers.status.offline')
    if (executingTask) return t('printers.status.printing')
    return t('printers.status.ready')
  }

  const getStatusClass = (cnc) => {
    const isWorking = cnc.isWorking !== undefined ? cnc.isWorking : cnc.Flags?.Connected
    const executingTask =cnc.Flags?.ExecutingTask
    
    if (!isWorking) return 'status-offline'
    if (executingTask) return 'status-printing'
    return 'status-online'
  }

  const getPrinterIcon = (cnc) => {
    const printerType = (cnc.CncType || cnc.printerType || cnc.MACHINE_TYPE || '').toUpperCase()
    
    switch (printerType) {
      case 'FDM':
      case 'FDM_PRINTER':
        return <Box className="h-5 w-5 text-primary-600 dark:text-primary-400" />
      case 'LASER':
        return <Sunset className="h-5 w-5 text-primary-600 dark:text-primary-400" />
      case 'SLA':
      case 'SLA_PRINTER':
        return <Droplet className="h-5 w-5 text-primary-600 dark:text-primary-400" />
      case 'SLS':
      case 'SLS_PRINTER':
        return <Layers className="h-5 w-5 text-primary-600 dark:text-primary-400" />
      case 'MILLING':
        return <Settings className="h-5 w-5 text-primary-600 dark:text-primary-400" />
      default:
        return <Printer className="h-5 w-5 text-primary-600 dark:text-primary-400" />
    }
  }

  const formatTime = (seconds) => {
    if (seconds === 0) return '--'
    const hours = Math.floor(seconds / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    return `${hours}${t('time.hours')} ${minutes}${t('time.minutes')}`
  }

  if (loading) {
    return (
      <div className="card">
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
          <span className="ml-3 text-gray-600 dark:text-gray-300">{t('status.loading')} {t('navigation.printers').toLowerCase()}...</span>
        </div>
      </div>
    )
  }

  return (
    <div className="card">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center space-x-3">
          <Printer className="h-6 w-6 text-primary-600 dark:text-primary-400" />
          <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">{t('navigation.printers')}</h2>
          <span className="bg-primary-100 text-primary-800 text-xs font-medium px-2.5 py-0.5 rounded-full dark:bg-primary-900/40 dark:text-primary-300">
            {cncs.length}
          </span>
        </div>
          <div className="flex items-center space-x-2 text-sm text-gray-500 dark:text-gray-400">
            <div className="flex items-center space-x-1">
              <div className="w-2 h-2 bg-success-500 rounded-full"></div>
              <span>{t('header.online')}: {cncs.filter(c => c.isWorking !== undefined ? c.isWorking : c.Flags?.Connected).length}</span>
            </div>
            <div className="flex items-center space-x-1">
              <div className="w-2 h-2 bg-warning-500 rounded-full animate-pulse"></div>
              <span>{t('header.printing')}: {cncs.filter(c => c.executingTask !== undefined ? c.executingTask : c.Flags?.ExecutingTask).length}</span>
            </div>
          </div>
      </div>

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="table-header">
              <th className="px-4 py-3 text-left rounded-l-lg">{t('navigation.printers')}</th>
              <th className="px-4 py-3 text-left">{t('printers.details.type')}</th>
              <th className="px-4 py-3 text-left">{t('printers.details.status')}</th>
              <th className="px-4 py-3 text-left">Position (X/Y/Z)</th>
              <th className="px-4 py-3 text-left">{t('printers.details.progress')}</th>
              <th className="px-4 py-3 text-left rounded-r-lg">{t('printers.details.connection')}</th>
            </tr>
          </thead>
          <tbody>
            {cncs.map((cnc, index) => (
              <tr
                key={cnc.uniqueKey}
                className="table-row border-b border-gray-100 hover:bg-primary-50 transition-colors duration-200 cursor-pointer dark:border-gray-700 dark:hover:bg-gray-800"
                onClick={() => showCncDetails(cnc)}
              >
                <td className="px-4 py-4">
                  <div className="flex items-center space-x-3">
                    <div className="p-2 bg-primary-100 rounded-lg dark:bg-primary-900/40">
                      {getPrinterIcon(cnc)}
                    </div>
                    <div>
                      <div className="font-medium text-gray-900 dark:text-gray-100">{cnc.CNCName || cnc.TARGET_MACHINE_NAME || 'Unnamed'}</div>
                    </div>
                  </div>
                </td>
                <td className="px-4 py-4">
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200">
                    {cnc.CncType || cnc.printerType || cnc.MACHINE_TYPE || 'Unknown'}
                  </span>
                </td>
                <td className="px-4 py-4">
                  <div className="flex items-center space-x-2">
                    {getStatusIcon(cnc)}
                    <span className={getStatusClass(cnc)}>
                      {getStatusText(cnc)}
                    </span>
                  </div>
                </td>
                <td className="px-4 py-4">
                  <div className="text-sm text-gray-700 dark:text-gray-300">
                    {typeof cnc?.position?.X === 'number' || typeof cnc?.position?.Y === 'number' || typeof cnc?.position?.Z === 'number'
                      ? `${cnc.position?.X ?? '-'} / ${cnc.position?.Y ?? '-'} / ${cnc.position?.Z ?? '-'}`
                      : '--'}
                  </div>
                </td>
                <td className="px-4 py-4">
                  {cnc.executingTask ? (
                    <div className="space-y-1">
                      <div className="flex items-center justify-between text-sm">
                        <span>{cnc.progress || 0}%</span>
                        <Clock className="h-4 w-4 text-gray-400" />
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2 dark:bg-gray-700">
                        <div
                          className="bg-success-500 h-2 rounded-full transition-all duration-500"
                          style={{ width: `${cnc.progress || 0}%` }}
                        />
                      </div>
                      <div className="text-xs text-gray-500 dark:text-gray-400">
                        {formatTime(cnc.timeRemaining || 0)}
                      </div>
                    </div>
                  ) : (
                    <span className="text-gray-400">--</span>
                  )}
                </td>
                <td className="px-4 py-4">
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900/40 dark:text-blue-300">
                    {cnc.typeOfConnection || 'Unknown'}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Empty State */}
      {cncs.length === 0 && !loading && (
        <div className="text-center py-12">
          <Printer className="h-12 w-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2 dark:text-gray-100">
            {t('printers.empty.title')}
          </h3>
          <p className="text-gray-500 dark:text-gray-400">
            {t('printers.empty.subtitle')}
          </p>
        </div>
      )}
    </div>
  )
}
export default PrintersTable

