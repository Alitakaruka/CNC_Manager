import wsClient from './WebSocketClient'

export default async function SendGCode(GCodeStr, UniqueKey) {
  if (!GCodeStr || !UniqueKey || UniqueKey === "Unknown") {
    throw new Error("Неверные параметры: G-code или ключ CNC станка отсутствуют")
  }

  // Try WS command
  if (wsClient.connected) {
    try {
      const res = await wsClient.request('command', { gcode: GCodeStr, uniqueKey: UniqueKey })
      return typeof res === 'string' ? res : 'OK'
    } catch (e) {
      // fallthrough to HTTP
    }
  }

  // HTTP fallback
  try {
    const query = `GCode=${encodeURIComponent(GCodeStr)}&uniqueKey=${encodeURIComponent(UniqueKey)}`
    const response = await fetch(`/api/sendGCode?${query}`, { method: "POST", headers: { 'Content-Type': 'application/json' } })
    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`Ошибка сервера: ${errorText}`)
    }
    const result = await response.text()
    return result
  } catch (error) {
    console.error('G-code error:', error)
    throw new Error(`Ошибка отправки G-code: ${error.message}`)
  }
}

export function moveAxis(axis, stepValue, direction = 1) {
  if (isNaN(stepValue) || stepValue <= 0) {
    throw new Error("Шаг должен быть положительным числом")
  }
  const value = direction * stepValue
  const gcode = `G91\nG1 ${axis}${value.toFixed(2)} F3000\nG90`
  return gcode
}

export function getHomeCommand(axis = null) {
  if (axis) return `G28 ${axis.toUpperCase()}0`
  return 'G28'
}

export function getDisableMotorsCommand() { return 'M84' }

export function getTemperatureCommand(nozzle = null, bed = null) {
  let commands = []
  if (nozzle !== null) commands.push(`M104 S${nozzle}`)
  if (bed !== null) commands.push(`M140 S${bed}`)
  return commands.join('\n')
}
