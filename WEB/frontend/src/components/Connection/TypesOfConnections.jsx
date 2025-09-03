import { motion } from 'framer-motion'
import { Monitor, Wifi, Usb } from 'lucide-react'

const ConnectionTypes = {
  COM: "COM",
  IP: "IP", 
  USB: "USB",
}

const PlaceHolders = {
  COM: "COM3",
  IP: "192.168.1.100:8080",
  USB: "USB001",
}

const connectionIcons = {
  COM: Monitor,
  IP: Wifi,
  USB: Usb
}

export default ConnectionTypes

export function ConnectionBuilder({ BuilderType, ConnectionRef, value, onChange }) {
  const IconComponent = connectionIcons[BuilderType]

  const handleInputChange = (e) => {
    const newValue = e.target.value
    ConnectionRef.current.connection = newValue
    if (onChange) {
      onChange(newValue)
    }
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      className="relative"
    >
      <div className="relative">
        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <IconComponent className="h-5 w-5 text-gray-400" />
        </div>
        <input
          type="text"
          className="input-field pl-10 pr-4 py-3 text-lg"
          placeholder={PlaceHolders[BuilderType]}
          value={value || ''}
          onChange={handleInputChange}
          autoFocus
        />
      </div>
      
      {/* Connection type indicator */}
      <div className="mt-2 flex items-center space-x-2">
        <div className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
          BuilderType === 'COM' ? 'bg-blue-100 text-blue-800' :
          BuilderType === 'IP' ? 'bg-green-100 text-green-800' :
          'bg-purple-100 text-purple-800'
        }`}>
          <IconComponent className="h-3 w-3 mr-1" />
          {BuilderType}
        </div>
        <span className="text-xs text-gray-500">
          {BuilderType === 'COM' && 'Последовательный порт'}
          {BuilderType === 'IP' && 'Сетевое подключение'}
          {BuilderType === 'USB' && 'USB подключение'}
        </span>
      </div>
    </motion.div>
  )
}