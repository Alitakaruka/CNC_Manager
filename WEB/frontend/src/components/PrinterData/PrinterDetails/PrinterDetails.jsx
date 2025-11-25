import { useState, useRef } from 'react'
import { motion } from 'framer-motion'
import {
  X,
  Printer,
  Thermometer,
  Zap,
  Box,
  Sunset,

  Play,
  Pause,
  Square,
  Upload,
  FileText,
  Settings,
  Gauge,
  Move,
  ArrowUp,
  ArrowDown,
  ArrowLeft,
  ArrowRight,
  RotateCcw,
  Home,
  Target,
  Droplet,
  Layers
} from 'lucide-react'
import MainPrinterData from '../MainPrinterData/MainPrinterData'
import PrinterGabarites from '../Gabarites/gabarites'
import PrinterCommands from '../Commands/PrinterCommands'
import toast from 'react-hot-toast'
import { useLocalization } from '../../../hooks/useLocalization.jsx'
import ConnectionHook from '../../../hooks/ConnectionHook'
import wsClient from '../../../hooks/WebSocketClient'

export default function Details({ PrinterData, SetDetailsIsOpen }) {
  const fileRef = useRef(null)
  const [isUploading, setIsUploading] = useState(false)
  const [isReconnecting, setIsReconnecting] = useState(false)
  const { t } = useLocalization()

  if (!PrinterData) {
    return null
  }

  const isWorking = PrinterData.isWorking !== undefined ? PrinterData.isWorking : PrinterData.Flags?.Connected
  const controlsDisabled = !isWorking

  const handleStartTask = async () => {
    if (controlsDisabled) {
      return
    }
    if (!fileRef.current?.files?.[0]) {
      toast.error(t('notifications.fileRequired'))
      return
    }

    if (!PrinterData.uniqueKey) {
      toast.error(t('notifications.printerNotSelected'))
      return
    }

    setIsUploading(true)

    try {
      const file = fileRef.current.files[0]
      
      // Try WebSocket first
    // Try WebSocket first
  if (wsClient.connected) {
  try {
    // Read file as text (G-code files are text)
    const fileData = await new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => resolve(e.target.result);
      reader.onerror = reject;
      reader.readAsText(file);
    });

    // --- Encode to Base64 ---
    const fileBase64 = btoa(unescape(encodeURIComponent(fileData)));
    // -------------------------

    const result = await wsClient.request('executeTask', {
      uniqueKey: PrinterData.uniqueKey,
      fileName: file.name,
      fileData: fileBase64   // <-- теперь строка Base64
    });

    toast.success(t('notifications.printStarted'));
    console.log('Ответ сервера:', result);
    setIsUploading(false);
    return;
  } catch (wsError) {
    console.warn('WebSocket executeTask failed, falling back to HTTP:', wsError);
    // Fall through to HTTP fallback
  }
}

      // HTTP fallback
      const formData = new FormData()
      formData.append('PrintFile', file)
      const query = `uniqueKey=${encodeURIComponent(PrinterData.uniqueKey)}`
      const response = await fetch(`/api/ExecuteTask?${query}`, {
        method: 'POST',
        body: formData
      })

      if (!response.ok) {
        throw new Error(t('notifications.printError'))
      }

      const result = await response.text()
      toast.success(t('notifications.printStarted'))
      console.log('Ответ сервера:', result)
    } catch (error) {
      toast.error(`${t('common.error')}: ${error.message}`)
    } finally {
      setIsUploading(false)
    }
  }

  const connectionType = `${PrinterData.typeOfConnection || PrinterData.TypeOfConnection || ''}`.trim()
  const connectionData = `${PrinterData.connectionData || PrinterData.ConnectionData || ''}`.trim()
  const canReconnect = Boolean(connectionType && connectionData)

  const handleReconnect = async () => {
    if (!canReconnect) {
      toast.error(t('notifications.printerNotSelected'))
      return
    }
    setIsReconnecting(true)
    try {
      await ConnectionHook(connectionType, connectionData)
      toast.success(t('header.connected'))
    } catch (e) {
      toast.error(`${t('common.error')}: ${e.message}`)
    } finally {
      setIsReconnecting(false)
    }
  }

  const getStatusColor = () => {
    const cnc = PrinterData
    const isWorking = cnc.isWorking !== undefined ? cnc.isWorking : cnc.Flags?.Connected
    const executingTask = cnc.Flags?.executingTask

    if (!isWorking) return 'text-danger-500'
    if (executingTask) return 'text-success-500'
    return 'text-primary-500'
  }

  const getStatusText = () => {
    const cnc = PrinterData
    const isWorking = cnc.isWorking !== undefined ? cnc.isWorking : cnc.Flags?.Connected
    const executingTask = cnc.Flags?.ExecutingTask

    if (!isWorking) return t('printers.status.offline')
    if (executingTask) return t('printers.status.printing')
    return t('printers.status.ready')
  }

  const getPrinterIcon = () => {
    const printerType = (PrinterData.CncType || PrinterData.printerType || PrinterData.MACHINE_TYPE || '').toUpperCase()

    switch (printerType) {
      case 'FDM':
      case 'FDM_PRINTER':
        return <Box className="h-6 w-6 text-primary-600 dark:text-primary-400" />
      case 'LASER':
        return <Sunset className="h-6 w-6 text-primary-600 dark:text-primary-400" />
      case 'SLA':
      case 'SLA_PRINTER':
        return <Droplet className="h-6 w-6 text-primary-600 dark:text-primary-400" />
      case 'SLS':
      case 'SLS_PRINTER':
        return <Layers className="h-6 w-6 text-primary-600 dark:text-primary-400" />
      case 'MILLING':
        return <Settings className="h-6 w-6 text-primary-600 dark:text-primary-400" />
      default:
        return <Printer className="h-6 w-6 text-primary-600 dark:text-primary-400" />
    }
  }

  return (
    <motion.div
      className="card h-fit"
      initial={{ opacity: 0, x: 100, scale: 0.95 }}
      animate={{ opacity: 1, x: 0, scale: 1 }}
      exit={{ opacity: 0, x: 100, scale: 0.95 }}
      transition={{ type: "spring", damping: 25, stiffness: 300 }}
    >
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center space-x-3">
          <div className="p-2 bg-primary-100 rounded-lg dark:bg-primary-900/40">
            {getPrinterIcon()}
          </div>
          <div>
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
              {PrinterData.CNCName || PrinterData.printerName || PrinterData.TARGET_MACHINE_NAME || 'Unnamed CNC'}
            </h3>
            <div className="flex items-center space-x-2">
              <span className={`text-sm font-medium ${getStatusColor()}`}>
                {getStatusText()}
              </span>
              <span className="text-xs text-gray-500 dark:text-gray-400">
                {PrinterData.uniqueKey || PrinterData.UniqueKey}
              </span>
            </div>
              <motion.button
            onClick={handleReconnect}
            className={`px-3 py-2 rounded-lg text-sm font-medium flex items-center space-x-2 ${controlsDisabled ? 'btn-primary' : 'btn-secondary'}`}
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            disabled={isReconnecting || !canReconnect}
            title={t('printers.details.reconnect')}
          >
            {isReconnecting ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                <span>{t('status.connecting')}</span>
              </>
            ) : (
              <>
                <RotateCcw className="h-4 w-4" />
                <span>{t('printers.details.reconnect')}</span>
              </>
            )}
          </motion.button>
          </div>
        </div>


        <div className="grix items-center space-x-2">
          <motion.button
            onClick={SetDetailsIsOpen}
            className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors duration-200 dark:hover:bg-gray-800"
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            title={t('common.close')}
          >
            <X className="h-5 w-5" />
          </motion.button>
        </div>
      </div>

      {/* Temperature Status */}
      {(PrinterData.CncType || PrinterData.printerType) === '3D PRINTER' && (
        <div className="grid grid-cols-2 gap-4 mb-6">
          <motion.div
            className="p-4 bg-orange-50 border border-orange-200 rounded-lg dark:bg-orange-900/20 dark:border-orange-800"
            whileHover={{ scale: 1.02 }}
          >
            <div className="flex items-center space-x-2 mb-2">
              <Thermometer className="h-5 w-5 text-orange-600" />
              <span className="text-sm font-medium text-orange-800 dark:text-orange-300">{t('printers.details.nozzle')}</span>
            </div>
            <div className="text-2xl font-bold text-orange-900 dark:text-orange-200">
              {PrinterData.nozzleTemp || PrinterData.NowTempNozzle || 0}°C
            </div>
          </motion.div>

          <motion.div
            className="p-4 bg-blue-50 border border-blue-200 rounded-lg dark:bg-blue-900/20 dark:border-blue-800"
            whileHover={{ scale: 1.02 }}
          >
            <div className="flex items-center space-x-2 mb-2">
              <Zap className="h-5 w-5 text-blue-600" />
              <span className="text-sm font-medium text-blue-800 dark:text-blue-300">{t('printers.details.bed')}</span>
            </div>
            <div className="text-2xl font-bold text-blue-900 dark:text-blue-200">
              {PrinterData.bedTemp || PrinterData.NowTempBed || 0}°C
            </div>
          </motion.div>
        </div>
      )}

      {/* CNC Job Controls */}
      <div className={`space-y-4 mb-6 ${controlsDisabled ? 'opacity-50 pointer-events-none' : ''}`}>
        <h4 className="text-sm font-semibold text-gray-700 flex items-center space-x-2 dark:text-gray-200">
          <Play className="h-4 w-4" />
          <span>{t('printers.controls.printControl')}</span>
        </h4>

        {/* File Upload */}
        <div className="space-y-3">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('printers.controls.gcodeFile')}
          </label>
          <div className="flex items-center space-x-3">
            <input
              type="file"
              ref={fileRef}
              accept=".gcode"
              className="hidden"
              id="fileInput"
              disabled={controlsDisabled}
            />
            <label
              htmlFor="fileInput"
              className={`flex-1 flex items-center justify-center px-4 py-2 border-2 border-dashed border-gray-300 rounded-lg transition-colors duration-200 cursor-pointer dark:border-gray-700 dark:hover:bg-gray-800 ${controlsDisabled ? 'cursor-not-allowed opacity-60' : 'hover:border-primary-400 hover:bg-primary-50'}`}
            >
              <Upload className="h-5 w-5 text-gray-400 mr-2" />
              <span className="text-sm text-gray-600 dark:text-gray-300">
                {fileRef.current?.files?.[0]?.name || t('printers.controls.selectFile')}
              </span>
            </label>
          </div>
        </div>

        {/* Start Task Button */}
        <motion.button
          onClick={handleStartTask}
          disabled={controlsDisabled || !fileRef.current?.files?.[0] || isUploading}
          className={`w-full py-3 px-4 rounded-lg font-medium transition-all duration-200 flex items-center justify-center space-x-2 ${controlsDisabled || !fileRef.current?.files?.[0] || isUploading
              ? 'bg-gray-300 text-gray-500 cursor-not-allowed dark:bg-gray-700 dark:text-gray-400'
              : 'btn-success'
            }`}
          whileHover={!controlsDisabled && fileRef.current?.files?.[0] && !isUploading ? { scale: 1.02 } : {}}
          whileTap={!controlsDisabled && fileRef.current?.files?.[0] && !isUploading ? { scale: 0.98 } : {}}
        >
          {isUploading ? (
            <>
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
              <span>{t('printers.controls.upload')}</span>
            </>
          ) : (
            <>
              <Play className="h-4 w-4" />
              <span>{t('printers.controls.startPrint')}</span>
            </>
          )}
        </motion.button>
      </div>

      {/* CNC Station Information */}
      <div className="space-y-4 mb-6">
        <h4 className="text-sm font-semibold text-gray-700 flex items-center space-x-2 dark:text-gray-200">
          <Settings className="h-4 w-4" />
          <span>{t('printers.details.title')}</span>
        </h4>

        <div className="grid grid-cols-3 gap-3 text-sm">
          <div>
            <span className="text-gray-500 dark:text-gray-400">{t('printers.details.type')}:</span>
            <div className="font-medium dark:text-gray-100">{PrinterData.CncType || PrinterData.printerType || PrinterData.MACHINE_TYPE || 'Unknown'}</div>
          </div>
          <div>
            <span className="text-gray-500 dark:text-gray-400">{t('printers.details.connection')}:</span>
            <div className="font-medium dark:text-gray-100">{PrinterData.typeOfConnection || PrinterData.ConnectionData || 'Unknown'}</div>
          </div>
          <div>
            <span className="text-gray-500 dark:text-gray-400">{t('printers.details.status')}:</span>
            <div className={`font-medium ${getStatusColor()}`}>
              {getStatusText()}
            </div>
          </div>
        </div>
      </div>

      {/* CNC Station Commands */}
      <div className="space-y-4">
        <PrinterCommands uniqueKey={PrinterData.uniqueKey} disabled={controlsDisabled} />
      </div>
    </motion.div>
  )
}