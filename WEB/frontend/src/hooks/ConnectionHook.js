
export default async function useConnectPrinter(TypeOfConnection = "", ConnectionData = "") {
  console.log("TypeOfConnection:", TypeOfConnection)
  console.log("ConnectionData:", ConnectionData)
  
  if (!TypeOfConnection || !ConnectionData) {
    throw new Error("Данные подключения не могут быть пустыми!")
  }

  // Валидация COM порта
  if (TypeOfConnection === "COM" && !/^COM\d+$/i.test(ConnectionData)) {
    throw new Error("Неверный формат COM порта (например, COM3)")
  }

  // Валидация IP адреса
  if (TypeOfConnection === "IP") {
    const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(:\d+)?$/
    if (!ipRegex.test(ConnectionData)) {
      throw new Error("Неверный формат IP адреса")
    }
  }

  // Try WS first
  try {
    const { wsClient } = await import('./WebSocketClient')
    if (wsClient?.connected) {
      const result = await wsClient.request('connect', { TypeOfConnection, ConnectionData })
      return typeof result === 'string' ? result : 'OK'
    }
  } catch {}

  // HTTP fallback
  try {
    const response = await fetch("/connect", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ TypeOfConnection, ConnectionData })
    })
    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(errorText || "Ошибка подключения")
    }
    const result = await response.text()
    console.log("Подключение успешно:", result)
    return result
  } catch (error) {
    console.error("Connection error:", error)
    throw new Error(`Ошибка подключения: ${error.message}`)
  }
}
