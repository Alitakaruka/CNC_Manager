import { useRef, useState } from 'react'
import { motion } from 'framer-motion'
import { 
  Move, 
  ArrowUp, 
  Minus,
  Plus,
  ArrowDown, 
  ArrowLeft, 
  ArrowRight, 
  RotateCcw, 
  Home, 
  Target,
  Send,
  Zap,
  Square
} from 'lucide-react'
import SendGCode from '../../../hooks/Gcode'
import toast from 'react-hot-toast'
import { useLocalization } from '../../../hooks/useLocalization.jsx'

export default function PrinterCommands({ uniqueKey }) {
  const [gcodeCommand, setGcodeCommand] = useState('')
  const [stepValue, setStepValue] = useState(10)
  const [isSending, setIsSending] = useState(false)
  const { t } = useLocalization()

  const handleGcodeSend = async () => {
    if (!gcodeCommand.trim()) {
      toast.error(t('notifications.gcodeRequired'))
      return
    }

    if (!uniqueKey) {
      toast.error(t('notifications.printerNotSelected'))
      return
    }

    setIsSending(true)

    try {
      const result = await SendGCode(gcodeCommand, uniqueKey)
      if (result && result.startsWith('Error:')) {
        toast.error(result)
      } else {
        toast.success(t('printers.controls.send'))
        setGcodeCommand('')
      }
    } catch (error) {
      toast.error(`${t('common.error')}: ${error.message}`)
    } finally {
      setIsSending(false)
    }
  }

  const moveAxis = async (axis, direction = 1) => {
    if (!uniqueKey) {
      toast.error(t('notifications.printerNotSelected'))
      return
    }

    // Формируем команду движения для CNC станка

    const value = direction * stepValue
    const gcode = `G1 ${axis}${value.toFixed(2)}`
    
    setIsSending(true)
    try {
      await SendGCode(gcode, uniqueKey)
      toast.success(t('printers.commands.movementCompleted'))
    } catch (error) {
      toast.error(`${t('common.error')}: ${error.message}`)
    } finally {
      setIsSending(false)
    }
  }

  const sendQuickCommand = async (command) => {
    if (!uniqueKey) {
      toast.error(t('notifications.printerNotSelected'))
      return
    }

    // Отправка команды на CNC станок
    setIsSending(true)
    try {
      await SendGCode(command, uniqueKey)
      toast.success(t('printers.commands.commandExecuted'))
    } catch (error) {
      toast.error(`${t('common.error')}: ${error.message}`)
    } finally {
      setIsSending(false)
    }
  }

  const quickCommands = [
    { command: 'G28', label: t('printers.commands.homeAll'), icon: Home, color: 'bg-blue-500 hover:bg-blue-600' },
    { command: 'M84', label: t('printers.commands.disableMotors'), icon: Square, color: 'bg-red-500 hover:bg-red-600' },
  ]

  return (
    <div className="space-y-6">
      {/* Quick Commands */}
      <div className="space-y-3">
        <h5 className="text-sm font-medium text-gray-700 flex items-center space-x-2">
          <Zap className="h-4 w-4" />
          <span>{t('printers.controls.quickCommands')}</span>
        </h5>
        
        <div className="grid grid-cols-2 gap-2">
          {quickCommands.map((cmd, index) => (
            <motion.button
              key={cmd.command}
              onClick={() => sendQuickCommand(cmd.command)}
              disabled={isSending}
              className={`${cmd.color} text-white py-2 px-3 rounded-lg text-sm font-medium transition-colors duration-200 flex items-center justify-center space-x-2`}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1 }}
            >
              <cmd.icon className="h-4 w-4" />
              <span>{cmd.label}</span>
            </motion.button>
          ))}
        </div>
      </div>

      {/* Custom G-code */}
      <div className="space-y-3">
        <h5 className="text-sm font-medium text-gray-700 flex items-center space-x-2">
          <Send className="h-4 w-4" />
          <span>{t('printers.controls.customGcode')}</span>
        </h5>
        
        <div className="flex space-x-2">
          <input
            type="text"
            value={gcodeCommand}
            onChange={(e) => setGcodeCommand(e.target.value)}
            placeholder="G28"
            className="flex-1 input-field text-sm"
            onKeyDown={(e) => e.key === 'Enter' && handleGcodeSend()}
          />
          <motion.button
            onClick={handleGcodeSend}
            disabled={!gcodeCommand.trim() || isSending}
            className={`px-4 py-2 rounded-lg font-medium transition-all duration-200 flex items-center space-x-2 ${
              !gcodeCommand.trim() || isSending
                ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                : 'btn-primary'
            }`}
            whileHover={gcodeCommand.trim() && !isSending ? { scale: 1.05 } : {}}
            whileTap={gcodeCommand.trim() && !isSending ? { scale: 0.95 } : {}}
          >
            {isSending ? (
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
            ) : (
              <Send className="h-4 w-4" />
            )}
            <span className="hidden sm:inline">{t('printers.controls.send')}</span>
          </motion.button>
        </div>
      </div>

      {/* Loading State */}
      {isSending && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="flex items-center justify-center py-2 text-sm text-gray-600"
        >
          <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-primary-600 mr-2"></div>
          {t('printers.controls.sending')}
        </motion.div>
      )}
    </div>
  )
}