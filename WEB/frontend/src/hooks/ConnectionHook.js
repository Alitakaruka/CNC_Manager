
async function ensureWsReady(wsClient, timeoutMs = 5000) {
  if (!wsClient) {
    throw new Error("WS client недоступен")
  }
  if (wsClient.connected) {
    return
  }

  try {
    wsClient.connect()
  } catch (err) {
    throw new Error(err?.message || "Не удалось инициировать WS соединение")
  }

  if (wsClient.connected) {
    return
  }

  await new Promise((resolve, reject) => {
    let settled = false
    const cleanup = () => {
      settled = true
      clearTimeout(timer)
      offOpen?.()
      offClose?.()
    }

    const handleOpen = () => {
      if (settled) return
      cleanup()
      resolve()
    }

    const handleClose = () => {
      if (settled) return
      cleanup()
      reject(new Error("WS connection closed"))
    }

    const offOpen = wsClient.on('open', handleOpen)
    const offClose = wsClient.on('close', handleClose)

    const timer = setTimeout(() => {
      if (settled) return
      cleanup()
      reject(new Error("WS connection timeout"))
    }, timeoutMs)
  })
}

export default async function ConnectCNC(TypeOfConnection = "", ConnectionData = "") {
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
    await ensureWsReady(wsClient)

    const result = await wsClient.request('connect', { TypeOfConnection, ConnectionData })
    console.log(result)

    return typeof result === 'string' ? result : 'OK'
  } catch (error) {
    console.error("WS connection error:", error)
    throw new Error(error?.message || "Ошибка подключения через WebSocket")
  }
}

export async function ReconnectCNC(UniqueKey = "") {
  console.log("UniqueKey:", UniqueKey)

  if (UniqueKey === ""){
      throw new Error("Unique key is empty!")
  }

  // Try WS first
  try {
    const { wsClient } = await import('./WebSocketClient')
    await ensureWsReady(wsClient)

    const result = await wsClient.request('reconnect', {UniqueKey})
    console.log(result)

    return typeof result === 'string' ? result : 'OK'
  } catch (error) {
    console.error("WS connection error:", error)
    throw new Error(error?.message || "Ошибка подключения через WebSocket")
  }
}
